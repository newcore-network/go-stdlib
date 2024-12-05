package libs

import (
	"fmt"
	"testing"

	"github.com/newcore-network/libs/configuration"
	"github.com/newcore-network/libs/database"
	"github.com/newcore-network/libs/database/drivers"
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
