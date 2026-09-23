package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/outscale/goutils/sdk/sanitize"
)

type tflogWrapper struct {
	sanitize bool
}

func NewTflogWrapper(sanitize bool) *tflogWrapper {
	return &tflogWrapper{
		sanitize: sanitize,
	}
}

func (t *tflogWrapper) RequestHttp(ctx context.Context, req *http.Request) {
	reqStr := req.Method + " " + req.URL.String()

	if req.GetBody != nil {
		var bodyReader io.ReadCloser
		var err error
		if t.sanitize {
			req = sanitize.HTTPRequest(req)
			// req is already cloned, we can take the body directly
			bodyReader = req.Body
		} else {
			bodyReader, err = req.GetBody()
		}
		if err == nil && bodyReader != nil {
			bodyBytes, readErr := io.ReadAll(bodyReader)
			_ = bodyReader.Close()
			if readErr == nil && len(bodyBytes) > 0 {
				bodyStr := string(bodyBytes)
				var jsonData any
				if json.Unmarshal(bodyBytes, &jsonData) == nil {
					if indentJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
						bodyStr = string(indentJSON)
					}
				}
				reqStr += "\n\n" + bodyStr
			}
		}
	}

	tflog.Debug(ctx, "SDK HTTP request", map[string]any{"req": reqStr})
}

func (t *tflogWrapper) ResponseHttp(ctx context.Context, resp *http.Response, d time.Duration) {
	if t.sanitize {
		resp = sanitize.HTTPResponse(resp)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	respStr := resp.Request.Method + " " + resp.Request.URL.String() +
		"\nStatus: " + strconv.Itoa(resp.StatusCode) + " (" + d.String() + ")"

	if len(bodyBytes) > 0 {
		bodyStr := string(bodyBytes)
		var jsonData any
		if json.Unmarshal(bodyBytes, &jsonData) == nil {
			if prettyJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
				bodyStr = string(prettyJSON)
			}
		}
		respStr += "\n\n" + bodyStr
	}

	fields := map[string]any{"resp": respStr}
	if resp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "SDK HTTP response error", fields)
	} else {
		tflog.Debug(ctx, "SDK HTTP response", fields)
	}
}

func (t *tflogWrapper) Request(ctx context.Context, req any) {
	if t.sanitize {
		req = sanitize.Sanitize(req)
	}
	tflog.Trace(ctx, "SDK request", map[string]any{
		"body": req,
	})
}

func (t *tflogWrapper) Response(ctx context.Context, resp any) {
	if t.sanitize {
		resp = sanitize.Sanitize(resp)
	}
	if jsonBytes, err := json.Marshal(resp); err == nil {
		var jsonData any
		if json.Unmarshal(jsonBytes, &jsonData) == nil {
			tflog.Trace(ctx, "SDK response", map[string]any{
				"body": jsonData,
			})
			return
		}
	}

	tflog.Trace(ctx, "SDK response", map[string]any{
		"body": resp,
	})
}

func (t *tflogWrapper) Error(ctx context.Context, err error) {
	errMsg := err.Error()
	if t.sanitize {
		errMsg = sanitize.String(errMsg)
	}

	tflog.Error(ctx, "SDK error", map[string]any{
		"error": errMsg,
	})
}
