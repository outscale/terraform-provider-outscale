package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type tflogWrapper struct {
}

func NewTflogWrapper() *tflogWrapper {
	return &tflogWrapper{}
}

func removeNulls(f any) any {
	switch fields := f.(type) {
	case map[string]any:
		cleaned := make(map[string]any)
		for k, v := range fields {
			if v == nil {
				continue
			}
			cleaned[k] = removeNulls(v)
		}
		return cleaned
	case []any:
		for i, v := range fields {
			fields[i] = removeNulls(v)
		}
		return fields
	default:
		return fields
	}
}

func (t *tflogWrapper) RequestHttp(ctx context.Context, req *http.Request) {
	fields := map[string]any{
		"method": req.Method,
		"url":    req.URL.String(),
	}

	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err == nil {
			bodyBytes, _ := io.ReadAll(bodyReader)
			if len(bodyBytes) > 0 {
				fields["body"] = string(bodyBytes)
				var jsonData any
				if json.Unmarshal(bodyBytes, &jsonData) == nil {
					// cleaned := removeNulls(jsonData)
					cleaned := jsonData
					if indentJSON, err := json.MarshalIndent(cleaned, "", "  "); err == nil {
						fields["body"] = string(indentJSON)
					}
				}
			}
		}
	}

	tflog.Debug(ctx, "SDK HTTP request", fields)
}

func (t *tflogWrapper) ResponseHttp(ctx context.Context, resp *http.Response, d time.Duration) {
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	fields := map[string]any{
		"status_code": resp.StatusCode,
		"duration":    d.String(),
	}

	if len(bodyBytes) > 0 {
		fields["body"] = string(bodyBytes)
		var jsonData any
		if json.Unmarshal(bodyBytes, &jsonData) == nil {
			if prettyJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
				fields["body"] = string(prettyJSON)
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "SDK HTTP response error", fields)
	} else {
		tflog.Debug(ctx, "SDK HTTP response", fields)
	}
}

func (t *tflogWrapper) Request(ctx context.Context, req any) {
	tflog.Trace(ctx, "SDK request", map[string]any{
		"body": req,
	})
}

func (t *tflogWrapper) Response(ctx context.Context, resp any) {
	tflog.Trace(ctx, "SDK response", map[string]any{
		"body": resp,
	})
}

func (t *tflogWrapper) Error(ctx context.Context, err error) {
	tflog.Error(ctx, "SDK error", map[string]any{
		"error": err.Error(),
	})
}
