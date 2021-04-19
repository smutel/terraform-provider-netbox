package netbox

import (
	"errors"
	"fmt"
	"net"

	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

const VMInterfaceType string = "virtualization.vminterface"
const CustomFieldsRegex = "custom_fields.%"

func expandToStringSlice(v []interface{}) []string {
	s := make([]string, len(v))
	for i, val := range v {
		if strVal, ok := val.(string); ok {
			s[i] = strVal
		}
	}

	return s
}

func expandToInt64Slice(v []interface{}) []int64 {
	s := make([]int64, len(v))
	for i, val := range v {
		if strVal, ok := val.(int64); ok {
			s[i] = strVal
		}
	}

	return s
}

func convertTagsToNestedTags(tags []interface{}) []*models.NestedTag {
	nestedTags := []*models.NestedTag{}

	for _, tag := range tags {
		t := tag.(map[string]interface{})

		tagName := t["name"].(string)
		tagSlug := t["slug"].(string)

		nestedTag := models.NestedTag{Name: &tagName, Slug: &tagSlug}
		nestedTags = append(nestedTags, &nestedTag)
	}

	return nestedTags
}

func convertNestedTagsToTags(tags []*models.NestedTag) []map[string]string {
	var tfTags []map[string]string

	for _, t := range tags {
		tag := map[string]string{}
		tag["name"] = *t.Name
		tag["slug"] = *t.Slug

		tfTags = append(tfTags, tag)
	}

	return tfTags
}

func getAssociatedVMForInterface(client *netboxclient.NetBoxAPI, interfaceID int64) (*models.VirtualMachineWithConfigContext, error) {
	ifResp, err := client.Virtualization.VirtualizationInterfacesRead(
		virtualization.NewVirtualizationInterfacesReadParams().WithID(interfaceID), nil,
	)
	if err != nil {
		return nil, err
	}

	intf := ifResp.Payload
	if intf.VirtualMachine == nil || intf.VirtualMachine.ID == 0 {
		return nil, errors.New("interface is not assigned to a virtual machine")
	}

	vmResp, err := client.Virtualization.VirtualizationVirtualMachinesRead(
		virtualization.NewVirtualizationVirtualMachinesReadParams().WithID(intf.VirtualMachine.ID), nil,
	)
	if err != nil {
		return nil, err
	}
	return vmResp.Payload, nil
}

func updatePrimaryStatus(client *netboxclient.NetBoxAPI, addr *models.IPAddress) error {
	if addr.AssignedObjectType != VMInterfaceType {
		return errors.New("needs to be assigned to a virtual machine interface to update the primary address")
	}

	vm, err := getAssociatedVMForInterface(client, *addr.AssignedObjectID)
	if err != nil {
		return err
	}

	params := &models.WritableVirtualMachineWithConfigContext{}
	params.Name = vm.Name
	params.Cluster = &vm.Cluster.ID

	cidr, _, err := net.ParseCIDR(*addr.Address)
	if err != nil {
		return fmt.Errorf("invalid address in model: %w", err)
	}
	if cidr.To4() != nil {
		params.PrimaryIp4 = &addr.ID
	} else {
		params.PrimaryIp6 = &addr.ID
	}

	update := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().
		WithData(params).
		WithID(vm.ID)
	_, err = client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(update, nil)
	return err
}

// custom_fields have multiple data type returns based on field type
// but terraform only supports map[string]string, so we convert all to strings
func convertCustomFieldsFromAPIToTerraform(customFields interface{}) map[string]string {
	toReturn := make(map[string]string)
	switch t := customFields.(type) {

	case map[string]interface{}:
		for key, value := range t {
			var strValue string
			if value != nil {
				switch v := value.(type) {
				default:
					strValue = fmt.Sprintf("%v", v)
				case map[string]interface{}:
					strValue = fmt.Sprintf("%v", v["value"])
				}
			}

			toReturn[key] = strValue
		}
	}

	return toReturn
}

func convertCustomFieldsFromTerraformToAPICreate(customFields map[string]interface{}) map[string]interface{} {
	toReturn := make(map[string]interface{})
	for key, value := range customFields {
		toReturn[key] = value

		// special handling for booleans, as they are the only parameter not supplied as string to netbox
		if value == "true" {
			toReturn[key] = 1
		} else if value == "false" {
			toReturn[key] = 0
		}
	}

	return toReturn
}

// all custom fields need to be submitted to the netbox API for updates
func convertCustomFieldsFromTerraformToAPIUpdate(stateCustomFields, resourceCustomFields interface{}) map[string]interface{} {
	toReturn := make(map[string]interface{})

	// netbox needs explicit empty string to remove old values
	// first we fill all existing fields from the state with an empty string
	for key := range stateCustomFields.(map[string]interface{}) {
		toReturn[key] = nil
	}

	// then we override the values that still exist in the terraform code with their respective value
	for key, value := range resourceCustomFields.(map[string]interface{}) {
		toReturn[key] = value
		// https://github.com/smutel/terraform-provider-netbox/issues/32
		// under certain circumstances terraform puts empty string into `resourceCustomFields` when parent block is removed
		// we simply replace these occurences with `nil` to not negatively affect the boolean field where empty string is equal to false
		if value == "" {
			toReturn[key] = nil
		}

		// special handling for booleans, as they are the only parameter not supplied as string to netbox
		if value == "true" {
			toReturn[key] = 1
		} else if value == "false" {
			toReturn[key] = 0
		}
	}

	return toReturn
}

// replaces the input with a single whitespace if it is empty
func hackEmptyString(in string) string {
	if in == "" {
		return " "
	}
	return in
}
