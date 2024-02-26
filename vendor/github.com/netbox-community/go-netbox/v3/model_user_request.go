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
	"time"
)

// checks if the UserRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserRequest{}

// UserRequest Extends the built-in ModelSerializer to enforce calling full_clean() on a copy of the associated instance during validation. (DRF does not do this by default; see https://github.com/encode/django-rest-framework/issues/3144)
type UserRequest struct {
	// Required. 150 characters or fewer. Letters, digits and @/./+/-/_ only.
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	// Designates whether the user can log into this admin site.
	IsStaff *bool `json:"is_staff,omitempty"`
	// Designates whether this user should be treated as active. Unselect this instead of deleting accounts.
	IsActive             *bool      `json:"is_active,omitempty"`
	DateJoined           *time.Time `json:"date_joined,omitempty"`
	Groups               []int32    `json:"groups,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _UserRequest UserRequest

// NewUserRequest instantiates a new UserRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserRequest(username string, password string) *UserRequest {
	this := UserRequest{}
	this.Username = username
	this.Password = password
	return &this
}

// NewUserRequestWithDefaults instantiates a new UserRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserRequestWithDefaults() *UserRequest {
	this := UserRequest{}
	return &this
}

// GetUsername returns the Username field value
func (o *UserRequest) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *UserRequest) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *UserRequest) SetUsername(v string) {
	o.Username = v
}

// GetPassword returns the Password field value
func (o *UserRequest) GetPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Password
}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
func (o *UserRequest) GetPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Password, true
}

// SetPassword sets field value
func (o *UserRequest) SetPassword(v string) {
	o.Password = v
}

// GetFirstName returns the FirstName field value if set, zero value otherwise.
func (o *UserRequest) GetFirstName() string {
	if o == nil || IsNil(o.FirstName) {
		var ret string
		return ret
	}
	return *o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRequest) GetFirstNameOk() (*string, bool) {
	if o == nil || IsNil(o.FirstName) {
		return nil, false
	}
	return o.FirstName, true
}

// HasFirstName returns a boolean if a field has been set.
func (o *UserRequest) HasFirstName() bool {
	if o != nil && !IsNil(o.FirstName) {
		return true
	}

	return false
}

// SetFirstName gets a reference to the given string and assigns it to the FirstName field.
func (o *UserRequest) SetFirstName(v string) {
	o.FirstName = &v
}

// GetLastName returns the LastName field value if set, zero value otherwise.
func (o *UserRequest) GetLastName() string {
	if o == nil || IsNil(o.LastName) {
		var ret string
		return ret
	}
	return *o.LastName
}

// GetLastNameOk returns a tuple with the LastName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRequest) GetLastNameOk() (*string, bool) {
	if o == nil || IsNil(o.LastName) {
		return nil, false
	}
	return o.LastName, true
}

// HasLastName returns a boolean if a field has been set.
func (o *UserRequest) HasLastName() bool {
	if o != nil && !IsNil(o.LastName) {
		return true
	}

	return false
}

// SetLastName gets a reference to the given string and assigns it to the LastName field.
func (o *UserRequest) SetLastName(v string) {
	o.LastName = &v
}

// GetEmail returns the Email field value if set, zero value otherwise.
func (o *UserRequest) GetEmail() string {
	if o == nil || IsNil(o.Email) {
		var ret string
		return ret
	}
	return *o.Email
}

// GetEmailOk returns a tuple with the Email field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRequest) GetEmailOk() (*string, bool) {
	if o == nil || IsNil(o.Email) {
		return nil, false
	}
	return o.Email, true
}

// HasEmail returns a boolean if a field has been set.
func (o *UserRequest) HasEmail() bool {
	if o != nil && !IsNil(o.Email) {
		return true
	}

	return false
}

// SetEmail gets a reference to the given string and assigns it to the Email field.
func (o *UserRequest) SetEmail(v string) {
	o.Email = &v
}

// GetIsStaff returns the IsStaff field value if set, zero value otherwise.
func (o *UserRequest) GetIsStaff() bool {
	if o == nil || IsNil(o.IsStaff) {
		var ret bool
		return ret
	}
	return *o.IsStaff
}

// GetIsStaffOk returns a tuple with the IsStaff field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRequest) GetIsStaffOk() (*bool, bool) {
	if o == nil || IsNil(o.IsStaff) {
		return nil, false
	}
	return o.IsStaff, true
}

// HasIsStaff returns a boolean if a field has been set.
func (o *UserRequest) HasIsStaff() bool {
	if o != nil && !IsNil(o.IsStaff) {
		return true
	}

	return false
}

// SetIsStaff gets a reference to the given bool and assigns it to the IsStaff field.
func (o *UserRequest) SetIsStaff(v bool) {
	o.IsStaff = &v
}

// GetIsActive returns the IsActive field value if set, zero value otherwise.
func (o *UserRequest) GetIsActive() bool {
	if o == nil || IsNil(o.IsActive) {
		var ret bool
		return ret
	}
	return *o.IsActive
}

// GetIsActiveOk returns a tuple with the IsActive field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRequest) GetIsActiveOk() (*bool, bool) {
	if o == nil || IsNil(o.IsActive) {
		return nil, false
	}
	return o.IsActive, true
}

// HasIsActive returns a boolean if a field has been set.
func (o *UserRequest) HasIsActive() bool {
	if o != nil && !IsNil(o.IsActive) {
		return true
	}

	return false
}

// SetIsActive gets a reference to the given bool and assigns it to the IsActive field.
func (o *UserRequest) SetIsActive(v bool) {
	o.IsActive = &v
}

// GetDateJoined returns the DateJoined field value if set, zero value otherwise.
func (o *UserRequest) GetDateJoined() time.Time {
	if o == nil || IsNil(o.DateJoined) {
		var ret time.Time
		return ret
	}
	return *o.DateJoined
}

// GetDateJoinedOk returns a tuple with the DateJoined field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRequest) GetDateJoinedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.DateJoined) {
		return nil, false
	}
	return o.DateJoined, true
}

// HasDateJoined returns a boolean if a field has been set.
func (o *UserRequest) HasDateJoined() bool {
	if o != nil && !IsNil(o.DateJoined) {
		return true
	}

	return false
}

// SetDateJoined gets a reference to the given time.Time and assigns it to the DateJoined field.
func (o *UserRequest) SetDateJoined(v time.Time) {
	o.DateJoined = &v
}

// GetGroups returns the Groups field value if set, zero value otherwise.
func (o *UserRequest) GetGroups() []int32 {
	if o == nil || IsNil(o.Groups) {
		var ret []int32
		return ret
	}
	return o.Groups
}

// GetGroupsOk returns a tuple with the Groups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserRequest) GetGroupsOk() ([]int32, bool) {
	if o == nil || IsNil(o.Groups) {
		return nil, false
	}
	return o.Groups, true
}

// HasGroups returns a boolean if a field has been set.
func (o *UserRequest) HasGroups() bool {
	if o != nil && !IsNil(o.Groups) {
		return true
	}

	return false
}

// SetGroups gets a reference to the given []int32 and assigns it to the Groups field.
func (o *UserRequest) SetGroups(v []int32) {
	o.Groups = v
}

func (o UserRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["username"] = o.Username
	toSerialize["password"] = o.Password
	if !IsNil(o.FirstName) {
		toSerialize["first_name"] = o.FirstName
	}
	if !IsNil(o.LastName) {
		toSerialize["last_name"] = o.LastName
	}
	if !IsNil(o.Email) {
		toSerialize["email"] = o.Email
	}
	if !IsNil(o.IsStaff) {
		toSerialize["is_staff"] = o.IsStaff
	}
	if !IsNil(o.IsActive) {
		toSerialize["is_active"] = o.IsActive
	}
	if !IsNil(o.DateJoined) {
		toSerialize["date_joined"] = o.DateJoined
	}
	if !IsNil(o.Groups) {
		toSerialize["groups"] = o.Groups
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *UserRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"username",
		"password",
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

	varUserRequest := _UserRequest{}

	err = json.Unmarshal(data, &varUserRequest)

	if err != nil {
		return err
	}

	*o = UserRequest(varUserRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "username")
		delete(additionalProperties, "password")
		delete(additionalProperties, "first_name")
		delete(additionalProperties, "last_name")
		delete(additionalProperties, "email")
		delete(additionalProperties, "is_staff")
		delete(additionalProperties, "is_active")
		delete(additionalProperties, "date_joined")
		delete(additionalProperties, "groups")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableUserRequest struct {
	value *UserRequest
	isSet bool
}

func (v NullableUserRequest) Get() *UserRequest {
	return v.value
}

func (v *NullableUserRequest) Set(val *UserRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableUserRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableUserRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserRequest(val *UserRequest) *NullableUserRequest {
	return &NullableUserRequest{value: val, isSet: true}
}

func (v NullableUserRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
