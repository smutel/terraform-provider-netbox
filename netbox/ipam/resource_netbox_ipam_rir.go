package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

func ResourceNetboxIpamRIR() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a rir (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamRIRCreate,
		ReadContext:   resourceNetboxIpamRIRRead,
		UpdateContext: resourceNetboxIpamRIRUpdate,
		DeleteContext: resourceNetboxIpamRIRDelete,
		// Exists:        resourceNetboxIpamRIRExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"aggregate_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of aggregates with this rir (ipam module).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this rir (ipam module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this rir was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this rir (ipam module).",
			},
			"is_private": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Date when this rir was created.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this rir was created.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this rir (ipam module).",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this rir (ipam module).",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this rir (ipam module).",
			},
		},
	}
}

var rirRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"slug",
	"tag",
}

func resourceNetboxIpamRIRCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.RIR{
		CustomFields: &customFields,
		Description:  d.Get("description").(string),
		IsPrivate:    d.Get("is_private").(bool),
		Name:         &name,
		Slug:         &slug,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	resource := ipam.NewIpamRirsCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamRirsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxIpamRIRRead(ctx, d, m)
}

func resourceNetboxIpamRIRRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamRirsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamRirsList(params, nil)
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
	if err = d.Set("aggregate_count", resource.AggregateCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("created", resource.Created.String()); err != nil {
		return diag.FromErr(err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)
	if err = d.Set("custom_field", customFields); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("is_private", resource.IsPrivate); err != nil {
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
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamRIRUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modiefiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.RIR{}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
		modiefiedFields["description"] = params.Description
	}
	if d.HasChange("is_private") {
		params.IsPrivate = d.Get("is_private").(bool)
		modiefiedFields["is_private"] = params.IsPrivate
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("slug") {
		slug := d.Get("slug").(string)
		params.Slug = &slug
	}
	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}

	resource := ipam.NewIpamRirsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamRirsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modiefiedFields, rirRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamRIRRead(ctx, d, m)
}

func resourceNetboxIpamRIRDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamRIRExists(d, m)
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

	resource := ipam.NewIpamRirsDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamRirsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamRIRExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamRirsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamRirsList(params, nil)
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
