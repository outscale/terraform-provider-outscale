package logging

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type smithyTFLogMiddleware struct{}

type readSeekCloser struct {
	*bytes.Reader
}

func (*readSeekCloser) Close() error {
	return nil
}

var _ middleware.DeserializeMiddleware = (*smithyTFLogMiddleware)(nil)

func WithSmithyTFLogMiddleware() config.LoadOptionsFunc {
	return func(options *config.LoadOptions) error {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			return stack.Deserialize.Add(&smithyTFLogMiddleware{}, middleware.After)
		})
		return nil
	}
}

func (*smithyTFLogMiddleware) ID() string {
	return "TFLogRequestResponseLogger"
}

func (m *smithyTFLogMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (out middleware.DeserializeOutput, metadata middleware.Metadata, err error) {
	request, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unexpected HTTP request type %T", in.Request)
	}

	httpRequest := request.Build(ctx)

	err = m.logRequest(ctx, httpRequest)
	if err != nil {
		return out, metadata, err
	}

	request, err = request.SetStream(httpRequest.Body)
	if err != nil {
		return out, metadata, err
	}
	in.Request = request

	start := time.Now()
	out, metadata, err = next.HandleDeserialize(ctx, in)
	elapsed := time.Since(start)
	if err != nil {
		return out, metadata, err
	}

	response, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		return out, metadata, fmt.Errorf("unexpected HTTP response type %T", out.RawResponse)
	}
	err = m.logResponse(ctx, response.Response, elapsed)
	if err != nil {
		return out, metadata, err
	}

	return out, metadata, nil
}

func (*smithyTFLogMiddleware) logRequest(ctx context.Context, request *http.Request) error {
	var body []byte
	if request.Body != nil && request.Body != http.NoBody {
		switch middleware.GetOperationName(ctx) {
		case "PutObject", "UploadPart":
			body = redactedBody(request.ContentLength, request.Header.Get("Content-Type"))
		default:
			var err error
			body, err = io.ReadAll(request.Body)
			if err != nil {
				return err
			}
			_ = request.Body.Close()

			request.Body = &readSeekCloser{
				Reader: bytes.NewReader(body),
			}
		}
	}

	headers := request.Header.Clone()
	if request.Host != "" {
		headers.Set("Host", request.Host)
	}
	addContentHeaders(headers, request.ContentLength, request.TransferEncoding)
	fields := map[string]any{
		"req": formatHTTP(request.Method+" "+request.URL.String(), headers, body),
	}
	addOperationField(ctx, fields)

	tflog.Debug(ctx, "SDK HTTP request", fields)
	return nil
}

func (*smithyTFLogMiddleware) logResponse(ctx context.Context, response *http.Response, duration time.Duration) error {
	var body []byte
	if response.Body != nil && response.Body != http.NoBody {
		switch middleware.GetOperationName(ctx) {
		case "GetObject":
			body = redactedBody(response.ContentLength, response.Header.Get("Content-Type"))
		default:
			var err error
			body, err = io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			_ = response.Body.Close()

			response.Body = io.NopCloser(bytes.NewReader(body))
		}
	}

	headers := response.Header.Clone()
	addContentHeaders(headers, response.ContentLength, response.TransferEncoding)

	responseLog := formatHTTP(fmt.Sprintf("%s %s\nStatus: %d (%s)", response.Request.Method, response.Request.URL, response.StatusCode, duration), headers, body)
	fields := map[string]any{
		"resp": responseLog,
	}
	addOperationField(ctx, fields)

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		tflog.Error(ctx, "SDK HTTP response error", fields)
		return nil
	}

	tflog.Debug(ctx, "SDK HTTP response", fields)

	return nil
}

func redactedBody(contentLength int64, contentType string) []byte {
	result := fmt.Sprintf("[Redacted: %d bytes", contentLength)

	if contentType != "" {
		result += ", Type: " + contentType
	}

	return []byte(result + "]")
}

func addContentHeaders(headers http.Header, contentLength int64, transferEncoding []string) {
	if contentLength >= 0 {
		headers.Set("Content-Length", strconv.FormatInt(contentLength, 10))
	}
	if len(transferEncoding) > 0 {
		headers["Transfer-Encoding"] = transferEncoding
	}
}

func addOperationField(ctx context.Context, fields map[string]any) {
	if operation := middleware.GetOperationName(ctx); operation != "" {
		fields["sdk_operation"] = operation
	}
}

func formatHTTP(summary string, headers http.Header, body []byte) string {
	var output strings.Builder
	output.WriteString(summary)
	output.WriteByte('\n')
	_ = headers.Write(&output)

	if len(body) > 0 {
		output.WriteByte('\n')
		output.WriteString(formatXMLBody(body))
	}

	return strings.ReplaceAll(output.String(), "\r\n", "\n")
}

func formatXMLBody(body []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(body)) //nolint:gosec
	var output bytes.Buffer
	encoder := xml.NewEncoder(&output)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return string(body)
		}

		// Remove "xmlns" namespace from XML elements for readability
		switch element := token.(type) {
		case xml.StartElement:
			element.Name.Space = ""
			token = element
		case xml.EndElement:
			element.Name.Space = ""
			token = element
		}

		if err := encoder.EncodeToken(token); err != nil {
			return string(body)
		}
	}
	if err := encoder.Flush(); err != nil {
		return string(body)
	}

	return strings.TrimSpace(output.String())
}
