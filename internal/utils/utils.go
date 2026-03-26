package utils

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const (
	SuffixConfigFilePath string = "/.osc/config.json"
)

func InterfaceSliceToStringSlicePtr(slice []any) *[]string {
	result := InterfaceSliceToStringSlice(slice)
	return &result
}

func SetToStringSlice(set *schema.Set) []string {
	return InterfaceSliceToStringSlice(set.List())
}

func SetToStringSlicePtr(set *schema.Set) *[]string {
	return InterfaceSliceToStringSlicePtr(set.List())
}

func SetToSuperStringSlice[T ~string](set *schema.Set) []T {
	return SliceToSuperStringSlice[T](set.List())
}

func SliceToSuperStringSlice[T ~string](slice []any) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		val, ok := v.(string)
		if ok && val != "" {
			result = append(result, T(val))
		}
	}
	return result
}

func InterfaceSliceToStringSlice(slice []any) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		val, ok := v.(string)
		if ok && val != "" {
			result = append(result, val)
		}
	}
	return result
}

func InterfaceSliceToStringList(slice []any) *[]string {
	res := InterfaceSliceToStringSlice(slice)
	return &res
}

func StringSlicePtrToInterfaceSlice(list *[]string) []any {
	if list == nil {
		return make([]any, 0)
	}
	vs := make([]any, 0, len(*list))
	for _, v := range *list {
		vs = append(vs, v)
	}
	return vs
}

func UnknownDataSourceFilterError(filterName string) error {
	return fmt.Errorf("datasource filter '%s' is not implemented in the provider or not supported by the api", filterName)
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
	for i := range src {
		if src[i] != nil {
			if n, err := strconv.Atoi(ptr.From(src[i])); err != nil {
				dst[i] = new(int64(n))
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

func SliceToTftypesValueSlice(src []string) (basetypes.SetValue, diag.Diagnostics) {
	nSet := lo.Map(src, func(str string, _ int) attr.Value {
		return types.StringValue(str)
	})
	setValue, diags := types.SetValue(basetypes.StringType{}, nSet)
	if diags != nil {
		return setValue, diags
	}
	return setValue, diags
}

// StringSliceToInt32Slice converts []string to []int32 ...
func StringSliceToInt32Slice(src []string) (res []int32) {
	for _, str := range src {
		res = append(res, cast.ToInt32(str))
	}
	return
}

func StringSliceToIntSlice(src []string) (res []int) {
	for _, str := range src {
		res = append(res, cast.ToInt(str))
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
	log.Printf("\n[WARN] %s %s not found, probably deleted manually, removing from state\n", name, id)
}

func IsResponseEmpty(len int, name, id string) bool {
	if len == 0 {
		LogManuallyDeleted(name, id)
		return true
	}
	return false
}

func RandIntRange(min, max int) int {
	return min + rand.IntN(max-min) //nolint:gosec
}

func ParsingfilterToDateFormat(filterName, value string) (time.Time, error) {
	var err error
	var filterDate iso8601.Time

	if value != "" {
		if filterDate, err = iso8601.ParseString(value); err != nil {
			return filterDate.Time, fmt.Errorf("%s value should be 'ISO 8601' format ('2017-06-14' or '2017-06-14T00:00:00Z, ...) %w", filterName, err)
		}
	}
	return filterDate.Time, nil
}

func StringSliceToTimeSlice(filterValues []string, filterName string) ([]time.Time, error) {
	var sliceDates []time.Time
	for val := range filterValues {
		valDate, err := ParsingfilterToDateFormat(filterName, filterValues[val])
		if err != nil {
			return sliceDates, err
		}
		sliceDates = append(sliceDates, valDate)
	}
	return sliceDates, nil
}

func FiltersTimesToStringSlice(filterValues []string, filterName string) ([]string, error) {
	var sliceString []string
	for val := range filterValues {
		valDate, err := ParsingfilterToDateFormat(filterName, filterValues[val])
		if err != nil {
			return sliceString, err
		}
		sliceString = append(sliceString, valDate.String())
	}
	return sliceString, nil
}

func I32toa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func GetRegion() string {
	region := os.Getenv("OUTSCALE_REGION")
	if region == "" {
		region = os.Getenv("OSC_REGION")
	}
	return region
}
