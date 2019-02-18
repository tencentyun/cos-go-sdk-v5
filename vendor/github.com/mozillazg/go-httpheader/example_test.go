package httpheader_test

import (
	"fmt"
	"net/http"

	"github.com/mozillazg/go-httpheader"
)

func ExampleHeader() {
	type Options struct {
		ContentType  string `header:"Content-Type"`
		Length       int
		XArray       []string `header:"X-Array"`
		TestHide     string   `header:"-"`
		IgnoreEmpty  string   `header:"X-Empty,omitempty"`
		IgnoreEmptyN string   `header:"X-Empty-N,omitempty"`
		CustomHeader http.Header
	}

	opt := Options{
		ContentType:  "application/json",
		Length:       2,
		XArray:       []string{"test1", "test2"},
		TestHide:     "hide",
		IgnoreEmptyN: "n",
		CustomHeader: http.Header{
			"X-Test-1": []string{"233"},
			"X-Test-2": []string{"666"},
		},
	}
	h, _ := httpheader.Header(opt)
	fmt.Println(h["Content-Type"])
	fmt.Println(h["Length"])
	fmt.Println(h["X-Array"])
	_, ok := h["TestHide"]
	fmt.Println(ok)
	_, ok = h["X-Empty"]
	fmt.Println(ok)
	fmt.Println(h["X-Empty-N"])
	fmt.Println(h["X-Test-1"])
	fmt.Println(h["X-Test-2"])
	// Output:
	// [application/json]
	// [2]
	// [test1 test2]
	// false
	// false
	// [n]
	// [233]
	// [666]
}
