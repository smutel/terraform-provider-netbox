package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxIpamVlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamVlanCreate,
		Read:   resourceNetboxIpamVlanRead,
		Update: resourceNetboxIpamVlanUpdate,
		Delete: resourceNetboxIpamVlanDelete,
		Exists: resourceNetboxIpamVlanExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"custom_field": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"text", "integer", "boolean",
								"date", "url", "selection", "multiple"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"vlan_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"role_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"active", "reserved",
					"deprecated"}, false),
			},
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceNetboxIpamVlanCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
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
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))
	return resourceNetboxIpamVlanRead(d, m)
}

func resourceNetboxIpamVlanRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamVlansListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVlansList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return err
			}

			var description interface{}
			if resource.Description == "" {
				description = nil
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return err
			}

			if resource.Group == nil {
				if err = d.Set("vlan_group_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("vlan_group_id", resource.Group.ID); err != nil {
					return err
				}
			}

			if err = d.Set("name", resource.Name); err != nil {
				return err
			}

			if resource.Role == nil {
				if err = d.Set("role_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("role_id", resource.Role.ID); err != nil {
					return err
				}
			}

			if resource.Site == nil {
				if err = d.Set("site_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("site_id", resource.Site.ID); err != nil {
					return err
				}
			}

			if resource.Status == nil {
				if err = d.Set("status", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("status", resource.Status.Value); err != nil {
					return err
				}
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return err
			}

			if resource.Tenant == nil {
				if err = d.Set("tenant_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("tenant_id", resource.Tenant.ID); err != nil {
					return err
				}
			}

			if err = d.Set("vlan_id", resource.Vid); err != nil {
				return err
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamVlanUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
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
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("vlan_group_id") {
		groupID := int64(d.Get("vlan_group_id").(int))
		if groupID != 0 {
			params.Group = &groupID
		}
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
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamVlansPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxIpamVlanRead(d, m)
}

func resourceNetboxIpamVlanDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamVlanExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource := ipam.NewIpamVlansDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamVlansDelete(resource, nil); err != nil {
		return err
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
