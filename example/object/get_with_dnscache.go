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
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("http://test-1259654469.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	r := &dnscache.Resolver{}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID:  os.Getenv("COS_SECRETID"),
			// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv("COS_SECRETKEY"),
			// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
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
