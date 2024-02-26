/*
NetBox REST API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 3.7.1 (3.7)
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package netbox

import (
	"encoding/json"
	"fmt"
)

// checks if the PowerPortTemplateRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PowerPortTemplateRequest{}

// PowerPortTemplateRequest Extends the built-in ModelSerializer to enforce calling full_clean() on a copy of the associated instance during validation. (DRF does not do this by default; see https://github.com/encode/django-rest-framework/issues/3144)
type PowerPortTemplateRequest struct {
	DeviceType NullableNestedDeviceTypeRequest `json:"device_type,omitempty"`
	ModuleType NullableNestedModuleTypeRequest `json:"module_type,omitempty"`
	// {module} is accepted as a substitution for the module bay position when attached to a module type.
	Name string `json:"name"`
	// Physical label
	Label *string                      `json:"label,omitempty"`
	Type  NullablePowerPortRequestType `json:"type,omitempty"`
	// Maximum power draw (watts)
	MaximumDraw NullableInt32 `json:"maximum_draw,omitempty"`
	// Allocated power draw (watts)
	AllocatedDraw        NullableInt32 `json:"allocated_draw,omitempty"`
	Description          *string       `json:"description,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PowerPortTemplateRequest PowerPortTemplateRequest

// NewPowerPortTemplateRequest instantiates a new PowerPortTemplateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPowerPortTemplateRequest(name string) *PowerPortTemplateRequest {
	this := PowerPortTemplateRequest{}
	this.Name = name
	return &this
}

// NewPowerPortTemplateRequestWithDefaults instantiates a new PowerPortTemplateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPowerPortTemplateRequestWithDefaults() *PowerPortTemplateRequest {
	this := PowerPortTemplateRequest{}
	return &this
}

// GetDeviceType returns the DeviceType field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PowerPortTemplateRequest) GetDeviceType() NestedDeviceTypeRequest {
	if o == nil || IsNil(o.DeviceType.Get()) {
		var ret NestedDeviceTypeRequest
		return ret
	}
	return *o.DeviceType.Get()
}

// GetDeviceTypeOk returns a tuple with the DeviceType field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PowerPortTemplateRequest) GetDeviceTypeOk() (*NestedDeviceTypeRequest, bool) {
	if o == nil {
		return nil, false
	}
	return o.DeviceType.Get(), o.DeviceType.IsSet()
}

// HasDeviceType returns a boolean if a field has been set.
func (o *PowerPortTemplateRequest) HasDeviceType() bool {
	if o != nil && o.DeviceType.IsSet() {
		return true
	}

	return false
}

// SetDeviceType gets a reference to the given NullableNestedDeviceTypeRequest and assigns it to the DeviceType field.
func (o *PowerPortTemplateRequest) SetDeviceType(v NestedDeviceTypeRequest) {
	o.DeviceType.Set(&v)
}

// SetDeviceTypeNil sets the value for DeviceType to be an explicit nil
func (o *PowerPortTemplateRequest) SetDeviceTypeNil() {
	o.DeviceType.Set(nil)
}

// UnsetDeviceType ensures that no value is present for DeviceType, not even an explicit nil
func (o *PowerPortTemplateRequest) UnsetDeviceType() {
	o.DeviceType.Unset()
}

// GetModuleType returns the ModuleType field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PowerPortTemplateRequest) GetModuleType() NestedModuleTypeRequest {
	if o == nil || IsNil(o.ModuleType.Get()) {
		var ret NestedModuleTypeRequest
		return ret
	}
	return *o.ModuleType.Get()
}

// GetModuleTypeOk returns a tuple with the ModuleType field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PowerPortTemplateRequest) GetModuleTypeOk() (*NestedModuleTypeRequest, bool) {
	if o == nil {
		return nil, false
	}
	return o.ModuleType.Get(), o.ModuleType.IsSet()
}

// HasModuleType returns a boolean if a field has been set.
func (o *PowerPortTemplateRequest) HasModuleType() bool {
	if o != nil && o.ModuleType.IsSet() {
		return true
	}

	return false
}

// SetModuleType gets a reference to the given NullableNestedModuleTypeRequest and assigns it to the ModuleType field.
func (o *PowerPortTemplateRequest) SetModuleType(v NestedModuleTypeRequest) {
	o.ModuleType.Set(&v)
}

// SetModuleTypeNil sets the value for ModuleType to be an explicit nil
func (o *PowerPortTemplateRequest) SetModuleTypeNil() {
	o.ModuleType.Set(nil)
}

// UnsetModuleType ensures that no value is present for ModuleType, not even an explicit nil
func (o *PowerPortTemplateRequest) UnsetModuleType() {
	o.ModuleType.Unset()
}

// GetName returns the Name field value
func (o *PowerPortTemplateRequest) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *PowerPortTemplateRequest) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *PowerPortTemplateRequest) SetName(v string) {
	o.Name = v
}

// GetLabel returns the Label field value if set, zero value otherwise.
func (o *PowerPortTemplateRequest) GetLabel() string {
	if o == nil || IsNil(o.Label) {
		var ret string
		return ret
	}
	return *o.Label
}

// GetLabelOk returns a tuple with the Label field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PowerPortTemplateRequest) GetLabelOk() (*string, bool) {
	if o == nil || IsNil(o.Label) {
		return nil, false
	}
	return o.Label, true
}

// HasLabel returns a boolean if a field has been set.
func (o *PowerPortTemplateRequest) HasLabel() bool {
	if o != nil && !IsNil(o.Label) {
		return true
	}

	return false
}

// SetLabel gets a reference to the given string and assigns it to the Label field.
func (o *PowerPortTemplateRequest) SetLabel(v string) {
	o.Label = &v
}

// GetType returns the Type field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PowerPortTemplateRequest) GetType() PowerPortRequestType {
	if o == nil || IsNil(o.Type.Get()) {
		var ret PowerPortRequestType
		return ret
	}
	return *o.Type.Get()
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PowerPortTemplateRequest) GetTypeOk() (*PowerPortRequestType, bool) {
	if o == nil {
		return nil, false
	}
	return o.Type.Get(), o.Type.IsSet()
}

// HasType returns a boolean if a field has been set.
func (o *PowerPortTemplateRequest) HasType() bool {
	if o != nil && o.Type.IsSet() {
		return true
	}

	return false
}

// SetType gets a reference to the given NullablePowerPortRequestType and assigns it to the Type field.
func (o *PowerPortTemplateRequest) SetType(v PowerPortRequestType) {
	o.Type.Set(&v)
}

// SetTypeNil sets the value for Type to be an explicit nil
func (o *PowerPortTemplateRequest) SetTypeNil() {
	o.Type.Set(nil)
}

// UnsetType ensures that no value is present for Type, not even an explicit nil
func (o *PowerPortTemplateRequest) UnsetType() {
	o.Type.Unset()
}

// GetMaximumDraw returns the MaximumDraw field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PowerPortTemplateRequest) GetMaximumDraw() int32 {
	if o == nil || IsNil(o.MaximumDraw.Get()) {
		var ret int32
		return ret
	}
	return *o.MaximumDraw.Get()
}

// GetMaximumDrawOk returns a tuple with the MaximumDraw field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PowerPortTemplateRequest) GetMaximumDrawOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.MaximumDraw.Get(), o.MaximumDraw.IsSet()
}

// HasMaximumDraw returns a boolean if a field has been set.
func (o *PowerPortTemplateRequest) HasMaximumDraw() bool {
	if o != nil && o.MaximumDraw.IsSet() {
		return true
	}

	return false
}

// SetMaximumDraw gets a reference to the given NullableInt32 and assigns it to the MaximumDraw field.
func (o *PowerPortTemplateRequest) SetMaximumDraw(v int32) {
	o.MaximumDraw.Set(&v)
}

// SetMaximumDrawNil sets the value for MaximumDraw to be an explicit nil
func (o *PowerPortTemplateRequest) SetMaximumDrawNil() {
	o.MaximumDraw.Set(nil)
}

// UnsetMaximumDraw ensures that no value is present for MaximumDraw, not even an explicit nil
func (o *PowerPortTemplateRequest) UnsetMaximumDraw() {
	o.MaximumDraw.Unset()
}

// GetAllocatedDraw returns the AllocatedDraw field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PowerPortTemplateRequest) GetAllocatedDraw() int32 {
	if o == nil || IsNil(o.AllocatedDraw.Get()) {
		var ret int32
		return ret
	}
	return *o.AllocatedDraw.Get()
}

// GetAllocatedDrawOk returns a tuple with the AllocatedDraw field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PowerPortTemplateRequest) GetAllocatedDrawOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.AllocatedDraw.Get(), o.AllocatedDraw.IsSet()
}

// HasAllocatedDraw returns a boolean if a field has been set.
func (o *PowerPortTemplateRequest) HasAllocatedDraw() bool {
	if o != nil && o.AllocatedDraw.IsSet() {
		return true
	}

	return false
}

// SetAllocatedDraw gets a reference to the given NullableInt32 and assigns it to the AllocatedDraw field.
func (o *PowerPortTemplateRequest) SetAllocatedDraw(v int32) {
	o.AllocatedDraw.Set(&v)
}

// SetAllocatedDrawNil sets the value for AllocatedDraw to be an explicit nil
func (o *PowerPortTemplateRequest) SetAllocatedDrawNil() {
	o.AllocatedDraw.Set(nil)
}

// UnsetAllocatedDraw ensures that no value is present for AllocatedDraw, not even an explicit nil
func (o *PowerPortTemplateRequest) UnsetAllocatedDraw() {
	o.AllocatedDraw.Unset()
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *PowerPortTemplateRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PowerPortTemplateRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *PowerPortTemplateRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *PowerPortTemplateRequest) SetDescription(v string) {
	o.Description = &v
}

func (o PowerPortTemplateRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PowerPortTemplateRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.DeviceType.IsSet() {
		toSerialize["device_type"] = o.DeviceType.Get()
	}
	if o.ModuleType.IsSet() {
		toSerialize["module_type"] = o.ModuleType.Get()
	}
	toSerialize["name"] = o.Name
	if !IsNil(o.Label) {
		toSerialize["label"] = o.Label
	}
	if o.Type.IsSet() {
		toSerialize["type"] = o.Type.Get()
	}
	if o.MaximumDraw.IsSet() {
		toSerialize["maximum_draw"] = o.MaximumDraw.Get()
	}
	if o.AllocatedDraw.IsSet() {
		toSerialize["allocated_draw"] = o.AllocatedDraw.Get()
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *PowerPortTemplateRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
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

	varPowerPortTemplateRequest := _PowerPortTemplateRequest{}

	err = json.Unmarshal(data, &varPowerPortTemplateRequest)

	if err != nil {
		return err
	}

	*o = PowerPortTemplateRequest(varPowerPortTemplateRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "device_type")
		delete(additionalProperties, "module_type")
		delete(additionalProperties, "name")
		delete(additionalProperties, "label")
		delete(additionalProperties, "type")
		delete(additionalProperties, "maximum_draw")
		delete(additionalProperties, "allocated_draw")
		delete(additionalProperties, "description")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePowerPortTemplateRequest struct {
	value *PowerPortTemplateRequest
	isSet bool
}

func (v NullablePowerPortTemplateRequest) Get() *PowerPortTemplateRequest {
	return v.value
}

func (v *NullablePowerPortTemplateRequest) Set(val *PowerPortTemplateRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePowerPortTemplateRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePowerPortTemplateRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePowerPortTemplateRequest(val *PowerPortTemplateRequest) *NullablePowerPortTemplateRequest {
	return &NullablePowerPortTemplateRequest{value: val, isSet: true}
}

func (v NullablePowerPortTemplateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePowerPortTemplateRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
