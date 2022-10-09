package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"

	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

func handler(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	ciphertext := m.Get("Ciphertext")
	kmsRegion := m.Get("KMSRegion")
	if len(ciphertext) == 0 || len(kmsRegion) == 0 {
		fmt.Fprint(w, "Ciphertext or KMSRegion is empty")
		return
	}

	c, err := kms.NewClientWithSecretId(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), kmsRegion)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	req := kms.NewDecryptRequest()
	req.CiphertextBlob = &ciphertext
	rsp, err := c.Decrypt(req)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	s, err := base64.StdEncoding.DecodeString(*rsp.Response.Plaintext)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	fmt.Fprint(w, string(s))
}

func GetLocalIp() string {
	localIP := "127.0.0.1"
	addrSlice, err := net.InterfaceAddrs()
	if nil == err {
		for _, addr := range addrSlice {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if nil != ipnet.IP.To4() {
					localIP = ipnet.IP.String()
					break
				}
			}
		}
	}
	return localIP
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(GetLocalIp()+":8082", nil)
}
