package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func logStatus(err error) {
	if err == nil {
		return
	}
	if cos.IsNotFoundError(err) {
		// WARN
		fmt.Println("WARN: Resource is not existed")
	} else if e, ok := cos.IsCOSError(err); ok {
		fmt.Printf("ERROR: Code: %v\n", e.Code)
		fmt.Printf("ERROR: Message: %v\n", e.Message)
		fmt.Printf("ERROR: Resource: %v\n", e.Resource)
		fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
		// ERROR
	} else {
		fmt.Printf("ERROR: %v\n", err)
		// ERROR
	}
}

func main() {
	u, _ := url.Parse("https://bj-1259654469.cos.ap-beijing.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	domain := "www.qq.com"
	domainPemPath := domain + ".pem"
	domainKeyPath := domain + ".key"

	fd, err := os.Open(domainPemPath)
	if err != nil {
		panic(err)
	}
	pem, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	fd.Close()
	fd, err = os.Open(domainKeyPath)
	if err != nil {
		panic(err)
	}
	key, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	fd.Close()

	opt := &cos.BucketPutDomainCertificateOptions{
		CertificateInfo: &cos.BucketDomainCertificateInfo{
			CertType: "CustomCert",
			CustomCert: &cos.BucketDomainCustomCert{
				Cert:       string(pem),
				PrivateKey: string(key),
			},
		},
		DomainList: []string{
			domain,
		},
	}

	_, err = c.Bucket.PutDomainCertificate(context.Background(), opt)
	logStatus(err)

	gopt := &cos.BucketGetDomainCertificateOptions{
		DomainName: domain,
	}
	res, _, err := c.Bucket.GetDomainCertificate(context.Background(), gopt)
	logStatus(err)
	fmt.Printf("%+v\n", res)

	dopt := &cos.BucketDeleteDomainCertificateOptions{
		DomainName: domain,
	}
	_, err = c.Bucket.DeleteDomainCertificate(context.Background(), dopt)
	logStatus(err)

}
