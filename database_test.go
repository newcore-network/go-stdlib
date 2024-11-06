package libs

import (
	"fmt"
	"testing"

	"github.com/styerr-development/libs/database"
	"github.com/styerr-development/libs/database/drivers"
)

func TestConnectionPassed(t *testing.T) {
	postgres := &drivers.PostgresConnection{}
	conn, err := database.NewConnection(postgres)
	if err != nil {
		fmt.Printf("error: %v", err)
		t.FailNow()
	}
	fmt.Println(conn)
}
