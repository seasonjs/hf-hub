package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type RepoType int

const (
	Model RepoType = iota
	Dataset
	Space
)

type Cache struct {
	path string
}

func NewCache(path string) *Cache {
	//if !filepath.IsAbs(path) {
	//	path, err := filepath.Abs(path)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return &Cache{path: path}, nil
	//}
	return &Cache{path: path}
}

func DefaultCache() (*Cache, error) {
	homePath := os.Getenv("HF_HOME")
	if len(homePath) == 0 {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		homePath = filepath.Join(homePath, ".cache", "huggingface")
	}
	cachePath := filepath.Join(homePath, "hub")
	return NewCache(cachePath), nil
}

func (c *Cache) Path() string {
	return c.path
}

func (c *Cache) TokenPath() string {
	return filepath.Join(filepath.Dir(c.path), "token")
}

func (c *Cache) Token() (string, error) {
	tokenPath := c.TokenPath()
	file, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	token := strings.TrimSpace(string(file))
	if len(token) == 0 {
		return "", errors.New("token file empty")
	}
	return token, nil
}

func (c *Cache) Repo(rep *Repo) *CacheRepo {
	return NewCacheRepo(c.Clone(), rep)
}

func (c *Cache) Model(modelId string) *CacheRepo {
	return c.Repo(NewRepo(modelId, Model))
}

func (c *Cache) Dataset(modelId string) *CacheRepo {
	return c.Repo(NewRepo(modelId, Dataset))
}

func (c *Cache) Space(modelId string) *CacheRepo {
	return c.Repo(NewRepo(modelId, Space))
}

func (c *Cache) TempPath() (string, error) {
	path := filepath.Join(c.path, "tmp")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	path = filepath.Join(path, randStr(7))
	return path, nil
}

func (c *Cache) Clone() *Cache {
	newCache := NewCache(c.path)
	return newCache
}

type CacheRepo struct {
	cache *Cache
	repo  *Repo
}

func NewCacheRepo(cache *Cache, repo *Repo) *CacheRepo {
	return &CacheRepo{repo: repo, cache: cache}
}

func (r *CacheRepo) Get(filename string) (string, error) {
	commitPath := r.refPath()
	commitHash, err := os.ReadFile(commitPath)
	if err != nil {
		return "", nil
	}
	path := r.PointerPath(string(commitHash))
	path = filepath.Join(path, filename)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return "", err
	}
	return path, nil
}

func (r *CacheRepo) path() string {
	return filepath.Join(r.cache.path, r.repo.FolderName())
}

func (r *CacheRepo) refPath() string {
	path := r.path()
	return filepath.Join(path, "refs", r.repo.Revision())
}

func (r *CacheRepo) PointerPath(etag string) string {
	return filepath.Join(r.path(), "blobs", etag)
}

type Repo struct {
	repoId   string
	repoType RepoType
	revision string
}

func NewRepo(repoId string, repoType RepoType) *Repo {
	return NewRepoWithRevision(repoId, repoType, "main")
}

func NewRepoWithRevision(repoId string, repoType RepoType, revision string) *Repo {
	return &Repo{repoId: repoId, repoType: repoType, revision: revision}
}

func NewModelRepo(repoId string) *Repo {
	return NewRepo(repoId, Model)
}

func NewDatasetRepo(repoId string) *Repo {
	return NewRepo(repoId, Dataset)
}

func NewSpaceRepo(repoId string) *Repo {
	return NewRepo(repoId, Space)
}

func (r *Repo) Clone() *Repo {
	newRepo := *r
	return &newRepo
}

// FolderName  The normalized folder nameof the repo within the cache directory
func (r *Repo) FolderName() string {
	var prefix string
	switch r.repoType {
	case Model:
		prefix = "model"
	case Dataset:
		prefix = "dataset"
	case Space:
		prefix = "space"
	}
	result := fmt.Sprintf("%s--%s", prefix, r.repoId)
	result = strings.ReplaceAll(result, "/", "--")
	return result
}

func (r *Repo) Revision() string {
	return r.revision
}

func (r *Repo) UrlRevision() string {
	return strings.ReplaceAll(r.revision, "/", "%2F")
}

func (r *Repo) ApiUrl() string {
	var prefix string
	switch r.repoType {
	case Model:
		prefix = "model"
	case Dataset:
		prefix = "dataset"
	case Space:
		prefix = "space"
	}
	return fmt.Sprintf("%s/%s/revision/%s", prefix, r.repoId, r.UrlRevision())
}

type ApiBuilder struct {
	endpoint         string
	cache            *Cache
	urlTemplate      string
	token            string
	maxFiles         uint64
	chunkSize        uint64
	parallelFailures uint64
	maxRetries       uint64
	progress         bool
}

func NewApiBuilder() (*ApiBuilder, error) {
	cache, err := DefaultCache()
	if err != nil {
		return nil, err
	}
	apiBuilder := &ApiBuilder{}
	return apiBuilder.FromCache(cache)
}

func (b *ApiBuilder) FromCache(cache *Cache) (*ApiBuilder, error) {
	token, err := cache.Token()
	if err != nil {
		return nil, err
	}
	return &ApiBuilder{
		endpoint: "https://huggingface.co",
		//"{endpoint}/{repo_id}/resolve/{revision}/{filename}"
		urlTemplate: "{{.Endpoint}}/{{.RepoId}}/resolve/{{.Revision}}/{{.Filename}}",
		cache:       cache,
		token:       token,
		progress:    true,
	}, nil
}

func (b *ApiBuilder) WithProgress(progress bool) *ApiBuilder {
	b.progress = progress
	return b
}

func (b *ApiBuilder) WithCacheDir(cacheDir string) *ApiBuilder {
	cache := NewCache(cacheDir)
	b.cache = cache
	return b
}

func (b *ApiBuilder) WithToken(token string) *ApiBuilder {
	b.token = token
	return b
}

func (b *ApiBuilder) BuildHeaders() http.Header {
	headers := make(http.Header)
	userAgent := fmt.Sprintf("unkown/None; %s/%s; rust/unknown", "hf-hub", "v0.0.1")
	headers.Add("User-Agent", userAgent)
	if len(b.token) > 0 {
		headers.Add("Authorization", fmt.Sprintf("Bearer %s", b.token))
	}
	return headers
}

func (b *ApiBuilder) Build() *Api {
	headers := b.BuildHeaders()
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}
	return &Api{
		endpoint:    b.endpoint,
		urlTemplate: b.urlTemplate,
		cache:       b.cache,
		headers:     headers,
		client:      client,
		progress:    b.progress,
	}

}

type Metadata struct {
}

type Api struct {
	endpoint    string
	urlTemplate string
	cache       *Cache
	headers     http.Header
	client      *http.Client
	progress    bool
}

func NewApi() (*Api, error) {
	builder, err := NewApiBuilder()
	if err != nil {
		return nil, err
	}
	api := builder.Build()
	return api, nil
}

func (a *Api) Repo(rep *Repo) *ApiRepo {
	return NewApiRepo(a.Clone(), rep)
}

func (a *Api) Model(modelId string) *ApiRepo {
	return NewApiRepo(a.Clone(), NewModelRepo(modelId))
}

func (a *Api) Dataset(datasetId string) *ApiRepo {
	return NewApiRepo(a.Clone(), NewDatasetRepo(datasetId))
}

func (a *Api) Space(spaceId string) *ApiRepo {
	return NewApiRepo(a.Clone(), NewSpaceRepo(spaceId))
}

func (a *Api) Clone() *Api {
	newApi := *a
	newCache := *a.cache
	newClient := *a.client

	newApi.headers = newApi.headers.Clone()
	newApi.cache = &newCache
	newApi.client = &newClient
	return &newApi
}

type ApiRepo struct {
	api  *Api
	repo *Repo
}

func NewApiRepo(api *Api, rep *Repo) *ApiRepo {
	return &ApiRepo{
		api:  api,
		repo: rep,
	}
}

type UrlTemplateEntry struct {
	Endpoint string
	RepoId   string
	Revision string
	Filename string
}

func (r *ApiRepo) Url(filename string) (string, error) {
	t, err := template.New("test").Parse(r.api.urlTemplate)
	if err != nil {
		return "", nil
	}
	var buf bytes.Buffer
	//"{{.Endpoint}}/{{.RepoId}}/resolve/{{.Revision}}/{{.Filename}}",
	err = t.Execute(&buf, UrlTemplateEntry{
		Endpoint: r.api.endpoint,
		RepoId:   r.repo.repoId,
		Revision: r.repo.revision,
		Filename: filename,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (r *ApiRepo) Get(filename string) (string, error) {
	path, err := r.api.cache.Repo(r.repo.Clone()).Get(filename)
	if err != nil {
		return r.Download(filename)
	}
	return path, nil
}

func (r *ApiRepo) Download(filename string) (string, error) {
	//url, err := r.Url(filename)
	//if err != nil {
	//	return "", err
	//}
	panic("todo")
}
