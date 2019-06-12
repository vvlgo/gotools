package example_test

import (
	"fmt"
	"github.com/vvlgo/gotools/jwttoken"
	"testing"
)

func TestToken(t *testing.T) {
	s, err := jwttoken.Sign("test", 100, "admin", "123")
	if err != nil {
		panic(err)
	}
	unsign, err := jwttoken.Unsign(s, "123")
	if err != nil {
		panic(err)
	}
	fmt.Println(unsign)
}
