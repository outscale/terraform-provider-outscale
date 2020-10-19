package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	oscgo "github.com/outscale/osc-sdk-go/osc"
	"github.com/spf13/cast"

	"github.com/aws/aws-sdk-go/aws"
)

// PrintToJSON method helper to debug responses
func PrintToJSON(v interface{}, msg string) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	fmt.Print("\n\n[DEBUG] ", msg, string(pretty))
}

func ToJSONString(v interface{}) string {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	return string(pretty)
}

func GetErrorResponse(err error) error {
	if e, ok := err.(oscgo.GenericOpenAPIError); ok {
		if errorResponse, oker := e.Model().(oscgo.ErrorResponse); oker {
			return fmt.Errorf("%s %s", err, ToJSONString(errorResponse))
		}
	}
	return err
}

// StringSliceToPtrInt64Slice ...
func StringSliceToPtrInt64Slice(src []*string) []*int64 {
	dst := make([]*int64, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			if n, err := strconv.Atoi(aws.StringValue(src[i])); err != nil {
				dst[i] = aws.Int64(int64(n))
			}
		}
	}
	return dst
}

// StringSliceToInt64Slice converts []string to []int64 ...
func StringSliceToInt64Slice(src []string) (res []int64) {
	for _, str := range src {
		res = append(res, cast.ToInt64(str))
	}
	return
}

// StringSliceToInt32Slice converts []string to []int32 ...
func StringSliceToInt32Slice(src []string) (res []int32) {
	for _, str := range src {
		res = append(res, cast.ToInt32(str))
	}
	return
}
