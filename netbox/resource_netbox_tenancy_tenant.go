package netbox

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/tenancy"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/customfield"
)

func resourceNetboxTenancyTenant() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a tenant (tenancy module) within Netbox.",
		CreateContext: resourceNetboxTenancyTenantCreate,
		ReadContext:   resourceNetboxTenancyTenantRead,
		UpdateContext: resourceNetboxTenancyTenantUpdate,
		DeleteContext: resourceNetboxTenancyTenantDelete,
		Exists:        resourceNetboxTenancyTenantExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Comments for this tenant (tenancy module).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this tenant (tenancy module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The description for this tenant (tenancy module).",
			},
			"tenant_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the group where this tenant (tenancy module) is attached to.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 30),
				Description:  "The name for this tenant (tenancy module).",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
				Description: "The slug for this tenant (tenancy module).",
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
				Description: "Existing tag to associate to this tenant (tenancy module).",
			},
		},
	}
}

func resourceNetboxTenancyTenantCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	comments := d.Get("comments").(string)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
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
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyTenantRead(ctx, d, m)
}

func resourceNetboxTenancyTenantRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyTenantsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyTenantsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {

			var comments interface{}
			if resource.Comments == "" {
				comments = nil
			} else {
				comments = resource.Comments
			}

			if err = d.Set("comments", comments); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return diag.FromErr(err)
			}

			var description interface{}
			if resource.Description == "" {
				description = nil
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return diag.FromErr(err)
			}

			if resource.Group == nil {
				if err = d.Set("tenant_group_id", 0); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("tenant_group_id", resource.Group.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("name", resource.Name); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("slug", resource.Slug); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxTenancyTenantUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableTenant{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	if d.HasChange("comments") {
		if comments, exist := d.GetOk("comments"); exist {
			params.Comments = comments.(string)
		} else {
			params.Comments = " "
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			params.Description = " "
		}
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
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyTenantsPartialUpdate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxTenancyTenantRead(ctx, d, m)
}

func resourceNetboxTenancyTenantDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyTenantExists(d, m)
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

	p := tenancy.NewTenancyTenantsDeleteParams().WithID(id)
	if _, err := client.Tenancy.TenancyTenantsDelete(p, nil); err != nil {
		return diag.FromErr(err)
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
