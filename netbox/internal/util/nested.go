package util

import (
	"github.com/smutel/go-netbox/v3/netbox/models"
)

func GetClusterStatusValue(nested *models.ClusterStatus) *string {
	if nested == nil {
		return nil
	}

	return nested.Value
}

func GetCustomFieldUIVisibilityValue(nested *models.CustomFieldUIVisibility) *string {
	if nested == nil {
		return nil
	}

	return nested.Value
}

func GetIPAddressAssignedObject(nested *models.IPAddress) (*string, *int64) {
	if nested == nil {
		return nil, nil
	}

	return nested.AssignedObjectType, nested.AssignedObjectID
}

func GetIPAddressFamilyLabel(nested *models.IPAddressFamily) *string {
	if nested == nil {
		return nil
	}

	return nested.Label
}

func GetIPAddressRoleValue(nested *models.IPAddressRole) *string {
	if nested == nil {
		return nil
	}

	return nested.Value
}

func GetIPAddressStatusValue(nested *models.IPAddressStatus) *string {
	if nested == nil {
		return nil
	}

	return nested.Value
}

func GetNestedIPAddressAddress(nested *models.NestedIPAddress) *string {
	if nested == nil {
		return nil
	}

	return nested.Address
}

func GetNestedIPAddressID(nested *models.NestedIPAddress) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedClusterGroupID(nested *models.NestedClusterGroup) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedDeviceID(nested *models.NestedDevice) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedManufacturerID(nested *models.NestedManufacturer) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedPlatformID(nested *models.NestedPlatform) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedRegionID(nested *models.NestedRegion) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedSiteID(nested *models.NestedSite) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedSiteGroupID(nested *models.NestedSiteGroup) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedTenantID(nested *models.NestedTenant) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedRoleID(nested *models.NestedDeviceRole) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedVlanID(nested *models.NestedVLAN) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedVrfID(nested *models.NestedVRF) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedLocationID(nested *models.NestedLocation) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetNestedRackRoleID(nested *models.NestedRackRole) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}

func GetRackStatusValue(nested *models.RackStatus) *string {
	if nested == nil {
		return nil
	}

	return nested.Value
}

func GetRackTypeValue(nested *models.RackType) *string {
	if nested == nil {
		return nil
	}

	return nested.Value
}

func GetRackOuterUnit(nested *models.RackOuterUnit) *string {
	if nested == nil {
		return nil
	}

	return nested.Value
}

func GetNestedRegionParentID(nested *models.NestedRegion) *int64 {
	if nested == nil {
		return nil
	}

	return &nested.ID
}
