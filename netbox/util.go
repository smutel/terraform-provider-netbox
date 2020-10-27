package netbox

import (
	"fmt"

	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

type InfosForPrimary struct {
	vm_id             int64
	vm_name           string
	vm_cluster_id     int64
	vm_primary_ip4_id int64
}

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
				vm_id_str := fmt.Sprintf("%d", i.VirtualMachine.ID)
				paramsVm := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(&vm_id_str)
				vms, err := client.Virtualization.VirtualizationVirtualMachinesList(paramsVm, nil)

				if err != nil {
					return info, err
				}

				if *vms.Payload.Count != 1 {
					return info, fmt.Errorf("Cannot set an interface as primary when the " +
						"interface is not attached to a VM.")
				}

				info.vm_id = i.VirtualMachine.ID
				info.vm_name = *vms.Payload.Results[0].Name
				info.vm_cluster_id = vms.Payload.Results[0].Cluster.ID

				if vms.Payload.Results[0].PrimaryIp4 != nil {
					info.vm_primary_ip4_id = vms.Payload.Results[0].PrimaryIp4.ID
				} else {
					info.vm_primary_ip4_id = 0
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

	if info.vm_name != "" {
		params := &models.WritableVirtualMachineWithConfigContext{}
		params.Name = &info.vm_name
		params.Cluster = &info.vm_cluster_id
		params.PrimaryIp4 = &id
		vm := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)
		vm.SetID(info.vm_id)
		_, err := client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
			vm, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
 * func diffSlices(oldSlice []string, newSlice []string) []string {
 *   var diff []string
 *
 *   for _, x := range oldSlice {
 *     found := false
 *     for _, y := range newSlice {
 *       if x == y {
 *         found = true
 *       }
 *     }
 *
 *     if found == false {
 *       diff = append(diff, x)
 *     }
 *   }
 *
 *   return diff
 * }
 */

/*
 * func convertCFToAPI(customFields []interface{}) (cf map[string]interface{},
 *   e error) {
 *
 *   customFieldsAPI := make(map[string]interface{})
 *   for _, customFieldRaw := range customFields {
 *     customField := customFieldRaw.(map[string]interface{})
 *     customFieldName := customField["name"].(string)
 *     customFieldType := customField["kind"].(string)
 *     customFieldValue := customField["value"].(string)
 *
 *     if customFieldType == "string" {
 *       customFieldsAPI[customFieldName] = customFieldValue
 *     } else if customFieldType == "int" {
 *       cfIntValue, err := strconv.ParseInt(customFieldValue, 10, 64)
 *       if err != nil {
 *         return nil, err
 *       }
 *       customFieldsAPI[customFieldName] = cfIntValue
 *     } else if customFieldType == "bool" {
 *       cfBoolValue, err := strconv.ParseBool(customFieldValue)
 *       if err != nil {
 *         return nil, err
 *       }
 *       customFieldsAPI[customFieldName] = cfBoolValue
 *     }
 *   }
 *
 *   return customFieldsAPI, nil
 * }
 */

/*
 * func convertAPIToCF(customFields interface{}) (cf []map[string]string) {
 *   var cfAPI map[string]string
 *   cfAPISize := 0
 *
 *   for _, v := range customFields.(map[string]interface{}) {
 *     if v != nil {
 *       cfAPISize++
 *     }
 *   }
 *
 *   customFieldsAPI := make([]map[string]string, cfAPISize)
 *
 *   i := 0
 *   for k, v := range customFields.(map[string]interface{}) {
 *     cfAPI = make(map[string]string)
 *     cfAPI["name"] = k
 *
 *     if v != nil {
 *       switch v.(type) {
 *       case json.Number:
 *         cfAPI["value"] = v.(json.Number).String()
 *         cfAPI["kind"] = "int"
 *       case bool:
 *         cfAPI["value"] = strconv.FormatBool(v.(bool))
 *         cfAPI["kind"] = "bool"
 *       default:
 *         cfAPI["value"] = v.(string)
 *         cfAPI["kind"] = "string"
 *       }
 *       customFieldsAPI[i] = cfAPI
 *       i++
 *     }
 *   }
 *
 *   return customFieldsAPI
 * }
 */
