package tenancy

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxTenancyContactAssignment() *schema.Resource {
	return &schema.Resource{
		Description:   "Link a contact to another resource within Netbox.",
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
				Type:     schema.TypeInt,
				Required: true,
				Description: "ID of the contact to link to " +
					"this contact assignment.",
			},
			"contact_role_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "The role of the contact for " +
					"this contact assignment.",
			},
			"content_type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Type of the object where the contact " +
					"will be linked.",
			},
			"object_id": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "ID of the object where the contact " +
					"will be linked.",
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "primary",
				ValidateFunc: validation.StringInSlice([]string{"primary",
					"secondary", "tertiary", "inactive"}, false),
				Description: "Priority of this contact among primary, " +
					"secondary and tertiary (primary by default).",
			},
		},
	}
}

func resourceNetboxTenancyContactAssignmentCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	contactID := d.Get("contact_id").(int)
	contactRoleID := d.Get("contact_role_id").(int)
	contentType := d.Get("content_type").(string)
	objectID := int64(d.Get("object_id").(int))
	priority := d.Get("priority").(string)

	newResource := netbox.NewWritableContactAssignmentRequestWithDefaults()
	newResource.SetObjectId(objectID)
	newResource.SetObjectType(contentType)

	if contactID != 0 {
		b, err := brief.GetBriefContactRequestFromID(ctx, client, contactID)
		if err != nil {
			return err
		}
		newResource.SetContact(*b)
	}

	if contactRoleID != 0 {
		b, err := brief.GetBriefContactRoleRequestFromID(ctx, client,
			contactRoleID)
		if err != nil {
			return err
		}
		newResource.SetRole(*b)
	}

	p, err := netbox.NewContactAssignmentPriorityValueFromValue(priority)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetPriority(*p)

	resourceCreated, response, err :=
		client.TenancyAPI.TenancyContactAssignmentsCreate(
			ctx).WritableContactAssignmentRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxTenancyContactAssignmentRead(ctx, d, m)
}

func resourceNetboxTenancyContactAssignmentRead(ctx context.Context,
	d *schema.ResourceData,
	m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.TenancyAPI.TenancyContactAssignmentsRetrieve(ctx,
			int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("contact_id", resource.GetContact().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", resource.GetObjectType()); err != nil {
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

func resourceNetboxTenancyContactAssignmentUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableContactAssignmentRequestWithDefaults()

	// Required parameters
	resource.SetObjectType(d.Get("content_type").(string))
	b, errDiag := brief.GetBriefContactRequestFromID(ctx, client,
		d.Get("contact_id").(int))
	if errDiag != nil {
		return errDiag
	}
	resource.SetContact(*b)

	if d.HasChange("contact_role_id") {
		contactRoleID := d.Get("contact_role_id").(int)
		if contactRoleID != 0 {
			b, errDiag := brief.GetBriefContactRoleRequestFromID(ctx, client,
				contactRoleID)
			if errDiag != nil {
				return errDiag
			}
			resource.SetRole(*b)
		}

	}

	if d.HasChange("object_id") {
		resource.SetObjectId(int64(d.Get("object_id").(int)))
	}

	if d.HasChange("priority") {
		p, err := netbox.NewContactAssignmentPriorityValueFromValue(
			d.Get("priority").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetPriority(*p)
	}

	if _, response, err :=
		client.TenancyAPI.TenancyContactAssignmentsUpdate(ctx,
			int32(resourceID)).WritableContactAssignmentRequest(
			*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxTenancyContactAssignmentRead(ctx, d, m)
}

func resourceNetboxTenancyContactAssignmentDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxTenancyContactAssignmentExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}

	if response, err :=
		client.TenancyAPI.TenancyContactAssignmentsDestroy(ctx,
			int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxTenancyContactAssignmentExists(d *schema.ResourceData,
	m any) (b bool,
	e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.TenancyAPI.TenancyContactAssignmentsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
