package ipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxIpamIPRange() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an ip range within Netbox.",
		CreateContext: resourceNetboxIpamIPRangeCreate,
		ReadContext:   resourceNetboxIpamIPRangeRead,
		UpdateContext: resourceNetboxIpamIPRangeUpdate,
		DeleteContext: resourceNetboxIpamIPRangeDelete,
		Exists:        resourceNetboxIpamIPRangeExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this ip range.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this ip range was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The description of this prefix.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this site was last updated.",
			},
			"start_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The first address of the ip range",
			},
			"end_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The last address of the ip range",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of addresses in the ip range",
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the role attached to this prefix.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"active",
					"reserved", "deprecated"}, false),
				Description: "Status among active, reserved, deprecated " +
					"(active by default).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "ID of the tenant where this ip range " +
					"is attached.",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vrf attached to this ip range.",
			},
		},
	}
}

func resourceNetboxIpamIPRangeCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	startAddress := d.Get("start_address").(string)
	endAddress := d.Get("end_address").(string)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	description := d.Get("description").(string)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritableIPRangeRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetEndAddress(endAddress)
	newResource.SetStartAddress(startAddress)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	s, err := netbox.NewPatchedWritableIPRangeRequestStatusFromValue(status)
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

	if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	if vrfID := d.Get("vrf_id").(int); vrfID != 0 {
		b, err := brief.GetBriefVRFRequestFromID(ctx, client, vrfID)
		if err != nil {
			return err
		}
		newResource.SetVrf(*b)
	}

	_, response, err :=
		client.IpamAPI.IpamIpRangesCreate(ctx).WritableIPRangeRequest(
			*newResource).Execute()

	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamIPRangeRead(ctx, d, m)
}

func resourceNetboxIpamIPRangeRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.IpamAPI.IpamIpRangesRetrieve(ctx,
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
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields,
		resource.GetCustomFields())

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

	if err = d.Set("start_address", resource.GetStartAddress()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("end_address", resource.GetEndAddress()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("size", resource.GetSize()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("role_id", resource.GetRole().Id); err != nil {
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

	if err = d.Set("vrf_id", resource.GetVrf().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamIPRangeUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableIPRangeRequestWithDefaults()

	// Required parameters
	startAddress := d.Get("start_address").(string)
	endAddress := d.Get("end_address").(string)
	resource.SetStartAddress(startAddress)
	resource.SetEndAddress(endAddress)

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

	if d.HasChange("status") {
		s, err := netbox.NewPatchedWritableIPRangeRequestStatusFromValue(
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

	if _, response, err := client.IpamAPI.IpamIpRangesUpdate(ctx,
		int32(resourceID)).WritableIPRangeRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamIPRangeRead(ctx, d, m)
}

func resourceNetboxIpamIPRangeDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamIPRangeExists(d, m)
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

	if response, err := client.IpamAPI.IpamIpRangesDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamIPRangeExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamIpRangesRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
