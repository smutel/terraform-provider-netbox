// Copyright (c)
// SPDX-License-Identifier: MIT

package ipam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxIpamPrefix() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a prefix within Netbox.",
		CreateContext: resourceNetboxIpamPrefixCreate,
		ReadContext:   resourceNetboxIpamPrefixRead,
		UpdateContext: resourceNetboxIpamPrefixUpdate,
		DeleteContext: resourceNetboxIpamPrefixDelete,
		Exists:        resourceNetboxIpamPrefixExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this prefix.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this prefix was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The description of this prefix.",
			},
			"is_pool": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,
				Description: "Define if this object is a pool " +
					"(false by default).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this site was last updated.",
			},
			"prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, util.Const256),
				ExactlyOneOf: []string{"prefix", "parent_prefix"},
				Description: "The prefix (IP address/mask) used " +
					"for this prefix. Required if parent_prefix is not set.",
			},
			"parent_prefix": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Description: "Parent prefix and length used for new prefix. " +
					"Required if prefix is not set",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Id of parent prefix",
						},
						"prefix_length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Length of new prefix",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, util.Const128)),
						},
					},
				},
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the role attached to this prefix.",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the site where this prefix is located.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"container",
					"active", "reserved", "deprecated"}, false),
				Description: "Status among container, active, reserved," +
					" deprecated (active by default).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this prefix is attached.",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vlan where this prefix is attached.",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vrf attached to this prefix.",
			},
		},
	}
}

func resourceNetboxIpamPrefixCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	var prefix string
	var prefixid int32
	update := false

	if stateprefix, ok := d.GetOk("prefix"); ok {
		prefix = stateprefix.(string)
	} else if pprefix, ok := d.GetOk("parent_prefix"); ok {
		set := pprefix.(*schema.Set)
		mappreffix := set.List()[0].(map[string]any)
		parentPrefix32, err := safecast.ToInt32(mappreffix["prefix"].(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		p, errDiag := getNewAvailablePrefix(ctx, client, parentPrefix32)
		if errDiag != nil {
			return errDiag
		}
		prefix = p.Prefix
		update = true
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
		nil, resourceCustomFields)
	description := d.Get("description").(string)
	isPool := d.Get("is_pool").(bool)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritablePrefixRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetIsPool(isPool)
	newResource.SetPrefix(prefix)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	s, err := netbox.NewPatchedWritablePrefixRequestStatusFromValue(status)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetStatus(*s)

	if roleID := d.Get("role_id").(int); roleID != 0 {
		b, err := brief.GetBriefVlanRoleRequestFromID(ctx, client, roleID)
		if err != nil {
			return err
		}
		newResource.SetRole(*b)
	}

	if siteID := d.Get("site_id").(int); siteID != 0 {
		b, err := brief.GetBriefSiteRequestFromID(ctx, client, siteID)
		if err != nil {
			return err
		}
		newResource.SetSite(*b)
	}

	if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	if vlanID := d.Get("vlan_id").(int); vlanID != 0 {
		b, err := brief.GetBriefVLANRequestFromID(ctx, client, vlanID)
		if err != nil {
			return err
		}
		newResource.SetVlan(*b)
	}

	if vrfID := d.Get("vrf_id").(int); vrfID != 0 {
		b, err := brief.GetBriefVRFRequestFromID(ctx, client, vrfID)
		if err != nil {
			return err
		}
		newResource.SetVrf(*b)
	}

	var response *http.Response
	if !update {
		_, response, err = client.IpamAPI.IpamPrefixesCreate(
			ctx).WritablePrefixRequest(*newResource).Execute()
	} else {
		_, response, err = client.IpamAPI.IpamPrefixesUpdate(
			ctx, prefixid).WritablePrefixRequest(*newResource).Execute()
	}
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamPrefixRead(ctx, d, m)
}

func resourceNetboxIpamPrefixRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.IpamAPI.IpamPrefixesRetrieve(
		ctx, int32(resourceID)).Execute()

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

	if err = d.Set("is_pool", resource.GetIsPool()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("prefix", resource.GetPrefix()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("role_id", resource.GetRole().Id); err != nil {
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

	if err = d.Set("vlan_id", resource.GetVlan().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vrf_id", resource.GetVrf().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamPrefixUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritablePrefixRequestWithDefaults()

	// Required parameters
	resource.SetPrefix(d.Get("prefix").(string))

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

	resource.SetIsPool(d.Get("is_pool").(bool))

	if d.HasChange("role_id") {
		if roleID := d.Get("role_id").(int); roleID != 0 {
			b, err := brief.GetBriefVlanRoleRequestFromID(ctx, client, roleID)
			if err != nil {
				return err
			}
			resource.SetRole(*b)
		} else {
			resource.SetRoleNil()
		}
	}

	if d.HasChange("site_id") {
		if siteID := d.Get("site_id").(int); siteID != 0 {
			b, err := brief.GetBriefSiteRequestFromID(ctx, client, siteID)
			if err != nil {
				return err
			}
			resource.SetSite(*b)
		} else {
			resource.SetSiteNil()
		}
	}

	if d.HasChange("status") {
		s, err := netbox.NewPatchedWritablePrefixRequestStatusFromValue(
			d.Get("status").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetStatus(*s)
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
			b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
			if err != nil {
				return err
			}
			resource.SetTenant(*b)
		} else {
			resource.SetTenantNil()
		}
	}

	if d.HasChange("vlan_id") {
		if vlanID := d.Get("vlan_id").(int); vlanID != 0 {
			b, err := brief.GetBriefVLANRequestFromID(ctx, client, vlanID)
			if err != nil {
				return err
			}
			resource.SetVlan(*b)
		} else {
			resource.SetVlanNil()
		}
	}

	if d.HasChange("vrf_id") {
		if vrfID := d.Get("vrf_id").(int); vrfID != 0 {
			b, err := brief.GetBriefVRFRequestFromID(ctx, client, vrfID)
			if err != nil {
				return err
			}
			resource.SetVrf(*b)
		} else {
			resource.SetVrfNil()
		}
	}

	if _, response, err := client.IpamAPI.IpamPrefixesUpdate(ctx,
		int32(resourceID)).WritablePrefixRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamPrefixRead(ctx, d, m)
}

func resourceNetboxIpamPrefixDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamPrefixExists(d, m)
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

	if response, err := client.IpamAPI.IpamPrefixesDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamPrefixExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, httpRequest, err := client.IpamAPI.IpamPrefixesRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && httpRequest.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && httpRequest.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
