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

// ClusterStatusLabel the model 'ClusterStatusLabel'
type ClusterStatusLabel string

// List of Cluster_status_label
const (
	CLUSTERSTATUSLABEL_PLANNED         ClusterStatusLabel = "Planned"
	CLUSTERSTATUSLABEL_STAGING         ClusterStatusLabel = "Staging"
	CLUSTERSTATUSLABEL_ACTIVE          ClusterStatusLabel = "Active"
	CLUSTERSTATUSLABEL_DECOMMISSIONING ClusterStatusLabel = "Decommissioning"
	CLUSTERSTATUSLABEL_OFFLINE         ClusterStatusLabel = "Offline"
)

// All allowed values of ClusterStatusLabel enum
var AllowedClusterStatusLabelEnumValues = []ClusterStatusLabel{
	"Planned",
	"Staging",
	"Active",
	"Decommissioning",
	"Offline",
}

func (v *ClusterStatusLabel) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ClusterStatusLabel(value)
	for _, existing := range AllowedClusterStatusLabelEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ClusterStatusLabel", value)
}

// NewClusterStatusLabelFromValue returns a pointer to a valid ClusterStatusLabel
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewClusterStatusLabelFromValue(v string) (*ClusterStatusLabel, error) {
	ev := ClusterStatusLabel(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ClusterStatusLabel: valid values are %v", v, AllowedClusterStatusLabelEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ClusterStatusLabel) IsValid() bool {
	for _, existing := range AllowedClusterStatusLabelEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Cluster_status_label value
func (v ClusterStatusLabel) Ptr() *ClusterStatusLabel {
	return &v
}

type NullableClusterStatusLabel struct {
	value *ClusterStatusLabel
	isSet bool
}

func (v NullableClusterStatusLabel) Get() *ClusterStatusLabel {
	return v.value
}

func (v *NullableClusterStatusLabel) Set(val *ClusterStatusLabel) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterStatusLabel) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterStatusLabel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterStatusLabel(val *ClusterStatusLabel) *NullableClusterStatusLabel {
	return &NullableClusterStatusLabel{value: val, isSet: true}
}

func (v NullableClusterStatusLabel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterStatusLabel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
