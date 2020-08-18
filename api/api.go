package api

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/guomio/go-template/tools"
	"github.com/imroc/req"
)

const (
	queryTag  = "json"
	fieldName = "file"
	get       = "GET"
	post      = "POST"
	delete    = "DELETE"
	put       = "PUT"
)

// TransformRequestFunc TransformRequest 类型
type TransformRequestFunc = func(data interface{}) interface{}

// TransformResponseFunc TransformResponse 类型
type TransformResponseFunc = func(data interface{})

// HeadersFunc headers 类型
type HeadersFunc = func() req.Header

func defaultHeaders() req.Header {
	return req.Header{}
}

func defaultTransformRequest(data interface{}) interface{} {
	return data
}

func defaultTransformResponse(data interface{}) {
}

// API 配置项
type API struct {
	BaseURL           string
	TransformRequest  TransformRequestFunc
	TransformResponse TransformResponseFunc
	Headers           HeadersFunc
}

// Request 参数
type Request struct {
	API
	Method string
	URL    string
	Data   interface{}
	Params interface{}
	Header interface{}
	File   string
	Extra  []interface{}
}

// NewAPI 新建实例
func NewAPI(opt *API) *API {
	if opt.Headers == nil {
		opt.Headers = defaultHeaders
	}
	if opt.TransformRequest == nil {
		opt.TransformRequest = defaultTransformRequest
	}
	if opt.TransformResponse == nil {
		opt.TransformResponse = defaultTransformResponse
	}

	return opt
}

// ReflectHeader header struct 转 map[string]string
func ReflectHeader(header interface{}) req.Header {
	h := make(req.Header)

	if header == nil {
		return h
	}

	elem := reflect.ValueOf(header).Elem()
	elemType := elem.Type()
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		key := field.Name
		if tag := field.Tag.Get(queryTag); tag != "" {
			key = tag
		}
		h[key] = elem.Field(i).String()
	}
	return h
}

// QueryURL URL 添加 query 参数
func QueryURL(URL string, query interface{}) string {
	if query == nil {
		return URL
	}
	qs := []string{}
	param := ReflectHeader(query)

	for k, v := range param {
		qs = append(qs, k+"="+v)
	}

	if strings.Contains(URL, "?") {
		return URL + "&" + strings.Join(qs, "&")
	}
	return URL + "?" + strings.Join(qs, "&")
}

// Field 上传文件
func Field(pattern, FieldName string) []req.FileUpload {
	matches := []string{}
	uploads := []req.FileUpload{}

	m, err := filepath.Glob(pattern)
	if err != nil {
		return uploads
	}
	matches = append(matches, m...)

	if len(matches) == 0 {
		return uploads
	}
	for _, match := range matches {
		if s, e := os.Stat(match); e != nil || s.IsDir() {
			continue
		}
		file, _ := os.Open(match)
		uploads = append(uploads, req.FileUpload{
			File:      file,
			FileName:  filepath.Base(match),
			FieldName: FieldName,
		})
	}

	return uploads
}

// File FieldName 默认 file
func File(pattern string) []req.FileUpload {
	return Field(pattern, fieldName)
}

// Request 发送请求
func (a *API) Request(inter interface{}, r *Request) (err error) {
	URL := QueryURL(tools.CombineURLs(tools.EmptyToString(r.BaseURL, a.BaseURL), r.URL), r.Params)

	var res *req.Resp

	transformRequest := a.TransformRequest

	if r.TransformRequest != nil {
		transformRequest = r.TransformRequest
	}

	transformResponse := a.TransformResponse

	if r.TransformResponse != nil {
		transformResponse = r.TransformResponse
	}

	rs := []interface{}{}

	headers := a.Headers()

	if r.Headers != nil {
		headers = r.Headers()
	}

	if r.Header != nil {
		header := ReflectHeader(r.Header)

		for k, v := range header {
			headers[k] = v
		}

		rs = append(rs, headers)
	}

	if r.Data != nil {
		rs = append(rs, req.BodyJSON(transformRequest(r.Data)))
	}

	if r.File != "" {
		rs = append(rs, File(r.File))
	}

	if r.Extra != nil {
		rs = append(rs, r.Extra...)
	}

	switch r.Method {
	case get:
		res, err = req.Get(URL, rs...)
		break
	case delete:
		res, err = req.Delete(URL, rs...)
		break
	case post:
		res, err = req.Post(URL, rs...)
		break
	case put:
		res, err = req.Put(URL, rs...)
		break
	default:
		err = errors.New("unknow http method")
	}

	if err != nil {
		return
	}

	err = res.ToJSON(inter)

	if err != nil {
		return
	}

	transformResponse(inter)

	return
}

// Get 请求
func (a *API) Get(inter interface{}, r *Request) error {
	r.Method = get
	return a.Request(inter, r)
}

// Post 请求
func (a *API) Post(inter interface{}, r *Request) error {
	r.Method = post
	return a.Request(inter, r)
}

// Put 请求
func (a *API) Put(inter interface{}, r *Request) error {
	r.Method = put
	return a.Request(inter, r)
}

// Delete 请求
func (a *API) Delete(inter interface{}, r *Request) error {
	r.Method = delete
	return a.Request(inter, r)
}
