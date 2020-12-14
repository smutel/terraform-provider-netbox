package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/tenancy"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxTenancyTenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxTenancyTenantCreate,
		Read:   resourceNetboxTenancyTenantRead,
		Update: resourceNetboxTenancyTenantUpdate,
		Delete: resourceNetboxTenancyTenantDelete,
		Exists: resourceNetboxTenancyTenantExists,

		Schema: map[string]*schema.Schema{
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  " ",
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// terraform default behavior sees a difference between null and an empty string
				// therefore we override the default, because null in terraform results in empty string in netbox
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// function is called for each member of map
					// including additional call on the amount of entries
					// we ignore the count, because the actual state always returns the amount of existing custom_fields and all are optional in terraform
					if k == CustomFieldsRegex {
						return true
					}
					return old == new
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      " ",
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"tenant_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 30),
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
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
		},
	}
}

func resourceNetboxTenancyTenantCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	comments := d.Get("comments").(string)
	resourceCustomFields := d.Get("custom_fields").(map[string]interface{})
	customFields := convertCustomFieldsFromTerraformToAPICreate(resourceCustomFields)
	description := d.Get("description").(string)
	groupID := int64(d.Get("tenant_group_id").(int))
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableTenant{
		Comments:     comments,
		CustomFields: &customFields,
		Description:  description,
		Name:         &name,
		Slug:         &slug,
		Tags:         convertTagsToNestedTags(tags),
	}

	if groupID != 0 {
		newResource.Group = &groupID
	}

	resource := tenancy.NewTenancyTenantsCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyTenantsCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyTenantRead(d, m)
}

func resourceNetboxTenancyTenantRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyTenantsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyTenantsList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			var comments string

			if resource.Comments == "" {
				comments = " "
			} else {
				comments = resource.Comments
			}

			if err = d.Set("comments", comments); err != nil {
				return err
			}

			customFields := convertCustomFieldsFromAPIToTerraform(resource.CustomFields)

			if err = d.Set("custom_fields", customFields); err != nil {
				return err
			}

			var description string

			if resource.Description == "" {
				description = " "
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return err
			}

			if resource.Group == nil {
				if err = d.Set("tenant_group_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("tenant_group_id", resource.Group.ID); err != nil {
					return err
				}
			}

			if err = d.Set("name", resource.Name); err != nil {
				return err
			}

			if err = d.Set("slug", resource.Slug); err != nil {
				return err
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return err
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxTenancyTenantUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableTenant{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	if d.HasChange("comments") {
		comments := d.Get("comments").(string)
		params.Comments = comments
	}

	if d.HasChange("custom_fields") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_fields")
		customFields := convertCustomFieldsFromTerraformToAPIUpdate(stateCustomFields, resourceCustomFields)
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		params.Description = description
	}

	if d.HasChange("tenant_group_id") {
		groupID := int64(d.Get("tenant_group_id").(int))
		params.Group = &groupID
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	resource := tenancy.NewTenancyTenantsPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyTenantsPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxTenancyTenantRead(d, m)
}

func resourceNetboxTenancyTenantDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyTenantExists(d, m)
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

	p := tenancy.NewTenancyTenantsDeleteParams().WithID(id)
	if _, err := client.Tenancy.TenancyTenantsDelete(p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxTenancyTenantExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := tenancy.NewTenancyTenantsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyTenantsList(params, nil)
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
