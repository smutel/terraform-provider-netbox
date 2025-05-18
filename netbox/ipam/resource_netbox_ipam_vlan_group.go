package ipam

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

func ResourceNetboxIpamVlanGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a vlan group within Netbox.",
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
				Description: "The content type of this vlan group.",
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
				ValidateFunc: validation.StringLenBetween(0, util.Const100),
				Description:  "The description of this vlan group.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this resource was last updated.",
			},
			"max_vid": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      util.Const4094,
				ValidateFunc: validation.IntBetween(1, util.Const4094),
				Description:  "Highest permissible ID of a child vlan.",
			},
			"min_vid": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, util.Const4094),
				Description:  "Lowest permissible ID of a child vlan.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const50),
				Description:  "The name for this vlan group.",
			},
			"scope": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "ID of the scope object for " +
								"this vlan group.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"dcim.location", "dcim.rack", "dcim.region",
								"dcim.site", "dcim.sitegroup",
								"virtualization.cluster",
								"virtualization.clustergroup"}, false),
							Description: "Type of the scope object. " +
								"Must me one of " +
								"\"dcim.location\", \"dcim.rack\", " +
								"\"dcim.region\", " +
								"\"dcim.site\", \"dcim.sitegroup\", " +
								"\"virtualization.cluster\", " +
								"\"virtualization.clustergroup\".",
						},
					},
				},
				Description: "Scope of this vlan group.",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
				Description: "The slug for this vlan group.",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this vlan group.",
			},
			"vlan_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of vlans assigned to this vlan group.",
			},
		},
	}
}

func resourceNetboxIpamVlanGroupCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewVLANGroupRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetName(d.Get("name").(string))
	newResource.SetSlug(d.Get("slug").(string))
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	maxVid32, err := safecast.ToInt32(d.Get("max_vid").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetMaxVid(maxVid32)

	minVid32, err := safecast.ToInt32(d.Get("min_vid").(int))
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetMinVid(minVid32)

	if scopes := d.Get("scope").(*schema.Set).List(); len(scopes) == 1 {
		scopeIntf := scopes[0].(map[string]any)
		scopeType := scopeIntf["type"].(string)
		scopeID32, err := safecast.ToInt32(scopeIntf["id"].(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetScopeId(scopeID32)
		newResource.SetScopeType(scopeType)
	}

	_, response, err := client.IpamAPI.IpamVlanGroupsCreate(
		ctx).VLANGroupRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.IpamAPI.IpamVlanGroupsRetrieve(ctx,
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

	if err = d.Set("max_vid", resource.GetMaxVid()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("min_vid", resource.GetMinVid()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	var scopes []map[string]any
	if resource.Scope != nil {
		scopes = []map[string]any{
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

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vlan_count", resource.GetVlanCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamVlanGroupUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewVLANGroupRequestWithDefaults()

	// Required parameters
	resource.SetName(d.Get("name").(string))
	resource.SetSlug(d.Get("slug").(string))

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(),
			resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("max_vid") {
		maxVid32, err := safecast.ToInt32(d.Get("max_vid").(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetMaxVid(maxVid32)
	}

	if d.HasChange("min_vid") {
		minVid32, err := safecast.ToInt32(d.Get("min_vid").(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetMinVid(minVid32)
	}

	if d.HasChange("scope") {
		if scopes := d.Get("scope").(*schema.Set).List(); len(scopes) == 1 {
			scopeIntf := scopes[0].(map[string]any)
			scopeType := scopeIntf["type"].(string)
			scopeID32, err := safecast.ToInt32(scopeIntf["id"].(int))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			resource.SetScopeId(scopeID32)
			resource.SetScopeType(scopeType)
		} else {
			resource.SetScopeIdNil()
			resource.SetScopeTypeNil()
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.IpamAPI.IpamVlanGroupsUpdate(ctx,
		int32(resourceID)).VLANGroupRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamVlanGroupExists(d, m)
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

	if response, err := client.IpamAPI.IpamVlanGroupsDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamVlanGroupExists(d *schema.ResourceData,
	m any) (b bool, e error) {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamVlanGroupsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
