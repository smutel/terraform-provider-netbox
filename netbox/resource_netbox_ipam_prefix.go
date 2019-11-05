package netbox

import (
	//"log"
	//"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
	pkgerrors "github.com/pkg/errors"
)

func resourceNetboxIpamPrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamPrefixCreate,
		Read:   resourceNetboxIpamPrefixRead,
		Update: resourceNetboxIpamPrefixUpdate,
		Delete: resourceNetboxIpamPrefixDelete,
		Exists: resourceNetboxIpamPrefixExists,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 256),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"container", "active",
					"reserved", "deprecated"}, false),
			},
			"vrf_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"is_pool": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,
			},
			"site_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceNetboxIpamPrefixCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)

	prefixPrefix := d.Get("prefix").(string)
	prefixStatus := d.Get("status").(string)
	prefixVrfID := int64(d.Get("vrf_id").(int))
	prefixRoleID := int64(d.Get("role_id").(int))
	prefixDescription := d.Get("description").(string)
	prefixIsPool := d.Get("is_pool").(bool)
	prefixSiteID := int64(d.Get("site_id").(int))
	prefixVlanID := int64(d.Get("vlan_id").(int))
	prefixTenantID := int64(d.Get("tenant_id").(int))
	prefixTags := d.Get("tags").(*schema.Set).List()

	newPrefix := &models.WritablePrefix{
		Prefix:      &prefixPrefix,
		Status:      prefixStatus,
		Tags:        expandToStringSlice(prefixTags),
		Description: prefixDescription,
		IsPool:      prefixIsPool,
	}

	if prefixVrfID != 0 {
		newPrefix.Vrf = &prefixVrfID
	}

	if prefixRoleID != 0 {
		newPrefix.Role = &prefixRoleID
	}

	if prefixSiteID != 0 {
		newPrefix.Site = &prefixSiteID
	}

	if prefixVlanID != 0 {
		newPrefix.Vlan = &prefixVlanID
	}

	if prefixTenantID != 0 {
		newPrefix.Tenant = &prefixTenantID
	}

	p := ipam.NewIpamPrefixesCreateParams().WithData(newPrefix)

	prefixCreated, err := client.Ipam.IpamPrefixesCreate(p, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(prefixCreated.Payload.ID, 10))

	return resourceNetboxIpamPrefixRead(d, m)
}

func resourceNetboxIpamPrefixRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)

	prefixID := d.Id()
	params := ipam.NewIpamPrefixesListParams().WithIDIn(&prefixID)
	prefixs, err := client.Ipam.IpamPrefixesList(params, nil)
	if err != nil {
		return err
	}

	for _, prefix := range prefixs.Payload.Results {
		if strconv.FormatInt(prefix.ID, 10) == d.Id() {
			if err = d.Set("prefix", prefix.Prefix); err != nil {
				return err
			}

			if err = d.Set("status", *prefix.Status.Value); err != nil {
				return err
			}

			if err = d.Set("description", prefix.Description); err != nil {
				return err
			}

			if err = d.Set("tags", prefix.Tags); err != nil {
				return err
			}

			if err = d.Set("is_pool", prefix.IsPool); err != nil {
				return err
			}

			if prefix.Vrf == nil {
				if err = d.Set("vrf_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("vrf_id", prefix.Vrf.ID); err != nil {
					return err
				}
			}

			if prefix.Role == nil {
				if err = d.Set("role_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("role_id", prefix.Role.ID); err != nil {
					return err
				}
			}

			if prefix.Site == nil {
				if err = d.Set("site_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("site_id", prefix.Site.ID); err != nil {
					return err
				}
			}

			if prefix.Vlan == nil {
				if err = d.Set("vlan_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("vlan_id", prefix.Vlan.ID); err != nil {
					return err
				}
			}

			if prefix.Tenant == nil {
				if err = d.Set("tenant_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("tenant_id", prefix.Tenant.ID); err != nil {
					return err
				}
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamPrefixUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)
	updatedParams := &models.WritablePrefix{}

	prefixPrefix := d.Get("prefix").(string)
	updatedParams.Prefix = &prefixPrefix

	prefixTags := d.Get("tags").(*schema.Set).List()
	updatedParams.Tags = expandToStringSlice(prefixTags)

	updatedParams.IsPool = d.Get("is_pool").(bool)

	if d.HasChange("status") {
		updatedParams.Status = d.Get("status").(string)
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updatedParams.Description = description
	}

	if d.HasChange("vrf_id") {
		vrfID := int64(d.Get("vrf_id").(int))
		if vrfID != 0 {
			updatedParams.Vrf = &vrfID
		}
	}

	if d.HasChange("role_id") {
		prefixRoleID := int64(d.Get("role_id").(int))
		if prefixRoleID != 0 {
			updatedParams.Role = &prefixRoleID
		}
	}

	if d.HasChange("site_id") {
		siteID := int64(d.Get("site_id").(int))
		if siteID != 0 {
			updatedParams.Site = &siteID
		}
	}

	if d.HasChange("vlan_id") {
		vlanID := int64(d.Get("vlan_id").(int))
		if vlanID != 0 {
			updatedParams.Vlan = &vlanID
		}
	}

	if d.HasChange("tenant_id") {
		prefixTenantID := int64(d.Get("tenant_id").(int))
		if prefixTenantID != 0 {
			updatedParams.Tenant = &prefixTenantID
		}
	}

	p := ipam.NewIpamPrefixesPartialUpdateParams().WithData(updatedParams)

	tenantID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return pkgerrors.New("Unable to convert tenant ID into int64")
	}

	p.SetID(tenantID)

	_, err = client.Ipam.IpamPrefixesPartialUpdate(p, nil)
	if err != nil {
		return err
	}

	return resourceNetboxIpamPrefixRead(d, m)
}

func resourceNetboxIpamPrefixDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)

	resourceExists, err := resourceNetboxIpamPrefixExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	prefixID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return pkgerrors.New("Unable to convert prefix group ID into int64")
	}

	p := ipam.NewIpamPrefixesDeleteParams().WithID(prefixID)
	if _, err := client.Ipam.IpamPrefixesDelete(p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxIpamPrefixExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBox)
	prefixExist := false

	prefixID := d.Id()
	params := ipam.NewIpamPrefixesListParams().WithIDIn(&prefixID)
	prefixs, err := client.Ipam.IpamPrefixesList(params, nil)
	if err != nil {
		return prefixExist, err
	}

	for _, prefix := range prefixs.Payload.Results {
		if strconv.FormatInt(prefix.ID, 10) == d.Id() {
			prefixExist = true
		}
	}

	return prefixExist, nil
}
