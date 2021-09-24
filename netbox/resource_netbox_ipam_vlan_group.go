package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxIpamVlanGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamVlanGroupCreate,
		Read:   resourceNetboxIpamVlanGroupRead,
		Update: resourceNetboxIpamVlanGroupUpdate,
		Delete: resourceNetboxIpamVlanGroupDelete,
		Exists: resourceNetboxIpamVlanGroupExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
		},
	}
}

func resourceNetboxIpamVlanGroupCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	groupName := d.Get("name").(string)
	groupSlug := d.Get("slug").(string)

	newResource := &models.VLANGroup{
		Name: &groupName,
		Slug: &groupSlug,
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
	params := &models.VLANGroup{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	resource := ipam.NewIpamVlanGroupsPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
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
		return fmt.Errorf("Unable to convert ID into int64")
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
