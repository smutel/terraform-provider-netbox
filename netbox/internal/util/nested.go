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

func GetNestedIPAddressAddress(nested *models.NestedIPAddress) *string {
	if nested == nil {
		return nil
	}

	return nested.Address
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
