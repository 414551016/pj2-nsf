package nsfContext_test

import (
	"fmt"
	"nsf/internal/db"
	"nsf/internal/nsfContext"
	"nsf/logger"
	"path/filepath"
	"sync"
	"testing"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
)

func isolatedContext(t *testing.T, path string) *nsfContext.NsfContext {
	t.Helper()
	previous := db.NSF_DB_PATH
	db.NSF_DB_PATH = path
	defer func() { db.NSF_DB_PATH = previous }()
	log := logger.NewNsfLogger(loggergoUtil.LogLevelString("error"))
	ctx := nsfContext.NewNsfContext(&nsfContext.NsfContextIE{DbType: "bbolt", NsfLogger: &log})
	if ctx == nil {
		t.Fatal("context initialization failed")
	}
	return ctx
}

func TestPersistenceAndPartialDeregistration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "recovery.db")
	ctx := isolatedContext(t, path)
	if err := ctx.RegisterServices("nf-a", "127.0.0.11", 8101, []string{"shared", "keep"}); err != nil {
		ctx.Release()
		t.Fatal(err)
	}
	if err := ctx.RegisterServices("nf-b", "127.0.0.12", 8102, []string{"shared"}); err != nil {
		ctx.Release()
		t.Fatal(err)
	}
	ctx.Release()
	ctx = isolatedContext(t, path)
	if name, addr, port, err := ctx.DiscoverService("keep"); err != nil || name != "nf-a" || addr != "127.0.0.11" || port != 8101 {
		ctx.Release()
		t.Fatalf("recovery failed: %s %s %d %v", name, addr, port, err)
	}
	if err := ctx.DeregisterServices("nf-a", []string{"shared"}); err != nil {
		ctx.Release()
		t.Fatal(err)
	}
	if name, _, _, err := ctx.DiscoverService("shared"); err != nil || name != "nf-b" {
		ctx.Release()
		t.Fatalf("other provider lost: %s %v", name, err)
	}
	if _, _, _, err := ctx.DiscoverService("keep"); err != nil {
		ctx.Release()
		t.Fatalf("unrelated service lost: %v", err)
	}
	if err := ctx.DeregisterServices("nf-b", []string{"shared"}); err != nil {
		ctx.Release()
		t.Fatal(err)
	}
	ctx.Release()
	ctx = isolatedContext(t, path)
	defer ctx.Release()
	if _, _, _, err := ctx.DiscoverService("shared"); err == nil {
		t.Fatal("deleted service reappeared after restart")
	}
	if _, _, _, err := ctx.DiscoverService("never-registered"); err == nil {
		t.Fatal("unknown service should return an error")
	}
}

func TestConcurrentServiceOperations(t *testing.T) {
	ctx := isolatedContext(t, filepath.Join(t.TempDir(), "concurrent.db"))
	defer ctx.Release()
	var wg sync.WaitGroup
	for i := range 24 {
		wg.Go(func() {
			name := fmt.Sprintf("nf-%d", i)
			service := fmt.Sprintf("service-%d", i)
			if err := ctx.RegisterServices(name, "127.0.0.1", 9000+i, []string{service}); err != nil {
				t.Error(err)
				return
			}
			got, _, port, err := ctx.DiscoverService(service)
			if err != nil || got != name || port != 9000+i {
				t.Errorf("incorrect concurrent result: %s %d %v", got, port, err)
			}
			if err := ctx.DeregisterServices(name, []string{service}); err != nil {
				t.Error(err)
			}
			if _, _, _, err := ctx.DiscoverService(service); err == nil {
				t.Errorf("%s still discoverable", service)
			}
		})
	}
	wg.Wait()
}
