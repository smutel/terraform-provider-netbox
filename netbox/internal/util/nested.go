package util

import (
	"github.com/smutel/go-netbox/v3/netbox/models"
)

func GetNestedIPAddressAddress(nested *models.NestedIPAddress) *string {
	if nested == nil {
		return nil
	}

	return nested.Address
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
