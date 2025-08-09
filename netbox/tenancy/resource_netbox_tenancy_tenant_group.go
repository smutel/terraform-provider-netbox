// Copyright (c)
// SPDX-License-Identifier: MIT

package tenancy

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxTenancyTenantGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a tenant group within Netbox.",
		CreateContext: resourceNetboxTenancyTenantGroupCreate,
		ReadContext:   resourceNetboxTenancyTenantGroupRead,
		UpdateContext: resourceNetboxTenancyTenantGroupUpdate,
		DeleteContext: resourceNetboxTenancyTenantGroupDelete,
		Exists:        resourceNetboxTenancyTenantGroupExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this tenant group.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this tenant group was created.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last date when this tenant group was updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name for this tenant group.",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,100}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,100}$"),
				Description: "The slug for this tenant group.",
			},
			"tag": &tag.TagSchema,
		},
	}
}

func resourceNetboxTenancyTenantGroupCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritableTenantGroupRequestWithDefaults()
	newResource.SetName(name)
	newResource.SetSlug(slug)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	_, response, err :=
		client.TenancyAPI.TenancyTenantGroupsCreate(
			ctx).WritableTenantGroupRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxTenancyTenantGroupRead(ctx, d, m)
}

func resourceNetboxTenancyTenantGroupRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx,
			int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxTenancyTenantGroupUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableTenantGroupRequestWithDefaults()

	// Required parameters
	name := d.Get("name").(string)
	resource.SetName(name)

	slug := d.Get("slug").(string)
	resource.SetSlug(slug)

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.TenancyAPI.TenancyTenantGroupsUpdate(ctx,
		int32(resourceID)).WritableTenantGroupRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxTenancyTenantGroupRead(ctx, d, m)
}

func resourceNetboxTenancyTenantGroupDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxTenancyTenantGroupExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}

	if response, err := client.TenancyAPI.TenancyTenantGroupsDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxTenancyTenantGroupExists(d *schema.ResourceData,
	m any) (b bool,
	e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.TenancyAPI.TenancyTenantGroupsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
