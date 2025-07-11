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

// checks if the NestedSiteGroupRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NestedSiteGroupRequest{}

// NestedSiteGroupRequest Represents an object related through a ForeignKey field. On write, it accepts a primary key (PK) value or a dictionary of attributes which can be used to uniquely identify the related object. This class should be subclassed to return a full representation of the related object on read.
type NestedSiteGroupRequest struct {
	Name                 string `json:"name"`
	Slug                 string `json:"slug" validate:"regexp=^[-a-zA-Z0-9_]+$"`
	AdditionalProperties map[string]interface{}
}

type _NestedSiteGroupRequest NestedSiteGroupRequest

// NewNestedSiteGroupRequest instantiates a new NestedSiteGroupRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNestedSiteGroupRequest(name string, slug string) *NestedSiteGroupRequest {
	this := NestedSiteGroupRequest{}
	this.Name = name
	this.Slug = slug
	return &this
}

// NewNestedSiteGroupRequestWithDefaults instantiates a new NestedSiteGroupRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNestedSiteGroupRequestWithDefaults() *NestedSiteGroupRequest {
	this := NestedSiteGroupRequest{}
	return &this
}

// GetName returns the Name field value
func (o *NestedSiteGroupRequest) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *NestedSiteGroupRequest) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *NestedSiteGroupRequest) SetName(v string) {
	o.Name = v
}

// GetSlug returns the Slug field value
func (o *NestedSiteGroupRequest) GetSlug() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Slug
}

// GetSlugOk returns a tuple with the Slug field value
// and a boolean to check if the value has been set.
func (o *NestedSiteGroupRequest) GetSlugOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Slug, true
}

// SetSlug sets field value
func (o *NestedSiteGroupRequest) SetSlug(v string) {
	o.Slug = v
}

func (o NestedSiteGroupRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NestedSiteGroupRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["slug"] = o.Slug

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *NestedSiteGroupRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"slug",
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
	varNestedSiteGroupRequest := _NestedSiteGroupRequest{}

	err = json.Unmarshal(data, &varNestedSiteGroupRequest)

	if err != nil {
		return err
	}

	*o = NestedSiteGroupRequest(varNestedSiteGroupRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "slug")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableNestedSiteGroupRequest struct {
	value *NestedSiteGroupRequest
	isSet bool
}

func (v NullableNestedSiteGroupRequest) Get() *NestedSiteGroupRequest {
	return v.value
}

func (v *NullableNestedSiteGroupRequest) Set(val *NestedSiteGroupRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableNestedSiteGroupRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableNestedSiteGroupRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNestedSiteGroupRequest(val *NestedSiteGroupRequest) *NullableNestedSiteGroupRequest {
	return &NullableNestedSiteGroupRequest{value: val, isSet: true}
}

func (v NullableNestedSiteGroupRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNestedSiteGroupRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
