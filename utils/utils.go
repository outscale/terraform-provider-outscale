package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
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

// StringSliceToInt64Slice ...
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
