package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/tenancy"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxTenancyContactAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxTenancyContactAssignmentCreate,
		Read:   resourceNetboxTenancyContactAssignmentRead,
		Update: resourceNetboxTenancyContactAssignmentUpdate,
		Delete: resourceNetboxTenancyContactAssignmentDelete,
		Exists: resourceNetboxTenancyContactAssignmentExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"contact_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"contact_role_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"object_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "primary",
				ValidateFunc: validation.StringInSlice([]string{"primary", "secondary", "tertiary",
					"inactive"}, false),
			},
		},
	}
}

func resourceNetboxTenancyContactAssignmentCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	fmt.Println(client)

	contactID := int64(d.Get("contact_id").(int))
	contactRoleID := int64(d.Get("contact_role_id").(int))
	contentType := d.Get("content_type").(string)
	objectID := int64(d.Get("object_id").(int))
	priority := d.Get("priority").(string)

	newResource := &models.WritableContactAssignment{
		Contact:     &contactID,
		ContentType: &contentType,
		ObjectID:    &objectID,
		Priority:    &priority,
		Role:        &contactRoleID,
	}

	resource := tenancy.NewTenancyContactAssignmentsCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyContactAssignmentsCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyContactAssignmentRead(d, m)
}

func resourceNetboxTenancyContactAssignmentRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyContactAssignmentsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactAssignmentsList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("contact_id", resource.Contact.ID); err != nil {
				return err
			}

			if err = d.Set("content_type", resource.ContentType); err != nil {
				return err
			}

			if err = d.Set("object_id", resource.ObjectID); err != nil {
				return err
			}

			if resource.Priority == nil {
				if err = d.Set("priority", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("priority", resource.Priority.Value); err != nil {
					return err
				}
			}

			if err = d.Set("contact_role_id", resource.Role.ID); err != nil {
				return err
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxTenancyContactAssignmentUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableContactAssignment{}

	// Required parameters
	contactID := int64(d.Get("contact_id").(int))
	params.Contact = &contactID

	contactRoleID := int64(d.Get("contact_role_id").(int))
	params.Role = &contactRoleID

	contentType := d.Get("content_type").(string)
	params.ContentType = &contentType

	objectID := int64(d.Get("object_id").(int))
	params.ObjectID = &objectID

	priority := d.Get("priority").(string)
	params.Priority = &priority

	resource := tenancy.NewTenancyContactAssignmentsPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyContactAssignmentsPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxTenancyContactAssignmentRead(d, m)
}

func resourceNetboxTenancyContactAssignmentDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyContactAssignmentExists(d, m)
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

	p := tenancy.NewTenancyContactAssignmentsDeleteParams().WithID(id)
	if _, err := client.Tenancy.TenancyContactAssignmentsDelete(p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxTenancyContactAssignmentExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := tenancy.NewTenancyContactAssignmentsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactAssignmentsList(params, nil)
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
