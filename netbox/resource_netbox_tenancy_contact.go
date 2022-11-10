package netbox

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/tenancy"
	"github.com/smutel/go-netbox/v3/netbox/models"
)

func resourceNetboxTenancyContact() *schema.Resource {
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
			"custom_field": &customFieldSchema,
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
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the existing tag.",
						},
						"slug": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Slug of the existing tag.",
						},
					},
				},
				Description: "Existing tag to associate to this contact (tenancy module).",
			},
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
	client := m.(*netboxclient.NetBoxAPI)

	address := d.Get("address").(string)
	comments := d.Get("comments").(string)
	groupID := int64(d.Get("contact_group_id").(int))
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	email := strfmt.Email(d.Get("email").(string))
	name := d.Get("name").(string)
	phone := d.Get("phone").(string)
	tags := d.Get("tag").(*schema.Set).List()
	title := d.Get("title").(string)

	newResource := &models.WritableContact{
		Address:      address,
		Comments:     comments,
		CustomFields: &customFields,
		Email:        email,
		Name:         &name,
		Phone:        phone,
		Tags:         convertTagsToNestedTags(tags),
		Title:        title,
	}

	if groupID != 0 {
		newResource.Group = &groupID
	}

	resource := tenancy.NewTenancyContactsCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyContactsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyContactRead(ctx, d, m)
}

func resourceNetboxTenancyContactRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyContactsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			var address interface{}
			if resource.Address == "" {
				address = nil
			} else {
				address = resource.Address
			}

			if err = d.Set("address", address); err != nil {
				return diag.FromErr(err)
			}

			var comments interface{}
			if resource.Comments == "" {
				comments = nil
			} else {
				comments = resource.Comments
			}

			if err = d.Set("comments", comments); err != nil {
				return diag.FromErr(err)
			}

			if resource.Group == nil {
				if err = d.Set("contact_group_id", 0); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("contact_group_id", resource.Group.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return diag.FromErr(err)
			}

			var email interface{}
			if resource.Email.String() == "" {
				email = nil
			} else {
				email = resource.Email.String()
			}

			if err = d.Set("email", email); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("name", resource.Name); err != nil {
				return diag.FromErr(err)
			}

			var phone interface{}
			if resource.Phone == "" {
				phone = nil
			} else {
				phone = resource.Phone
			}

			if err = d.Set("phone", phone); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			var title interface{}
			if resource.Title == "" {
				title = nil
			} else {
				title = resource.Title
			}

			if err = d.Set("title", title); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxTenancyContactUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	dropFields := []string{
		"created",
		"last_updated",
	}
	emptyFields := make(map[string]interface{})
	params := &models.WritableContact{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	if d.HasChange("address") {
		if address, exist := d.GetOk("address"); exist {
			params.Address = address.(string)
		} else {
			params.Address = " "
		}
	}

	if d.HasChange("comments") {
		if comments, exist := d.GetOk("comments"); exist {
			params.Comments = comments.(string)
		} else {
			params.Comments = " "
		}
	}

	if d.HasChange("contact_group_id") {
		groupID := int64(d.Get("contact_group_id").(int))
		if groupID == 0 {
			params.Group = nil
		} else {
			params.Group = &groupID
		}
	} else {
		dropFields = append(dropFields, "group")
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("email") {
		if email, exist := d.GetOk("email"); exist {
			params.Email = strfmt.Email(email.(string))
		} else {
			params.Email = strfmt.Email(" ")
		}
	}

	if d.HasChange("phone") {
		if phone, exist := d.GetOk("phone"); exist {
			params.Phone = phone.(string)
		} else {
			params.Phone = " "
		}
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	if d.HasChange("title") {
		if title, exist := d.GetOk("title"); exist {
			params.Title = title.(string)
		} else {
			params.Title = " "
		}
	}

	resource := tenancy.NewTenancyContactsPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyContactsPartialUpdate(resource, nil, newRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxTenancyContactRead(ctx, d, m)
}

func resourceNetboxTenancyContactDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyContactExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	p := tenancy.NewTenancyContactsDeleteParams().WithID(id)
	if _, err := client.Tenancy.TenancyContactsDelete(p, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxTenancyContactExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := tenancy.NewTenancyContactsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactsList(params, nil)
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
