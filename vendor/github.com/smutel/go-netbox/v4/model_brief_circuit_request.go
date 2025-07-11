/*
NetBox REST API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 4.0.11 (4.0)
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package netbox

import (
	"encoding/json"
	"fmt"
)

// checks if the BriefCircuitRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BriefCircuitRequest{}

// BriefCircuitRequest Adds support for custom fields and tags.
type BriefCircuitRequest struct {
	// Unique circuit ID
	Cid                  string  `json:"cid"`
	Description          *string `json:"description,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _BriefCircuitRequest BriefCircuitRequest

// NewBriefCircuitRequest instantiates a new BriefCircuitRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBriefCircuitRequest(cid string) *BriefCircuitRequest {
	this := BriefCircuitRequest{}
	this.Cid = cid
	return &this
}

// NewBriefCircuitRequestWithDefaults instantiates a new BriefCircuitRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBriefCircuitRequestWithDefaults() *BriefCircuitRequest {
	this := BriefCircuitRequest{}
	return &this
}

// GetCid returns the Cid field value
func (o *BriefCircuitRequest) GetCid() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Cid
}

// GetCidOk returns a tuple with the Cid field value
// and a boolean to check if the value has been set.
func (o *BriefCircuitRequest) GetCidOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Cid, true
}

// SetCid sets field value
func (o *BriefCircuitRequest) SetCid(v string) {
	o.Cid = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *BriefCircuitRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BriefCircuitRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *BriefCircuitRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *BriefCircuitRequest) SetDescription(v string) {
	o.Description = &v
}

func (o BriefCircuitRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BriefCircuitRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["cid"] = o.Cid
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *BriefCircuitRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"cid",
	}

	// defaultValueFuncMap captures the default values for required properties.
	// These values are used when required properties are missing from the payload.
	defaultValueFuncMap := map[string]func() interface{}{}
	var defaultValueApplied bool
	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if value, exists := allProperties[requiredProperty]; !exists || value == "" {
			if _, ok := defaultValueFuncMap[requiredProperty]; ok {
				allProperties[requiredProperty] = defaultValueFuncMap[requiredProperty]()
				defaultValueApplied = true
			}
		}
		if value, exists := allProperties[requiredProperty]; !exists || value == "" {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	if defaultValueApplied {
		data, err = json.Marshal(allProperties)
		if err != nil {
			return err
		}
	}
	varBriefCircuitRequest := _BriefCircuitRequest{}

	err = json.Unmarshal(data, &varBriefCircuitRequest)

	if err != nil {
		return err
	}

	*o = BriefCircuitRequest(varBriefCircuitRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "cid")
		delete(additionalProperties, "description")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableBriefCircuitRequest struct {
	value *BriefCircuitRequest
	isSet bool
}

func (v NullableBriefCircuitRequest) Get() *BriefCircuitRequest {
	return v.value
}

func (v *NullableBriefCircuitRequest) Set(val *BriefCircuitRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableBriefCircuitRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableBriefCircuitRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBriefCircuitRequest(val *BriefCircuitRequest) *NullableBriefCircuitRequest {
	return &NullableBriefCircuitRequest{value: val, isSet: true}
}

func (v NullableBriefCircuitRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBriefCircuitRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
