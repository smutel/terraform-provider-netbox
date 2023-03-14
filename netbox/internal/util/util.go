package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-openapi/strfmt"
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

// Convert URL in content_type
func ConvertURIContentType(uri strfmt.URI) string {
	uriSplit := strings.Split(uri.String(), "/")
	firstLevel := uriSplit[len(uriSplit)-4]
	secondLevel := uriSplit[len(uriSplit)-3]

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

func ToListofInts(in []interface{}) []int64 {
	out := make([]int64, len(in))
	for i := range in {
		out[i] = int64(in[i].(int))
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
	k = strings.Replace(k, "_", "", -1)
	return k
}
