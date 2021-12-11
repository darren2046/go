package golanglibs

import "testing"

func TestHttpGet(t *testing.T) {
	Print(Http.Get("http://ifconfig.me"))
}
