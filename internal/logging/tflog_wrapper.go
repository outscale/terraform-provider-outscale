package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type tflogWrapper struct{}

func NewTflogWrapper() *tflogWrapper {
	return &tflogWrapper{}
}

func (t *tflogWrapper) RequestHttp(ctx context.Context, req *http.Request) {
	reqStr := req.Method + " " + req.URL.String()

	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err == nil {
			bodyBytes, _ := io.ReadAll(bodyReader)
			if len(bodyBytes) > 0 {
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
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	respStr := resp.Request.Method + " " + resp.Request.URL.String() +
		"\nStatus: " + strconv.Itoa(resp.StatusCode) + " (" + d.String() + ")"

	if len(bodyBytes) > 0 {
		bodyStr := string(bodyBytes)
		var jsonData any
		if json.Unmarshal(bodyBytes, &jsonData) == nil {
			jsonData = maskSensitiveValues(jsonData)
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
	tflog.Trace(ctx, "SDK request", map[string]any{
		"body": req,
	})
}

func (t *tflogWrapper) Response(ctx context.Context, resp any) {
	if jsonBytes, err := json.Marshal(resp); err == nil {
		var jsonData any
		if json.Unmarshal(jsonBytes, &jsonData) == nil {
			tflog.Trace(ctx, "SDK response", map[string]any{
				"body": maskSensitiveValues(jsonData),
			})
			return
		}
	}

	tflog.Trace(ctx, "SDK response", map[string]any{
		"body": resp,
	})
}

func (t *tflogWrapper) Error(ctx context.Context, err error) {
	tflog.Error(ctx, "SDK error", map[string]any{
		"error": err.Error(),
	})
}

// Masks values of sensitive keys of JSON data.
// It supports masking the values based on the key name,
// and subvalues stored in a string (kubeconfig case).
func maskSensitiveValues(data any) any {
	mask := "(sensitive)"
	sensitiveKeys := []string{
		"client-certificate-data",
		"client-key-data",
		"certificate-authority-data",
		"token",
	}

	sensitiveValueRegexes := make([]*regexp.Regexp, len(sensitiveKeys))
	for i, key := range sensitiveKeys {
		pattern := regexp.QuoteMeta(key) + `:\s*\S+`
		sensitiveValueRegexes[i] = regexp.MustCompile(pattern)
	}

	switch v := data.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, value := range v {
			if slices.Contains(sensitiveKeys, key) {
				result[key] = mask
			} else if str, ok := value.(string); ok {
				maskedStr := str
				for _, regex := range sensitiveValueRegexes {
					maskedStr = regex.ReplaceAllStringFunc(maskedStr, func(match string) string {
						parts := strings.SplitN(match, ":", 2)
						if len(parts) == 2 {
							return parts[0] + ": " + mask
						}
						return mask
					})
				}
				result[key] = maskedStr
			} else {
				result[key] = maskSensitiveValues(value)
			}
		}
		return result
	case []any:
		for i, item := range v {
			v[i] = maskSensitiveValues(item)
		}
		return v
	default:
		return v
	}
}
