package netbox

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/virtualization"
	"github.com/smutel/go-netbox/netbox/models"
)

type InfosForPrimary struct {
	vmID           int64
	vmName         string
	vmClusterID    int64
	vmPrimaryIP4ID int64
	vmTags         []*models.NestedTag
}

const VMInterfaceType string = "virtualization.vminterface"

const CustomFieldBoolean = "boolean"

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

func purgeNestedTagsToNestedTags(vmTags []*models.NestedTag) []*models.NestedTag {
	nestedTags := []*models.NestedTag{}

	for _, tag := range vmTags {
		tagName := tag.Name
		tagSlug := tag.Slug

		nestedTag := models.NestedTag{Name: tagName, Slug: tagSlug}
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

func getInfoForPrimary(m interface{}, objectID int64) (InfosForPrimary, error) {
	client := m.(*netboxclient.NetBoxAPI)

	objectIDStr := fmt.Sprintf("%d", objectID)
	info := InfosForPrimary{}
	paramsInterface := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&objectIDStr)
	interfaces, err := client.Virtualization.VirtualizationInterfacesList(
		paramsInterface, nil)

	if err != nil {
		return info, err
	}

	for _, i := range interfaces.Payload.Results {
		if i.ID == objectID {
			if i.VirtualMachine != nil {
				vmIDStr := fmt.Sprintf("%d", i.VirtualMachine.ID)
				paramsVM := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(&vmIDStr)
				vms, err := client.Virtualization.VirtualizationVirtualMachinesList(
					paramsVM, nil)

				if err != nil {
					return info, err
				}

				if *vms.Payload.Count != 1 {
					return info, fmt.Errorf("Cannot set an interface as primary when " +
						"the interface is not attached to a VM.")
				}

				info.vmID = i.VirtualMachine.ID
				info.vmName = *vms.Payload.Results[0].Name
				info.vmClusterID = vms.Payload.Results[0].Cluster.ID
				info.vmTags = vms.Payload.Results[0].Tags

				if vms.Payload.Results[0].PrimaryIp4 != nil {
					info.vmPrimaryIP4ID = vms.Payload.Results[0].PrimaryIp4.ID
				} else {
					info.vmPrimaryIP4ID = 0
				}
			} else {
				return info, fmt.Errorf("Cannot set an interface as primary when the " +
					"interface is not attached to a VM.")
			}
		}
	}

	return info, nil
}

func updatePrimaryStatus(m interface{}, info InfosForPrimary, id int64) error {
	client := m.(*netboxclient.NetBoxAPI)

	if info.vmName != "" {
		params := &models.WritableVirtualMachineWithConfigContext{}
		params.Name = &info.vmName
		params.Cluster = &info.vmClusterID
		params.PrimaryIp4 = &id
		params.Tags = purgeNestedTagsToNestedTags(info.vmTags)
		vm := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)
		vm.SetID(info.vmID)
		_, err := client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
			vm, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertArrayInterfaceString(arrayInterface []interface{}) string {
	var arrayString []string

	for _, item := range arrayInterface {
		switch v := item.(type) {
		case string:
			arrayString = append(arrayString, v)
		case int:
			strV := strconv.FormatInt(int64(v), 10)
			arrayString = append(arrayString, strV)
		}

	}

	sort.Strings(arrayString)
	result := strings.Join(arrayString, ",")

	return result
}

// Pick the custom fields in the state file and update values with data from API
func updateCustomFieldsFromAPI(stateCustomFields, customFields interface{}) []map[string]string {
	var tfCms []map[string]string

	switch t := customFields.(type) {
	case map[string]interface{}:
		for _, stateCustomField := range stateCustomFields.([]interface{}) {
			for key, value := range t {
				if stateCustomField.(map[string]interface{})["name"].(string) == key {
					var strValue string

					cm := map[string]string{}
					cm["name"] = key
					cm["type"] = stateCustomField.(map[string]interface{})["type"].(string)

					if value != nil {
						switch v := value.(type) {
						case []interface{}:
							strValue = convertArrayInterfaceString(v)
						default:
							strValue = fmt.Sprintf("%v", v)
						}

						if strValue == "1" && cm["type"] == CustomFieldBoolean {
							strValue = "true"
						} else if strValue == "0" && cm["type"] == CustomFieldBoolean {
							strValue = "false"
						}

						cm["value"] = strValue
					} else {
						cm["value"] = ""
					}

					tfCms = append(tfCms, cm)
				}
			}
		}
	}

	return tfCms
}

// Convert custom field regarding his type
func convertCustomFieldsFromTerraformToAPI(stateCustomFields []interface{}, customFields []interface{}) map[string]interface{} {
	toReturn := make(map[string]interface{})

	for _, stateCf := range stateCustomFields {
		stateCustomField := stateCf.(map[string]interface{})
		toReturn[stateCustomField["name"].(string)] = nil
	}

	for _, cf := range customFields {
		customField := cf.(map[string]interface{})

		cfName := customField["name"].(string)
		cfType := customField["type"].(string)
		cfValue := customField["value"].(string)

		if len(cfValue) > 0 {
			if cfType == "integer" {
				cfValueInt, _ := strconv.Atoi(cfValue)
				toReturn[cfName] = cfValueInt
			} else if cfType == CustomFieldBoolean {
				if cfValue == "true" {
					toReturn[cfName] = true
				} else if cfValue == "false" {
					toReturn[cfName] = false
				}
			} else if cfType == "multiple" {
				cfValueArray := strings.Split(cfValue, ",")
				sort.Strings(cfValueArray)
				toReturn[cfName] = cfValueArray
			} else {
				toReturn[cfName] = cfValue
			}
		} else {
			toReturn[cfName] = nil
		}
	}

	return toReturn
}
