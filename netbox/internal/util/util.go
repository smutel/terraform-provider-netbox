// Copyright (c)
// SPDX-License-Identifier: MIT

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

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/smutel/go-netbox/v4"
)

const cr string = "\n"
const empty string = ""

const Const10 = 10
const Const19 = 19
const Const21 = 21
const Const23 = 23
const Const32 = 32
const Const50 = 50
const Const64 = 64
const Const100 = 100
const Const128 = 128
const Const200 = 200
const Const201 = 201
const Const254 = 254
const Const256 = 256
const Const404 = 404
const Const1000 = 1000
const Const4094 = 4094
const Const16384 = 16384
const Const65535 = 65535
const Const65536 = 65536
const Const4294967295 = 4294967295

// const Name string = "name"
// const Type string = "type"
// const Tag string = "tag"
// const Slug string = "slug"
// const IP4ID string = "primary_ip4_id"
// const IP6ID string = "primary_ip6_id"
// const Enabled sring = "enabled"

type GenericResponse struct {
	ID     int32          `json:"id"`
	Others map[string]any `json:"-"`
}

func ConvertASNsToInts(asns []netbox.ASN) []int32 {
	var tfASNs []int32

	for _, t := range asns {
		tfASNs = append(tfASNs, t.GetId())
	}

	return tfASNs
}

func ConvertAPIVlansToVlans(vlans []netbox.VLAN) []int32 {
	var tfVlans []int32

	for _, t := range vlans {
		vlan := t.GetId()

		tfVlans = append(tfVlans, vlan)
	}

	return tfVlans
}

func ConvertURLContentType(url string) string {
	urlSplit := strings.Split(url, "/")
	firstLevel := urlSplit[len(urlSplit)-4]
	secondLevel := urlSplit[len(urlSplit)-3]

	slRegexp := regexp.MustCompile(`s$`)
	secondLevel = strings.ReplaceAll(secondLevel, "-", empty)
	secondLevel = slRegexp.ReplaceAllString(secondLevel, empty)

	contentType := firstLevel + "." + secondLevel
	return contentType
}

func ExpandToInt64Slice(v []any) ([]int64, error) {
	s := make([]int64, len(v))
	for i, val := range v {
		//nolint:revive
		if strVal, ok := val.(int); ok {
			s[i] = int64(strVal)
		} else {
			return nil, fmt.Errorf("could not convert values to int")
		}
	}

	return s, nil
}

func ExpandToInt32Slice(v []any) ([]int32, error) {
	s := make([]int32, len(v))
	for i, val := range v {
		//nolint:revive
		if intVal, ok := val.(int); ok {
			intVal32, err := safecast.ToInt32(intVal)
			if err == nil {
				s[i] = intVal32
			}
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

func GetLocalContextData(data any) (*string, error) {
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

func ToListofInts(in []any) []int32 {
	out := make([]int32, len(in))
	for i := range in {
		in32, err := safecast.ToInt32(in[i].(int))
		if err == nil {
			out[i] = in32
		}
	}
	return out
}

func ToListofStrings(in []any) []string {
	out := make([]string, len(in))
	for i := range in {
		out[i] = in[i].(string)
	}
	return out
}

func ToInterfaceArrayToInteraceArrayArray(in []any) [][]any {
	out := make([][]any, len(in))
	for i, extraChoice := range in {
		e := extraChoice.(map[string]any)

		extraChoiceArray := make([]any, 2)
		extraChoiceArray[0] = e["value"].(string)
		extraChoiceArray[1] = e["label"].(string)

		out[i] = extraChoiceArray
	}
	return out
}

func ToInterfaceArrayArrayToInteraceArray(in [][]any) []any {
	out := make([]any, len(in))

	for i, extraChoiceArray := range in {
		e := make(map[string]any)
		e["value"] = extraChoiceArray[0]
		e["label"] = extraChoiceArray[1]
		out[i] = e
	}

	return out
}

func TrimString(val any) string {
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

	msg := empty
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		msg += "File: " + path.Base(file) + cr
		msg += "Function: " + path.Base(runtime.FuncForPC(pc).Name()) + cr
		msg += fmt.Sprintf("%s%d"+cr, "Line: ", line)
	}

	if err != nil {
		msg += "Error: " + err.Error() + cr
	} else {
		msg += "Error: unknown error from Netbox" + cr
	}

	if response != nil {
		bytes, readErr := io.ReadAll(response.Body)
		msg += "URL: " + response.Request.URL.String() + cr
		msg += "Method: " + response.Request.Method + cr
		if readErr == nil {
			msg += "Response: " + string(bytes) + cr
		} else {
			msg += "Response: error (" + readErr.Error() + ")" + cr
		}
	} else {
		msg += "No additional information from Netbox" + cr
	}

	providerDiag := diag.Diagnostic{
		Summary: "Error with Netbox API",
		Detail:  msg,
	}
	diags = append(diags, providerDiag)

	return diags
}

func UnmarshalID(body io.ReadCloser) (int32, error) {
	byteValue, err := io.ReadAll(body)
	var gr GenericResponse

	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(byteValue, &gr); err != nil {
		return 0, err
	}

	return gr.ID, nil
}
