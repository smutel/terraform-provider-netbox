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

// checks if the BriefVLANGroup type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BriefVLANGroup{}

// BriefVLANGroup Adds support for custom fields and tags.
type BriefVLANGroup struct {
	Id                   int32   `json:"id"`
	Url                  string  `json:"url"`
	Display              string  `json:"display"`
	Name                 string  `json:"name"`
	Slug                 string  `json:"slug" validate:"regexp=^[-a-zA-Z0-9_]+$"`
	Description          *string `json:"description,omitempty"`
	VlanCount            *int64  `json:"vlan_count,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _BriefVLANGroup BriefVLANGroup

// NewBriefVLANGroup instantiates a new BriefVLANGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBriefVLANGroup(id int32, url string, display string, name string, slug string) *BriefVLANGroup {
	this := BriefVLANGroup{}
	this.Id = id
	this.Url = url
	this.Display = display
	this.Name = name
	this.Slug = slug
	return &this
}

// NewBriefVLANGroupWithDefaults instantiates a new BriefVLANGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBriefVLANGroupWithDefaults() *BriefVLANGroup {
	this := BriefVLANGroup{}
	return &this
}

// GetId returns the Id field value
func (o *BriefVLANGroup) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *BriefVLANGroup) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *BriefVLANGroup) SetId(v int32) {
	o.Id = v
}

// GetUrl returns the Url field value
func (o *BriefVLANGroup) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *BriefVLANGroup) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *BriefVLANGroup) SetUrl(v string) {
	o.Url = v
}

// GetDisplay returns the Display field value
func (o *BriefVLANGroup) GetDisplay() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Display
}

// GetDisplayOk returns a tuple with the Display field value
// and a boolean to check if the value has been set.
func (o *BriefVLANGroup) GetDisplayOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Display, true
}

// SetDisplay sets field value
func (o *BriefVLANGroup) SetDisplay(v string) {
	o.Display = v
}

// GetName returns the Name field value
func (o *BriefVLANGroup) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *BriefVLANGroup) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *BriefVLANGroup) SetName(v string) {
	o.Name = v
}

// GetSlug returns the Slug field value
func (o *BriefVLANGroup) GetSlug() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Slug
}

// GetSlugOk returns a tuple with the Slug field value
// and a boolean to check if the value has been set.
func (o *BriefVLANGroup) GetSlugOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Slug, true
}

// SetSlug sets field value
func (o *BriefVLANGroup) SetSlug(v string) {
	o.Slug = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *BriefVLANGroup) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BriefVLANGroup) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *BriefVLANGroup) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *BriefVLANGroup) SetDescription(v string) {
	o.Description = &v
}

// GetVlanCount returns the VlanCount field value if set, zero value otherwise.
func (o *BriefVLANGroup) GetVlanCount() int64 {
	if o == nil || IsNil(o.VlanCount) {
		var ret int64
		return ret
	}
	return *o.VlanCount
}

// GetVlanCountOk returns a tuple with the VlanCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BriefVLANGroup) GetVlanCountOk() (*int64, bool) {
	if o == nil || IsNil(o.VlanCount) {
		return nil, false
	}
	return o.VlanCount, true
}

// HasVlanCount returns a boolean if a field has been set.
func (o *BriefVLANGroup) HasVlanCount() bool {
	if o != nil && !IsNil(o.VlanCount) {
		return true
	}

	return false
}

// SetVlanCount gets a reference to the given int64 and assigns it to the VlanCount field.
func (o *BriefVLANGroup) SetVlanCount(v int64) {
	o.VlanCount = &v
}

func (o BriefVLANGroup) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BriefVLANGroup) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["url"] = o.Url
	toSerialize["display"] = o.Display
	toSerialize["name"] = o.Name
	toSerialize["slug"] = o.Slug
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.VlanCount) {
		toSerialize["vlan_count"] = o.VlanCount
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *BriefVLANGroup) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"url",
		"display",
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
	varBriefVLANGroup := _BriefVLANGroup{}

	err = json.Unmarshal(data, &varBriefVLANGroup)

	if err != nil {
		return err
	}

	*o = BriefVLANGroup(varBriefVLANGroup)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "url")
		delete(additionalProperties, "display")
		delete(additionalProperties, "name")
		delete(additionalProperties, "slug")
		delete(additionalProperties, "description")
		delete(additionalProperties, "vlan_count")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableBriefVLANGroup struct {
	value *BriefVLANGroup
	isSet bool
}

func (v NullableBriefVLANGroup) Get() *BriefVLANGroup {
	return v.value
}

func (v *NullableBriefVLANGroup) Set(val *BriefVLANGroup) {
	v.value = val
	v.isSet = true
}

func (v NullableBriefVLANGroup) IsSet() bool {
	return v.isSet
}

func (v *NullableBriefVLANGroup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBriefVLANGroup(val *BriefVLANGroup) *NullableBriefVLANGroup {
	return &NullableBriefVLANGroup{value: val, isSet: true}
}

func (v NullableBriefVLANGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBriefVLANGroup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
