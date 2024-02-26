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

// checks if the PatchedWritableServiceRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PatchedWritableServiceRequest{}

// PatchedWritableServiceRequest Adds support for custom fields and tags.
type PatchedWritableServiceRequest struct {
	Device         NullableInt32                          `json:"device,omitempty"`
	VirtualMachine NullableInt32                          `json:"virtual_machine,omitempty"`
	Name           *string                                `json:"name,omitempty"`
	Ports          []int32                                `json:"ports,omitempty"`
	Protocol       *PatchedWritableServiceRequestProtocol `json:"protocol,omitempty"`
	// The specific IP addresses (if any) to which this service is bound
	Ipaddresses          []int32                `json:"ipaddresses,omitempty"`
	Description          *string                `json:"description,omitempty"`
	Comments             *string                `json:"comments,omitempty"`
	Tags                 []NestedTagRequest     `json:"tags,omitempty"`
	CustomFields         map[string]interface{} `json:"custom_fields,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PatchedWritableServiceRequest PatchedWritableServiceRequest

// NewPatchedWritableServiceRequest instantiates a new PatchedWritableServiceRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchedWritableServiceRequest() *PatchedWritableServiceRequest {
	this := PatchedWritableServiceRequest{}
	return &this
}

// NewPatchedWritableServiceRequestWithDefaults instantiates a new PatchedWritableServiceRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchedWritableServiceRequestWithDefaults() *PatchedWritableServiceRequest {
	this := PatchedWritableServiceRequest{}
	return &this
}

// GetDevice returns the Device field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PatchedWritableServiceRequest) GetDevice() int32 {
	if o == nil || IsNil(o.Device.Get()) {
		var ret int32
		return ret
	}
	return *o.Device.Get()
}

// GetDeviceOk returns a tuple with the Device field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PatchedWritableServiceRequest) GetDeviceOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Device.Get(), o.Device.IsSet()
}

// HasDevice returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasDevice() bool {
	if o != nil && o.Device.IsSet() {
		return true
	}

	return false
}

// SetDevice gets a reference to the given NullableInt32 and assigns it to the Device field.
func (o *PatchedWritableServiceRequest) SetDevice(v int32) {
	o.Device.Set(&v)
}

// SetDeviceNil sets the value for Device to be an explicit nil
func (o *PatchedWritableServiceRequest) SetDeviceNil() {
	o.Device.Set(nil)
}

// UnsetDevice ensures that no value is present for Device, not even an explicit nil
func (o *PatchedWritableServiceRequest) UnsetDevice() {
	o.Device.Unset()
}

// GetVirtualMachine returns the VirtualMachine field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PatchedWritableServiceRequest) GetVirtualMachine() int32 {
	if o == nil || IsNil(o.VirtualMachine.Get()) {
		var ret int32
		return ret
	}
	return *o.VirtualMachine.Get()
}

// GetVirtualMachineOk returns a tuple with the VirtualMachine field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PatchedWritableServiceRequest) GetVirtualMachineOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.VirtualMachine.Get(), o.VirtualMachine.IsSet()
}

// HasVirtualMachine returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasVirtualMachine() bool {
	if o != nil && o.VirtualMachine.IsSet() {
		return true
	}

	return false
}

// SetVirtualMachine gets a reference to the given NullableInt32 and assigns it to the VirtualMachine field.
func (o *PatchedWritableServiceRequest) SetVirtualMachine(v int32) {
	o.VirtualMachine.Set(&v)
}

// SetVirtualMachineNil sets the value for VirtualMachine to be an explicit nil
func (o *PatchedWritableServiceRequest) SetVirtualMachineNil() {
	o.VirtualMachine.Set(nil)
}

// UnsetVirtualMachine ensures that no value is present for VirtualMachine, not even an explicit nil
func (o *PatchedWritableServiceRequest) UnsetVirtualMachine() {
	o.VirtualMachine.Unset()
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PatchedWritableServiceRequest) SetName(v string) {
	o.Name = &v
}

// GetPorts returns the Ports field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetPorts() []int32 {
	if o == nil || IsNil(o.Ports) {
		var ret []int32
		return ret
	}
	return o.Ports
}

// GetPortsOk returns a tuple with the Ports field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetPortsOk() ([]int32, bool) {
	if o == nil || IsNil(o.Ports) {
		return nil, false
	}
	return o.Ports, true
}

// HasPorts returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasPorts() bool {
	if o != nil && !IsNil(o.Ports) {
		return true
	}

	return false
}

// SetPorts gets a reference to the given []int32 and assigns it to the Ports field.
func (o *PatchedWritableServiceRequest) SetPorts(v []int32) {
	o.Ports = v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetProtocol() PatchedWritableServiceRequestProtocol {
	if o == nil || IsNil(o.Protocol) {
		var ret PatchedWritableServiceRequestProtocol
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetProtocolOk() (*PatchedWritableServiceRequestProtocol, bool) {
	if o == nil || IsNil(o.Protocol) {
		return nil, false
	}
	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasProtocol() bool {
	if o != nil && !IsNil(o.Protocol) {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given PatchedWritableServiceRequestProtocol and assigns it to the Protocol field.
func (o *PatchedWritableServiceRequest) SetProtocol(v PatchedWritableServiceRequestProtocol) {
	o.Protocol = &v
}

// GetIpaddresses returns the Ipaddresses field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetIpaddresses() []int32 {
	if o == nil || IsNil(o.Ipaddresses) {
		var ret []int32
		return ret
	}
	return o.Ipaddresses
}

// GetIpaddressesOk returns a tuple with the Ipaddresses field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetIpaddressesOk() ([]int32, bool) {
	if o == nil || IsNil(o.Ipaddresses) {
		return nil, false
	}
	return o.Ipaddresses, true
}

// HasIpaddresses returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasIpaddresses() bool {
	if o != nil && !IsNil(o.Ipaddresses) {
		return true
	}

	return false
}

// SetIpaddresses gets a reference to the given []int32 and assigns it to the Ipaddresses field.
func (o *PatchedWritableServiceRequest) SetIpaddresses(v []int32) {
	o.Ipaddresses = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *PatchedWritableServiceRequest) SetDescription(v string) {
	o.Description = &v
}

// GetComments returns the Comments field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetComments() string {
	if o == nil || IsNil(o.Comments) {
		var ret string
		return ret
	}
	return *o.Comments
}

// GetCommentsOk returns a tuple with the Comments field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetCommentsOk() (*string, bool) {
	if o == nil || IsNil(o.Comments) {
		return nil, false
	}
	return o.Comments, true
}

// HasComments returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasComments() bool {
	if o != nil && !IsNil(o.Comments) {
		return true
	}

	return false
}

// SetComments gets a reference to the given string and assigns it to the Comments field.
func (o *PatchedWritableServiceRequest) SetComments(v string) {
	o.Comments = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetTags() []NestedTagRequest {
	if o == nil || IsNil(o.Tags) {
		var ret []NestedTagRequest
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetTagsOk() ([]NestedTagRequest, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []NestedTagRequest and assigns it to the Tags field.
func (o *PatchedWritableServiceRequest) SetTags(v []NestedTagRequest) {
	o.Tags = v
}

// GetCustomFields returns the CustomFields field value if set, zero value otherwise.
func (o *PatchedWritableServiceRequest) GetCustomFields() map[string]interface{} {
	if o == nil || IsNil(o.CustomFields) {
		var ret map[string]interface{}
		return ret
	}
	return o.CustomFields
}

// GetCustomFieldsOk returns a tuple with the CustomFields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PatchedWritableServiceRequest) GetCustomFieldsOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.CustomFields) {
		return map[string]interface{}{}, false
	}
	return o.CustomFields, true
}

// HasCustomFields returns a boolean if a field has been set.
func (o *PatchedWritableServiceRequest) HasCustomFields() bool {
	if o != nil && !IsNil(o.CustomFields) {
		return true
	}

	return false
}

// SetCustomFields gets a reference to the given map[string]interface{} and assigns it to the CustomFields field.
func (o *PatchedWritableServiceRequest) SetCustomFields(v map[string]interface{}) {
	o.CustomFields = v
}

func (o PatchedWritableServiceRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PatchedWritableServiceRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Device.IsSet() {
		toSerialize["device"] = o.Device.Get()
	}
	if o.VirtualMachine.IsSet() {
		toSerialize["virtual_machine"] = o.VirtualMachine.Get()
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Ports) {
		toSerialize["ports"] = o.Ports
	}
	if !IsNil(o.Protocol) {
		toSerialize["protocol"] = o.Protocol
	}
	if !IsNil(o.Ipaddresses) {
		toSerialize["ipaddresses"] = o.Ipaddresses
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

func (o *PatchedWritableServiceRequest) UnmarshalJSON(data []byte) (err error) {
	varPatchedWritableServiceRequest := _PatchedWritableServiceRequest{}

	err = json.Unmarshal(data, &varPatchedWritableServiceRequest)

	if err != nil {
		return err
	}

	*o = PatchedWritableServiceRequest(varPatchedWritableServiceRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "device")
		delete(additionalProperties, "virtual_machine")
		delete(additionalProperties, "name")
		delete(additionalProperties, "ports")
		delete(additionalProperties, "protocol")
		delete(additionalProperties, "ipaddresses")
		delete(additionalProperties, "description")
		delete(additionalProperties, "comments")
		delete(additionalProperties, "tags")
		delete(additionalProperties, "custom_fields")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePatchedWritableServiceRequest struct {
	value *PatchedWritableServiceRequest
	isSet bool
}

func (v NullablePatchedWritableServiceRequest) Get() *PatchedWritableServiceRequest {
	return v.value
}

func (v *NullablePatchedWritableServiceRequest) Set(val *PatchedWritableServiceRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchedWritableServiceRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchedWritableServiceRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchedWritableServiceRequest(val *PatchedWritableServiceRequest) *NullablePatchedWritableServiceRequest {
	return &NullablePatchedWritableServiceRequest{value: val, isSet: true}
}

func (v NullablePatchedWritableServiceRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchedWritableServiceRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
