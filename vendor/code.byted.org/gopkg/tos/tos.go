package tos

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	TosAccessHeader       = "X-Tos-Access"
	TosRestoreHeader      = "X-Tos-Restore"
	TosStorageClassHeader = "X-Tos-Storage-Class"
	TosCopySrcHeader      = "X-Tos-Copy-Source"
	TosCopyDstHeader      = "X-Tos-Copy-Dst"
	TosCopyAccessHeader   = "X-Tos-Copy-Access"
	StorageClassStandard  = "Standard"
	StorageClassArchive   = "Archive"
	HttpHeaderUserAgent   = "User-Agent"
	defaultServiceName    = "toutiao.tos.tosapi"
	MinPartSize           = 5 * 1024 * 1024
	Version               = "v1.0.4"

	TosCopyDstBucket = "dstBucket"
	TosCopyDstObject = "dstObject"
)

var DefaultReqTimeout = 10 * time.Second
var userAgentInfo string

func init() {
	initUserAgentInfo()
}

type options struct {
	Cluster  string
	Bucket   string
	Token    string
	IDC      string
	Endpoint string
}

type Option func(o *options)

/********** User can set Content-Type/MD5/Meta info Start **********/
const (
	HTTPHeaderContentType = "Content-Type"
)

type (
	ObjOptions struct {
		Value interface{}
	}

	ObjOption func(map[string]ObjOptions) error
)

func ContentType(value string) ObjOption {
	return setHeader(HTTPHeaderContentType, value)
}

func setHeader(key string, value interface{}) ObjOption {
	return func(params map[string]ObjOptions) error {
		if value == nil {
			return nil
		}
		params[key] = ObjOptions{value}
		return nil
	}
}

func handleOptions(header http.Header, options []ObjOption) error {
	params := map[string]ObjOptions{}
	for _, option := range options {
		if option != nil {
			if err := option(params); err != nil {
				return err
			}
		}
	}

	for k, v := range params {
		header.Set(k, v.Value.(string))
	}
	return nil
}

/********** User can set Content-Type/MD5/Meta info End   **********/

func WithCluster(cluster string) Option {
	return func(o *options) {
		o.Cluster = cluster
	}
}

//Endpoint example: tos-cn-north.byted.org
func WithEndpoint(endpoint string) Option {
	return func(o *options) {
		o.Endpoint = endpoint
	}
}

func WithBucket(bucket string) Option {
	return func(o *options) {
		o.Bucket = bucket
	}
}

func WithToken(token string) Option {
	return func(o *options) {
		o.Token = token
	}
}

func WithAuth(bucket, token string) Option {
	return func(o *options) {
		o.Bucket = bucket
		o.Token = token
	}
}

func WithIDC(idc string) Option {
	return func(o *options) {
		o.IDC = idc
	}
}

type Tos struct {
	opts       options
	httpClient *httpClient
}

// NewTos return a tos client instance
func NewTos(ops ...Option) (*Tos, error) {
	cli := &Tos{}
	for _, op := range ops {
		op(&cli.opts)
	}
	if cli.opts.Cluster == "" {
		cli.opts.Cluster = "default"
	}
	if cli.opts.Bucket == "" {
		return nil, errors.New("bucket not set")
	}

	httpClient, err := newHttpClient(cli.opts.Cluster, cli.opts.IDC, cli.opts.Endpoint)
	if err != nil {
		return nil, err
	}
	cli.httpClient = httpClient

	return cli, nil
}

func (t *Tos) makeuri(object string) string {
	if isEndpointValidDomain(t.opts.Endpoint) {
		return "http://" + t.opts.Bucket + "." + t.opts.Endpoint + "/" + object
	} else if t.opts.Endpoint != "" {
		return "http://" + t.opts.Endpoint + "/" + t.opts.Bucket + "/" + object
	}
	name := defaultServiceName
	if t.opts.IDC != "" {
		name += ".service." + t.opts.IDC
	}
	ret := "http://" + name + "/" + t.opts.Bucket
	if object == "" {
		return ret
	}
	return ret + "/" + object
}

type ObjectInfo struct {
	R       io.ReadCloser
	Size    int64
	MTime   time.Time
	Headers http.Header
}

func (t *Tos) doReq(ctx context.Context, method, o string, body io.Reader, ops ...ObjOption) (
	*http.Response, error) {
	url := t.makeuri(o)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		nbody, ok := body.(httpbody)
		if ok {
			req.ContentLength = nbody.ContentLength()
		}
	}

	/*  */
	handleOptions(req.Header, ops)
	/*  */
	return t.doHttpReq(ctx, req)
}

func (t *Tos) doHttpReq(ctx context.Context, req *http.Request) (*http.Response, error) {
	timeout := DefaultReqTimeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = -time.Since(deadline)
	}
	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&timeout=" + timeout.String()
	} else {
		req.URL.RawQuery += "timeout=" + timeout.String()
	}
	req.Header.Set(TosAccessHeader, t.opts.Token)
	req.Header.Set(HttpHeaderUserAgent, userAgent())
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			req.Host = connInfo.Conn.RemoteAddr().String()
		},
	}
	ctx = httptrace.WithClientTrace(ctx, trace)
	resp, err := t.httpClient.do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	resp.Request.Host = req.Host
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, DecodeErr(resp)
	}
	return resp, nil
}

func (t *Tos) GetObject(ctx context.Context, object string) (*ObjectInfo, error) {
	resp, err := t.doReq(ctx, "GET", object, nil)
	if err != nil {
		return nil, err
	}
	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	mtime, _ := time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
	return &ObjectInfo{
		R:       resp.Body,
		Size:    size,
		MTime:   mtime,
		Headers: resp.Header,
	}, nil
}

func (t *Tos) GetObjectFromOffset(ctx context.Context, o string, off int64) (*ObjectInfo, error) {
	req, _ := http.NewRequest("GET", t.makeuri(o), nil)
	req.Header.Set("Range", "bytes="+strconv.FormatInt(off, 10)+"-")
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 206 {
		return nil, errors.New("expect http 206")
	}
	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	mtime, _ := time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
	return &ObjectInfo{
		R:       resp.Body,
		Size:    size,
		MTime:   mtime,
		Headers: resp.Header,
	}, nil
}

func (t *Tos) GetObjectFromRange(ctx context.Context, o string, start, end int64) (*ObjectInfo, error) {
	req, _ := http.NewRequest("GET", t.makeuri(o), nil)
	req.Header.Set("Range", strings.Join([]string{"bytes=",
		strconv.FormatInt(start, 10), "-", strconv.FormatInt(end, 10)}, ""))

	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 206 {
		return nil, errors.New("expect http 206")
	}
	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	mtime, _ := time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
	return &ObjectInfo{
		R:       resp.Body,
		Size:    size,
		MTime:   mtime,
		Headers: resp.Header,
	}, nil
}

func (t *Tos) HttpForward(ctx context.Context, o string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" && r.Method != "HEAD" {
		return errors.New("method not allowed")
	}
	req, _ := http.NewRequest(r.Method, t.makeuri(o), nil)
	for k, v := range r.Header {
		req.Header[k] = v
	}
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	wheader := w.Header()
	for k, vv := range resp.Header {
		if _, ok := wheader[k]; ok {
			continue // not overwrite
		}
		for _, v := range vv {
			wheader.Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	return nil
}

// HeadObject return an object's meta info
func (t *Tos) HeadObject(ctx context.Context, object string) (*ObjectInfo, error) {
	resp, err := t.doReq(ctx, "HEAD", object, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // it is safe for HEAD method
	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	mtime, _ := time.Parse(http.TimeFormat, resp.Header.Get("Last-Modified"))
	return &ObjectInfo{
		R:       nil,
		Size:    size,
		MTime:   mtime,
		Headers: resp.Header,
	}, nil
}

type httpbody interface {
	Read(p []byte) (int, error)
	ContentLength() int64
}

type withContentLengthReader struct {
	R io.Reader
	N int64
}

func (r *withContentLengthReader) ContentLength() int64 {
	return r.N
}

func (r *withContentLengthReader) Read(p []byte) (int, error) {
	return r.R.Read(p)
}

// PutObject write an object into the storage server
func (t *Tos) PutObject(ctx context.Context, object string, size int64, r io.Reader, ops ...ObjOption) error {
	body := &withContentLengthReader{R: r, N: size}
	resp, err := t.doReq(ctx, "PUT", object, body, ops...)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// DelObject delete an object
func (t *Tos) DelObject(ctx context.Context, object string) error {
	resp, err := t.doReq(ctx, "DELETE", object, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (t *Tos) RestoreObject(ctx context.Context, object string) error {
	req, _ := http.NewRequest(http.MethodPost, t.makeuri(object)+"?restore", nil)
	_, err := t.doHttpReq(ctx, req)
	return err
}

// InitUpload init an upload session
func (t *Tos) InitUpload(ctx context.Context, object string, ops ...ObjOption) (string, error) {
	req, _ := http.NewRequest("POST", t.makeuri(object)+"?uploads", nil)

	// UploadPartCopyTo, should assign dst
	dstBucket, dstBucketExist := ctx.Value(TosCopyDstBucket).(string)
	dstObject, dstObjectExist := ctx.Value(TosCopyDstObject).(string)
	if dstBucketExist && dstObjectExist {
		req.Header.Set(TosCopyDstHeader, "/"+dstBucket+"/"+dstObject)
	}
	/*  */
	handleOptions(req.Header, ops)
	/*  */
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var d struct {
		Success int `json:"success"`
		Payload struct {
			UploadID string `json:"uploadID"`
		} `json:"payload"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&d); err != nil {
		return "", err
	}
	return d.Payload.UploadID, nil
}

type Part struct {
	PartID string `json:"partID"`
	Etag   string `json:"etag"`
}

// UploadPart upload a part, then get a result.
func (t *Tos) UploadPart(ctx context.Context, object, uploadID string, index int, data []byte) (*Part, error) {
	if index < 0 || index > 9999 {
		return nil, errors.New("index range[0, 10000]")
	}
	partID := strconv.Itoa(index)
	m := md5.New()
	uri := t.makeuri(object) + fmt.Sprintf("?partNumber=%s&uploadID=%s", partID, uploadID)
	req, _ := http.NewRequest("PUT", uri, io.TeeReader(bytes.NewReader(data), m))
	req.ContentLength = int64(len(data))
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	md5Hash := hex.EncodeToString(m.Sum(nil))
	if md5Hash != resp.Header.Get("X-Tos-MD5") {
		return nil, ErrChecksum
	}
	return &Part{
		PartID: partID,
		Etag:   resp.Header.Get("X-Tos-ETag"),
	}, nil
}

// CompleteUpload
func (t *Tos) CompleteUpload(ctx context.Context, object, uploadID string, parts []Part) error {
	var partList []string
	for _, part := range parts {
		if part.Etag != "" {
			partList = append(partList, part.PartID+":"+part.Etag)
		} else {
			partList = append(partList, part.PartID)
		}
	}
	body := bytes.NewBufferString(strings.Join(partList, ","))
	req, _ := http.NewRequest("POST", t.makeuri(object)+"?uploadID="+uploadID, body)

	// UploadPartCopyTo, should assign dst
	dstBucket, dstBucketExist := ctx.Value("dstBucket").(string)
	dstObject, dstObjectExist := ctx.Value("dstObject").(string)
	if dstBucketExist && dstObjectExist {
		req.Header.Set(TosCopyDstHeader, "/"+dstBucket+"/"+dstObject)
	}

	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ListParts return a list of all parts have uploaded
func (t *Tos) ListParts(ctx context.Context, object, uploadID string) ([]Part, error) {
	req, _ := http.NewRequest("GET", t.makeuri(object)+"?uploadID="+uploadID, nil)
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res struct {
		Success int `json:"success"`
		Payload struct {
			UploadID string `json:"uploadID"`
			Parts    []Part `json:"parts"`
		} `json:"payload"`
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res.Payload.Parts, nil
}

// AbortUpload abort an upload session with the uploadID
func (t *Tos) AbortUpload(ctx context.Context, object, uploadID string) error {
	req, _ := http.NewRequest("DELETE", t.makeuri(object)+"?uploadID="+uploadID, nil)
	res, err := t.doHttpReq(ctx, req)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}

type ListPrefixInput struct {
	Prefix     string
	Delimiter  string
	StartAfter string
	MaxKeys    int
}

type ListObject struct {
	Key          string `json:"key"`
	LastModified string `json:"lastModified"`
	Size         int64  `json:"size"`
}

type ListPrefixOutput struct {
	IsTruncated  bool         `json:"isTruncated"` // HasMore
	CommonPrefix []string     `json:"commonPrefix"`
	Objects      []ListObject `json:"objects"`
	StartAfter   string       `json:"startAfter"`
}

func (t *Tos) ListPrefix(ctx context.Context, input ListPrefixInput) (*ListPrefixOutput, error) {
	uv := url.Values{}
	uv.Set("prefix", input.Prefix)
	uv.Set("delimiter", input.Delimiter)
	uv.Set("start-after", input.StartAfter)
	uv.Set("max-keys", strconv.Itoa(input.MaxKeys))
	req, _ := http.NewRequest("GET", t.makeuri("")+"?"+uv.Encode(), nil)
	res, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var d struct {
		Success int              `json:"success"`
		Payload ListPrefixOutput `json:"payload"`
	}
	err = json.NewDecoder(res.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	return &d.Payload, nil
}

func (t *Tos) RemoveAll(ctx context.Context, prefix string) error {
	if prefix == "" || prefix == "/" {
		return errors.New("prefix not allowed")
	}
	hasmore := true
	startAfter := ""
	for hasmore {
		resp, err := t.ListPrefix(ctx, ListPrefixInput{Prefix: prefix, StartAfter: startAfter, MaxKeys: 100})
		if err != nil {
			return err
		}
		for _, commonPrefix := range resp.CommonPrefix {
			if err := t.RemoveAll(ctx, commonPrefix); err != nil {
				return err
			}
		}
		for _, o := range resp.Objects {
			if err := t.DelObject(ctx, o.Key); err != nil {
				return err
			}
		}
		hasmore, startAfter = resp.IsTruncated, resp.StartAfter
	}
	return nil
}

// CopyObject copy srcObject to a new object dstObject
// run in one bucket
func (t *Tos) CopyObject(ctx context.Context, srcObject, dstObject string) error {
	req, _ := http.NewRequest(http.MethodPost, t.makeuri(dstObject)+"?copyobject", nil)
	req.Header.Set(TosCopySrcHeader, "/"+t.opts.Bucket+"/"+srcObject)
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// CopyObjectFrom copy srcObject from srcBucket to localObject, which in the initial bucket
func (t *Tos) CopyObjectFrom(ctx context.Context, srcObject, dstObject, srcBucket, srcBucketToken string) error {
	req, _ := http.NewRequest(http.MethodPost, t.makeuri(dstObject)+"?copyobjectfrom", nil)
	req.Header.Set(TosCopySrcHeader, "/"+srcBucket+"/"+srcObject)
	req.Header.Set(TosCopyAccessHeader, srcBucketToken)
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// CopyObjectTo copy srcObject from initial bucket to dstBucket, named by dstObject
func (t *Tos) CopyObjectTo(ctx context.Context, srcObject, dstObject, dstBucket, dstBucketToken string) error {
	req, _ := http.NewRequest(http.MethodPost, t.makeuri(srcObject)+"?copyobjectto", nil)
	req.Header.Set(TosCopyDstHeader, "/"+dstBucket+"/"+dstObject)
	req.Header.Set(TosCopyAccessHeader, dstBucketToken)
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (t *Tos) UploadPartCopy(ctx context.Context, srcObject, dstObject, uploadID string, index int, startOffset, partSize uint64) (*Part, error) {
	if index < 0 || index > 9999 {
		return nil, errors.New("index range[0, 10000]")
	}
	partID := strconv.Itoa(index)
	uri := t.makeuri(dstObject) + fmt.Sprintf("?partNumber=%s&uploadID=%s&startOffset=%d&partSize=%d", partID, uploadID, startOffset, partSize)
	req, _ := http.NewRequest("PUT", uri, nil)
	req.Header.Set(TosCopySrcHeader, "/"+t.opts.Bucket+"/"+srcObject)
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Part{
		PartID: partID,
		Etag:   resp.Header.Get("X-Tos-ETag"),
	}, nil
}

func (t *Tos) UploadPartCopyFrom(ctx context.Context, srcObject, srcBucket, srcBucketToken, dstObject, uploadID string, index int, startOffset, partSize uint64) (*Part, error) {
	if index < 0 || index > 9999 {
		return nil, errors.New("index range[0, 10000]")
	}
	partID := strconv.Itoa(index)
	uri := t.makeuri(dstObject) + fmt.Sprintf("?partNumber=%s&uploadID=%s&startOffset=%d&partSize=%d", partID, uploadID, startOffset, partSize)
	req, _ := http.NewRequest("PUT", uri, nil)
	req.Header.Set(TosCopySrcHeader, "/"+srcBucket+"/"+srcObject)
	req.Header.Set(TosCopyAccessHeader, srcBucketToken)
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Part{
		PartID: partID,
		Etag:   resp.Header.Get("X-Tos-ETag"),
	}, nil
}

func (t *Tos) UploadPartCopyTo(ctx context.Context, srcObject, dstObject, dstBucket, dstBucketToken, uploadID string, index int, startOffset, partSize uint64) (*Part, error) {
	if index < 0 || index > 9999 {
		return nil, errors.New("index range[0, 10000]")
	}
	partID := strconv.Itoa(index)
	uri := t.makeuri(srcObject) + fmt.Sprintf("?partNumber=%s&uploadID=%s&startOffset=%d&partSize=%d", partID, uploadID, startOffset, partSize)
	req, _ := http.NewRequest("PUT", uri, nil)
	req.Header.Set(TosCopyDstHeader, "/"+dstBucket+"/"+dstObject)
	req.Header.Set(TosCopyAccessHeader, dstBucketToken)
	resp, err := t.doHttpReq(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &Part{
		PartID: partID,
		Etag:   resp.Header.Get("X-Tos-ETag"),
	}, nil
}

func userAgent() string {
	return userAgentInfo
}

func initUserAgentInfo() {
	name := runtime.GOOS
	release := "-"
	machine := runtime.GOARCH
	if out, err := exec.Command("uname", "-s").CombinedOutput(); err == nil {
		name = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command("uname", "-r").CombinedOutput(); err == nil {
		release = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command("uname", "-m").CombinedOutput(); err == nil {
		machine = string(bytes.TrimSpace(out))
	}

	userAgentInfo = fmt.Sprintf("tos-sdk-go/%s (%s/%s/%s;%s)", Version, name,
		release, machine, runtime.Version())
}
