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
	Head     func(uri string, args ...interface{}) *httpResp
	PostFile func(uri string, filePath string, args ...interface{}) *httpResp
	PostRaw  func(uri string, body string, args ...interface{}) *httpResp
	PostJSON func(uri string, json interface{}, args ...interface{}) *httpResp
	Post     func(uri string, args ...interface{}) *httpResp
	Get      func(uri string, args ...interface{}) *httpResp
	PutJSON  func(uri string, json interface{}, args ...interface{}) *httpResp
	Put      func(uri string, args ...interface{}) *httpResp
	PutRaw   func(uri string, body string, args ...interface{}) *httpResp
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
	Timeout             int
	ReadBodySize        int
	DoNotFollowRedirect bool
	HttpProxy           string
	TimeoutRetryTimes   int
	InsecureSkipVerify  bool
}

type httpRequestStruct struct {
	uri               string
	header            http.Header // request header
	param             req.Param
	readBodySize      int
	timeoutRetryTimes int
}

func initHttpRequest(uri string, args ...interface{}) *httpRequestStruct {
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
			if vv.Timeout != 0 {
				req.SetTimeout(getTimeDuration(vv.Timeout))
			} else {
				req.SetTimeout(getTimeDuration(10))
			}
			readBodySize = vv.ReadBodySize
			if vv.DoNotFollowRedirect {
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			if vv.HttpProxy != "" {
				u, err := url.Parse(vv.HttpProxy)
				Panicerr(err)
				client.Transport = &http.Transport{DisableKeepAlives: true, Proxy: http.ProxyURL(u)}
			}
			timeoutRetryTimes = vv.TimeoutRetryTimes
			if vv.InsecureSkipVerify {
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
	}
}

func (m *httpRequestStruct) doRequest(reqfunc func(uri string, header http.Header, param req.Param) (*req.Resp, error)) *httpResp {
	var timeouttimes int = 0
	var respBody []byte
	headers := make(map[string]string)
	var hresp *http.Response
	for {
		err := func() error {
			resp, err := reqfunc(m.uri, m.header, m.param)
			if err != nil {
				return err
			}

			hresp = resp.Response()
			for k, v := range hresp.Header {
				headers[k] = String(" ").Join(v).Get()
			}

			defer hresp.Body.Close()

			if m.readBodySize == 0 {
				respBody, err = ioutil.ReadAll(hresp.Body)
				if err != nil {
					return err
				}
			} else {
				buffer := make([]byte, m.readBodySize)
				num, err := io.ReadAtLeast(hresp.Body, buffer, m.readBodySize)
				if err != nil {
					return err
				}

				respBody = buffer[:num]
			}
			return nil
		}()

		if err != nil {
			if String("context deadline exceeded").In(err.Error()) || String("Timeout exceeded").In(err.Error()) {
				timeouttimes += 1
				if m.timeoutRetryTimes != -1 && timeouttimes >= m.timeoutRetryTimes {
					Panicerr(err)
				} else {
					continue
				}
			} else {
				Panicerr(err)
			}
		}

		break
	}

	return &httpResp{
		Content:    string(respBody),
		Headers:    headers,
		StatusCode: hresp.StatusCode,
		URL:        hresp.Request.URL.String(),
	}
}

func httpHead(uri string, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Head(uri, header, param)
	})
}

func httpPostFile(uri string, filePath string, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Post(uri, req.File(filePath), header, param)
	})
}

func httpPostRaw(uri string, body string, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Post(uri, body, header, param)
	})
}

func httpPostJSON(uri string, json interface{}, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Post(uri, req.BodyJSON(&json), header, param)
	})
}

func httpPost(uri string, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Post(uri, header, param)
	})
}

func httpGet(uri string, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Get(uri, header, param)
	})
}

func httpPutJSON(uri string, json interface{}, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Put(uri, req.BodyJSON(&json), header, param)
	})
}

func httpPut(uri string, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Put(uri, header, param)
	})
}

func httpPutRaw(uri string, body string, args ...interface{}) *httpResp {
	return initHttpRequest(uri, args...).doRequest(func(uri string, header http.Header, param req.Param) (*req.Resp, error) {
		return req.Put(uri, body, header, param)
	})
}
