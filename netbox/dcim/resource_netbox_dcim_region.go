// Copyright (c)
// SPDX-License-Identifier: MIT

package dcim

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxDcimRegion() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a region within Netbox.",
		CreateContext: resourceNetboxDcimRegionCreate,
		ReadContext:   resourceNetboxDcimRegionRead,
		UpdateContext: resourceNetboxDcimRegionUpdate,
		DeleteContext: resourceNetboxDcimRegionDelete,
		Exists:        resourceNetboxDcimRegionExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this region.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this region was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"depth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Depth of this region.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, util.Const200),
				Description:  "Description of this region.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last date when this region was updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name of this region.",
			},
			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the parent for this region.",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,100}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,100}$"),
				Description: "The slug, unique identifier used in URLs," +
					"of this region.",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this region.",
			},
		},
	}
}

func resourceNetboxDcimRegionCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()
	parentID := d.Get("parent_id").(int)

	newResource := netbox.NewWritableRegionRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetName(name)
	newResource.SetSlug(slug)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	if parentID != 0 {
		parentID32, err := safecast.ToInt32(parentID)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetParent(parentID32)
	}

	_, response, err := client.DcimAPI.DcimRegionsCreate(
		ctx).WritableRegionRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxDcimRegionRead(ctx, d, m)
}

func resourceNetboxDcimRegionRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.DcimAPI.DcimRegionsRetrieve(ctx,
		int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(
		resourceCustomFields,
		resource.GetCustomFields())

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("depth", resource.GetDepth()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
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

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(
		resource.GetTags())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimRegionUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableRegionRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))
	resource.SetSlug(d.Get("slug").(string))

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields :=
			customfield.ConvertCustomFieldsFromTerraformToAPI(
				stateCustomFields.(*schema.Set).List(),
				resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("parent_id") {
		if parentID, exist := d.GetOk("parent_id"); exist {
			parentID32, err := safecast.ToInt32(parentID.(int))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			resource.SetParent(parentID32)
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.DcimAPI.DcimRegionsUpdate(ctx,
		int32(resourceID)).WritableRegionRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimRegionRead(ctx, d, m)
}

func resourceNetboxDcimRegionDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimRegionExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}

	if response, err := client.DcimAPI.DcimRegionsDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimRegionExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimRegionsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
