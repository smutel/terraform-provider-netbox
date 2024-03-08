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

// checks if the WritableRearPortTemplateRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WritableRearPortTemplateRequest{}

// WritableRearPortTemplateRequest Extends the built-in ModelSerializer to enforce calling full_clean() on a copy of the associated instance during validation. (DRF does not do this by default; see https://github.com/encode/django-rest-framework/issues/3144)
type WritableRearPortTemplateRequest struct {
	DeviceType NullableDeviceTypeRequest `json:"device_type,omitempty"`
	ModuleType NullableModuleTypeRequest `json:"module_type,omitempty"`
	// {module} is accepted as a substitution for the module bay position when attached to a module type.
	Name string `json:"name"`
	// Physical label
	Label                *string            `json:"label,omitempty"`
	Type                 FrontPortTypeValue `json:"type"`
	Color                *string            `json:"color,omitempty"`
	Positions            *int32             `json:"positions,omitempty"`
	Description          *string            `json:"description,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _WritableRearPortTemplateRequest WritableRearPortTemplateRequest

// NewWritableRearPortTemplateRequest instantiates a new WritableRearPortTemplateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWritableRearPortTemplateRequest(name string, type_ FrontPortTypeValue) *WritableRearPortTemplateRequest {
	this := WritableRearPortTemplateRequest{}
	this.Name = name
	this.Type = type_
	return &this
}

// NewWritableRearPortTemplateRequestWithDefaults instantiates a new WritableRearPortTemplateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWritableRearPortTemplateRequestWithDefaults() *WritableRearPortTemplateRequest {
	this := WritableRearPortTemplateRequest{}
	return &this
}

// GetDeviceType returns the DeviceType field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WritableRearPortTemplateRequest) GetDeviceType() DeviceTypeRequest {
	if o == nil || IsNil(o.DeviceType.Get()) {
		var ret DeviceTypeRequest
		return ret
	}
	return *o.DeviceType.Get()
}

// GetDeviceTypeOk returns a tuple with the DeviceType field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WritableRearPortTemplateRequest) GetDeviceTypeOk() (*DeviceTypeRequest, bool) {
	if o == nil {
		return nil, false
	}
	return o.DeviceType.Get(), o.DeviceType.IsSet()
}

// HasDeviceType returns a boolean if a field has been set.
func (o *WritableRearPortTemplateRequest) HasDeviceType() bool {
	if o != nil && o.DeviceType.IsSet() {
		return true
	}

	return false
}

// SetDeviceType gets a reference to the given NullableDeviceTypeRequest and assigns it to the DeviceType field.
func (o *WritableRearPortTemplateRequest) SetDeviceType(v DeviceTypeRequest) {
	o.DeviceType.Set(&v)
}

// SetDeviceTypeNil sets the value for DeviceType to be an explicit nil
func (o *WritableRearPortTemplateRequest) SetDeviceTypeNil() {
	o.DeviceType.Set(nil)
}

// UnsetDeviceType ensures that no value is present for DeviceType, not even an explicit nil
func (o *WritableRearPortTemplateRequest) UnsetDeviceType() {
	o.DeviceType.Unset()
}

// GetModuleType returns the ModuleType field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WritableRearPortTemplateRequest) GetModuleType() ModuleTypeRequest {
	if o == nil || IsNil(o.ModuleType.Get()) {
		var ret ModuleTypeRequest
		return ret
	}
	return *o.ModuleType.Get()
}

// GetModuleTypeOk returns a tuple with the ModuleType field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WritableRearPortTemplateRequest) GetModuleTypeOk() (*ModuleTypeRequest, bool) {
	if o == nil {
		return nil, false
	}
	return o.ModuleType.Get(), o.ModuleType.IsSet()
}

// HasModuleType returns a boolean if a field has been set.
func (o *WritableRearPortTemplateRequest) HasModuleType() bool {
	if o != nil && o.ModuleType.IsSet() {
		return true
	}

	return false
}

// SetModuleType gets a reference to the given NullableModuleTypeRequest and assigns it to the ModuleType field.
func (o *WritableRearPortTemplateRequest) SetModuleType(v ModuleTypeRequest) {
	o.ModuleType.Set(&v)
}

// SetModuleTypeNil sets the value for ModuleType to be an explicit nil
func (o *WritableRearPortTemplateRequest) SetModuleTypeNil() {
	o.ModuleType.Set(nil)
}

// UnsetModuleType ensures that no value is present for ModuleType, not even an explicit nil
func (o *WritableRearPortTemplateRequest) UnsetModuleType() {
	o.ModuleType.Unset()
}

// GetName returns the Name field value
func (o *WritableRearPortTemplateRequest) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *WritableRearPortTemplateRequest) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *WritableRearPortTemplateRequest) SetName(v string) {
	o.Name = v
}

// GetLabel returns the Label field value if set, zero value otherwise.
func (o *WritableRearPortTemplateRequest) GetLabel() string {
	if o == nil || IsNil(o.Label) {
		var ret string
		return ret
	}
	return *o.Label
}

// GetLabelOk returns a tuple with the Label field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WritableRearPortTemplateRequest) GetLabelOk() (*string, bool) {
	if o == nil || IsNil(o.Label) {
		return nil, false
	}
	return o.Label, true
}

// HasLabel returns a boolean if a field has been set.
func (o *WritableRearPortTemplateRequest) HasLabel() bool {
	if o != nil && !IsNil(o.Label) {
		return true
	}

	return false
}

// SetLabel gets a reference to the given string and assigns it to the Label field.
func (o *WritableRearPortTemplateRequest) SetLabel(v string) {
	o.Label = &v
}

// GetType returns the Type field value
func (o *WritableRearPortTemplateRequest) GetType() FrontPortTypeValue {
	if o == nil {
		var ret FrontPortTypeValue
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *WritableRearPortTemplateRequest) GetTypeOk() (*FrontPortTypeValue, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *WritableRearPortTemplateRequest) SetType(v FrontPortTypeValue) {
	o.Type = v
}

// GetColor returns the Color field value if set, zero value otherwise.
func (o *WritableRearPortTemplateRequest) GetColor() string {
	if o == nil || IsNil(o.Color) {
		var ret string
		return ret
	}
	return *o.Color
}

// GetColorOk returns a tuple with the Color field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WritableRearPortTemplateRequest) GetColorOk() (*string, bool) {
	if o == nil || IsNil(o.Color) {
		return nil, false
	}
	return o.Color, true
}

// HasColor returns a boolean if a field has been set.
func (o *WritableRearPortTemplateRequest) HasColor() bool {
	if o != nil && !IsNil(o.Color) {
		return true
	}

	return false
}

// SetColor gets a reference to the given string and assigns it to the Color field.
func (o *WritableRearPortTemplateRequest) SetColor(v string) {
	o.Color = &v
}

// GetPositions returns the Positions field value if set, zero value otherwise.
func (o *WritableRearPortTemplateRequest) GetPositions() int32 {
	if o == nil || IsNil(o.Positions) {
		var ret int32
		return ret
	}
	return *o.Positions
}

// GetPositionsOk returns a tuple with the Positions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WritableRearPortTemplateRequest) GetPositionsOk() (*int32, bool) {
	if o == nil || IsNil(o.Positions) {
		return nil, false
	}
	return o.Positions, true
}

// HasPositions returns a boolean if a field has been set.
func (o *WritableRearPortTemplateRequest) HasPositions() bool {
	if o != nil && !IsNil(o.Positions) {
		return true
	}

	return false
}

// SetPositions gets a reference to the given int32 and assigns it to the Positions field.
func (o *WritableRearPortTemplateRequest) SetPositions(v int32) {
	o.Positions = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *WritableRearPortTemplateRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WritableRearPortTemplateRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *WritableRearPortTemplateRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *WritableRearPortTemplateRequest) SetDescription(v string) {
	o.Description = &v
}

func (o WritableRearPortTemplateRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WritableRearPortTemplateRequest) ToMap() (map[string]interface{}, error) {
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
	toSerialize["type"] = o.Type
	if !IsNil(o.Color) {
		toSerialize["color"] = o.Color
	}
	if !IsNil(o.Positions) {
		toSerialize["positions"] = o.Positions
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *WritableRearPortTemplateRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"type",
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

	varWritableRearPortTemplateRequest := _WritableRearPortTemplateRequest{}

	err = json.Unmarshal(data, &varWritableRearPortTemplateRequest)

	if err != nil {
		return err
	}

	*o = WritableRearPortTemplateRequest(varWritableRearPortTemplateRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "device_type")
		delete(additionalProperties, "module_type")
		delete(additionalProperties, "name")
		delete(additionalProperties, "label")
		delete(additionalProperties, "type")
		delete(additionalProperties, "color")
		delete(additionalProperties, "positions")
		delete(additionalProperties, "description")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableWritableRearPortTemplateRequest struct {
	value *WritableRearPortTemplateRequest
	isSet bool
}

func (v NullableWritableRearPortTemplateRequest) Get() *WritableRearPortTemplateRequest {
	return v.value
}

func (v *NullableWritableRearPortTemplateRequest) Set(val *WritableRearPortTemplateRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableWritableRearPortTemplateRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableWritableRearPortTemplateRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWritableRearPortTemplateRequest(val *WritableRearPortTemplateRequest) *NullableWritableRearPortTemplateRequest {
	return &NullableWritableRearPortTemplateRequest{value: val, isSet: true}
}

func (v NullableWritableRearPortTemplateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWritableRearPortTemplateRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
