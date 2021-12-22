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

type httpRequestStruct struct {
	uri               string
	header            http.Header // request header
	param             req.Param
	readBodySize      int
	timeoutRetryTimes int
	timeouttimes      int
	headers           map[string]string // response header
	hresp             *http.Response
	respBody          []byte
}

func getHttpRequest(uri string, args ...interface{}) *httpRequestStruct {
	if !String(uri).StartsWith("http://") && !String(uri).StartsWith("https://") {
		uri = "http://" + uri
	}

	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	req.SetClient(client)

	req.SetTimeout(getTimeDuration(10))

	header := make(http.Header)
	header.Set("User-Agent", uarand.GetRandom())

	var readBodySize int
	var timeoutRetryTimes int = 0
	param := make(req.Param)
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
	return &httpRequestStruct{
		uri:               uri,
		header:            header,
		param:             param,
		readBodySize:      readBodySize,
		timeoutRetryTimes: timeoutRetryTimes,
		timeouttimes:      0,
	}
}

func (m *httpRequestStruct) responseHandler(resp *req.Resp, err error) bool {
	m.headers = make(map[string]string)
	if err != nil {
		// lg.trace(err)
		if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
			m.timeouttimes += 1
			if m.timeoutRetryTimes != -1 && m.timeouttimes >= m.timeoutRetryTimes {
				Panicerr(err)
			} else {
				return false
			}
		} else {
			Panicerr(err)
		}
	}

	m.hresp = resp.Response()
	for k, v := range m.hresp.Header {
		m.headers[k] = String(" ").Join(v).Get()
	}

	defer m.hresp.Body.Close()

	if m.readBodySize == 0 {
		m.respBody, err = ioutil.ReadAll(m.hresp.Body)
		if err != nil {
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				m.timeouttimes += 1
				if m.timeoutRetryTimes != -1 && m.timeouttimes >= m.timeoutRetryTimes {
					Panicerr(err)
				} else {
					return false
				}
			} else {
				Panicerr(err)
			}
		}
	} else {
		buffer := make([]byte, m.readBodySize)
		num, err := io.ReadAtLeast(m.hresp.Body, buffer, m.readBodySize)
		if err != nil {
			// lg.trace(err)
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				m.timeouttimes += 1
				if m.timeoutRetryTimes != -1 && m.timeouttimes >= m.timeoutRetryTimes {
					Panicerr(err)
				} else {
					return false
				}
			} else {
				if !String("unexpected EOF").In(err.Error()) {
					Panicerr(err)
				}
			}
		}

		m.respBody = buffer[:num]
	}
	return true
}

func httpHead(uri string, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)

	for {
		resp, err := req.Head(httpRequest.uri, httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

func httpPostFile(uri string, filePath string, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Post(uri, req.File(filePath), httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

func httpPostRaw(uri string, body string, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Post(uri, body, httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

func httpPostJSON(uri string, json interface{}, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Post(uri, req.BodyJSON(&json), httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

// httpPost(url, HttpHeader{}, HttpParam{}) {
func httpPost(uri string, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Post(uri, httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

// httpGet(url, HttpHeader{}, HttpParam{}) {
func httpGet(uri string, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Get(uri, httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

func httpPutJSON(uri string, json interface{}, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Put(uri, req.BodyJSON(&json), httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

// httpPost(url, HttpHeader{}, HttpParam{}) {
func httpPut(uri string, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Put(uri, httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}

func httpPutRaw(uri string, body string, args ...interface{}) httpResp {
	httpRequest := getHttpRequest(uri, args...)
	for {
		resp, err := req.Put(uri, body, httpRequest.header, httpRequest.param)
		if !httpRequest.responseHandler(resp, err) {
			continue
		}
		break
	}

	return httpResp{
		Content:    string(httpRequest.respBody),
		Headers:    httpRequest.headers,
		StatusCode: httpRequest.hresp.StatusCode,
		URL:        httpRequest.hresp.Request.URL.String(),
	}
}
