package tenancy

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxTenancyContactAssignment() *schema.Resource {
	return &schema.Resource{
		Description:   "Link a contact (tenancy module) to another resource within Netbox.",
		CreateContext: resourceNetboxTenancyContactAssignmentCreate,
		ReadContext:   resourceNetboxTenancyContactAssignmentRead,
		UpdateContext: resourceNetboxTenancyContactAssignmentUpdate,
		DeleteContext: resourceNetboxTenancyContactAssignmentDelete,
		Exists:        resourceNetboxTenancyContactAssignmentExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"contact_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the contact to link to this contact assignment (tenancy module).",
			},
			"contact_role_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The role of the contact for this contact assignment (tenancy module).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the object where the contact will be linked.",
			},
			"object_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the object where the contact will be linked.",
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "primary",
				ValidateFunc: validation.StringInSlice([]string{"primary", "secondary", "tertiary",
					"inactive"}, false),
				Description: "Priority of this contact among primary, secondary and tertiary (primary by default).",
			},
		},
	}
}

func resourceNetboxTenancyContactAssignmentCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	contactID := int32(d.Get("contact_id").(int))
	contactRoleID := int32(d.Get("contact_role_id").(int))
	contentType := d.Get("content_type").(string)
	objectID := int64(d.Get("object_id").(int))
	priority := d.Get("priority").(string)

	newResource := netbox.NewWritableContactAssignmentRequestWithDefaults()
	newResource.SetContact(contactID)
	newResource.SetContentType(contentType)
	newResource.SetObjectId(objectID)
	newResource.SetRole(contactRoleID)

	p, err := netbox.NewContactAssignmentPriorityValueFromValue(priority)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetPriority(*p)

	resourceCreated, response, err := client.TenancyAPI.TenancyContactAssignmentsCreate(ctx).WritableContactAssignmentRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxTenancyContactAssignmentRead(ctx, d, m)
}

func resourceNetboxTenancyContactAssignmentRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.TenancyAPI.TenancyContactAssignmentsRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("contact_id", resource.GetContact().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("object_id", resource.GetObjectId()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("priority", resource.GetPriority().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("contact_role_id", resource.GetRole().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxTenancyContactAssignmentUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableContactAssignmentRequestWithDefaults()

	// Required parameters
	resource.SetContact(int32(d.Get("contact_id").(int)))

	if d.HasChange("contact_role_id") {
		resource.SetRole(int32(d.Get("contact_role_id").(int)))
	}

	if d.HasChange("content_type") {
		resource.SetContentType(d.Get("content_type").(string))
	}

	if d.HasChange("object_id") {
		resource.SetObjectId(int64(d.Get("object_id").(int)))
	}

	if d.HasChange("priority") {
		p, err := netbox.NewContactAssignmentPriorityValueFromValue(d.Get("priority").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetPriority(*p)
	}

	if _, response, err := client.TenancyAPI.TenancyContactAssignmentsUpdate(ctx, int32(resourceID)).WritableContactAssignmentRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxTenancyContactAssignmentRead(ctx, d, m)
}

func resourceNetboxTenancyContactAssignmentDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxTenancyContactAssignmentExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}

	if response, err := client.TenancyAPI.TenancyContactAssignmentsDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxTenancyContactAssignmentExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.TenancyAPI.TenancyContactAssignmentsRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
