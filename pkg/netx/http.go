package netx

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/51xpx/pkgx/pkg/jsonx"
)

var (
	errStatusCodeNot200 = errors.New("Status code != 200")
)

var (
	maxWaitTime = time.Second * 15
	retryMax    = 3
)

// RetryClient retriable client
type RetryClient struct {
	RetryMax     int
	RetryWaitMax time.Duration
	client       *http.Client
}

// NewRetryClient new retry client
func NewRetryClient() *RetryClient {
	client := retryablehttp.NewClient()
	client.RetryMax = retryMax
	client.RetryWaitMax = maxWaitTime
	return &RetryClient{
		RetryMax:     retryMax,
		RetryWaitMax: maxWaitTime,
		client:       client.StandardClient(),
	}
}

// Get .
func (c *RetryClient) Get(url string, ctx *gin.Context) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// if app.Config.ThirdService.JaegerTrace.IsOpen == 1 {
	if true {
		tracer, _ := ctx.Get("Tracer")
		parentSpanContext, _ := ctx.Get("ParentSpanContext")

		span := opentracing.StartSpan(
			url,
			opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)),
			opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			ext.SpanKindRPCClient,
			// ext.HTTPUrl,
			// ext.SpanKindRPCClient,
		)

		span.Finish()

		injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if injectErr != nil {
			log.Fatalf("%s: Couldn't inject headers", injectErr)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errStatusCodeNot200
	}

	return string(body), nil
}

// Post .
func (c *RetryClient) Post(url, contentType string, bodyContent io.Reader) (string, error) {
	resp, err := c.client.Post(url, contentType, bodyContent)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errStatusCodeNot200
	}

	return string(body), nil
}

// PostJSON .
func (c *RetryClient) PostJSON(url string, params map[string]interface{}) (string, error) {
	jsonContent, err := jsonx.Marshal(params)
	if err != nil {
		return "", err
	}

	return c.Post(url, "application/Json", bytes.NewBuffer(jsonContent))
}

// PostForm .
func (c *RetryClient) PostForm(url string, bodyContent url.Values) (string, error) {
	resp, err := c.client.PostForm(url, bodyContent)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errStatusCodeNot200
	}

	return string(body), nil
}
