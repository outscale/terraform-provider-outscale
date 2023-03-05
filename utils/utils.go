package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/spf13/cast"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// PrintToJSON method helper to debug responses
const (
	randMin         float32 = 1.0
	randMax         float32 = 20.0
	MinPort         int     = 1
	MaxPort         int     = 65535
	MinIops         int     = 100
	MaxIops         int     = 13000
	DefaultIops     int     = 150
	MaxSize         int     = 14901
	InvalidState    string  = "InvalidState"
	VolumeIOPSError string  = `
- The "iops" parameter can only be set if "io1" volume type is created.
- "Standard" volume types have a default value of 150 iops.
- For "gp2" volume types, iops value depend on your volume size.
`
	TestCaPem string = `-----BEGIN CERTIFICATE-----
MIIFzzCCA7egAwIBAgIUG4MiEOc008Sc0VFUYiiPud2MINEwDQYJKoZIhvcNAQEL
BQAwdzELMAkGA1UEBhMCRlIxDjAMBgNVBAgMBVBhcmlzMREwDwYDVQQKDAhvdXRz
Y2FsZTEYMBYGA1UEAwwPT1VBVFRBUkEgdGhpZXJ5MSswKQYJKoZIhvcNAQkBFhx0
aGllcnkub3VhdHRhcmFAb3V0c2NhbGUuY29tMB4XDTIyMDgwODA4MTA0MVoXDTIz
MDgwODA4MTA0MVowdzELMAkGA1UEBhMCRlIxDjAMBgNVBAgMBVBhcmlzMREwDwYD
VQQKDAhvdXRzY2FsZTEYMBYGA1UEAwwPT1VBVFRBUkEgdGhpZXJ5MSswKQYJKoZI
hvcNAQkBFhx0aGllcnkub3VhdHRhcmFAb3V0c2NhbGUuY29tMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAmw8MV5gj8nkSYwHSk9BSFahuLwr4hg/4eAn+
+2vlb3efl0NqIsZg9GF9blI9gFbf47U6QN3bqzigCODBnRkfgMRtHWqyzDRzsIJG
RamD73L6goat/Eg9Gm35L03rWX0fYfxXA4gebVvjrqiPYAWnGrCIWuZPAlWYe8DL
7SBCdr9r6or2y5fpILkm+2Ngem6yJvMFUCKO1cr/K3h3UDMZfEgvkXxFDG4g/og+
Hpgmezx6h1GTgzflxgQRXLTDjo9pHpqMHU38jDIYb5ne0F3iCWL9DAXI/LZ40RLy
jY/yzIPZzBXZANhwMwS0R/AL0dkCl/09+8GvCgQmN2ftxOFE56tasbYT4Z/Ibk6h
XvkbgcExmifkKYo+5iBs4N6UydQazvQrjJOBtuPFro0YfmCnkQK4W0WaPH0ujitJ
UWQ5fsKqlI0V/dkFFF3C/wkuhAXTWreoHm78ikqgt2hLl8Jpcn2e7J6yv/odD0dL
uZ3QZxmdjzTtOXi1ORkOC6jkVJM7mxacSHelMNT5THGgpq9w8iYvsqCcJWOwebAi
70nXNHs8CRr2h1tXJxVJUgXh7gHhlgWV73DpcNq8VsmDHdgob2waWF1X2f/1J2gP
zNGexNXG3MCkYJZFJoKbizCl7AD6zkom7bM5b/kk7m3DYCQKBf9rz2JjGiyRhliJ
5eIXUuMCAwEAAaNTMFEwHQYDVR0OBBYEFD0uPSRhYoH5zdwc9oy5YVDfgYrzMB8G
A1UdIwQYMBaAFD0uPSRhYoH5zdwc9oy5YVDfgYrzMA8GA1UdEwEB/wQFMAMBAf8w
DQYJKoZIhvcNAQELBQADggIBAFcIa1IEzoUvKZjDNOH7RqnefykwIIDfkIt2xysu
D2uTsh4dNN6EKf3bF3xk7d3p6vq6tqsR2/m0ERsgnWY48uterqzP7uE7Fw+OcR7i
V6PnMuxb/ihOTFFB6iNQuuJfKUYlqN8HGXRXXQbqHnpy/U8lC8tMqyxfykM+GZmP
VoR4040ne2+3l8LASMXi1o0PNwBUnQen5oQc090Girfrt/j6n2ggupaV8ee8Nppp
JV2gXAGZirUqBZfjUsNx0tF5W19JSQrYYAZtfl5QD2e5djAxJpS+J9myHsEQjiTj
b/UpFUOJVnmlvUurmajaBG54ybY6rV4Ai/IHb8M8XmFAaDxpzgbqucFR8F0Fatg4
IRYvAy7+juFZEjzuytoNA6KNbbZ//6bxUuTAum9DJ2b5S7/qwHk/ajlvd6ihPi9I
Fd/94aZvBzRUe/KQ0XfhQ2rE5qv35uwuwUWtVLnVC6C2zKrBSNgjuWbv7ObJ6F/S
oZozgmcdXa0aEtYlJgYXJcGsco2gkKwgBs3AutjGZL1sp/EnjdKffFLLIh6zGKsr
IGal1JuTKfbI41vbQoayesGWcvpQXASqZftRmtaYE1mWUW8z5dGaezms/EUHKRrh
NeSlzb9gysR43vScnnOZN4rdsj/hAc91HynLSBEr05LGKGOw+sFfBOkOlZaTUEvy
iRHX
-----END CERTIFICATE-----
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
	region := fmt.Sprintf("%s", os.Getenv("TESTACC_REGION"))
	if region == "" {
		region = "eu-west-2"
	}
	return region
}
