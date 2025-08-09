// Copyright (c)
// SPDX-License-Identifier: MIT

package dcim

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxDcimManufacturer() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a manufacturer within Netbox.",
		CreateContext: resourceNetboxDcimManufacturerCreate,
		ReadContext:   resourceNetboxDcimManufacturerRead,
		UpdateContext: resourceNetboxDcimManufacturerUpdate,
		DeleteContext: resourceNetboxDcimManufacturerDelete,
		Exists:        resourceNetboxDcimManufacturerExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this manufacturer.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this manufacturer was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, util.Const100),
				Description:  "The description of this manufacturer.",
			},
			"devicetype_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of device types of this manufacturer.",
			},
			"inventoryitem_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The number of inventory items of this " +
					"manufacturer.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this manufacturer was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name of this manufacturer.",
			},
			"platform_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of platforms of this manufacturer.",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The slug of this manufacturer.",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this manufacturer.",
			},
		},
	}
}

func resourceNetboxDcimManufacturerCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewManufacturerRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetName(name)
	newResource.SetSlug(slug)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	_, response, err := client.DcimAPI.DcimManufacturersCreate(
		ctx).ManufacturerRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxDcimManufacturerRead(ctx, d, m)
}

func resourceNetboxDcimManufacturerRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.DcimAPI.DcimManufacturersRetrieve(ctx,
			int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type",
		util.ConvertURLContentType(resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(
		resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("devicetype_count",
		resource.GetDevicetypeCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("inventoryitem_count",
		resource.GetInventoryitemCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("platform_count",
		resource.GetPlatformCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimManufacturerUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewManufacturerRequestWithDefaults()

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
		if description, exist := d.GetOk("description"); exist {
			resource.SetDescription(description.(string))
		} else {
			resource.SetDescription("")
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.DcimAPI.DcimManufacturersUpdate(ctx,
		int32(resourceID)).ManufacturerRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimManufacturerRead(ctx, d, m)
}

func resourceNetboxDcimManufacturerDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimManufacturerExists(d, m)
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

	if response, err := client.DcimAPI.DcimManufacturersDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimManufacturerExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimManufacturersRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
