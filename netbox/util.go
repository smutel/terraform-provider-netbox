package netbox

import (
	"github.com/netbox-community/go-netbox/netbox/models"
)

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
