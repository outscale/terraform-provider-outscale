package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
)

// PrintToJSON method helper to debug responses
func PrintToJSON(v interface{}, msg string) {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	fmt.Print("\n\n[DEBUG] ", msg, string(pretty))
}

// DebugRequest ...
func DebugRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string("####################"))
	fmt.Println(string("###### REQUEST #######"))
	fmt.Println(string(requestDump))
}

// DebugResponse ...
func DebugResponse(req *http.Response) {
	requestDump, err := httputil.DumpResponse(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string("####################"))
	fmt.Println(string("###### RESPONSE #######"))
	fmt.Println(string(requestDump))
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
