// Copyright (c)
// SPDX-License-Identifier: MIT

package virtualization

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxVirtualizationCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a cluster within Netbox.",
		CreateContext: resourceNetboxVirtualizationClusterCreate,
		ReadContext:   resourceNetboxVirtualizationClusterRead,
		UpdateContext: resourceNetboxVirtualizationClusterUpdate,
		DeleteContext: resourceNetboxVirtualizationClusterDelete,
		Exists:        resourceNetboxVirtualizationClusterExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this cluster.",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   util.TrimString,
				Description: "Comments for this cluster.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this cluster was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of devices in this cluster.",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The cluster group of this cluster.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this cluster was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name of this cluster.",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The site of this cluster.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"offline",
					"active", "planned", "staging", "decommissioning"}, false),
				Description: "The status among offline, active, planned, " +
					"staging or decommissioning (active by default).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the tenant where this cluster is attached.",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Default:     nil,
				Description: "Type of this cluster.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this cluster.",
			},
			"virtualmachine_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of virtual machines in this cluster.",
			},
		},
	}
}

func resourceNetboxVirtualizationClusterCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	groupID := d.Get("group_id").(int)
	name := d.Get("name").(string)
	typeID := d.Get("type_id").(int)
	siteID := d.Get("site_id").(int)
	tenantID := d.Get("tenant_id").(int)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritableClusterRequestWithDefaults()
	newResource.SetComments(d.Get("comments").(string))
	newResource.SetCustomFields(customFields)
	newResource.SetName(name)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	b, err := brief.GetBriefClusterTypeRequestFromID(ctx, client, typeID)
	if err != nil {
		return err
	}
	newResource.SetType(*b)

	if groupID != 0 {
		b, err := brief.GetBriefClusterGroupRequestFromID(ctx, client, groupID)
		if err != nil {
			return err
		}
		newResource.SetGroup(*b)
	}

	s, errDiag := netbox.NewClusterStatusValueFromValue(
		d.Get("status").(string))
	if errDiag != nil {
		return util.GenerateErrorMessage(nil, errDiag)
	}
	newResource.SetStatus(*s)

	if siteID != 0 {
		b, err := brief.GetBriefSiteRequestFromID(ctx, client, siteID)
		if err != nil {
			return err
		}
		newResource.SetSite(*b)
	}

	if tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	_, response, errDiag :=
		client.VirtualizationAPI.VirtualizationClustersCreate(
			ctx).WritableClusterRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && errDiag != nil {
		return util.GenerateErrorMessage(response, errDiag)
	}

	var resourceID int32
	if resourceID, errDiag = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, errDiag)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxVirtualizationClusterRead(ctx, d, m)
}

func resourceNetboxVirtualizationClusterRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.VirtualizationAPI.VirtualizationClustersRetrieve(ctx,
			int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("comments", resource.GetComments()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields,
		resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("device_count", resource.GetDeviceCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("group_id", resource.GetGroup().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("site_id", resource.GetSite().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("status", resource.GetStatus().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tenant_id", resource.GetTenant().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("type_id", resource.GetType().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("virtualmachine_count",
		resource.GetVirtualmachineCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxVirtualizationClusterUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableClusterRequestWithDefaults()

	// Required parameters
	resource.SetName(d.Get("name").(string))
	typeID := d.Get("type_id").(int)
	b, errDiag := brief.GetBriefClusterTypeRequestFromID(ctx, client, typeID)
	if errDiag != nil {
		return errDiag
	}
	resource.SetType(*b)

	if d.HasChange("comments") {
		resource.SetComments(d.Get("comments").(string))
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields :=
			customfield.ConvertCustomFieldsFromTerraformToAPI(
				stateCustomFields.(*schema.Set).List(),
				resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("group_id") {
		groupID := d.Get("group_id").(int)
		if groupID != 0 {
			b, err := brief.GetBriefClusterGroupRequestFromID(
				ctx, client, groupID)
			if err != nil {
				return err
			}
			resource.SetGroup(*b)
		} else {
			resource.SetGroupNil()
		}
	}

	if d.HasChange("status") {
		s, err := netbox.NewClusterStatusValueFromValue(
			d.Get("status").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetStatus(*s)
	}

	if d.HasChange("site_id") {
		siteID := d.Get("site_id").(int)
		if siteID != 0 {
			b, err := brief.GetBriefSiteRequestFromID(ctx, client, siteID)
			if err != nil {
				return err
			}
			resource.SetSite(*b)
		} else {
			resource.SetSiteNil()
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		tenantID := d.Get("tenant_id").(int)
		if tenantID != 0 {
			b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
			if err != nil {
				return err
			}
			resource.SetTenant(*b)
		} else {
			resource.SetTenantNil()
		}
	}

	if _, response, err :=
		client.VirtualizationAPI.VirtualizationClustersUpdate(ctx,
			int32(resourceID)).WritableClusterRequest(
			*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxVirtualizationClusterRead(ctx, d, m)
}

func resourceNetboxVirtualizationClusterDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxVirtualizationClusterExists(d, m)
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

	if response, err :=
		client.VirtualizationAPI.VirtualizationClustersDestroy(ctx,
			int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxVirtualizationClusterExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err :=
		client.VirtualizationAPI.VirtualizationClustersRetrieve(nil,
			int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
