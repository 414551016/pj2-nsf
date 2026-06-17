package nsfContext_test

import (
	"fmt"
	"nsf/internal/db"
	"nsf/internal/nsfContext"
	"nsf/logger"
	"os"
	"testing"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
)

/*
The TestMain function sets up the testing environment for the nsfContext package.
It initializes a new NsfContext instance with a specified database path and logger configuration.
After running the tests, it releases the context and cleans up by removing the test database file.
This ensures that each test run starts with a clean state and that resources are properly managed. 
*/

var nsfCtx *nsfContext.NsfContext

func TestMain(m *testing.M) {
	db.NSF_DB_PATH = "/tmp/nsf_db_context_test.db"

	logger := logger.NewNsfLogger(loggergoUtil.LogLevelString("error"))

	nsfCtx = nsfContext.NewNsfContext(&nsfContext.NsfContextIE{
		DbType: "bbolt",

		NsfLogger: &logger,
	})

	exitCode := m.Run()

	nsfCtx.Release()

	if err := os.Remove(db.NSF_DB_PATH); err != nil {
		fmt.Printf("Failed to remove test DB file: %v\n", err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
