package ipam

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxIpamPrefix() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a prefix (ipam module) within Netbox.",
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
				Description: "The content type of this prefix (ipam module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The description of this prefix (ipam module).",
			},
			"is_pool": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     nil,
				Description: "Define if this object is a pool (false by default).",
			},
			"prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 256),
				ExactlyOneOf: []string{"prefix", "parent_prefix"},
				Description:  "The prefix (IP address/mask) used for this prefix (ipam module). Required if parent_prefix is not set.",
			},
			"parent_prefix": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "Parent prefix and length used for new prefix. Required if prefix is not set",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Id of parent prefix",
						},
						"prefix_length": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "Length of new prefix",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 128)),
						},
					},
				},
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the role attached to this prefix (ipam module).",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the site where this prefix (ipam module) is located.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"container", "active",
					"reserved", "deprecated"}, false),
				Description: "Status among container, active, reserved, deprecated (active by default).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this prefix (ipam module) is attached.",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vlan where this prefix (ipam module) is attached.",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vrf attached to this prefix (ipam module).",
			},
		},
	}
}

func resourceNetboxIpamPrefixCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	var prefix string
	var prefixid int32
	update := false

	if stateprefix, ok := d.GetOk("prefix"); ok {
		prefix = stateprefix.(string)
	} else if pprefix, ok := d.GetOk("parent_prefix"); ok {
		set := pprefix.(*schema.Set)
		mappreffix := set.List()[0].(map[string]interface{})
		parentPrefix := int32(mappreffix["prefix"].(int))
		p, err := getNewAvailablePrefix(ctx, client, parentPrefix)
		if err != nil {
			return err
		}
		prefix = p.Prefix
		update = true
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	isPool := d.Get("is_pool").(bool)
	roleID := int32(d.Get("role_id").(int))
	siteID := int32(d.Get("site_id").(int))
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int32(d.Get("tenant_id").(int))
	vlanID := int32(d.Get("vlan_id").(int))
	vrfID := int32(d.Get("vrf_id").(int))

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

	if roleID != 0 {
		newResource.SetRole(roleID)
	}

	if siteID != 0 {
		newResource.SetSite(siteID)
	}

	if tenantID != 0 {
		newResource.SetTenant(tenantID)
	}

	if vlanID != 0 {
		newResource.SetVlan(vlanID)
	}

	if vrfID != 0 {
		newResource.SetVrf(vrfID)
	}

	var resource *netbox.Prefix
	var response *http.Response
	if !update {
		resource, response, err = client.IpamAPI.IpamPrefixesCreate(ctx).WritablePrefixRequest(*newResource).Execute()
	} else {
		resource, response, err = client.IpamAPI.IpamPrefixesUpdate(ctx, prefixid).WritablePrefixRequest(*newResource).Execute()
	}
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resource.GetId() == 0 {
		return diag.FromErr(errors.New("Bug Netbox - TO BE FIXED"))
	}

	prefixid = resource.GetId()

	d.SetId(strconv.FormatInt(int64(prefixid), 10))

	return resourceNetboxIpamPrefixRead(ctx, d, m)
}

func resourceNetboxIpamPrefixRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.IpamAPI.IpamPrefixesRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("is_pool", resource.GetIsPool()); err != nil {
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

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
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

func resourceNetboxIpamPrefixUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritablePrefixRequestWithDefaults()

	// Required parameters
	resource.SetPrefix(d.Get("prefix").(string))

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
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
		roleID := int32(d.Get("role_id").(int))
		if roleID != 0 {
			resource.SetRole(roleID)
		} else {
			resource.SetRoleNil()
		}
	}

	if d.HasChange("site_id") {
		siteID := int32(d.Get("role_id").(int))
		if siteID != 0 {
			resource.SetSite(siteID)
		} else {
			resource.SetSiteNil()
		}
	}

	if d.HasChange("status") {
		s, err := netbox.NewPatchedWritablePrefixRequestStatusFromValue(d.Get("status").(string))
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
		tenantID := int32(d.Get("tenant_id").(int))
		if tenantID != 0 {
			resource.SetTenant(tenantID)
		} else {
			resource.SetTenantNil()
		}
	}

	if d.HasChange("vlan_id") {
		vlanID := int32(d.Get("vlan_id").(int))
		if vlanID != 0 {
			resource.SetVlan(vlanID)
		} else {
			resource.SetVlanNil()
		}
	}

	if d.HasChange("vrf_id") {
		vrfID := int32(d.Get("vrf_id").(int))
		if vrfID != 0 {
			resource.SetVrf(vrfID)
		} else {
			resource.SetVrfNil()
		}
	}

	if _, response, err := client.IpamAPI.IpamPrefixesUpdate(ctx, int32(resourceID)).WritablePrefixRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamPrefixRead(ctx, d, m)
}

func resourceNetboxIpamPrefixDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamPrefixExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int64"))
	}

	if response, err := client.IpamAPI.IpamPrefixesDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamPrefixExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamPrefixesRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
