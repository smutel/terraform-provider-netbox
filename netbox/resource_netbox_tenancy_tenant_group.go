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
)

func resourceNetboxTenancyTenantGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a tenant group (tenancy module) within Netbox.",
		CreateContext: resourceNetboxTenancyTenantGroupCreate,
		ReadContext:   resourceNetboxTenancyTenantGroupRead,
		UpdateContext: resourceNetboxTenancyTenantGroupUpdate,
		DeleteContext: resourceNetboxTenancyTenantGroupDelete,
		Exists:        resourceNetboxTenancyTenantGroupExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this tenant group (tenancy module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
				Description:  "The name for this tenant group (tenancy module).",
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_]{1,50}$"),
					"Must be like ^[-a-zA-Z0-9_]{1,50}$"),
				Description: "The slug for this tenant group (tenancy module).",
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
				Description: "Existing tag to associate to this tenant group (tenancy module).",
			},
		},
	}
}

func resourceNetboxTenancyTenantGroupCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	groupName := d.Get("name").(string)
	groupSlug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableTenantGroup{
		Name: &groupName,
		Slug: &groupSlug,
		Tags: convertTagsToNestedTags(tags),
	}

	resource := tenancy.NewTenancyTenantGroupsCreateParams().WithData(newResource)

	resourceCreated, err := client.Tenancy.TenancyTenantGroupsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxTenancyTenantGroupRead(ctx, d, m)
}

func resourceNetboxTenancyTenantGroupRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := tenancy.NewTenancyTenantGroupsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyTenantGroupsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
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

func resourceNetboxTenancyTenantGroupUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableTenantGroup{}

	// Required parameters
	slug := d.Get("slug").(string)
	params.Slug = &slug

	name := d.Get("name").(string)
	params.Name = &name

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	resource := tenancy.NewTenancyTenantGroupsPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Tenancy.TenancyTenantGroupsPartialUpdate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxTenancyTenantGroupRead(ctx, d, m)
}

func resourceNetboxTenancyTenantGroupDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxTenancyTenantGroupExists(d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert tenant ID into int64")
	}

	resource := tenancy.NewTenancyTenantGroupsDeleteParams().WithID(resourceID)
	if _, err := client.Tenancy.TenancyTenantGroupsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxTenancyTenantGroupExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := tenancy.NewTenancyTenantGroupsListParams().WithID(&resourceID)
	resources, err := client.Tenancy.TenancyTenantGroupsList(params, nil)
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
