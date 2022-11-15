package tenancy

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
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
)

func ResourceNetboxTenancyContactGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a contact group (tenancy module) within Netbox.",
		CreateContext: resourceNetboxTenancyContactGroupCreate,
		ReadContext:   resourceNetboxTenancyContactGroupRead,
		UpdateContext: resourceNetboxTenancyContactGroupUpdate,
		DeleteContext: resourceNetboxTenancyContactGroupDelete,
		Exists:        resourceNetboxTenancyContactGroupExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this contact group (tenancy module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "Description for this contact group (tenancy module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name for this contact group (tenancy module).",
			},
			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the contact group parent of this one.",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
				Description: "The slug for this contact group (tenancy module).",
			},
			"tag": &tag.TagSchema,
		},
	}
}

func resourceNetboxTenancyContactGroupCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
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
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	if parentID != 0 {
		newResource.Parent = &parentID
	}

	resource := tenancy.NewTenancyContactGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyContactGroupsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyContactGroupRead(ctx, d, m)
}

func resourceNetboxTenancyContactGroupRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyContactGroupsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyContactGroupsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
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

			if err = d.Set("name", resource.Name); err != nil {
				return diag.FromErr(err)
			}

			if resource.Parent == nil {
				if err = d.Set("parent_id", 0); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("parent_id", resource.Parent.ID); err != nil {
					return diag.FromErr(err)
				}
			}

			if err = d.Set("slug", resource.Slug); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxTenancyContactGroupUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableContactGroup{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name

	slug := d.Get("slug").(string)
	params.Slug = &slug

	parentID := int64(d.Get("parent_id").(int))
	if parentID != 0 {
		params.Parent = &parentID
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

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = tag.ConvertTagsToNestedTags(tags)

	resource := tenancy.NewTenancyContactGroupsPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyContactGroupsPartialUpdate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxTenancyContactGroupRead(ctx, d, m)
}

func resourceNetboxTenancyContactGroupDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyContactGroupExists(d, m)
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

	p := tenancy.NewTenancyContactGroupsDeleteParams().WithID(id)
	if _, err := client.Tenancy.TenancyContactGroupsDelete(p, nil); err != nil {
		return diag.FromErr(err)
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
