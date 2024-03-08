package dcim

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v3"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxDcimRegion() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a region (dcim module) within Netbox.",
		CreateContext: resourceNetboxDcimRegionCreate,
		ReadContext:   resourceNetboxDcimRegionRead,
		UpdateContext: resourceNetboxDcimRegionUpdate,
		DeleteContext: resourceNetboxDcimRegionDelete,
		Exists:        resourceNetboxDcimRegionExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this region (dcim module) was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"depth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Depth of this region (dcim module).",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
				Description:  "Description of this region (dcim module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this region was last updated (dcim module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this region (dcim module).",
			},
			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the parent for this region (dcim module).",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this site (dcim module).",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this region (dcim module).",
			},
		},
	}
}

func resourceNetboxDcimRegionCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()
	parentID := int32(d.Get("parent_id").(int))

	newResource := netbox.NewWritableRegionRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetName(name)
	newResource.SetSlug(slug)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	if parentID != 0 {
		newResource.SetParent(parentID)
	}

	resourceCreated, response, err := client.DcimAPI.DcimRegionsCreate(ctx).WritableRegionRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resourceCreated.GetId() == 0 {
		return util.GenerateErrorMessage(response, errors.New("Bug Netbox - TO BE FIXED"))
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxDcimRegionRead(ctx, d, m)
}

func resourceNetboxDcimRegionRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.DcimAPI.DcimRegionsRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("depth", resource.GetDepth()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("parent_id", resource.GetParent().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.GetTags())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimRegionUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	resource := netbox.NewWritableRegionRequestWithDefaults()

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("name") {
		resource.SetName(d.Get("name").(string))
	}

	if d.HasChange("parent_id") {
		if parentID, exist := d.GetOk("parent_id"); exist {
			resource.SetParent(int32(parentID.(int)))
		} else {
			resource.SetParentNil()
		}
	}

	if d.HasChange("slug") {
		resource.SetSlug(d.Get("slug").(string))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.DcimAPI.DcimRegionsUpdate(ctx, int32(resourceID)).WritableRegionRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimRegionRead(ctx, d, m)
}

func resourceNetboxDcimRegionDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimRegionExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	if response, err := client.DcimAPI.DcimRegionsDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimRegionExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimRegionsRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
