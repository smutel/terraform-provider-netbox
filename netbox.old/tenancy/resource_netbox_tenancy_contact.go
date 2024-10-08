package tenancy

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxTenancyContact() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a contact (tenancy module) within Netbox.",
		CreateContext: resourceNetboxTenancyContactCreate,
		ReadContext:   resourceNetboxTenancyContactRead,
		UpdateContext: resourceNetboxTenancyContactUpdate,
		DeleteContext: resourceNetboxTenancyContactDelete,
		Exists:        resourceNetboxTenancyContactExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The address for this contact (tenancy module).",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Comments for this contact (tenancy module).",
			},
			"contact_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the group where this contact (tenancy module) belongs to.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this contact (tenancy module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					email := val.(string)
					if len(email) >= 254 {
						errs = append(errs, fmt.Errorf("Length of %q must be lower than 254, got: %d", key, len(email)))
					} else if !strfmt.IsEmail(email) {
						errs = append(errs, fmt.Errorf("%q is not a valid email", key))
					}
					return
				},
				Description: "The e-mail for this contact (tenancy module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name for this contact (tenancy module).",
			},
			"phone": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The phone for this contact (tenancy module).",
			},
			"tag": &tag.TagSchema,
			"title": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The title for this contact (tenancy module).",
			},
		},
	}
}

func resourceNetboxTenancyContactCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	address := d.Get("address").(string)
	comments := d.Get("comments").(string)
	groupID := int32(d.Get("contact_group_id").(int))
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	email := d.Get("email").(string)
	name := d.Get("name").(string)
	phone := d.Get("phone").(string)
	tags := d.Get("tag").(*schema.Set).List()
	title := d.Get("title").(string)

	newResource := netbox.NewWritableContactRequestWithDefaults()
	newResource.SetAddress(address)
	newResource.SetComments(comments)
	newResource.SetCustomFields(customFields)
	newResource.SetEmail(email)
	newResource.SetName(name)
	newResource.SetPhone(phone)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	newResource.SetTitle(title)

	if groupID != 0 {
		newResource.SetGroup(groupID)
	}

	resourceCreated, response, err := client.TenancyAPI.TenancyContactsCreate(ctx).ContactRequest

	WritableContactRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxTenancyContactRead(ctx, d, m)
}

func resourceNetboxTenancyContactRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.TenancyAPI.TenancyContactsRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("address", resource.GetAddress()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("comments", resource.GetComments()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("contact_group_id", resource.GetGroup().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("email", resource.GetEmail()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("phone", resource.GetPhone()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("title", resource.GetTitle()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxTenancyContactUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableContactRequestWithDefaults()

	// Required parameters
	name := d.Get("name").(string)
	resource.SetName(name)

	if d.HasChange("address") {
		resource.SetAddress(d.Get("address").(string))
	}

	if d.HasChange("comments") {
		resource.SetComments(d.Get("comments").(string))
	}

	if d.HasChange("contact_group_id") {
		groupID := int32(d.Get("contact_group_id").(int))
		if groupID == 0 {
			resource.SetGroupNil()
		} else {
			resource.SetGroup(groupID)
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("email") {
		resource.SetEmail(d.Get("email").(string))
	}

	if d.HasChange("phone") {
		resource.SetEmail(d.Get("phone").(string))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("title") {
		resource.SetEmail(d.Get("title").(string))
	}

	if _, response, err := client.TenancyAPI.TenancyContactsUpdate(ctx, int32(resourceID)).ContactRequest

  WritableContactRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxTenancyContactRead(ctx, d, m)
}

func resourceNetboxTenancyContactDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxTenancyContactExists(d, m)
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

	if response, err := client.TenancyAPI.TenancyContactsDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxTenancyContactExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.TenancyAPI.TenancyContactsRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
