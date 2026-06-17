package model_test

import (
	"nsf/model"
	"testing"

	"github.com/free-ran-ue/util"
	"github.com/google/go-cmp/cmp"
)

const (
	CONFIG_FILE = "../config.yaml"
)

var targetNsfConfig = model.NsfConfig{
	Nsf: model.NsfIE{
		SbiAddress: "0.0.0.0",
		SbiPort:    8000,
		DbType:     "bbolt",
	},
	Logger: model.LoggerIE{
		Level: "info",
	},
}

func TestNsfConfig(t *testing.T) {
	nsfConfig := model.NsfConfig{}
	if err := util.LoadFromYaml(CONFIG_FILE, &nsfConfig); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if diff := cmp.Diff(targetNsfConfig, nsfConfig); diff != "" {
		t.Errorf("NsfConfig mismatch (-want +got):\n%s", diff)
	}
}
