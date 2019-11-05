package netbox

import (
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
	pkgerrors "github.com/pkg/errors"
)

func resourceNetboxIpamVlanGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamVlanGroupCreate,
		Read:   resourceNetboxIpamVlanGroupRead,
		Update: resourceNetboxIpamVlanGroupUpdate,
		Delete: resourceNetboxIpamVlanGroupDelete,
		Exists: resourceNetboxIpamVlanGroupExists,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
			},
			"site_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNetboxIpamVlanGroupCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)

	vlanGroupName := d.Get("name").(string)
	vlanGroupSlug := d.Get("slug").(string)
	vlanGroupSiteID := int64(d.Get("site_id").(int))

	newVlanGroup := &models.WritableVLANGroup{
		Name: &vlanGroupName,
		Slug: &vlanGroupSlug,
	}

	if vlanGroupSiteID != 0 {
		newVlanGroup.Site = &vlanGroupSiteID
	}

	p := ipam.NewIpamVlanGroupsCreateParams().WithData(newVlanGroup)

	vlanGroupCreated, err := client.Ipam.IpamVlanGroupsCreate(p, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(vlanGroupCreated.Payload.ID, 10))
	return resourceNetboxIpamVlanGroupRead(d, m)
}

func resourceNetboxIpamVlanGroupRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)

	vlanGroupID := d.Id()
	params := ipam.NewIpamVlanGroupsListParams().WithID(&vlanGroupID)
	vlanGroups, err := client.Ipam.IpamVlanGroupsList(params, nil)
	if err != nil {
		return err
	}

	for _, vlanGroup := range vlanGroups.Payload.Results {
		if strconv.FormatInt(vlanGroup.ID, 10) == d.Id() {
			if err = d.Set("name", vlanGroup.Name); err != nil {
				return err
			}

			if err = d.Set("slug", vlanGroup.Slug); err != nil {
				return err
			}

			if vlanGroup.Site == nil {
				if err = d.Set("site_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("site_id", vlanGroup.Site.ID); err != nil {
					return err
				}
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamVlanGroupUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBox)
	updatedParams := &models.WritableVLANGroup{}

	name := d.Get("name").(string)
	updatedParams.Name = &name

	slug := d.Get("slug").(string)
	updatedParams.Slug = &slug

	if d.HasChange("site_id") {
		siteID := int64(d.Get("site_id").(int))
		if siteID != 0 {
			updatedParams.Site = &siteID
		}
	}

	p := ipam.NewIpamVlanGroupsPartialUpdateParams().WithData(
		updatedParams)

	tenantID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return pkgerrors.New("Unable to convert tenant ID into int64")
	}

	p.SetID(tenantID)

	_, err = client.Ipam.IpamVlanGroupsPartialUpdate(p, nil)
	if err != nil {
		return err
	}

	return resourceNetboxIpamVlanGroupRead(d, m)
}

func resourceNetboxIpamVlanGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBox)

	resourceExists, err := resourceNetboxIpamVlanGroupExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	vlanGroupID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return pkgerrors.New("Unable to convert vlan group ID into int64")
	}

	p := ipam.NewIpamVlanGroupsDeleteParams().WithID(vlanGroupID)
	if _, err := client.Ipam.IpamVlanGroupsDelete(p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxIpamVlanGroupExists(d *schema.ResourceData, m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBox)
	vlanGroupExist := false

	vlanGroupID := d.Id()
	params := ipam.NewIpamVlanGroupsListParams().WithID(&vlanGroupID)
	vlanGroups, err := client.Ipam.IpamVlanGroupsList(params, nil)
	if err != nil {
		return vlanGroupExist, err
	}

	for _, vlanGroup := range vlanGroups.Payload.Results {
		if strconv.FormatInt(vlanGroup.ID, 10) == d.Id() {
			vlanGroupExist = true
		}
	}

	return vlanGroupExist, nil
}
