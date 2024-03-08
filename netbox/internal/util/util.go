package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/smutel/go-netbox/v3/netbox/models"
)

func ConvertNestedASNsToASNs(asns []*models.NestedASN) []int64 {
	var tfASNs []int64

	for _, t := range asns {
		asn := t.ID

		tfASNs = append(tfASNs, asn)
	}

	return tfASNs
}

func ConvertNestedIPsToIPs(asns []*models.NestedIPAddress) []int64 {
	var tfASNs []int64

	for _, t := range asns {
		asn := t.ID

		tfASNs = append(tfASNs, asn)
	}

	return tfASNs
}

func ConvertNestedVlansToVlans(vlans []*models.NestedVLAN) []int64 {
	var tfVlans []int64

	for _, t := range vlans {
		vlan := t.ID

		tfVlans = append(tfVlans, vlan)
	}

	return tfVlans
}

func ConvertURLContentType(url string) string {
	urlSplit := strings.Split(url, "/")
	firstLevel := urlSplit[len(urlSplit)-4]
	secondLevel := urlSplit[len(urlSplit)-3]

	slRegexp := regexp.MustCompile(`s$`)
	secondLevel = strings.ReplaceAll(secondLevel, "-", "")
	secondLevel = slRegexp.ReplaceAllString(secondLevel, "")

	contentType := firstLevel + "." + secondLevel
	return contentType
}

func ExpandToInt64Slice(v []interface{}) ([]int64, error) {
	s := make([]int64, len(v))
	for i, val := range v {
		if strVal, ok := val.(int); ok {
			s[i] = int64(strVal)
		} else {
			return nil, fmt.Errorf("could not convert values to int")
		}
	}

	return s, nil
}

func ExpandToInt32Slice(v []interface{}) ([]int32, error) {
	s := make([]int32, len(v))
	for i, val := range v {
		if strVal, ok := val.(int); ok {
			s[i] = int32(strVal)
		} else {
			return nil, fmt.Errorf("could not convert values to int")
		}
	}

	return s, nil
}

func Float2stringptr(vcpus *float64) *string {
	if vcpus == nil {
		return nil
	}

	vcpustr := fmt.Sprintf("%v", *vcpus)
	return &vcpustr
}

func GetLocalContextData(data interface{}) (*string, error) {
	if data == nil {
		return nil, nil
	}
	localContextDataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	jsonstring := string(localContextDataJSON)
	return &jsonstring, nil
}

func ToListofInts(in []interface{}) []int32 {
	out := make([]int32, len(in))
	for i := range in {
		out[i] = int32(in[i].(int))
	}
	return out
}

func ToListofStrings(in []interface{}) []string {
	out := make([]string, len(in))
	for i := range in {
		out[i] = in[i].(string)
	}
	return out
}

func TrimString(val interface{}) string {
	return strings.TrimSpace(val.(string))
}

func FieldNameToStructName(k string) string {
	r := []rune(k)
	r[0] = unicode.ToUpper(r[0])
	k = string(r)

	split := strings.Split(k, "_")
	k = split[0]
	for i := 1; i < len(split); i++ {
		r2 := []rune(split[i])
		r2[0] = unicode.ToUpper(r2[0])
		k = k + string(r2)
	}

	return k
}

func GenerateErrorMessage(response *http.Response, err error) diag.Diagnostics {
	var diags diag.Diagnostics

	msg := ""
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		msg += "File: " + path.Base(file) + "\n"
		msg += "Function: " + path.Base(runtime.FuncForPC(pc).Name()) + "\n"
		msg += fmt.Sprintf("%s%d\n", "Line: ", line)
	}

	if err != nil {
		msg += "Error: " + err.Error() + "\n"
	} else {
		msg += "Error: unknown error from Netbox\n"
	}

	if response != nil {
		bytes, readErr := io.ReadAll(response.Body)
		msg += "URL: " + response.Request.URL.String() + "\n"
		msg += "Method: " + response.Request.Method + "\n"
		if readErr == nil {
			msg += "Response: " + string(bytes) + "\n"
		} else {
			msg += "Response: error (" + readErr.Error() + ")\n"
		}
	} else {
		msg += "No additional information from Netbox\n"
	}

	providerDiag := diag.Diagnostic{
		Summary: "Error with Netbox API",
		Detail:  msg,
	}
	diags = append(diags, providerDiag)

	return diags
}
