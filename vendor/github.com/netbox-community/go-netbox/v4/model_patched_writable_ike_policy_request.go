/*
NetBox REST API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 4.0.3 (4.0)
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package netbox

import (
	"encoding/json"
)

// checks if the PatchedWritableIKEPolicyRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedWritableIKEPolicyRequest{}

// PatchedWritableIKEPolicyRequest Adds support for custom fields and tags.
type PatchedWritableIKEPolicyRequest struct {
	Name                 *string                                 `json:"name,omitempty"`
	Description          *string                                 `json:"description,omitempty"`
	Version              *PatchedWritableIKEPolicyRequestVersion `json:"version,omitempty"`
	Mode                 *PatchedWritableIKEPolicyRequestMode    `json:"mode,omitempty"`
	Proposals            []int32                                 `json:"proposals,omitempty"`
	PresharedKey         *string                                 `json:"preshared_key,omitempty"`
	Comments             *string                                 `json:"comments,omitempty"`
	Tags                 []NestedTagRequest                      `json:"tags,omitempty"`
	CustomFields         map[string]interface{}                  `json:"custom_fields,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PatchedWritableIKEPolicyRequest PatchedWritableIKEPolicyRequest

// NewPatchedWritableIKEPolicyRequest instantiates a new PatchedWritableIKEPolicyRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedWritableIKEPolicyRequest() *PatchedWritableIKEPolicyRequest {
	this := PatchedWritableIKEPolicyRequest{}
	return &this
}

// NewPatchedWritableIKEPolicyRequestWithDefaults instantiates a new PatchedWritableIKEPolicyRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedWritableIKEPolicyRequestWithDefaults() *PatchedWritableIKEPolicyRequest {
	this := PatchedWritableIKEPolicyRequest{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PatchedWritableIKEPolicyRequest) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *PatchedWritableIKEPolicyRequest) SetDescription(v string) {
	o.Description = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetVersion() PatchedWritableIKEPolicyRequestVersion {
	if o == nil || IsNil(o.Version) {
		var ret PatchedWritableIKEPolicyRequestVersion
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetVersionOk() (*PatchedWritableIKEPolicyRequestVersion, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given PatchedWritableIKEPolicyRequestVersion and assigns it to the Version field.
func (o *PatchedWritableIKEPolicyRequest) SetVersion(v PatchedWritableIKEPolicyRequestVersion) {
	o.Version = &v
}

// GetMode returns the Mode field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetMode() PatchedWritableIKEPolicyRequestMode {
	if o == nil || IsNil(o.Mode) {
		var ret PatchedWritableIKEPolicyRequestMode
		return ret
	}
	return *o.Mode
}

// GetModeOk returns a tuple with the Mode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetModeOk() (*PatchedWritableIKEPolicyRequestMode, bool) {
	if o == nil || IsNil(o.Mode) {
		return nil, false
	}
	return o.Mode, true
}

// HasMode returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasMode() bool {
	if o != nil && !IsNil(o.Mode) {
		return true
	}

	return false
}

// SetMode gets a reference to the given PatchedWritableIKEPolicyRequestMode and assigns it to the Mode field.
func (o *PatchedWritableIKEPolicyRequest) SetMode(v PatchedWritableIKEPolicyRequestMode) {
	o.Mode = &v
}

// GetProposals returns the Proposals field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetProposals() []int32 {
	if o == nil || IsNil(o.Proposals) {
		var ret []int32
		return ret
	}
	return o.Proposals
}

// GetProposalsOk returns a tuple with the Proposals field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetProposalsOk() ([]int32, bool) {
	if o == nil || IsNil(o.Proposals) {
		return nil, false
	}
	return o.Proposals, true
}

// HasProposals returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasProposals() bool {
	if o != nil && !IsNil(o.Proposals) {
		return true
	}

	return false
}

// SetProposals gets a reference to the given []int32 and assigns it to the Proposals field.
func (o *PatchedWritableIKEPolicyRequest) SetProposals(v []int32) {
	o.Proposals = v
}

// GetPresharedKey returns the PresharedKey field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetPresharedKey() string {
	if o == nil || IsNil(o.PresharedKey) {
		var ret string
		return ret
	}
	return *o.PresharedKey
}

// GetPresharedKeyOk returns a tuple with the PresharedKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetPresharedKeyOk() (*string, bool) {
	if o == nil || IsNil(o.PresharedKey) {
		return nil, false
	}
	return o.PresharedKey, true
}

// HasPresharedKey returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasPresharedKey() bool {
	if o != nil && !IsNil(o.PresharedKey) {
		return true
	}

	return false
}

// SetPresharedKey gets a reference to the given string and assigns it to the PresharedKey field.
func (o *PatchedWritableIKEPolicyRequest) SetPresharedKey(v string) {
	o.PresharedKey = &v
}

// GetComments returns the Comments field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetComments() string {
	if o == nil || IsNil(o.Comments) {
		var ret string
		return ret
	}
	return *o.Comments
}

// GetCommentsOk returns a tuple with the Comments field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetCommentsOk() (*string, bool) {
	if o == nil || IsNil(o.Comments) {
		return nil, false
	}
	return o.Comments, true
}

// HasComments returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasComments() bool {
	if o != nil && !IsNil(o.Comments) {
		return true
	}

	return false
}

// SetComments gets a reference to the given string and assigns it to the Comments field.
func (o *PatchedWritableIKEPolicyRequest) SetComments(v string) {
	o.Comments = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetTags() []NestedTagRequest {
	if o == nil || IsNil(o.Tags) {
		var ret []NestedTagRequest
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetTagsOk() ([]NestedTagRequest, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []NestedTagRequest and assigns it to the Tags field.
func (o *PatchedWritableIKEPolicyRequest) SetTags(v []NestedTagRequest) {
	o.Tags = v
}

// GetCustomFields returns the CustomFields field value if set, zero value otherwise.
func (o *PatchedWritableIKEPolicyRequest) GetCustomFields() map[string]interface{} {
	if o == nil || IsNil(o.CustomFields) {
		var ret map[string]interface{}
		return ret
	}
	return o.CustomFields
}

// GetCustomFieldsOk returns a tuple with the CustomFields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableIKEPolicyRequest) GetCustomFieldsOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.CustomFields) {
		return map[string]interface{}{}, false
	}
	return o.CustomFields, true
}

// HasCustomFields returns a boolean if a field has been set.
func (o *PatchedWritableIKEPolicyRequest) HasCustomFields() bool {
	if o != nil && !IsNil(o.CustomFields) {
		return true
	}

	return false
}

// SetCustomFields gets a reference to the given map[string]interface{} and assigns it to the CustomFields field.
func (o *PatchedWritableIKEPolicyRequest) SetCustomFields(v map[string]interface{}) {
	o.CustomFields = v
}

func (o PatchedWritableIKEPolicyRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedWritableIKEPolicyRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.Mode) {
		toSerialize["mode"] = o.Mode
	}
	if !IsNil(o.Proposals) {
		toSerialize["proposals"] = o.Proposals
	}
	if !IsNil(o.PresharedKey) {
		toSerialize["preshared_key"] = o.PresharedKey
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

func (o *PatchedWritableIKEPolicyRequest) UnmarshalJSON(data []byte) (err error) {
	varPatchedWritableIKEPolicyRequest := _PatchedWritableIKEPolicyRequest{}

	err = json.Unmarshal(data, &varPatchedWritableIKEPolicyRequest)

	if err != nil {
		return err
	}

	*o = PatchedWritableIKEPolicyRequest(varPatchedWritableIKEPolicyRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "description")
		delete(additionalProperties, "version")
		delete(additionalProperties, "mode")
		delete(additionalProperties, "proposals")
		delete(additionalProperties, "preshared_key")
		delete(additionalProperties, "comments")
		delete(additionalProperties, "tags")
		delete(additionalProperties, "custom_fields")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePatchedWritableIKEPolicyRequest struct {
	value *PatchedWritableIKEPolicyRequest
	isSet bool
}

func (v NullablePatchedWritableIKEPolicyRequest) Get() *PatchedWritableIKEPolicyRequest {
	return v.value
}

func (v *NullablePatchedWritableIKEPolicyRequest) Set(val *PatchedWritableIKEPolicyRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedWritableIKEPolicyRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedWritableIKEPolicyRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedWritableIKEPolicyRequest(val *PatchedWritableIKEPolicyRequest) *NullablePatchedWritableIKEPolicyRequest {
	return &NullablePatchedWritableIKEPolicyRequest{value: val, isSet: true}
}

func (v NullablePatchedWritableIKEPolicyRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedWritableIKEPolicyRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
