package golanglibs

import "net/url"

type urlComponents struct {
	schema   string
	host     string
	port     string
	user     string
	pass     string
	path     string
	query    string
	fragment string
}

type urlStruct struct {
	url string
}

func getUrl(url string) *urlStruct {
	return &urlStruct{url: url}
}

func (u *urlStruct) Parse() *urlComponents {
	uu, err := url.Parse(u.url)
	panicerr(err)

	pass, _ := uu.User.Password()

	var port string

	if uu.Port() == "" {
		if uu.Scheme == "https" {
			port = "443"
		}
		if uu.Scheme == "http" {
			port = "80"
		}
	} else {
		port = uu.Port()
	}

	return &urlComponents{
		schema:   uu.Scheme,
		host:     uu.Hostname(),
		port:     port,
		user:     uu.User.Username(),
		pass:     pass,
		path:     uu.Path,
		query:    uu.RawQuery,
		fragment: uu.Fragment,
	}
}

func (u *urlStruct) Encode() string {
	return url.QueryEscape(u.url)
}

func (u *urlStruct) Decode() string {
	str, err := url.QueryUnescape(u.url)
	panicerr(err)
	return str
}
