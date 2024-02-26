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

// DataSourceStatusLabel the model 'DataSourceStatusLabel'
type DataSourceStatusLabel string

// List of DataSource_status_label
const (
	DATASOURCESTATUSLABEL_NEW       DataSourceStatusLabel = "New"
	DATASOURCESTATUSLABEL_QUEUED    DataSourceStatusLabel = "Queued"
	DATASOURCESTATUSLABEL_SYNCING   DataSourceStatusLabel = "Syncing"
	DATASOURCESTATUSLABEL_COMPLETED DataSourceStatusLabel = "Completed"
	DATASOURCESTATUSLABEL_FAILED    DataSourceStatusLabel = "Failed"
)

// All allowed values of DataSourceStatusLabel enum
var AllowedDataSourceStatusLabelEnumValues = []DataSourceStatusLabel{
	"New",
	"Queued",
	"Syncing",
	"Completed",
	"Failed",
}

func (v *DataSourceStatusLabel) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DataSourceStatusLabel(value)
	for _, existing := range AllowedDataSourceStatusLabelEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DataSourceStatusLabel", value)
}

// NewDataSourceStatusLabelFromValue returns a pointer to a valid DataSourceStatusLabel
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewDataSourceStatusLabelFromValue(v string) (*DataSourceStatusLabel, error) {
	ev := DataSourceStatusLabel(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for DataSourceStatusLabel: valid values are %v", v, AllowedDataSourceStatusLabelEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v DataSourceStatusLabel) IsValid() bool {
	for _, existing := range AllowedDataSourceStatusLabelEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to DataSource_status_label value
func (v DataSourceStatusLabel) Ptr() *DataSourceStatusLabel {
	return &v
}

type NullableDataSourceStatusLabel struct {
	value *DataSourceStatusLabel
	isSet bool
}

func (v NullableDataSourceStatusLabel) Get() *DataSourceStatusLabel {
	return v.value
}

func (v *NullableDataSourceStatusLabel) Set(val *DataSourceStatusLabel) {
	v.value = val
	v.isSet = true
}

func (v NullableDataSourceStatusLabel) IsSet() bool {
	return v.isSet
}

func (v *NullableDataSourceStatusLabel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDataSourceStatusLabel(val *DataSourceStatusLabel) *NullableDataSourceStatusLabel {
	return &NullableDataSourceStatusLabel{value: val, isSet: true}
}

func (v NullableDataSourceStatusLabel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDataSourceStatusLabel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
