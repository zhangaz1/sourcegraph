package httpclientutil

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/neelance/parallel"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sourcegraph/sourcegraph/internal/metrics"
	"github.com/sourcegraph/sourcegraph/internal/trace/ot"
	"golang.org/x/net/context/ctxhttp"
)

type Client struct {
	httpClient  *http.Client
	httpLimiter *parallel.Run
	userAgent   string
	options     ClientOptions
}

type ClientOptions struct {
	Name                string
	MaxIdleConnsPerHost int
	CategoryFunc        func(*url.URL) string
}

func NewClient(options ClientOptions) *Client {
	httpClient := newMeteredHTTPClient(
		options.MaxIdleConnsPerHost,
		options.CategoryFunc,
	)

	return &Client{
		httpClient:  httpClient,
		httpLimiter: parallel.NewRun(500), // TODO - configure
		userAgent:   filepath.Base(os.Args[0]),
		options:     options,
	}
}

func newMeteredHTTPClient(maxIdleConnsPerHost int, categoryFunc func(*url.URL) string) *http.Client {
	requestMeter := metrics.NewRequestMeter(
		"precise_code_intel_bundle_manager",                                   // TODO
		"Total number of requests sent to precise code intel bundle manager.", // TODO
	)

	return &http.Client{
		Transport: &ot.Transport{
			RoundTripper: requestMeter.Transport(
				&http.Transport{
					MaxIdleConnsPerHost: maxIdleConnsPerHost,
				},
				categoryFunc,
			),
		},
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (_ *http.Response, err error) {
	span, ctx := ot.StartSpanFromContext(ctx, fmt.Sprintf("%s.do", c.options.Name))
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.SetTag("err", err.Error())
		}
		span.Finish()
	}()

	req.Header.Set("User-Agent", c.userAgent)
	req = req.WithContext(ctx)

	if c.httpLimiter != nil {
		span.LogKV("event", "Waiting on HTTP limiter")
		c.httpLimiter.Acquire()
		defer c.httpLimiter.Release()
		span.LogKV("event", "Acquired HTTP limiter")
	}

	req, ht := nethttp.TraceRequest(
		span.Tracer(),
		req,
		nethttp.OperationName("Precise Code Intel Bundle Manager Client"), // TODO
		nethttp.ClientTrace(false),
	)
	defer ht.Finish()

	resp, err := ctxhttp.Do(req.Context(), c.httpClient, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
