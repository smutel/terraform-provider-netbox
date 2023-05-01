package customfield

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
func UpdateCustomFieldsFromAPI(stateCustomFields, customFields interface{}) []map[string]string {
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
func ConvertCustomFieldsFromTerraformToAPI(stateCustomFields []interface{}, customFields []interface{}) map[string]interface{} {
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
				toReturn[cfName] = cfValueInt
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
