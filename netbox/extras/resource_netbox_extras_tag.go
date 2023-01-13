package extras

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/extras"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

func ResourceNetboxExtrasTag() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a tag (extra module) within Netbox.",
		CreateContext: resourceNetboxExtrasTagCreate,
		ReadContext:   resourceNetboxExtrasTagRead,
		UpdateContext: resourceNetboxExtrasTagUpdate,
		DeleteContext: resourceNetboxExtrasTagDelete,
		Exists:        resourceNetboxExtrasTagExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this tag (extra module).",
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "9e9e9e",
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 6),
					validation.StringMatch(
						regexp.MustCompile("^[0-9a-f]{1,6}$"),
						"^[0-9a-f]{1,6})$")),
				Description: "The color of this tag. Default is grey (#9e9e9e).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this tag was created.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this tag.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this tag was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this tag.",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this tag.",
			},
			"tagged_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of items tagged with this tag.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this tag.",
			},
		},
	}
}

var tagRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"slug",
}

func resourceNetboxExtrasTagCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)

	newResource := &models.Tag{
		Name:        &name,
		Slug:        &slug,
		Description: d.Get("description").(string),
		Color:       d.Get("color").(string),
	}

	resource := extras.NewExtrasTagsCreateParams().WithData(newResource)

	resourceCreated, err := client.Extras.ExtrasTagsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxExtrasTagRead(ctx, d, m)
}

func resourceNetboxExtrasTagRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := extras.NewExtrasTagsListParams().WithID(&resourceID)
	resources, err := client.Extras.ExtrasTagsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]

	if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("color", resource.Color); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("created", resource.Created.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("slug", resource.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tagged_items", resource.TaggedItems); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxExtrasTagUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modifiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	params := &models.Tag{}

	if d.HasChange("color") {
		params.Color = d.Get("color").(string)
	}
	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
		modifiedFields["description"] = params.Description
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("slug") {
		slug := d.Get("slug").(string)
		params.Slug = &slug
	}

	resource := extras.NewExtrasTagsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Extras.ExtrasTagsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, tagRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxExtrasTagRead(ctx, d, m)
}

func resourceNetboxExtrasTagDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxExtrasTagExists(d, m)
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

	resource := extras.NewExtrasTagsDeleteParams().WithID(id)
	if _, err := client.Extras.ExtrasTagsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxExtrasTagExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := extras.NewExtrasTagsListParams().WithID(&resourceID)
	resources, err := client.Extras.ExtrasTagsList(params, nil)
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
