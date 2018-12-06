package cos

import (
	"fmt"
	"testing"
)

func Test_calSHA1Digest(t *testing.T) {
	want := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	got := fmt.Sprintf("%x", calSHA1Digest([]byte("test")))
	if got != want {

		t.Errorf("calSHA1Digest request sha1: %+v, want %+v", got, want)
	}
}

func Test_calMD5Digest(t *testing.T) {
	want := "098f6bcd4621d373cade4e832627b4f6"
	got := fmt.Sprintf("%x", calMD5Digest([]byte("test")))
	if got != want {

		t.Errorf("calMD5Digest request md5: %+v, want %+v", got, want)
	}
}
