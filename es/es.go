package util

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/gocarina/gocsv"
	"github.com/iancoleman/orderedmap"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Doc[Source, Aggregations any] struct {
	Took     int64 `json:"token"`
	TimedOut bool  `json:"timed_out"`
	Shards   struct {
		Total      int64 `json:"total"`
		Successful int64 `json:"successful"`
		Skipped    int64 `json:"skipped"`
		Failed     int64 `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64     `json:"max_score"`
		Hits     Hit[Source] `json:"hits"`
	} `json:"hits"`
	Aggregations Aggregations `json:"aggregations"`
}

type Hit[T any] struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
	Sort   []any   `json:"sort"`
}

type DoFunc func(*http.Request) (*http.Response, error)

// ElasticSearchRequest 主要封装 ElasticSearch 调用方便使用
// 整个过程划分为三个阶段
// 1. 构造 ElasticSearch 请求
// 2. 发送 HTTP 请求，获取响应
// 3. 将响应信息已特定方式返回给使用者
// 为了方便调用者使用，在 1 阶段所有产生的错误先临时记录下来，在 2 阶段执行前再返回，这样能够让代码非常简洁
type ElasticSearchRequest struct {
	url  string
	req  *http.Request
	resp *http.Response
	body []byte
	err  error
	do   DoFunc
}

// Search 快速构造一个搜索请求
func Search(index string, do DoFunc) *ElasticSearchRequest {
	url := path.Join(index, "_search")
	return Post(url, do)
}

// Update 快速构造一个更新请求
func Update(index, id string, do DoFunc) *ElasticSearchRequest {
	url := path.Join(index, "_update", id)
	return Post(url, do)
}

func DeleteDoc(index, id string, do DoFunc) *ElasticSearchRequest {
	url := path.Join(index, "_doc", id)
	return Delete(url, do)
}

func SQL(do DoFunc) *ElasticSearchRequest {
	return Post("_sql?format=csv", do)
}

// Post 快速构造一个 Post 请求
func Post(url string, do DoFunc) *ElasticSearchRequest {
	return NewElasticSearchRequest(url, http.MethodPost, do)
}

func Delete(url string, do DoFunc) *ElasticSearchRequest {
	return NewElasticSearchRequest(url, http.MethodDelete, do)
}

// NewElasticSearchRequest 构建 ElasticSearch 请求
func NewElasticSearchRequest(url, method string, do DoFunc) *ElasticSearchRequest {
	return NewElasticSearchRequestWithContext(context.Background(), url, method, do)
}

// NewElasticSearchRequestWithContext 构建 ElasticSearch 请求，并支持用户传递 context
func NewElasticSearchRequestWithContext(ctx context.Context, url, method string, do DoFunc) *ElasticSearchRequest {
	e := &ElasticSearchRequest{
		url:  url,
		resp: &http.Response{},
		do:   do,
	}
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		e.addError(err)
		return e
	}
	e.req = req
	return e
}

func (e *ElasticSearchRequest) addError(err error) {
	if e.err == nil {
		e.err = err
	}
}

// Body 将 data 放入到 http.Request 中
func (e *ElasticSearchRequest) Body(data any) *ElasticSearchRequest {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		e.req.Body = io.NopCloser(bf)
		e.req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bf), nil
		}
		e.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		e.req.Body = io.NopCloser(bf)
		e.req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bf), nil
		}
		e.req.ContentLength = int64(len(t))
	default:
		e.addError(fmt.Errorf("unsupported body data type: %s", t))
	}
	return e
}

// JSONBody 传输的变量进行 json 序列化，并添加对应的 content type
// 对于 ElasticSearch 的场景中，很多请求体都是用户自己构建的，对于 string 类型就不进行 json 序列化了
func (e *ElasticSearchRequest) JSONBody(a any) *ElasticSearchRequest {
	e.req.Header.Set("Content-Type", "application/json")
	switch a.(type) {
	case string:
		return e.Body(a)
	}
	b, err := json.Marshal(a)
	if err != nil {
		e.addError(err)
		return e
	}
	return e.Body(b)
}

// UpdateBody 给需要更新的结构体添加 doc 外壳
func (e *ElasticSearchRequest) UpdateBody(a any) *ElasticSearchRequest {
	doc := struct {
		Doc any `json:"doc"`
	}{Doc: a}
	return e.JSONBody(doc)
}

// SearchAfterBody 在原有 query 中附加 search_after 字段，方便开发者使用
// 因为在 ElasticSearch 请求的内容 json 字段的顺序是有含义的，比如排序（如果一个对象中写多个的话就可能有问题），但是 golang 的 map json 序列化是无序的，从而会导致报错
// 为了避免出现语义错误，使用了有序 map 来解决这个问题，当然如果使用者感觉比较繁琐，可直接用 Search 即可
func (e *ElasticSearchRequest) SearchAfterBody(query string, searchAfter []any) *ElasticSearchRequest {
	if len(searchAfter) == 0 {
		return e.JSONBody(query)
	}
	m := orderedmap.New()
	if err := json.UnmarshalFromString(query, &m); err != nil {
		e.addError(err)
		return e
	}
	if _, ok := m.Get("search_after"); ok {
		e.addError(errors.New("search_after already exists"))
		return e
	}
	m.Set("search_after", searchAfter)
	query, err := json.MarshalToString(m)
	if err != nil {
		e.addError(err)
		return e
	}
	return e.JSONBody(query)
}

func (e *ElasticSearchRequest) SQLBody(sql string) *ElasticSearchRequest {
	query := struct {
		Query string `json:"query"`
	}{
		Query: sql,
	}
	return e.JSONBody(query)
}

// DoRequst 发送请求，目前这个还是依赖全局的 ElasticSearch 客户端
func (e *ElasticSearchRequest) DoRequst() (*http.Response, error) {
	if e.err != nil {
		return nil, e.err
	}
	return e.do(e.req)
}

// Bytes 发送请求并读取返回值
func (e *ElasticSearchRequest) Bytes() ([]byte, error) {
	if e.body != nil {
		return e.body, nil
	}
	resp, err := e.DoRequst()
	if err != nil {
		return nil, err
	}
	e.resp = resp
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		var err error
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}
	e.body, err = io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return e.body, nil
}

// ToJSON 发送请求并将返回值反序列化为 any
func (e *ElasticSearchRequest) ToJSON(a any) error {
	b, err := e.Bytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, a); err != nil {
		return err
	}
	return nil
}

// ToError 发送请求并对于返回结果进行检查，若返回相关 ElasticSearch 的错误，也会报错
func (e *ElasticSearchRequest) ToError() error {
	b, err := e.Bytes()
	if err != nil {
		return err
	}
	errorField := json.Get(b, "error").ToString()
	if len(errorField) != 0 {
		return errors.New(string(b))
	}
	if json.Get(e.body, "errors").ToBool() {
		return errors.New(string(b))
	}
	return nil
}

// ToHits 将返回的反序列化为 []Hit[T] 数组
func (e *ElasticSearchRequest) ToHits(a any) error {
	if err := e.ToError(); err != nil {
		return err
	}
	b, err := e.Bytes()
	if err != nil {
		return err
	}
	hitsField := json.Get(b, "hits", "hits").ToString()
	if err := json.UnmarshalFromString(hitsField, a); err != nil {
		return err
	}
	return nil
}

// ToBytes 返回正常请求返回的字节流
func (e *ElasticSearchRequest) ToBytes() ([]byte, error) {
	if err := e.ToError(); err != nil {
		return nil, err
	}
	return e.Bytes()
}

func (e *ElasticSearchRequest) ToCSV(a any) error {
	b, err := e.ToBytes()
	if err != nil {
		return err
	}
	return gocsv.UnmarshalBytes(b, a)
}
