/*
NetBox REST API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 4.0.3 (4.0)
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package netbox

import (
	"encoding/json"
	"fmt"
)

// DcimConsolePortTemplatesListTypeParameter the model 'DcimConsolePortTemplatesListTypeParameter'
type DcimConsolePortTemplatesListTypeParameter string

// List of dcim_console_port_templates_list_type_parameter
const (
	DCIMCONSOLEPORTTEMPLATESLISTTYPEPARAMETER_OTHER  DcimConsolePortTemplatesListTypeParameter = "Other"
	DCIMCONSOLEPORTTEMPLATESLISTTYPEPARAMETER_SERIAL DcimConsolePortTemplatesListTypeParameter = "Serial"
	DCIMCONSOLEPORTTEMPLATESLISTTYPEPARAMETER_USB    DcimConsolePortTemplatesListTypeParameter = "USB"
)

// All allowed values of DcimConsolePortTemplatesListTypeParameter enum
var AllowedDcimConsolePortTemplatesListTypeParameterEnumValues = []DcimConsolePortTemplatesListTypeParameter{
	"Other",
	"Serial",
	"USB",
}

func (v *DcimConsolePortTemplatesListTypeParameter) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DcimConsolePortTemplatesListTypeParameter(value)
	for _, existing := range AllowedDcimConsolePortTemplatesListTypeParameterEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DcimConsolePortTemplatesListTypeParameter", value)
}

// NewDcimConsolePortTemplatesListTypeParameterFromValue returns a pointer to a valid DcimConsolePortTemplatesListTypeParameter
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewDcimConsolePortTemplatesListTypeParameterFromValue(v string) (*DcimConsolePortTemplatesListTypeParameter, error) {
	ev := DcimConsolePortTemplatesListTypeParameter(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for DcimConsolePortTemplatesListTypeParameter: valid values are %v", v, AllowedDcimConsolePortTemplatesListTypeParameterEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v DcimConsolePortTemplatesListTypeParameter) IsValid() bool {
	for _, existing := range AllowedDcimConsolePortTemplatesListTypeParameterEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to dcim_console_port_templates_list_type_parameter value
func (v DcimConsolePortTemplatesListTypeParameter) Ptr() *DcimConsolePortTemplatesListTypeParameter {
	return &v
}

type NullableDcimConsolePortTemplatesListTypeParameter struct {
	value *DcimConsolePortTemplatesListTypeParameter
	isSet bool
}

func (v NullableDcimConsolePortTemplatesListTypeParameter) Get() *DcimConsolePortTemplatesListTypeParameter {
	return v.value
}

func (v *NullableDcimConsolePortTemplatesListTypeParameter) Set(val *DcimConsolePortTemplatesListTypeParameter) {
	v.value = val
	v.isSet = true
}

func (v NullableDcimConsolePortTemplatesListTypeParameter) IsSet() bool {
	return v.isSet
}

func (v *NullableDcimConsolePortTemplatesListTypeParameter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDcimConsolePortTemplatesListTypeParameter(val *DcimConsolePortTemplatesListTypeParameter) *NullableDcimConsolePortTemplatesListTypeParameter {
	return &NullableDcimConsolePortTemplatesListTypeParameter{value: val, isSet: true}
}

func (v NullableDcimConsolePortTemplatesListTypeParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDcimConsolePortTemplatesListTypeParameter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
