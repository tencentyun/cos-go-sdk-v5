package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 本示例演示 rapid bucket 的 session 鉴权用法。
//
// 与普通桶不同, rapid bucket的数据面请求需要先通过 CreateSession 换取临时 session
// 凭证再签名。cos.RapidCredentialTransport 把这一步自动化: 你只需提供一个底层
// 基础凭证 transport(如 AuthorizationTransport / CVMCredentialTransport /
// StsCredentialTransport), 它会在对rapid bucket发数据面请求时自动 CreateSession
// 并缓存/复用 session 凭证, 底层凭证变更后自动重签。
//
// cos.NewRapidCredentialTransport 会组装 session 鉴权与基础凭证签名。内部网关负载
// 均衡（按需刷新 DNS 快照、逐请求选择网关 IP、失败节点暂时冷却）由 base 自身按请求
// 域名自动判断是否介入, 无需额外指定。
//
// 运行前请设置环境变量 SECRETID / SECRETKEY, 并把下面的 BucketURL 换成你自己的
// rapid bucket域名(形如 http://<bucket>-x--<appid>.cosrapid.<region>.myqcloud.com)。
// 网关负载均衡只在 HTTP, 无显式端口, URL 目标不是 IP 且未配置代理时介入;
// 不满足条件时会原样透传, 不会阻断请求。
func main() {
	u, _ := url.Parse("http://example-x--1259654469.cosrapid.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}

	// 最底层的网络传输层, 负责 TCP 连接池与网络 I/O。生产环境建议显式设置非零的
	// IdleConnTimeout: DNS 摘除节点后, 空闲旧连接需要由 http.Transport 定时回收。
	transport := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	// 基础凭证 transport, 用于签 CreateSession, 网络层跟非 rapid 场景写法完全一样
	// 直接设置。也可换成 CVMCredentialTransport / CredentialTransport / StsCredentialTransport。
	base := &cos.AuthorizationTransport{
		SecretID:  os.Getenv("SECRETID"),
		SecretKey: os.Getenv("SECRETKEY"),
		Transport: transport,
	}

	// 关键: 用 NewRapidCredentialTransport 包裹底层凭证, 启用 session 鉴权; base 内部
	// 判断到目标是 rapid bucket 时会自动在 transport 外面套一层内部网关负载均衡,
	// 并保证签名层位于选点层之前。CreateSession 与数据面请求会经过同一套网关选点逻辑。
	// 如果底层 transport 的自定义 DialContext 或 resolver 覆盖了目标地址, 以该网络层的
	// 实际拨号结果为准; 网关负载均衡的节点分布, 冷却和 DNS 扩缩容语义不再适用。
	c := cos.NewClient(b, &http.Client{
		Transport: cos.NewRapidCredentialTransport(base),
		Timeout:   30 * time.Second,
	})

	ctx := context.Background()

	// 1) 自动 session 的 PutObject / GetObject。
	//    首次对rapid bucket发数据面请求时, RapidCredentialTransport 会自动
	//    CreateSession 拿临时凭证再签名。
	key := "rapid-example/hello.txt"
	payload := []byte("hello rapid bucket")
	if _, err := c.Object.Put(ctx, key, bytes.NewReader(payload), nil); err != nil {
		fmt.Printf("PutObject error: %v\n", err)
		return
	}
	fmt.Println("PutObject OK (session 自动创建并签名)")

	resp, err := c.Object.Get(ctx, key, nil)
	if err != nil {
		fmt.Printf("GetObject error: %v\n", err)
		return
	}
	got, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("GetObject OK: %s\n", got)

	// 2) RenameObject: rapid bucket内原子重命名对象。
	//    name 为目标 key, 第二个参数为源 key(内部会补前导 '/' 并做 URL 编码)。
	dstKey := "rapid-example/hello-renamed.txt"
	if _, err := c.Object.RenameObject(ctx, dstKey, key, nil); err != nil {
		fmt.Printf("RenameObject error: %v\n", err)
		return
	}
	fmt.Printf("RenameObject OK: %s -> %s\n", key, dstKey)

	// 如需防覆盖, 可通过 ObjectRenameOptions.XCosForbidOverwrite 传入 "true":
	//   c.Object.RenameObject(ctx, dstKey, key, &cos.ObjectRenameOptions{XCosForbidOverwrite: "true"})

	// 3) GetSessionPresignedURL: 为rapid bucket对象签出带 session 凭证的预签名 URL。
	//    expired 决定 URL 的有效期(受 session 临时凭证寿命约束)。
	presigned, err := c.Object.GetSessionPresignedURL(ctx, http.MethodGet, dstKey, 10*time.Minute, nil)
	if err != nil {
		fmt.Printf("GetSessionPresignedURL error: %v\n", err)
		return
	}
	fmt.Printf("Presigned URL: %s\n", presigned.String())

	// 该预签名 URL 自带 session 凭证与 token(全在 query 里), 可直接用裸 http.Client 访问。
	pr, err := http.Get(presigned.String())
	if err != nil {
		fmt.Printf("get via presigned error: %v\n", err)
		return
	}
	pb, _ := ioutil.ReadAll(pr.Body)
	pr.Body.Close()
	fmt.Printf("GET via presigned URL: status=%d body=%s\n", pr.StatusCode, pb)
}
