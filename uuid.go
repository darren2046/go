package golanglibs

import (
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid"
)

type uuidStruct struct {
	Uuid4      func() string
	Shortuuid4 func() string
}

var Uuid uuidStruct

func init() {
	Uuid = uuidStruct{
		Uuid4:      uuid4,
		Shortuuid4: shortuuid4,
	}
}

func uuid4() string {
	return uuid.New().String()
}

func shortuuid4() string {
	return shortuuid.New()
}
