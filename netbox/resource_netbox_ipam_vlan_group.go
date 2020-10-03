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
			"site_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
			},
		},
	}
}

func resourceNetboxIpamVlanGroupCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	groupName := d.Get("name").(string)
	groupSiteID := int64(d.Get("site_id").(int))
	groupSlug := d.Get("slug").(string)

	newResource := &models.WritableVLANGroup{
		Name: &groupName,
		Slug: &groupSlug,
	}

	if groupSiteID != 0 {
		newResource.Site = &groupSiteID
	}

	resource := ipam.NewIpamVlanGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamVlanGroupsCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))
	return resourceNetboxIpamVlanGroupRead(d, m)
}

func resourceNetboxIpamVlanGroupRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamVlanGroupsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVlanGroupsList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("name", resource.Name); err != nil {
				return err
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

			if err = d.Set("slug", resource.Slug); err != nil {
				return err
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamVlanGroupUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableVLANGroup{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	if d.HasChange("site_id") {
		siteID := int64(d.Get("site_id").(int))
		if siteID != 0 {
			params.Site = &siteID
		}
	}

	resource := ipam.NewIpamVlanGroupsPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return pkgerrors.New("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamVlanGroupsPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxIpamVlanGroupRead(d, m)
}

func resourceNetboxIpamVlanGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamVlanGroupExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return pkgerrors.New("Unable to convert ID into int64")
	}

	resource := ipam.NewIpamVlanGroupsDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamVlanGroupsDelete(resource, nil); err != nil {
		return err
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
