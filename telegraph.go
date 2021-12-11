package golanglibs

import (
	"gitlab.com/toby3d/telegraph"
)

type telegraphStruct struct {
	account *telegraph.Account
}

type telegraphPageInfo struct {
	author  string
	title   string
	url     string
	content string
}

func getTelegraph(AuthorName string) *telegraphStruct {
	account, err := telegraph.CreateAccount(telegraph.Account{
		ShortName:  randomStr(6, "abcdefghijklmn1234567890"),
		AuthorName: AuthorName,
	})
	panicerr(err)
	return &telegraphStruct{
		account: account,
	}
}

func (m *telegraphStruct) post(title string, content string) *telegraphPageInfo {
	tcontent, err := telegraph.ContentFormat(content)
	panicerr(err)

	page, err := m.account.CreatePage(telegraph.Page{
		URL:        "this-is-a-test-url",
		Title:      title,
		AuthorName: m.account.AuthorName,
		Content:    tcontent,
	}, true)
	panicerr(err)
	return &telegraphPageInfo{
		author:  page.AuthorName,
		title:   page.Title,
		url:     page.URL,
		content: content,
	}
}
