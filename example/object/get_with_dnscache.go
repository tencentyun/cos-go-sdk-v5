package main

import (
	"context"
	"fmt"
	"github.com/rs/dnscache"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func log_status(err error) {
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
	u, _ := url.Parse("http://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	r := &dnscache.Resolver{}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   false,
				// dns 缓存, 详见库 https://github.com/rs/dnscache
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
						host, port, err := net.SplitHostPort(addr)
						if err != nil {
							return nil, err
						}
						ips, err := r.LookupHost(ctx, host)
						if err != nil {
							return nil, err
						}
						for _, ip := range ips {
							var dialer net.Dialer
							conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
							if err == nil {
								break
							}
						}
						return
					},
				},
			},
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 定期刷新缓存
	go func(ctx context.Context) {
		// 设置 dns 缓存周期为 1800 秒
		timeTickerChan := time.Tick(time.Second * 1800)
		for {
			select {
			case <-ctx.Done():
				return
			case <-timeTickerChan:
				// 刷新缓存，会重新发起 dns 请求，并清理上个周期以来未被使用的记录。
				r.Refresh(true)
			}
		}
	}(ctx)
	name := "test"
	for i := 0; i < 100; i++ {
		_, err := c.Object.Get(context.Background(), name, nil)
		log_status(err)
		time.Sleep(time.Second)
	}
}
