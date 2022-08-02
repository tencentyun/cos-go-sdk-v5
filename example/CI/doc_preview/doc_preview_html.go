package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
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
	u, _ := url.Parse("https://test-1234567890.cos.ap-chongqing.myqcloud.com")
	cu, _ := url.Parse("https://test-1234567890.ci.ap-chongqing.myqcloud.com")
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
				ResponseBody:   false,
			},
		},
	})

	key := "input/doc_preview.ppt"
	opt := &cos.DocPreviewHTMLOptions{
		DstType:  "html",
		SrcType:  "ppt",
		Copyable: "1",
		HtmlParams: &cos.HtmlParams{
			CommonOptions: &cos.HtmlCommonParams{
				IsShowTopArea: false,
			},
			PptOptions: &cos.HtmlPptParams{
				IsShowBottomStatusBar: true,
			},
		},
		Htmlwaterword:  "5pWw5o2u5LiH6LGhLeaWh+aho+mihOiniA==",
		Htmlfillstyle:  "cmdiYSgxMDIsMjA0LDI1NSwwLjMp", // rgba(102,204,255,0.3)
		Htmlfront:      "Ym9sZCAyNXB4IFNlcmlm",         // bold 25px Serif
		Htmlrotate:     "315",
		Htmlhorizontal: "50",
		Htmlvertical:   "100",
	}
	resp, err := c.CI.DocPreviewHTML(context.Background(), key, opt)
	log_status(err)
	fd, _ := os.OpenFile("doc_preview.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	io.Copy(fd, resp.Body)
	fd.Close()
}
