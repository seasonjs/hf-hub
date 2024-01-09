// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package api_test

import (
	"github.com/seasonjs/hf-hub/api"
	"io"
	"os"
	"strings"
	"testing"
)

func assert(t *testing.T, a any, b any) {
	if a != b {
		t.Error("assert failed", "a:", a, "b:", b)
		return
	}
	t.Log("assert ok")
}

func TestApi(t *testing.T) {
	builder, err := api.NewApiBuilder()
	if err != nil {
		t.Error(err)
		return
	}
	hapi := builder.WithCacheDir("../tmp").Build()

	downloadedPath, err := hapi.Model("julien-c/dummy-unknown").Download("config.json")
	if err != nil {
		t.Error(err)
		return
	}

	if _, err = os.Stat(downloadedPath); os.IsNotExist(err) {
		t.Error(err)
		return
	}

	sha256, err := api.GetSHA256FromFile(downloadedPath)
	if err != nil {
		t.Error(err)
		return
	}

	assert(t, sha256, "b908f2b7227d4d31a2105dfa31095e28d304f9bc938bfaaa57ee2cacf1f62d32")

	cachePath, err := hapi.Repo(api.NewRepo("julien-c/dummy-unknown", api.Model)).Get("config.json")
	if err != nil {
		t.Error(err)
		return
	}

	assert(t, cachePath, downloadedPath)

}

func TestApi_Dataset(t *testing.T) {
	builder, err := api.NewApiBuilder()
	if err != nil {
		t.Error(err)
		return
	}
	hapi := builder.
		WithCacheDir("../tmp").
		Build()

	repo := api.NewRepoWithRevision("wikitext", api.Dataset, "refs/convert/parquet")

	downloadedPath, err := hapi.Repo(repo).Download("wikitext-103-v1/test/0000.parquet")
	if err != nil {
		t.Error(err)
		return
	}

	if _, err = os.Stat(downloadedPath); os.IsNotExist(err) {
		t.Error(err)
		return
	}

	sha256, err := api.GetSHA256FromFile(downloadedPath)
	if err != nil {
		t.Error(err)
		return
	}

	assert(t, sha256, "abdfc9f83b1103b502924072460d4c92f277c9b49c313cef3e48cfcf7428e125")
}

func TestApi_Model(t *testing.T) {
	builder, err := api.NewApiBuilder()
	if err != nil {
		t.Error(err)
		return
	}
	hapi := builder.
		WithCacheDir("../tmp").
		Build()

	repo := api.NewRepoWithRevision("BAAI/bGe-reRanker-Base", api.Model, "refs/pr/5")

	downloadedPath, err := hapi.Repo(repo).Download("tokenizer.json")
	if err != nil {
		t.Error(err)
		return
	}

	if _, err = os.Stat(downloadedPath); os.IsNotExist(err) {
		t.Error(err)
		return
	}

	sha256, err := api.GetSHA256FromFile(downloadedPath)
	if err != nil {
		t.Error(err)
		return
	}

	assert(t, sha256, strings.ToLower("9EB652AC4E40CC093272BBBE0F55D521CF67570060227109B5CDC20945A4489E"))
}

func TestApiRepo_Info(t *testing.T) {
	builder, err := api.NewApiBuilder()
	if err != nil {
		t.Error(err)
		return
	}
	hapi := builder.
		WithCacheDir("../tmp").
		Build()

	repo := api.NewRepoWithRevision("wikitext", api.Dataset, "refs/convert/parquet")

	info, err := hapi.Repo(repo).Info()

	t.Log(info)
}

func TestApiRepo_InfoRequest(t *testing.T) {
	builder, err := api.NewApiBuilder()
	if err != nil {
		t.Error(err)
		return
	}
	hapi := builder.
		WithCacheDir("../tmp").
		Build()

	repo := api.NewRepoWithRevision("mcpotato/42-eicar-street", api.Model, "8b3861f6931c4026b0cd22b38dbc09e7668983ac")

	res, err := hapi.Repo(repo).InfoRequest()
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(data))

	assert(t, res.StatusCode, 200)
}
