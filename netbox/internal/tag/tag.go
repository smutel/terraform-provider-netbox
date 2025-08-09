// Copyright (c)
// SPDX-License-Identifier: MIT

package tag

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netbox "github.com/smutel/go-netbox/v4"
)

var TagSchema = schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the existing tag.",
			},
			"slug": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slug of the existing tag.",
			},
		},
	},
	Description: "Existing tag to associate to this resource.",
}

func ConvertTagsToNestedTagRequest(
	tags []any) []netbox.NestedTagRequest {

	nestedTags := []netbox.NestedTagRequest{}

	for _, tag := range tags {
		t := tag.(map[string]any)

		tagName := t["name"].(string)
		tagSlug := t["slug"].(string)

		nestedTag := netbox.NestedTagRequest{Name: tagName, Slug: tagSlug}
		nestedTags = append(nestedTags, nestedTag)
	}

	return nestedTags
}

func ConvertNestedTagRequestToTags(
	tags []netbox.NestedTag) []map[string]string {

	var tfTags []map[string]string

	for _, t := range tags {
		tag := map[string]string{}
		tag["name"] = t.Name
		tag["slug"] = t.Slug

		tfTags = append(tfTags, tag)
	}

	return tfTags
}
