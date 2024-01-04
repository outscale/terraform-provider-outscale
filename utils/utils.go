package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"
)

// PrintToJSON method helper to debug responses
const (
	randMin         float32 = 1.0
	randMax         float32 = 20.0
	MinPort         int     = 1
	MaxPort         int     = 65535
	MinIops         int     = 100
	MaxIops         int     = 13000
	DefaultIops     int32   = 150
	MaxSize         int     = 14901
	InvalidState    string  = "InvalidState"
	VolumeIOPSError string  = `
- The "iops" parameter can only be set if "io1" volume type is created.
- "Standard" volume types have a default value of 150 iops.
- For "gp2" volume types, iops value depend on your volume size.
`
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

func ParseStringToInt32(str string) int32 {
	parsed, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(parsed)
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

func LogManuallyDeleted(name, id string) {
	log.Printf("[WARN] %s %s not found, probably deleted manually, removing from state", name, id)
}

func IsResponseEmpty(len int, name, id string) bool {
	if len == 0 {
		LogManuallyDeleted(name, id)
		return true
	}
	return false
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

func CheckThrottling(httpResp *http.Response, err error) *resource.RetryError {
	rand.Seed(time.Now().UnixNano())
	if httpResp != nil {
		errCode := httpResp.StatusCode
		if errCode == http.StatusServiceUnavailable || errCode == http.StatusTooManyRequests ||
			errCode == http.StatusConflict || errCode == http.StatusFailedDependency {
			randTime := (rand.Float32()*(randMax-randMin) + randMin) * 1000
			time.Sleep(time.Duration(randTime) * time.Millisecond)
			return resource.RetryableError(err)
		}
	}
	return resource.NonRetryableError(err)
}

func RandIntRange(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandVpcCidr() string {
	var result string
	prefix := RandIntRange(16, 29)
	switch rand.Intn(3) {
	case 0:
		//10.0.0.0 - 10.255.255.255 (10/8 prefix)
		result = fmt.Sprintf("10.%d.0.0/%d", rand.Intn(256), prefix)
	case 1:
		//172.16.0.0 - 172.31.255.255 (172.16/12 prefix)
		result = fmt.Sprintf("172.%d.0.0/%d", RandIntRange(16, 32), prefix)
	case 2:
		//192.168.0.0 - 192.168.255.255 (192.168/16 prefix)
		result = fmt.Sprintf("192.168.0.0/%d", prefix)
	}
	return result
}
func InterfaceSliceToStringSlicePtr(slice []interface{}) *[]string {
	result := InterfaceSliceToStringSlice(slice)
	return &result
}

func SetToStringSlice(set *schema.Set) []string {
	return InterfaceSliceToStringSlice(set.List())
}

func SetToStringSlicePtr(set *schema.Set) *[]string {
	return InterfaceSliceToStringSlicePtr(set.List())
}

func InterfaceSliceToStringSlice(slice []interface{}) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		val, ok := v.(string)
		if ok && val != "" {
			result = append(result, v.(string))
		}
	}
	return result
}

func InterfaceSliceToStringList(slice []interface{}) *[]string {
	res := InterfaceSliceToStringSlice(slice)
	return &res
}

func StringSlicePtrToInterfaceSlice(list *[]string) []interface{} {
	if list == nil {
		return make([]interface{}, 0)
	}
	vs := make([]interface{}, 0, len(*list))
	for _, v := range *list {
		vs = append(vs, v)
	}
	return vs
}

func I32toa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func GetRegion() string {
	region := fmt.Sprintf("%s", os.Getenv("OUTSCALE_REGION"))
	if region == "" {
		region = fmt.Sprintf("%s", os.Getenv("OSC_REGION"))
	}
	return region
}

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func Strings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}

func GetEnvVariableValue(envVariables []string) string {

	for _, envVariable := range envVariables {
		if value := os.Getenv(envVariable); value != "" {
			return value
		}
	}
	return ""
}
