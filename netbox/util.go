package netbox

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/client/virtualization"
	"github.com/smutel/go-netbox/v3/netbox/models"
)

// Type of vm interface in Netbox
const vMInterfaceType string = "virtualization.vminterface"

// Boolean string for custom field
const customFieldBoolean = "boolean"
const customFieldJSON = "json"
const customFieldMultiSelectLegacy = "multiple"
const customFieldMultiSelect = "multiselect"
const customFieldObject = "object"
const customFieldMultiObject = "multiobject"

var customFieldSchema = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the existing custom field.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"text", "longtext", "integer", customFieldBoolean,
					"date", "url", customFieldJSON, "select", customFieldMultiSelect, customFieldObject, customFieldMultiObject, customFieldMultiSelectLegacy, "selection"}, false),
				Description: "Type of the existing custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject, selection (deprecated), multiple(deprecated)).",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value of the existing custom field.",
			},
		},
	},
	Description: "Existing custom fields to associate to this ressource.",
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

func getVMIDForInterface(m interface{}, objectID int64) (int64, error) {
	client := m.(*netboxclient.NetBoxAPI)

	objectIDStr := fmt.Sprintf("%d", objectID)
	paramsInterface := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&objectIDStr)
	interfaces, err := client.Virtualization.VirtualizationInterfacesList(
		paramsInterface, nil)

	if err != nil {
		return 0, err
	}

	for _, i := range interfaces.Payload.Results {
		if i.ID == objectID {
			if i.VirtualMachine != nil {
				return i.VirtualMachine.ID, nil
			}
		}
	}
	return 0, fmt.Errorf("virtual machine not found")
}

func isprimary(m interface{}, objectID *int64, ipID int64, ip4 bool) (bool, error) {
	client := m.(*netboxclient.NetBoxAPI)

	if objectID == nil {
		return false, nil
	}

	objectIDStr := strconv.FormatInt(*objectID, 10)
	paramsInterface := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&objectIDStr)
	interfaces, err := client.Virtualization.VirtualizationInterfacesList(
		paramsInterface, nil)

	if err != nil {
		return false, err
	}

	for _, i := range interfaces.Payload.Results {
		if i.ID == *objectID {
			if i.VirtualMachine != nil {
				vmIDStr := fmt.Sprintf("%d", i.VirtualMachine.ID)
				paramsVM := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(&vmIDStr)
				vms, err := client.Virtualization.VirtualizationVirtualMachinesList(
					paramsVM, nil)

				if err != nil {
					return false, err
				}

				if *vms.Payload.Count != 1 {
					return false, fmt.Errorf("Cannot set an interface as primary when " +
						"the interface is not attached to a VM.")
				}

				if ip4 && vms.Payload.Results[0].PrimaryIp4 != nil {
					return vms.Payload.Results[0].PrimaryIp4.ID == ipID, nil
				} else if !ip4 && vms.Payload.Results[0].PrimaryIp6 != nil {
					return vms.Payload.Results[0].PrimaryIp6.ID == ipID, nil
				} else {
					return false, nil
				}
			} else {
				return false, fmt.Errorf("Cannot set an interface as primary when the " +
					"interface is not attached to a VM.")
			}
		}
	}

	return false, nil
}

func updatePrimaryStatus(m interface{}, vmid, ipid int64, primary bool) error {
	client := m.(*netboxclient.NetBoxAPI)

	emptyFields := make(map[string]interface{})
	dropFields := []string{
		"created",
		"last_updated",
		"name",
		"cluster",
		"tags",
	}

	params := &models.WritableVirtualMachineWithConfigContext{}
	if primary {
		params.PrimaryIp4 = &ipid
	} else {
		params.PrimaryIp4 = nil
		emptyFields["primary_ip4"] = nil
	}
	vm := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)
	vm.SetID(vmid)
	_, err := client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
		vm, nil, newRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return err
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

func convertArrayInterfaceJSONString(arrayInterface []interface{}) string {
	var arrayString []interface{}

	for _, item := range arrayInterface {
		switch v := item.(type) {
		case map[string]interface{}:
			arrayString = append(arrayString, string(v["id"].(json.Number)))
		default:
			arrayString = append(arrayString, v)
		}
	}

	if len(arrayString) > 1 {
		switch arrayString[0].(type) {
		case json.Number:
			sort.Slice(arrayString, func(i, j int) bool {
				return arrayString[i].(json.Number) < arrayString[j].(json.Number)
			})
		case string:
			sort.Slice(arrayString, func(i, j int) bool {
				return arrayString[i].(string) < arrayString[j].(string)
			})
		default:
			panic("")
		}
	}
	jsonValue, _ := json.Marshal(arrayString)
	result := string(jsonValue)

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

					if value != nil && cm["type"] == customFieldJSON {
						jsonValue, _ := json.Marshal(value)
						cm["value"] = string(jsonValue)
					} else if value != nil {
						switch v := value.(type) {
						case []interface{}:
							if cm["type"] == customFieldMultiSelectLegacy {
								strValue = convertArrayInterfaceString(v)
							} else if cm["type"] == customFieldObject {
								strValue = string(value.([]interface{})[0].(map[string]interface{})["id"].(json.Number))
							} else {
								strValue = convertArrayInterfaceJSONString(v)
							}
						case map[string]interface{}:
							strValue = string(v["id"].(json.Number))
						default:
							strValue = fmt.Sprintf("%v", v)
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
			switch cfType {
			case "integer":
				cfValueInt, _ := strconv.Atoi(cfValue)
				toReturn[cfName] = cfValueInt
			case customFieldBoolean:
				if cfValue == "true" {
					toReturn[cfName] = true
				} else if cfValue == "false" {
					toReturn[cfName] = false
				}
			case customFieldMultiSelectLegacy:
				cfValueArray := strings.Split(cfValue, ",")
				sort.Strings(cfValueArray)
				toReturn[cfName] = cfValueArray
			case customFieldJSON:
				jsonMap := make(map[string]interface{})
				err := json.Unmarshal([]byte(cfValue), &jsonMap)
				if err != nil {
					continue
				}
				toReturn[cfName] = jsonMap
			case customFieldMultiSelect:
				var jsonList []interface{}
				err := json.Unmarshal([]byte(cfValue), &jsonList)
				if err != nil {
					continue
				}
				sort.Slice(jsonList, func(i, j int) bool {
					return jsonList[i].(string) < jsonList[j].(string)
				})
				toReturn[cfName] = jsonList
			case customFieldObject:
				cfValueInt, _ := strconv.Atoi(cfValue)
				toReturn[cfName] = []int{cfValueInt}
			case customFieldMultiObject:
				var jsonList []interface{}
				err := json.Unmarshal([]byte(cfValue), &jsonList)
				if err != nil {
					continue
				}
				sort.Slice(jsonList, func(i, j int) bool {
					return jsonList[i].(string) < jsonList[j].(string)
				})
				toReturn[cfName] = jsonList
			default:
				toReturn[cfName] = cfValue
			}
		} else {
			toReturn[cfName] = nil
		}
	}

	return toReturn
}

// Convert URL in content_type
func convertURIContentType(uri strfmt.URI) string {
	uriSplit := strings.Split(uri.String(), "/")
	firstLevel := uriSplit[len(uriSplit)-4]
	secondLevel := uriSplit[len(uriSplit)-3]

	slRegexp := regexp.MustCompile(`s$`)
	secondLevel = strings.ReplaceAll(secondLevel, "-", "")
	secondLevel = slRegexp.ReplaceAllString(secondLevel, "")

	contentType := firstLevel + "." + secondLevel
	return contentType
}

func getNewAvailableIPForIPRange(client *netboxclient.NetBoxAPI, id int64) (*models.IPAddress, error) {
	params := ipam.NewIpamIPRangesAvailableIpsCreateParams().WithID(id)
	params.Data = []*models.WritableAvailableIP{
		{},
	}
	list, err := client.Ipam.IpamIPRangesAvailableIpsCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

func getNewAvailableIPForPrefix(client *netboxclient.NetBoxAPI, id int64) (*models.IPAddress, error) {
	params := ipam.NewIpamPrefixesAvailableIpsCreateParams().WithID(id)
	params.Data = []*models.WritableAvailableIP{
		{},
	}
	list, err := client.Ipam.IpamPrefixesAvailableIpsCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

func getNewAvailablePrefix(client *netboxclient.NetBoxAPI, id int64, length int64) (*models.Prefix, error) {
	params := ipam.NewIpamPrefixesAvailablePrefixesCreateParams().WithID(id)
	params.Data = []*models.PrefixLength{
		{PrefixLength: &length},
	}
	list, err := client.Ipam.IpamPrefixesAvailablePrefixesCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

type requestModifier struct {
	origwriter      runtime.ClientRequestWriter
	overwritefields map[string]interface{}
	dropfields      []string
}

func (o requestModifier) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {
	err := o.origwriter.WriteToRequest(r, reg)
	if err != nil {
		return err
	}

	jsonString, err := json.Marshal(r.GetBodyParam())
	if err != nil {
		return err
	}

	var objmap map[string]interface{}
	err = json.Unmarshal(jsonString, &objmap)
	if err != nil {
		return err
	}
	for _, k := range o.dropfields {
		delete(objmap, k)
	}
	for k, v := range o.overwritefields {
		objmap[k] = v
	}

	err = r.SetBodyParam(objmap)
	return err
}

func newRequestModifierOperation(empty map[string]interface{}, drop []string) func(*runtime.ClientOperation) {
	return func(op *runtime.ClientOperation) {
		if len(empty) > 0 || len(drop) > 0 {
			tmp := requestModifier{
				origwriter:      op.Params,
				overwritefields: empty,
				dropfields:      drop,
			}
			op.Params = tmp
		}
	}
}
