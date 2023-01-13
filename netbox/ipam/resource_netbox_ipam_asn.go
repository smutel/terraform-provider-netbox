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

func ResourceNetboxIpamASN() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a asn (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamASNCreate,
		ReadContext:   resourceNetboxIpamASNRead,
		UpdateContext: resourceNetboxIpamASNUpdate,
		DeleteContext: resourceNetboxIpamASNDelete,
		// Exists:        resourceNetboxIpamASNExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"asn": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The asn number of this asn (ipam module).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this asn (ipam module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this asn was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this asn (ipam module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this rir was created.",
			},
			"provider_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of providers for this asn (ipam module).",
			},
			"rir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The rir for this asn (ipam module).",
			},
			"site_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of sites for this asn (ipam module).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The rir for this asn (ipam module).",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this asn (ipam module).",
			},
		},
	}
}

var asnRequiredFields = []string{
	"created",
	"last_updated",
	"asn",
	"rir",
	"tag",
}

func resourceNetboxIpamASNCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)

	asn := int64(d.Get("asn").(int))
	rirID := int64(d.Get("rir_id").(int))
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableASN{
		Asn:          &asn,
		CustomFields: &customFields,
		Description:  d.Get("description").(string),
		Rir:          &rirID,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	if tenantID := int64(d.Get("tenant_id").(int)); tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	resource := ipam.NewIpamAsnsCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamAsnsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxIpamASNRead(ctx, d, m)
}

func resourceNetboxIpamASNRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamAsnsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamAsnsList(params, nil)
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
	if err = d.Set("asn", resource.Asn); err != nil {
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
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("provider_count", resource.ProviderCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("rir_id", resource.Rir); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("site_count", resource.SiteCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tenant_id", util.GetNestedTenantID(resource.Tenant)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamASNUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modiefiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableASN{}

	if d.HasChange("asn") {
		asnID := int64(d.Get("asn").(int))
		params.Asn = &asnID
	}
	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
		modiefiedFields["description"] = params.Description
	}
	if d.HasChange("rir_id") {
		rirID := int64(d.Get("rir_id").(int))
		params.Rir = &rirID
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}

	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		params.Tenant = &tenantID
		modiefiedFields["tenant"] = tenantID
	}

	resource := ipam.NewIpamAsnsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamAsnsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modiefiedFields, asnRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamASNRead(ctx, d, m)
}

func resourceNetboxIpamASNDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamASNExists(d, m)
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

	resource := ipam.NewIpamAsnsDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamAsnsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamASNExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamAsnsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamAsnsList(params, nil)
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
