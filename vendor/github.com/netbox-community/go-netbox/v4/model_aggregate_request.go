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

// checks if the AggregateRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AggregateRequest{}

// AggregateRequest Adds support for custom fields and tags.
type AggregateRequest struct {
	Prefix               string                 `json:"prefix"`
	Rir                  RIRRequest             `json:"rir"`
	Tenant               NullableTenantRequest  `json:"tenant,omitempty"`
	DateAdded            NullableString         `json:"date_added,omitempty"`
	Description          *string                `json:"description,omitempty"`
	Comments             *string                `json:"comments,omitempty"`
	Tags                 []NestedTagRequest     `json:"tags,omitempty"`
	CustomFields         map[string]interface{} `json:"custom_fields,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _AggregateRequest AggregateRequest

// NewAggregateRequest instantiates a new AggregateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAggregateRequest(prefix string, rir RIRRequest) *AggregateRequest {
	this := AggregateRequest{}
	this.Prefix = prefix
	this.Rir = rir
	return &this
}

// NewAggregateRequestWithDefaults instantiates a new AggregateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAggregateRequestWithDefaults() *AggregateRequest {
	this := AggregateRequest{}
	return &this
}

// GetPrefix returns the Prefix field value
func (o *AggregateRequest) GetPrefix() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Prefix
}

// GetPrefixOk returns a tuple with the Prefix field value
// and a boolean to check if the value has been set.
func (o *AggregateRequest) GetPrefixOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Prefix, true
}

// SetPrefix sets field value
func (o *AggregateRequest) SetPrefix(v string) {
	o.Prefix = v
}

// GetRir returns the Rir field value
func (o *AggregateRequest) GetRir() RIRRequest {
	if o == nil {
		var ret RIRRequest
		return ret
	}

	return o.Rir
}

// GetRirOk returns a tuple with the Rir field value
// and a boolean to check if the value has been set.
func (o *AggregateRequest) GetRirOk() (*RIRRequest, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Rir, true
}

// SetRir sets field value
func (o *AggregateRequest) SetRir(v RIRRequest) {
	o.Rir = v
}

// GetTenant returns the Tenant field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *AggregateRequest) GetTenant() TenantRequest {
	if o == nil || IsNil(o.Tenant.Get()) {
		var ret TenantRequest
		return ret
	}
	return *o.Tenant.Get()
}

// GetTenantOk returns a tuple with the Tenant field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AggregateRequest) GetTenantOk() (*TenantRequest, bool) {
	if o == nil {
		return nil, false
	}
	return o.Tenant.Get(), o.Tenant.IsSet()
}

// HasTenant returns a boolean if a field has been set.
func (o *AggregateRequest) HasTenant() bool {
	if o != nil && o.Tenant.IsSet() {
		return true
	}

	return false
}

// SetTenant gets a reference to the given NullableTenantRequest and assigns it to the Tenant field.
func (o *AggregateRequest) SetTenant(v TenantRequest) {
	o.Tenant.Set(&v)
}

// SetTenantNil sets the value for Tenant to be an explicit nil
func (o *AggregateRequest) SetTenantNil() {
	o.Tenant.Set(nil)
}

// UnsetTenant ensures that no value is present for Tenant, not even an explicit nil
func (o *AggregateRequest) UnsetTenant() {
	o.Tenant.Unset()
}

// GetDateAdded returns the DateAdded field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *AggregateRequest) GetDateAdded() string {
	if o == nil || IsNil(o.DateAdded.Get()) {
		var ret string
		return ret
	}
	return *o.DateAdded.Get()
}

// GetDateAddedOk returns a tuple with the DateAdded field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AggregateRequest) GetDateAddedOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DateAdded.Get(), o.DateAdded.IsSet()
}

// HasDateAdded returns a boolean if a field has been set.
func (o *AggregateRequest) HasDateAdded() bool {
	if o != nil && o.DateAdded.IsSet() {
		return true
	}

	return false
}

// SetDateAdded gets a reference to the given NullableString and assigns it to the DateAdded field.
func (o *AggregateRequest) SetDateAdded(v string) {
	o.DateAdded.Set(&v)
}

// SetDateAddedNil sets the value for DateAdded to be an explicit nil
func (o *AggregateRequest) SetDateAddedNil() {
	o.DateAdded.Set(nil)
}

// UnsetDateAdded ensures that no value is present for DateAdded, not even an explicit nil
func (o *AggregateRequest) UnsetDateAdded() {
	o.DateAdded.Unset()
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *AggregateRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AggregateRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *AggregateRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *AggregateRequest) SetDescription(v string) {
	o.Description = &v
}

// GetComments returns the Comments field value if set, zero value otherwise.
func (o *AggregateRequest) GetComments() string {
	if o == nil || IsNil(o.Comments) {
		var ret string
		return ret
	}
	return *o.Comments
}

// GetCommentsOk returns a tuple with the Comments field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AggregateRequest) GetCommentsOk() (*string, bool) {
	if o == nil || IsNil(o.Comments) {
		return nil, false
	}
	return o.Comments, true
}

// HasComments returns a boolean if a field has been set.
func (o *AggregateRequest) HasComments() bool {
	if o != nil && !IsNil(o.Comments) {
		return true
	}

	return false
}

// SetComments gets a reference to the given string and assigns it to the Comments field.
func (o *AggregateRequest) SetComments(v string) {
	o.Comments = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *AggregateRequest) GetTags() []NestedTagRequest {
	if o == nil || IsNil(o.Tags) {
		var ret []NestedTagRequest
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AggregateRequest) GetTagsOk() ([]NestedTagRequest, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *AggregateRequest) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []NestedTagRequest and assigns it to the Tags field.
func (o *AggregateRequest) SetTags(v []NestedTagRequest) {
	o.Tags = v
}

// GetCustomFields returns the CustomFields field value if set, zero value otherwise.
func (o *AggregateRequest) GetCustomFields() map[string]interface{} {
	if o == nil || IsNil(o.CustomFields) {
		var ret map[string]interface{}
		return ret
	}
	return o.CustomFields
}

// GetCustomFieldsOk returns a tuple with the CustomFields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AggregateRequest) GetCustomFieldsOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.CustomFields) {
		return map[string]interface{}{}, false
	}
	return o.CustomFields, true
}

// HasCustomFields returns a boolean if a field has been set.
func (o *AggregateRequest) HasCustomFields() bool {
	if o != nil && !IsNil(o.CustomFields) {
		return true
	}

	return false
}

// SetCustomFields gets a reference to the given map[string]interface{} and assigns it to the CustomFields field.
func (o *AggregateRequest) SetCustomFields(v map[string]interface{}) {
	o.CustomFields = v
}

func (o AggregateRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AggregateRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["prefix"] = o.Prefix
	toSerialize["rir"] = o.Rir
	if o.Tenant.IsSet() {
		toSerialize["tenant"] = o.Tenant.Get()
	}
	if o.DateAdded.IsSet() {
		toSerialize["date_added"] = o.DateAdded.Get()
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Comments) {
		toSerialize["comments"] = o.Comments
	}
	if !IsNil(o.Tags) {
		toSerialize["tags"] = o.Tags
	}
	if !IsNil(o.CustomFields) {
		toSerialize["custom_fields"] = o.CustomFields
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *AggregateRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"prefix",
		"rir",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varAggregateRequest := _AggregateRequest{}

	err = json.Unmarshal(data, &varAggregateRequest)

	if err != nil {
		return err
	}

	*o = AggregateRequest(varAggregateRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "prefix")
		delete(additionalProperties, "rir")
		delete(additionalProperties, "tenant")
		delete(additionalProperties, "date_added")
		delete(additionalProperties, "description")
		delete(additionalProperties, "comments")
		delete(additionalProperties, "tags")
		delete(additionalProperties, "custom_fields")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableAggregateRequest struct {
	value *AggregateRequest
	isSet bool
}

func (v NullableAggregateRequest) Get() *AggregateRequest {
	return v.value
}

func (v *NullableAggregateRequest) Set(val *AggregateRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableAggregateRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableAggregateRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAggregateRequest(val *AggregateRequest) *NullableAggregateRequest {
	return &NullableAggregateRequest{value: val, isSet: true}
}

func (v NullableAggregateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAggregateRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
