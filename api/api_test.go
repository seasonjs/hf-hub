package api_test

import (
	"hf-hub/api"
	"os"
	"testing"
)

func Assest(t *testing.T, con bool) {
	if !con {
		t.Error("assest error")
	}
}

func TestApi(t *testing.T) {
	builder, err := api.NewApiBuilder()
	if err != nil {
		t.Error(err)
		return
	}
	hapi := builder.WithCacheDir("./tmp").Build()

	downloadedPath, err := hapi.Model("julien-c/dummy-unknown").Download("config.json")
	if err != nil {
		t.Error(err)
		return
	}

	if _, err = os.Stat(downloadedPath); os.IsNotExist(err) {
		t.Error(err)
		return
	}

}
