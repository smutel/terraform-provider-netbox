package tag

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/go-netbox/v3/netbox/models"
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

func ConvertTagsToNestedTags(tags []interface{}) []*models.NestedTag {
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

func ConvertTagsToNestedTagRequest(tags []interface{}) []netbox.NestedTagRequest {
	nestedTags := []netbox.NestedTagRequest{}

	for _, tag := range tags {
		t := tag.(map[string]interface{})

		tagName := t["name"].(string)
		tagSlug := t["slug"].(string)

		nestedTag := netbox.NestedTagRequest{Name: tagName, Slug: tagSlug}
		nestedTags = append(nestedTags, nestedTag)
	}

	return nestedTags
}

func ConvertNestedTagsToTags(tags []*models.NestedTag) []map[string]string {
	var tfTags []map[string]string

	for _, t := range tags {
		tag := map[string]string{}
		tag["name"] = *t.Name
		tag["slug"] = *t.Slug

		tfTags = append(tfTags, tag)
	}

	return tfTags
}

func ConvertNestedTagRequestToTags(tags []netbox.NestedTag) []map[string]string {
	var tfTags []map[string]string

	for _, t := range tags {
		tag := map[string]string{}
		tag["name"] = t.Name
		tag["slug"] = t.Slug

		tfTags = append(tfTags, tag)
	}

	return tfTags
}
