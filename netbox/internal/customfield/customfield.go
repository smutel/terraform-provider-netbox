package customfield

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

// Boolean string for custom field
const customFieldBoolean = "boolean"
const customFieldJSON = "json"
const customFieldMultiSelectLegacy = "multiple"
const customFieldMultiSelect = "multiselect"
const customFieldObject = "object"
const customFieldMultiObject = "multiobject"

var CustomFieldSchema = schema.Schema{
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
				ValidateFunc: validation.StringInSlice([]string{"text",
					"longtext", "integer", customFieldBoolean, "date", "url",
					customFieldJSON, "select", customFieldMultiSelect,
					customFieldObject, customFieldMultiObject,
					customFieldMultiSelectLegacy, "selection"}, false),
				Description: "Type of the existing custom field (text," +
					"longtext, integer, boolean, date, url, json, select, " +
					"multiselect, object, multiobject, selection " +
					"(deprecated), multiple(deprecated)).",
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

func convertArrayInterfaceString(arrayInterface []any) string {
	var arrayString []string

	for _, item := range arrayInterface {
		switch v := item.(type) {
		case string:
			arrayString = append(arrayString, v)
		case int:
			strV := strconv.FormatInt(int64(v), util.Const10)
			arrayString = append(arrayString, strV)
		}

	}

	sort.Strings(arrayString)
	result := strings.Join(arrayString, ",")

	return result
}

func convertArrayInterfaceJSONString(arrayInterface []any) string {
	var arrayString []any

	for _, item := range arrayInterface {
		switch v := item.(type) {
		case map[string]any:
			arrayString = append(arrayString,
				fmt.Sprintf("%.0f", v["id"].(float64)))
		default:
			arrayString = append(arrayString, v)
		}
	}

	if len(arrayString) > 1 {
		switch arrayString[0].(type) {
		case float64:
			sort.Slice(arrayString, func(i, j int) bool {
				return arrayString[i].(float64) < arrayString[j].(float64)
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

// Pick the custom fields in the state file
// and update values with data from API
func UpdateCustomFieldsFromAPI(stateCustomFields,
	customFields any) []map[string]string {
	var tfCms []map[string]string

	//nolint:revive
	switch t := customFields.(type) {
	case map[string]any:
		for _, stateCustomField := range stateCustomFields.([]any) {
			for key, value := range t {
				if stateCustomField.(map[string]any)["name"].(string) == key { //nolint:revive
					var strValue string

					cm := map[string]string{}
					cm["name"] = key
					cm["type"] =
						stateCustomField.(map[string]any)["type"].(string) //nolint:revive

					if value != nil && cm["type"] == customFieldJSON {
						jsonValue, _ := json.Marshal(value)
						cm["value"] = string(jsonValue)
					} else if value != nil {
						switch v := value.(type) {
						case []any:
							//nolint:revive
							if cm["type"] == customFieldMultiSelectLegacy {
								strValue = convertArrayInterfaceString(v)
							} else if cm["type"] == customFieldObject {
								strValue =
									string(
										value.([]any)[0].(map[string]any)["id"].(json.Number)) //nolint:revive
							} else {
								strValue = convertArrayInterfaceJSONString(v)
							}
						case map[string]any:
							strValue = fmt.Sprintf("%.0f", v["id"].(float64))
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
func ConvertCustomFieldsFromTerraformToAPI(stateCustomFields []any,
	customFields []any) map[string]any {

	toReturn := make(map[string]any)

	for _, stateCf := range stateCustomFields {
		stateCustomField := stateCf.(map[string]any)
		toReturn[stateCustomField["name"].(string)] = nil
	}

	for _, cf := range customFields {
		customField := cf.(map[string]any)

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
				jsonMap := make(map[string]any)
				err := json.Unmarshal([]byte(cfValue), &jsonMap)
				if err != nil {
					continue
				}
				toReturn[cfName] = jsonMap
			case customFieldMultiSelect:
				var jsonList []any
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
				toReturn[cfName] = cfValueInt
			case customFieldMultiObject:
				var jsonList []any
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
