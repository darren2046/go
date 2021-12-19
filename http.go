package golanglibs

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/corpix/uarand"
	"github.com/imroc/req"
)

type httpStruct struct {
	Head     func(uri string, args ...interface{}) httpResp
	PostFile func(uri string, filePath string, args ...interface{}) httpResp
	PostRaw  func(uri string, body string, args ...interface{}) httpResp
	PostJSON func(uri string, json interface{}, args ...interface{}) httpResp
	Post     func(uri string, args ...interface{}) httpResp
	Get      func(uri string, args ...interface{}) httpResp
	PutJSON  func(uri string, json interface{}, args ...interface{}) httpResp
	Put      func(uri string, args ...interface{}) httpResp
	PutRaw   func(uri string, body string, args ...interface{}) httpResp
}

var Http httpStruct

func init() {
	Http = httpStruct{
		Head:     httpHead,
		PostFile: httpPostFile,
		PostRaw:  httpPostRaw,
		PostJSON: httpPostJSON,
		Post:     httpPost,
		Get:      httpGet,
		PutJSON:  httpPutJSON,
		Put:      httpPut,
		PutRaw:   httpPutRaw,
	}
}

type HttpHeader map[string]string
type HttpParam map[string]interface{}

type httpResp struct {
	Headers    map[string]string
	Content    string
	StatusCode int
	URL        string
}

type HttpConfig struct {
	timeout             int
	readBodySize        int
	doNotFollowRedirect bool
	httpProxy           string
	timeoutRetryTimes   int
	insecureSkipVerify  bool
}

func httpHead(uri string, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Head(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

func httpPostFile(uri string, filePath string, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, req.File(filePath), header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

func httpPostRaw(uri string, body string, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, body, header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

func httpPostJSON(uri string, json interface{}, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, req.BodyJSON(&json), header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

// httpPost(url, HttpHeader{}, HttpParam{}) {
func httpPost(uri string, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Post(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

// httpGet(url, HttpHeader{}, HttpParam{}) {
func httpGet(uri string, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}

	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Get(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

func httpPutJSON(uri string, json interface{}, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Put(uri, req.BodyJSON(&json), header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

// httpPost(url, HttpHeader{}, HttpParam{}) {
func httpPut(uri string, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Put(uri, header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

func httpPutRaw(uri string, body string, args ...interface{}) httpResp {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	param := make(req.Param)
	var readBodySize int
	var timeoutRetryTimes int = 0
	for _, v := range args {
		switch vv := v.(type) {
		case HttpHeader:
			for k, vvv := range vv {
				header.Set(k, vvv)
			}
		case HttpParam:
			for k, vvv := range vv {
				param[k] = vvv
			}
		case HttpConfig:
			if vv.timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.readBodySize
			if vv.doNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.httpProxy != "" {
				u, err := url.Parse(vv.httpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.timeoutRetryTimes
			if vv.insecureSkipVerify {
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}
		}
	}

	var timeouttimes int = 0
	var resp *req.Resp
	var err error
	var respBody []byte
	var hresp *http.Response
	headers := make(map[string]string)
	// lg.trace("timeoutRetryTimes:", timeoutRetryTimes)
	for {
		resp, err = req.Put(uri, body, header, param)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		hresp = resp.Response()
		for k, v := range hresp.Header {
			headers[k] = String(" ").Join(v).Get()
		}

		defer hresp.Body.Close()

		if readBodySize == 0 {
			respBody, err = ioutil.ReadAll(hresp.Body)
			if err != nil {
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					Panicerr(err)
				}
			}
		} else {
			buffer := make([]byte, readBodySize)
			num, err := io.ReadAtLeast(hresp.Body, buffer, readBodySize)
			if err != nil {
				// lg.trace(err)
				if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
					timeouttimes += 1
					if timeoutRetryTimes != -1 && timeouttimes >= timeoutRetryTimes {
						Panicerr(err)
					} else {
						continue
					}
				} else {
					if !String("unexpected EOF").In(err.Error()) {
						Panicerr(err)
					}
				}
			}

			respBody = buffer[:num]
		}
		break
	}

	return httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}
