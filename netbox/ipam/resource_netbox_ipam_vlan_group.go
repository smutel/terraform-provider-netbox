package ipam

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
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

var vlanGroupRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"slug",
	"tags",
}

func resourceNetboxIpamVlanGroupCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	groupName := d.Get("name").(string)
	groupSlug := d.Get("slug").(string)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.VLANGroup{
		CustomFields: customFields,
		Description:  d.Get("description").(string),
		MaxVid:       int64(d.Get("max_vid").(int)),
		MinVid:       int64(d.Get("min_vid").(int)),
		Name:         &groupName,
		Slug:         &groupSlug,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	if scopes := d.Get("scope").(*schema.Set).List(); len(scopes) == 1 {
		scopeIntf := scopes[0].(map[string]interface{})
		scopeID := int64(scopeIntf["id"].(int))
		scopeType := scopeIntf["type"].(string)
		newResource.ScopeID = &scopeID
		newResource.ScopeType = &scopeType
	}

	resource := ipam.NewIpamVlanGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamVlanGroupsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))
	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamVlanGroupsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVlanGroupsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]

	if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("created", resource.Created.String()); err != nil {
		return diag.FromErr(err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

	if err = d.Set("custom_field", customFields); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("max_vid", resource.MaxVid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("min_vid", resource.MinVid); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}

	var scopes []map[string]interface{}
	if resource.Scope != nil {
		scopes = []map[string]interface{}{
			{
				"id":   resource.ScopeID,
				"type": resource.ScopeType,
			},
		}
	}
	if err = d.Set("scope", scopes); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("slug", resource.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vlan_count", resource.VlanCount); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamVlanGroupUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.VLANGroup{}
	modifiedFields := map[string]interface{}{}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		params.Description = description
		modifiedFields["description"] = description
	}
	if d.HasChange("max_vid") {
		params.MaxVid = int64(d.Get("max_vid").(int))
	}
	if d.HasChange("min_vid") {
		params.MinVid = int64(d.Get("min_vid").(int))
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("scope") {
		if scopes := d.Get("scope").(*schema.Set).List(); len(scopes) == 1 {
			scopeIntf := scopes[0].(map[string]interface{})
			scopeID := int64(scopeIntf["id"].(int))
			scopeType := scopeIntf["type"].(string)
			params.ScopeID = &scopeID
			params.ScopeType = &scopeType
			modifiedFields["scope_id"] = scopeID
			modifiedFields["scope_type"] = scopeType
		} else {
			modifiedFields["scope_id"] = nil
			modifiedFields["scope_type"] = nil
		}
	}
	if d.HasChange("slug") {
		slug := d.Get("slug").(string)
		params.Slug = &slug
	}
	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}

	resource := ipam.NewIpamVlanGroupsPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamVlanGroupsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, vlanGroupRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamVlanGroupRead(ctx, d, m)
}

func resourceNetboxIpamVlanGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamVlanGroupExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource := ipam.NewIpamVlanGroupsDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamVlanGroupsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamVlanGroupExists(d *schema.ResourceData, m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamVlanGroupsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVlanGroupsList(params, nil)
	if err != nil {
		return resourceExist, err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
