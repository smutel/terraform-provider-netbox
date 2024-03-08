package extras

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
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

func resourceNetboxExtrasTagCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	newResource := netbox.NewTagRequestWithDefaults()
	newResource.SetName(d.Get("name").(string))
	newResource.SetSlug(d.Get("slug").(string))
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetColor(d.Get("color").(string))

	resourceCreated, response, err := client.ExtrasAPI.ExtrasTagsCreate(ctx).TagRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resourceCreated.GetId() == 0 {
		return diag.FromErr(errors.New("Bug Netbox - TO BE FIXED"))
	}

	d.SetId(strconv.FormatInt(int64(resourceCreated.GetId()), 10))

	return resourceNetboxExtrasTagRead(ctx, d, m)
}

func resourceNetboxExtrasTagRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.ExtrasAPI.ExtrasTagsRetrieve(ctx, int32(resourceID)).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("color", resource.GetColor()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tagged_items", resource.GetObjectTypes()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxExtrasTagUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource := netbox.NewTagRequestWithDefaults()

	if d.HasChange("color") {
		resource.SetColor(d.Get("color").(string))
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("name") {
		resource.SetName(d.Get("name").(string))
	}

	if d.HasChange("slug") {
		resource.SetName(d.Get("slug").(string))
	}

	if _, response, err := client.ExtrasAPI.ExtrasTagsUpdate(ctx, int32(resourceID)).TagRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxExtrasTagRead(ctx, d, m)
}

func resourceNetboxExtrasTagDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxExtrasTagExists(d, m)
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

	if response, err := client.ExtrasAPI.ExtrasTagsDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxExtrasTagExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.ExtrasAPI.ExtrasTagsRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
