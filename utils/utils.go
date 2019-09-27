package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

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

// StringSliceToPtrInt64Slice ...
func StringSliceToInt64Slice(src []string) []int64 {
	dst := make([]int64, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != "" {
			if n, err := strconv.Atoi(src[i]); err != nil {
				dst[i] = int64(n)
			}
		}
	}
	return dst
}
