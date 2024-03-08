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

// checks if the Module type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Module{}

// Module Adds support for custom fields and tags.
type Module struct {
	Id                   int32           `json:"id"`
	Url                  string          `json:"url"`
	Display              string          `json:"display"`
	Device               Device          `json:"device"`
	ModuleBay            NestedModuleBay `json:"module_bay"`
	AdditionalProperties map[string]interface{}
}

type _Module Module

// NewModule instantiates a new Module object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewModule(id int32, url string, display string, device Device, moduleBay NestedModuleBay) *Module {
	this := Module{}
	this.Id = id
	this.Url = url
	this.Display = display
	this.Device = device
	this.ModuleBay = moduleBay
	return &this
}

// NewModuleWithDefaults instantiates a new Module object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewModuleWithDefaults() *Module {
	this := Module{}
	return &this
}

// GetId returns the Id field value
func (o *Module) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Module) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Module) SetId(v int32) {
	o.Id = v
}

// GetUrl returns the Url field value
func (o *Module) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *Module) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *Module) SetUrl(v string) {
	o.Url = v
}

// GetDisplay returns the Display field value
func (o *Module) GetDisplay() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Display
}

// GetDisplayOk returns a tuple with the Display field value
// and a boolean to check if the value has been set.
func (o *Module) GetDisplayOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Display, true
}

// SetDisplay sets field value
func (o *Module) SetDisplay(v string) {
	o.Display = v
}

// GetDevice returns the Device field value
func (o *Module) GetDevice() Device {
	if o == nil {
		var ret Device
		return ret
	}

	return o.Device
}

// GetDeviceOk returns a tuple with the Device field value
// and a boolean to check if the value has been set.
func (o *Module) GetDeviceOk() (*Device, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Device, true
}

// SetDevice sets field value
func (o *Module) SetDevice(v Device) {
	o.Device = v
}

// GetModuleBay returns the ModuleBay field value
func (o *Module) GetModuleBay() NestedModuleBay {
	if o == nil {
		var ret NestedModuleBay
		return ret
	}

	return o.ModuleBay
}

// GetModuleBayOk returns a tuple with the ModuleBay field value
// and a boolean to check if the value has been set.
func (o *Module) GetModuleBayOk() (*NestedModuleBay, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ModuleBay, true
}

// SetModuleBay sets field value
func (o *Module) SetModuleBay(v NestedModuleBay) {
	o.ModuleBay = v
}

func (o Module) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Module) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["url"] = o.Url
	toSerialize["display"] = o.Display
	toSerialize["device"] = o.Device
	toSerialize["module_bay"] = o.ModuleBay

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Module) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"url",
		"display",
		"device",
		"module_bay",
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

	varModule := _Module{}

	err = json.Unmarshal(data, &varModule)

	if err != nil {
		return err
	}

	*o = Module(varModule)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "url")
		delete(additionalProperties, "display")
		delete(additionalProperties, "device")
		delete(additionalProperties, "module_bay")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableModule struct {
	value *Module
	isSet bool
}

func (v NullableModule) Get() *Module {
	return v.value
}

func (v *NullableModule) Set(val *Module) {
	v.value = val
	v.isSet = true
}

func (v NullableModule) IsSet() bool {
	return v.isSet
}

func (v *NullableModule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableModule(val *Module) *NullableModule {
	return &NullableModule{value: val, isSet: true}
}

func (v NullableModule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableModule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
