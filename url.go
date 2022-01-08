package golanglibs

import "net/url"

type urlComponents struct {
	Schema   string
	Host     string
	Port     string
	User     string
	Pass     string
	Path     string
	Query    string
	Fragment string
}

type urlStruct struct {
	url string
}

func getUrl(url string) *urlStruct {
	return &urlStruct{url: url}
}

func (u *urlStruct) Parse() *urlComponents {
	uu, err := url.Parse(u.url)
	Panicerr(err)

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
		Schema:   uu.Scheme,
		Host:     uu.Hostname(),
		Port:     port,
		User:     uu.User.Username(),
		Pass:     pass,
		Path:     uu.Path,
		Query:    uu.RawQuery,
		Fragment: uu.Fragment,
	}
}

func (u *urlStruct) Encode() string {
	return url.QueryEscape(u.url)
}

func (u *urlStruct) Decode() string {
	str, err := url.QueryUnescape(u.url)
	Panicerr(err)
	return str
}
