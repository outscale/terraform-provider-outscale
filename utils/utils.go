package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// PrintToJSON method helper to debug responses
const (
	ResourceNotFound string  = "InvalidResource"
	ResourceConflict string  = "Conflict"
	InvalidState     string  = "InvalidState"
	Throttled        string  = "Request rate exceeded"
	randMin          float32 = 1.0
	randMax          float32 = 20.0
)

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

// StringSliceToFloat32Slice converts []string to []float32 ...
func StringSliceToFloat32Slice(src []string) (res []float32) {
	for _, str := range src {
		res = append(res, cast.ToFloat32(str))
	}
	return
}

func IsResponseEmptyOrMutiple(rLen int, resName string) error {
	if rLen == 0 {
		return fmt.Errorf("Unable to find %v", resName)
	}
	if rLen > 1 {
		return fmt.Errorf("Multiple %vs matched; use additional constraints to reduce matches to a single %v", resName, resName)
	}
	return nil
}

func CheckThrottling(err error) *resource.RetryError {
	rand.Seed(time.Now().UnixNano())
	if strings.Contains(err.Error(), Throttled) {
		randTime := (rand.Float32()*(randMax-randMin) + randMin) * 1000
		time.Sleep(time.Duration(randTime) * time.Millisecond)
		return resource.RetryableError(err)
	}
	return resource.NonRetryableError(err)
}
