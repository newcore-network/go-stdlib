package configuration

import (
	"fmt"
	"testing"
)

func TestConfiguration(t *testing.T) {
	conf := GetFromEnvFile(".env")
	fmt.Println(conf)
}
