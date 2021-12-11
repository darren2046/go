package golanglibs

import (
	"time"

	pquernatotp "github.com/pquerna/otp/totp"
)

type totpStruct struct {
	key string
}

func getTotp(key string) *totpStruct {
	return &totpStruct{key: key}
}

func (m *totpStruct) validate(pass string) bool {
	return pquernatotp.Validate(pass, m.key)
}

func (m *totpStruct) password() string {
	p, err := pquernatotp.GenerateCode(m.key, time.Now())
	panicerr(err)
	return p
}
