package golanglibs

import (
	"fmt"
	"testing"
	"time"

	"github.com/icrowley/fake"
)

func TestFakeNameChinese(t *testing.T) {
	// for range Range(50) {
	// 	fmt.Println(Funcs.FakeNameChinese())
	// }

	for range Range(50) {
		fake.Seed(time.Now().UnixMicro())
		fmt.Println("DomainName:", fake.DomainName())
		fmt.Println("EmailAddress:", fake.EmailAddress())
		fmt.Println("Phone:", fake.Phone())
		fmt.Println("Country: ", fake.Country())
		fmt.Println("State: ", fake.State())
		fmt.Println("City: ", fake.City())
		fmt.Println("Street: ", fake.Street())
		fmt.Println("Company: ", fake.Company())
	}
}
