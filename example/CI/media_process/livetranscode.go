package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

type URLToken struct {
	SessionToken string `url:"x-cos-security-token,omitempty" header:"-"`
}

type JwtTokens struct {
	// base info
	Type     string `json:"Type"`
	AppId    string `json:"AppId"`
	BucketId string `json:"BucketId"`
	Object   string `json:"Object"`
	Issuer   string `json:"Issuer"`
	// time info
	IssuedTimeStamp int64 `json:"IssuedTimeStamp"`
	ExpireTimeStamp int64 `json:"ExpireTimeStamp"`
	// other info
	Random int64 `json:"Random"`
	// times info
	UsageLimit int `json:"UsageLimit"`
	// secret info
	ProtectSchema     string `json:"ProtectSchema"`
	PublicKey         string `json:"PublicKey"`
	ProtectContentKey int    `json:"ProtectContentKey"`
}

func (token JwtTokens) Valid() error {
	return nil
}

// 生成jwt
func GenerateToken(appId string, bucketId string, objectKey string, secret []byte) (string, error) {
	t := time.Now()
	now := t.Unix()
	payLoad := JwtTokens{
		// 固定为 CosCiToken， 必填参数
		Type: "CosCiToken",
		// app id，必填参数
		AppId: appId,
		// 播放文件所在的BucketId， 必填参数
		BucketId: bucketId,
		// 播放文件名
		Object: url.QueryEscape(objectKey),
		// 固定为client，必填参数
		Issuer: "client",
		// token颁发时间戳，必填参数
		IssuedTimeStamp: now,
		// token过期时间戳，非必填参数，默认1天过期
		ExpireTimeStamp: t.Add(time.Hour * 24 * 6).Unix(),
		// token使用次数限制，非必填参数，默认限制100次
		UsageLimit: 20,
		// 保护模式，填写为 rsa1024 ，则表示使用 RSA 非对称加密的方式保护，公私钥对长度为 1024 bit
		ProtectSchema: "rsa1024",
		// 公钥。1024 bit 的 RSA 公钥，需使用 Base64 进行编码
		PublicKey: "xxx",
		// 是否加密解密密钥（播放时解密ts视频流的密钥），1表示对解密密钥加密，0表示不对解密密钥加密。
		ProtectContentKey: 0,
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payLoad)

	//使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString(secret)
}

// 验证环境url
func GetCIDomainURL(tak string, tsk string, token *URLToken, appId string, bucketId string, region string, objectKey string, playkey []byte) {
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
		return
	}
	// 生成token
	generateToken, _ := GenerateToken(appId, bucketId, objectKey, playkey)
	resultUrl := presignedURL.String() + "&tokenType=JwtToken&expires=3600&object=" + url.QueryEscape(objectKey) + "&token=" + generateToken
	fmt.Println(resultUrl)
}

// cos环境url
func GetCOSDomainURL(tak string, tsk string, token *URLToken, appId string, bucketId string, region string, objectKey string, playkey []byte) {
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

	// 获取预签名
	presignedURL, err := c.Object.GetPresignedURL3(ctx, http.MethodGet, objectKey, time.Hour, token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// 生成token
	generateToken, _ := GenerateToken(appId, bucketId, objectKey, playkey)
	resultUrl := presignedURL.String() + "&ci-process=pm3u8&expires=43200&&tokenType=JwtToken&token=" + generateToken
	fmt.Println(resultUrl)
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
		},
	}
	createJobRes, _, err := c.CI.CreateGeneratePlayListJob(context.Background(), createJobOpt)
	log_status(err)
	fmt.Printf("%+v\n", createJobRes.JobsDetail)
}

func main() {
	InvokeGeneratePlayListJob()
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
	// 替换为自己播放密钥，控制台可以查询
	var secret = []byte("aaaaaaaaaaa")
	GetCIDomainURL(tak, tsk, token, appId, bucketId, region, objectKey, secret)
	// GetCOSDomainURL(tak, tsk, token, appId, bucketId, region, objectKey, secret)
}
