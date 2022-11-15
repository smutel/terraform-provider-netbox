package netbox

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/customfield"
)

func resourceNetboxIpamVlan() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a vlan (ipam module) within Netbox.",
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
				Description: "The content type of this vlan (ipam module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The description of this vlan (ipam module).",
			},
			"vlan_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the group where this vlan (ipam module) belongs to.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The name for this vlan (ipam module).",
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the role attached to this vlan (ipam module).",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the site where this vlan (ipam module) is located.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"active", "reserved",
					"deprecated"}, false),
				Description: "The description of this vlan (ipam module).",
			},
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the existing tag.",
						},
						"slug": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Slug of the existing tag.",
						},
					},
				},
				Description: "Existing tag to associate to this vlan (ipam module).",
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this vlan (ipam module) is attached.",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the vlan (vlan tag).",
			},
		},
	}
}

func resourceNetboxIpamVlanCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	groupID := int64(d.Get("vlan_group_id").(int))
	name := d.Get("name").(string)
	roleID := int64(d.Get("role_id").(int))
	siteID := int64(d.Get("site_id").(int))
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))
	vid := int64(d.Get("vlan_id").(int))

	newResource := &models.WritableVLAN{
		CustomFields: &customFields,
		Description:  description,
		Name:         &name,
		Status:       status,
		Tags:         convertTagsToNestedTags(tags),
		Vid:          &vid,
	}

	if groupID != 0 {
		newResource.Group = &groupID
	}

	if roleID != 0 {
		newResource.Role = &roleID
	}

	if siteID != 0 {
		newResource.Site = &siteID
	}

	if tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	resource := ipam.NewIpamVlansCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamVlansCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))
	return resourceNetboxIpamVlanRead(ctx, d, m)
}

func resourceNetboxIpamVlanRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamVlansListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVlansList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return diag.FromErr(err)
			}

			var description interface{}
			if resource.Description == "" {
				description = nil
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return diag.FromErr(err)
			}

			if resource.Group == nil {
				if err = d.Set("vlan_group_id", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("vlan_group_id", resource.Group.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("name", resource.Name); err != nil {
				return diag.FromErr(err)
			}

			if resource.Role == nil {
				if err = d.Set("role_id", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("role_id", resource.Role.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if resource.Site == nil {
				if err = d.Set("site_id", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("site_id", resource.Site.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if resource.Status == nil {
				if err = d.Set("status", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("status", resource.Status.Value); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			if resource.Tenant == nil {
				if err = d.Set("tenant_id", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("tenant_id", resource.Tenant.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("vlan_id", resource.Vid); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamVlanUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	dropFields := []string{
		"created",
		"last_updated",
	}
	emptyFields := make(map[string]interface{})
	params := &models.WritableVLAN{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	vid := int64(d.Get("vlan_id").(int))
	params.Vid = &vid

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			params.Description = " "
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("vlan_group_id") {
		groupID := int64(d.Get("vlan_group_id").(int))
		if groupID != 0 {
			params.Group = &groupID
		}
	} else {
		dropFields = append(dropFields, "group")
	}

	if d.HasChange("role_id") {
		roleID := int64(d.Get("role_id").(int))
		if roleID != 0 {
			params.Role = &roleID
		}
	}

	if d.HasChange("site_id") {
		siteID := int64(d.Get("site_id").(int))
		if siteID != 0 {
			params.Site = &siteID
		}
	}

	if d.HasChange("status") {
		params.Status = d.Get("status").(string)
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		if tenantID != 0 {
			params.Tenant = &tenantID
		}
	}

	resource := ipam.NewIpamVlansPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamVlansPartialUpdate(resource, nil, newRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamVlanRead(ctx, d, m)
}

func resourceNetboxIpamVlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamVlanExists(d, m)
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

	resource := ipam.NewIpamVlansDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamVlansDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamVlanExists(d *schema.ResourceData, m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	vlanID := d.Id()
	params := ipam.NewIpamVlansListParams().WithID(&vlanID)
	vlans, err := client.Ipam.IpamVlansList(params, nil)

	if err != nil {
		return resourceExist, err
	}

	for _, vlan := range vlans.Payload.Results {
		if strconv.FormatInt(vlan.ID, 10) == d.Id() {
			resourceExist = true
		}
	}

	return resourceExist, nil
}
