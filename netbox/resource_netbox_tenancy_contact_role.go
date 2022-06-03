package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/tenancy"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxTenancyContactRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxTenancyContactRoleCreate,
		Read:   resourceNetboxTenancyContactRoleRead,
		Update: resourceNetboxTenancyContactRoleUpdate,
		Delete: resourceNetboxTenancyContactRoleDelete,
		Exists: resourceNetboxTenancyContactRoleExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
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

func resourceNetboxTenancyContactRoleCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.ContactRole{
		CustomFields: &customFields,
		Description:  description,
		Name:         &name,
		Slug:         &slug,
		Tags:         convertTagsToNestedTags(tags),
	}

	resource := tenancy.NewTenancyContactRolesCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyContactRolesCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyContactRoleRead(d, m)
}

func resourceNetboxTenancyContactRoleRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyContactRolesListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactRolesList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return err
			}

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

func resourceNetboxTenancyContactRoleUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.ContactRole{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			params.Description = " "
		}
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	resource := tenancy.NewTenancyContactRolesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyContactRolesPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxTenancyContactRoleRead(d, m)
}

func resourceNetboxTenancyContactRoleDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyContactRoleExists(d, m)
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

	p := tenancy.NewTenancyContactRolesDeleteParams().WithID(id)
	if _, err := client.Tenancy.TenancyContactRolesDelete(p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxTenancyContactRoleExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := tenancy.NewTenancyContactRolesListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactRolesList(params, nil)
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
