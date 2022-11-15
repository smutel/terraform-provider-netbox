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

func resourceNetboxIpamPrefix() *schema.Resource {
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
				Description: "Existing tag to associate to this prefix (ipam module).",
			},
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
	client := m.(*netboxclient.NetBoxAPI)

	var prefix string
	var prefixid *int64
	if stateprefix, ok := d.GetOk("prefix"); ok {
		prefix = stateprefix.(string)
	} else if pprefix, ok := d.GetOk("parent_prefix"); ok {
		set := pprefix.(*schema.Set)
		mappreffix := set.List()[0].(map[string]interface{})
		parentPrefix := int64(mappreffix["prefix"].(int))
		prefixlength := int64(mappreffix["prefix_length"].(int))
		p, err := getNewAvailablePrefix(client, parentPrefix, prefixlength)
		if err != nil {
			return diag.FromErr(err)
		}
		prefix = *p.Prefix
		prefixid = &p.ID
	} else {
		return diag.Errorf("exactly one of (prefix, parent_prefix) must be specified")
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	isPool := d.Get("is_pool").(bool)
	roleID := int64(d.Get("role_id").(int))
	siteID := int64(d.Get("site_id").(int))
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))
	vlanID := int64(d.Get("vlan_id").(int))
	vrfID := int64(d.Get("vrf_id").(int))

	newResource := &models.WritablePrefix{
		CustomFields: &customFields,
		Description:  description,
		IsPool:       isPool,
		Prefix:       &prefix,
		Status:       status,
		Tags:         convertTagsToNestedTags(tags),
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

	if vlanID != 0 {
		newResource.Vlan = &vlanID
	}

	if vrfID != 0 {
		newResource.Vrf = &vrfID
	}

	if prefixid == nil {
		resource := ipam.NewIpamPrefixesCreateParams().WithData(newResource)

		resourceCreated, err := client.Ipam.IpamPrefixesCreate(resource, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		prefixid = &resourceCreated.Payload.ID
	} else {
		resource := ipam.NewIpamPrefixesUpdateParams().WithID(*prefixid).WithData(newResource)

		resourceCreated, err := client.Ipam.IpamPrefixesUpdate(resource, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		prefixid = &resourceCreated.Payload.ID
	}

	d.SetId(strconv.FormatInt(*prefixid, 10))

	return resourceNetboxIpamPrefixRead(ctx, d, m)
}

func resourceNetboxIpamPrefixRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamPrefixesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamPrefixesList(params, nil)
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

			if err = d.Set("is_pool", resource.IsPool); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("prefix", resource.Prefix); err != nil {
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

			if resource.Vlan == nil {
				if err = d.Set("vlan_id", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("vlan_id", resource.Vlan.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if resource.Vrf == nil {
				if err = d.Set("vrf_id", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("vrf_id", resource.Vrf.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamPrefixUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritablePrefix{}

	// Required parameters
	prefix := d.Get("prefix").(string)
	params.Prefix = &prefix

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			params.Description = " "
		}
	}

	params.IsPool = d.Get("is_pool").(bool)

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

	if d.HasChange("vlan_id") {
		vlanID := int64(d.Get("vlan_id").(int))
		if vlanID != 0 {
			params.Vlan = &vlanID
		}
	}

	if d.HasChange("vrf_id") {
		vrfID := int64(d.Get("vrf_id").(int))
		if vrfID != 0 {
			params.Vrf = &vrfID
		}
	}

	resource := ipam.NewIpamPrefixesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamPrefixesPartialUpdate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamPrefixRead(ctx, d, m)
}

func resourceNetboxIpamPrefixDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamPrefixExists(d, m)
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

	resource := ipam.NewIpamPrefixesDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamPrefixesDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamPrefixExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamPrefixesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamPrefixesList(params, nil)
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
