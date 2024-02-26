/*
NetBox REST API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 3.7.1 (3.7)
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package netbox

import (
	"encoding/json"
)

// checks if the DeviceFace type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DeviceFace{}

// DeviceFace struct for DeviceFace
type DeviceFace struct {
	Value                *DeviceFaceValue `json:"value,omitempty"`
	Label                *DeviceFaceLabel `json:"label,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _DeviceFace DeviceFace

// NewDeviceFace instantiates a new DeviceFace object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDeviceFace() *DeviceFace {
	this := DeviceFace{}
	return &this
}

// NewDeviceFaceWithDefaults instantiates a new DeviceFace object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDeviceFaceWithDefaults() *DeviceFace {
	this := DeviceFace{}
	return &this
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *DeviceFace) GetValue() DeviceFaceValue {
	if o == nil || IsNil(o.Value) {
		var ret DeviceFaceValue
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DeviceFace) GetValueOk() (*DeviceFaceValue, bool) {
	if o == nil || IsNil(o.Value) {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *DeviceFace) HasValue() bool {
	if o != nil && !IsNil(o.Value) {
		return true
	}

	return false
}

// SetValue gets a reference to the given DeviceFaceValue and assigns it to the Value field.
func (o *DeviceFace) SetValue(v DeviceFaceValue) {
	o.Value = &v
}

// GetLabel returns the Label field value if set, zero value otherwise.
func (o *DeviceFace) GetLabel() DeviceFaceLabel {
	if o == nil || IsNil(o.Label) {
		var ret DeviceFaceLabel
		return ret
	}
	return *o.Label
}

// GetLabelOk returns a tuple with the Label field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DeviceFace) GetLabelOk() (*DeviceFaceLabel, bool) {
	if o == nil || IsNil(o.Label) {
		return nil, false
	}
	return o.Label, true
}

// HasLabel returns a boolean if a field has been set.
func (o *DeviceFace) HasLabel() bool {
	if o != nil && !IsNil(o.Label) {
		return true
	}

	return false
}

// SetLabel gets a reference to the given DeviceFaceLabel and assigns it to the Label field.
func (o *DeviceFace) SetLabel(v DeviceFaceLabel) {
	o.Label = &v
}

func (o DeviceFace) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DeviceFace) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Value) {
		toSerialize["value"] = o.Value
	}
	if !IsNil(o.Label) {
		toSerialize["label"] = o.Label
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *DeviceFace) UnmarshalJSON(data []byte) (err error) {
	varDeviceFace := _DeviceFace{}

	err = json.Unmarshal(data, &varDeviceFace)

	if err != nil {
		return err
	}

	*o = DeviceFace(varDeviceFace)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "value")
		delete(additionalProperties, "label")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDeviceFace struct {
	value *DeviceFace
	isSet bool
}

func (v NullableDeviceFace) Get() *DeviceFace {
	return v.value
}

func (v *NullableDeviceFace) Set(val *DeviceFace) {
	v.value = val
	v.isSet = true
}

func (v NullableDeviceFace) IsSet() bool {
	return v.isSet
}

func (v *NullableDeviceFace) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDeviceFace(val *DeviceFace) *NullableDeviceFace {
	return &NullableDeviceFace{value: val, isSet: true}
}

func (v NullableDeviceFace) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDeviceFace) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
