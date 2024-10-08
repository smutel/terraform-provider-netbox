package ipam

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxIpamVlanGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a vlan group (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamVlanGroupCreate,
		ReadContext:   resourceNetboxIpamVlanGroupRead,
		UpdateContext: resourceNetboxIpamVlanGroupUpdate,
		DeleteContext: resourceNetboxIpamVlanGroupDelete,
		Exists:        resourceNetboxIpamVlanGroupExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this vlan group (ipam module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this resource was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
				Description:  "The description of this vlan group (ipam module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this resource was last updated.",
			},
			"max_vid": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      4094,
				ValidateFunc: validation.IntBetween(1, 4094),
				Description:  "Highest permissible ID of a child vlan (ipam module).",
			},
			"min_vid": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 4094),
				Description:  "Lowest permissible ID of a child vlan (ipam module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The name for this vlan group (ipam module).",
			},
			"scope": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "ID of the scope object for this vlan group (ipam module).",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"dcim.location", "dcim.rack", "dcim.region", "dcim.site", "dcim.sitegroup", "virtualization.cluster", "virtualization.clustergroup"}, false),
							Description:  "Type of the scope object. Must me one of \"dcim.location\", \"dcim.rack\", \"dcim.region\", \"dcim.site\", \"dcim.sitegroup\", \"virtualization.cluster\", \"virtualization.clustergroup\".",
						},
					},
				},
				Description: "Scope of this vlan group (ipam module).",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
				Description: "The slug for this vlan group (ipam module).",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this vlan group (ipam module).",
			},
			"vlan_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of vlans assigned to this vlan group (ipam module).",
			},
		},
	}
}

func resourceNetboxIpamVlanGroupCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewVLANGroupRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetMaxVid(int32(d.Get("max_vid").(int)))
	newResource.SetMinVid(int32(d.Get("min_vid").(int)))
	newResource.SetName(d.Get("name").(string))
	newResource.SetSlug(d.Get("slug").(string))
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	if scopes := d.Get("scope").(*schema.Set).List(); len(scopes) == 1 {
		scopeIntf := scopes[0].(map[string]interface{})
		scopeID := int32(scopeIntf["id"].(int))
		scopeType := scopeIntf["type"].(string)
		newResource.SetScopeId(scopeID)
		newResource.SetScopeType(scopeType)
	}

	resourceCreated, response, err := client.IpamAPI.IpamVlanGroupsCreate(ctx).VLANGroupRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resourceCreated.GetId() == 0 {
		return diag.FromErr(errors.New("Bug Netbox - TO BE FIXED"))
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.IpamAPI.IpamVlanGroupsRetrieve(ctx, int32(resourceID)).Execute()

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

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
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

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("max_vid", resource.GetMaxVid()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("min_vid", resource.GetMinVid()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	var scopes []map[string]interface{}
	if resource.Scope != nil {
		scopes = []map[string]interface{}{
			{
				"id":   resource.GetScopeId(),
				"type": resource.GetScopeType(),
			},
		}
	}

	if err = d.Set("scope", scopes); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vlan_count", resource.GetVlanCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamVlanGroupUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewVLANGroupRequestWithDefaults()

	// Required parameters
	resource.SetName(d.Get("name").(string))

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("max_vid") {
		resource.SetMaxVid(int32(d.Get("max_vid").(int)))
	}

	if d.HasChange("min_vid") {
		resource.SetMinVid(int32(d.Get("min_vid").(int)))
	}

	if d.HasChange("scope") {
		if scopes := d.Get("scope").(*schema.Set).List(); len(scopes) == 1 {
			scopeIntf := scopes[0].(map[string]interface{})
			scopeID := int32(scopeIntf["id"].(int))
			scopeType := scopeIntf["type"].(string)
			resource.SetScopeId(scopeID)
			resource.SetScopeType(scopeType)
		} else {
			resource.SetScopeIdNil()
			resource.SetScopeTypeNil()
		}
	}

	if d.HasChange("slug") {
		resource.SetSlug(d.Get("slug").(string))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.IpamAPI.IpamVlanGroupsUpdate(ctx, int32(resourceID)).VLANGroupRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamVlanGroupExists(d, m)
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

	if response, err := client.IpamAPI.IpamVlanGroupsDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamVlanGroupExists(d *schema.ResourceData, m interface{}) (b bool,
	e error) {
	client := m.(*netbox.APIClient)
	resourceExist := false

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamVlanGroupsRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}

	return resourceExist, nil
}
