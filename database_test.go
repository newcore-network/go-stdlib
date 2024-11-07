package libs

import (
	"fmt"
	"testing"

	"github.com/styerr-development/libs/configuration"
	"github.com/styerr-development/libs/database"
	"github.com/styerr-development/libs/database/drivers"
)

func TestConnectionPassed(t *testing.T) {
	postgres := &drivers.PostgresConnection{}
	cfg := configuration.GeneralConfig{}
	conn, err := database.NewConnection(postgres, cfg)
	if err != nil {
		fmt.Printf("error: %v", err)
		t.FailNow()
	}
	fmt.Println(conn)
}
