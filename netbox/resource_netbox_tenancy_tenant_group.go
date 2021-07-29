package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/tenancy"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxTenancyTenantGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxTenancyTenantGroupCreate,
		Read:   resourceNetboxTenancyTenantGroupRead,
		Update: resourceNetboxTenancyTenantGroupUpdate,
		Delete: resourceNetboxTenancyTenantGroupDelete,
		Exists: resourceNetboxTenancyTenantGroupExists,
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

func resourceNetboxTenancyTenantGroupCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	groupName := d.Get("name").(string)
	groupSlug := d.Get("slug").(string)

	newResource := &models.WritableTenantGroup{
		Name: &groupName,
		Slug: &groupSlug,
	}

	resource := tenancy.NewTenancyTenantGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyTenantGroupsCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyTenantGroupRead(d, m)
}

func resourceNetboxTenancyTenantGroupRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyTenantGroupsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyTenantGroupsList(params, nil)
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

func resourceNetboxTenancyTenantGroupUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableTenantGroup{}

	// Required parameters
	slug := d.Get("slug").(string)
	params.Slug = &slug

	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}

	resource := tenancy.NewTenancyTenantGroupsPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyTenantGroupsPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxTenancyTenantGroupRead(d, m)
}

func resourceNetboxTenancyTenantGroupDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyTenantGroupExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert tenant ID into int64")
	}

	resource := tenancy.NewTenancyTenantGroupsDeleteParams().WithID(resourceID)
	if _, err := client.Tenancy.TenancyTenantGroupsDelete(resource, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxTenancyTenantGroupExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := tenancy.NewTenancyTenantGroupsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyTenantGroupsList(params, nil)
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
