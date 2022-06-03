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

func resourceNetboxTenancyContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxTenancyContactGroupCreate,
		Read:   resourceNetboxTenancyContactGroupRead,
		Update: resourceNetboxTenancyContactGroupUpdate,
		Delete: resourceNetboxTenancyContactGroupDelete,
		Exists: resourceNetboxTenancyContactGroupExists,
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
			"parent_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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

func resourceNetboxTenancyContactGroupCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	parentID := int64(d.Get("parent_id").(int))
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableContactGroup{
		CustomFields: &customFields,
		Description:  description,
		Name:         &name,
		Slug:         &slug,
		Tags:         convertTagsToNestedTags(tags),
	}

	if parentID != 0 {
		newResource.Parent = &parentID
	}

	resource := tenancy.NewTenancyContactGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyContactGroupsCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyContactGroupRead(d, m)
}

func resourceNetboxTenancyContactGroupRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyContactGroupsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactGroupsList(params, nil)
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

			if resource.Parent == nil {
				if err = d.Set("parent_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("parent_id", resource.Parent.ID); err != nil {
					return err
				}
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

func resourceNetboxTenancyContactGroupUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableContactGroup{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	parentID := int64(d.Get("parent_id").(int))
	params.Parent = &parentID

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

	resource := tenancy.NewTenancyContactGroupsPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyContactGroupsPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxTenancyContactGroupRead(d, m)
}

func resourceNetboxTenancyContactGroupDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyContactGroupExists(d, m)
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

	p := tenancy.NewTenancyContactGroupsDeleteParams().WithID(id)
	if _, err := client.Tenancy.TenancyContactGroupsDelete(p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxTenancyContactGroupExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := tenancy.NewTenancyContactGroupsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactGroupsList(params, nil)
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
