package ipam

import (
	"context"
	"errors"
	"fmt"
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

func ResourceNetboxIpamVlan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a vlan within Netbox.",
		CreateContext: resourceNetboxIpamVlanCreate,
		ReadContext:   resourceNetboxIpamVlanRead,
		UpdateContext: resourceNetboxIpamVlanUpdate,
		DeleteContext: resourceNetboxIpamVlanDelete,
		Exists:        resourceNetboxIpamVlanExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this vlan.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this vlan was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The description of this vlan.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last date when this vlan was updated.",
			},
			"vlan_group_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "ID of the group where this vlan belongs to.",
				ConflictsWith: []string{"site_id"},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description:  "The name for this vlan.",
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the role attached to this vlan.",
			},
			"site_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "ID of the site where this vlan is located.",
				ConflictsWith: []string{"vlan_group_id"},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"active",
					"reserved", "deprecated"}, false),
				Description: "The description of this vlan.",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this vlan is attached.",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the vlan (vlan tag).",
			},
		},
	}
}

func resourceNetboxIpamVlanCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
		nil, resourceCustomFields)
	name := d.Get("name").(string)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritableVLANRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetName(name)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	vid32, err := safecast.ToInt32(d.Get("vlan_id").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetVid(vid32)

	s, err := netbox.NewPatchedWritableVLANRequestStatusFromValue(status)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetStatus(*s)

	if groupID := d.Get("vlan_group_id").(int); groupID != 0 {
		b, err := brief.GetBriefVlanGroupRequestFromID(ctx, client, groupID)
		if err != nil {
			return err
		}
		newResource.SetGroup(*b)
	}

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

	_, response, err := client.IpamAPI.IpamVlansCreate(
		ctx).WritableVLANRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamVlanRead(ctx, d, m)
}

func resourceNetboxIpamVlanRead(ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.IpamAPI.IpamVlansRetrieve(ctx,
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

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(
		resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vlan_group_id", resource.GetGroup().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
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

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(
		resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tenant_id", resource.GetTenant().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vlan_id", resource.GetVid()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamVlanUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableVLANRequestWithDefaults()

	// Required parameters
	resource.SetName(d.Get("name").(string))

	vid32, err := safecast.ToInt32(d.Get("vlan_id").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	resource.SetVid(vid32)

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields :=
			customfield.ConvertCustomFieldsFromTerraformToAPI(
				stateCustomFields.(*schema.Set).List(),
				resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("vlan_group_id") {
		if groupID := d.Get("vlan_group_id").(int); groupID != 0 {
			b, err := brief.GetBriefVlanGroupRequestFromID(ctx, client, groupID)
			if err != nil {
				return err
			}
			resource.SetGroup(*b)
		} else {
			resource.SetGroupNil()
		}
	}

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
		s, err :=
			netbox.NewPatchedWritableVLANRequestStatusFromValue(
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

	if _, response, err := client.IpamAPI.IpamVlansUpdate(ctx,
		int32(resourceID)).WritableVLANRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamVlanRead(ctx, d, m)
}

func resourceNetboxIpamVlanDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamVlanExists(d, m)
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

	if response, err := client.IpamAPI.IpamVlansDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamVlanExists(d *schema.ResourceData,
	m any) (b bool, e error) {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamVlansRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
