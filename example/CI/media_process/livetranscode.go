package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

type URLToken struct {
	SessionToken string `url:"x-cos-security-token,omitempty" header:"-"`
}

// 生成jwt
func GenerateToken(appId string, bucketId string, objectKey string, secret []byte) (string, error) {
	t := time.Now()
	now := t.Unix()
	payLoad := jwt.MapClaims{
		// 固定为 CosCiToken， 必填参数
		"Type": "CosCiToken",
		// app id，必填参数
		"AppId": appId,
		// 播放文件所在的BucketId， 必填参数
		"BucketId": bucketId,
		// 播放文件名
		"Object": url.QueryEscape(objectKey),
		// 固定为client，必填参数
		"Issuer": "client",
		// token颁发时间戳，必填参数
		"IssuedTimeStamp": now,
		// token过期时间戳，非必填参数，默认1天过期
		"ExpireTimeStamp": t.Add(time.Hour * 24 * 6).Unix(),
		// token使用次数限制，非必填参数，默认限制100次
		"UsageLimit": 20,
		// 保护模式，填写为 rsa1024 ，则表示使用 RSA 非对称加密的方式保护，公私钥对长度为 1024 bit
		"ProtectSchema": "rsa1024",
		// 公钥。1024 bit 的 RSA 公钥，需使用 Base64 进行编码
		"PublicKey": "xxx",
		// 是否加密解密密钥（播放时解密ts视频流的密钥），1表示对解密密钥加密，0表示不对解密密钥加密。
		"ProtectContentKey": 0,
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payLoad)

	//使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString(secret)
}

// CI验证环境
func GetCIDomainVideoEncryptionURL(tak string, tsk string, token *URLToken, bucketId string, region string, objectKey string, jwtToken string) string {
	// 固定为getplaylist
	name := "getplaylist"

	u, _ := url.Parse("https://" + bucketId + ".ci." + region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     tak,
			SecretKey:    tsk,
			SessionToken: token.SessionToken,
		},
	})
	ctx := context.Background()

	// 获取预签名
	presignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodGet, name, tak, tsk, time.Hour, token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	resultUrl := presignedURL.String() + "&tokenType=JwtToken&expires=3600&object=" + url.QueryEscape(objectKey) + "&token=" + jwtToken
	return resultUrl
}

// COS环境
func GetCOSDomainVideoEncryptionURL(tak string, tsk string, token *URLToken, bucketId string, region string, objectKey string, jwtToken string) string {
	u, _ := url.Parse("https://" + bucketId + ".cos." + region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     tak,
			SecretKey:    tsk,
			SessionToken: token.SessionToken,
		},
	})
	ctx := context.Background()

	opt := &cos.PresignedURLOptions{
		Query:  &url.Values{},
		Header: &http.Header{},
	}
	opt.Query.Add("ci-process", "getplaylist")
	opt.Query.Add("signType", "cos")
	opt.Query.Add("expires", "43200")
	// opt.Query.Add("exper", "30") 试看时长
	opt.Query.Add("tokenType", "JwtToken")
	opt.Query.Add("token", jwtToken)

	var signHost bool = true
	// 获取预签名
	presignedURL, err := c.Object.GetPresignedURL2(ctx, http.MethodGet, objectKey, 10*time.Hour, opt, signHost)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}

	resultUrl := presignedURL.String()
	return resultUrl
}

// COS环境
func GetCOSDomainURL(tak string, tsk string, token *URLToken, appId string, bucketId string, region string, objectKey string) string {
	u, _ := url.Parse("https://" + bucketId + ".cos." + region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     tak,
			SecretKey:    tsk,
			SessionToken: token.SessionToken,
		},
	})
	ctx := context.Background()
	opt := &cos.PresignedURLOptions{
		Query:  &url.Values{},
		Header: &http.Header{},
	}
	opt.Query.Add("ci-process", "getplaylist")
	opt.Query.Add("signType", "cos")
	opt.Query.Add("expires", "43200")
	// opt.Query.Add("exper", "30") 试看时长
	var signHost bool = true
	// 获取预签名
	presignedURL, err := c.Object.GetPresignedURL2(ctx, http.MethodGet, objectKey, time.Hour, opt, signHost)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	resultUrl := presignedURL.String()
	return resultUrl
}

// CDN域名
func GetCDNDomainVideoEncryptionURL(cdn string, objectKey string, jwtToken string) string {
	url := cdn + "/" + objectKey
	resultUrl := url + "?ci-process=getplaylist&signType=no&expires=43200&&tokenType=JwtToken&token=" + jwtToken
	return resultUrl
}

// CDN域名
func GetCDNDomainURL(cdn string, objectKey string) string {
	url := cdn + "/" + objectKey
	resultUrl := url + "?ci-process=getplaylist&signType=no"
	return resultUrl
}

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

func getClient(cosHost, ciHost string) *cos.Client {
	u, _ := url.Parse(cosHost)
	cu, _ := url.Parse(ciHost)
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader: true,
				// Notice when put a large file and set need the request body, might happend out of memory error.
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	return c
}

func InvokeGeneratePlayListJob() {
	// 函数内替换为自己桶
	c := getClient("https://test-1234567890.cos.ap-chongqing.myqcloud.com", "https://test-1234567890.ci.ap-chongqing.myqcloud.com")
	createJobOpt := &cos.CreateGeneratePlayListJobOptions{
		Tag: "GeneratePlayList",
		Input: &cos.JobInput{
			Object: "a.mp4",
			// Vod: &cos.VodInfo{
			// 	FileId: "243791581857019308",
			// },
		},
		Operation: &cos.GeneratePlayListJobOperation{
			Output: &cos.JobOutput{
				Region: "ap-chongqing",
				Bucket: "test-1250000000",
				Object: "live/a.m3u8",
			},
			Transcode: &cos.LiveTanscode{
				Video: &cos.LiveTanscodeVideo{
					Codec:   "H.264",
					Width:   "960", // 设置480、720、960、1080
					Bitrate: "2000",
					Maxrate: "5000",
					Fps:     "30",
				},
				Container: &cos.Container{
					Format: "hls",
					ClipConfig: &cos.ClipConfig{
						Duration: "5",
					},
				},
				TransConfig: &cos.LiveTanscodeTransConfig{
					HlsEncrypt: &cos.HlsEncrypt{
						IsHlsEncrypt: true,
					},
					InitialClipNum: "2",
					CosTag:         "a=a&b=b",
				},
			},
			Watermark: []cos.Watermark{
				{
					Type:    "Text",
					LocMode: "Absolute",
					Dx:      "640",
					Pos:     "TopLeft",
					Text: &cos.Text{
						Text:         "helloworld",
						FontSize:     "25",
						FontType:     "simfang.ttf",
						FontColor:    "0xff0000",
						Transparency: "100",
					},
				},
			},
		},
	}
	createJobRes, _, err := c.CI.CreateGeneratePlayListJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

func main() {
	// InvokeGeneratePlayListJob()
	// 替换成您的密钥
	tak := os.Getenv("COS_SECRETID")
	tsk := os.Getenv("COS_SECRETKEY")
	token := &URLToken{
		SessionToken: "",
	}
	// 替换成您的桶名称
	appId := "1250000000"
	// 替换成您的桶名称
	bucketId := "test-1250000000"
	// 替换成您桶所在的region
	region := "ap-chongqing"
	// 替换成您需要播放的视频名称
	objectKey := "live/a.m3u8"

	playUrl := ""

	playUrl = GetCOSDomainURL(tak, tsk, token, appId, bucketId, region, objectKey)
	fmt.Println(playUrl)
	// 替换为自己cdn域名
	cdn := "http://abc.cdn.com"
	playUrl = GetCDNDomainURL(cdn, objectKey)
	fmt.Println(playUrl)
	// // 替换为自己播放密钥，控制台可以查询
	// var playkey = []byte("aaaaaaaaaaa")
	// // 生成token
	// jwtToken, _ := GenerateToken(appId, bucketId, objectKey, playkey)
	// playUrl = GetCOSDomainVideoEncryptionURL(tak, tsk, token, bucketId, region, objectKey, jwtToken)
	// fmt.Println(playUrl)
	// playUrl = GetCDNDomainVideoEncryptionURL(cdn, objectKey, jwtToken)
	// fmt.Println(playUrl)
}
