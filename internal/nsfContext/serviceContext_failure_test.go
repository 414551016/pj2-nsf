package nsfContext

import (
	"errors"
	"nsf/internal/db"
	"nsf/logger"
	"testing"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
)

type failingStore struct {
	db.DbIf
	failure error
}

func (f *failingStore) Exist(string) (bool, error)      { return false, nil }
func (f *failingStore) SaveAll(map[string]string) error { return f.failure }

func TestRegistrationReportsDatabaseFailure(t *testing.T) {
	failure := errors.New("storage unavailable")
	log := logger.NewNsfLogger(loggergoUtil.LogLevelString("error"))
	ctx := newServiceContext(&serviceContextIE{DbIf: &failingStore{failure: failure}, NsfLogger: &log})
	if err := ctx.registerServices("nf", "127.0.0.1", 8000, []string{"svc"}); !errors.Is(err, failure) {
		t.Fatalf("expected database failure, got %v", err)
	}
	ctx.release()
	if err := ctx.registerServices("nf", "127.0.0.1", 8000, []string{"svc"}); err == nil {
		t.Fatal("closed context accepted registration")
	}
}
