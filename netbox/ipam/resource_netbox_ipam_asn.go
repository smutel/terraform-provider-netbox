package ipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v3"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxIpamASN() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a asn (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamASNCreate,
		ReadContext:   resourceNetboxIpamASNRead,
		UpdateContext: resourceNetboxIpamASNUpdate,
		DeleteContext: resourceNetboxIpamASNDelete,
		Exists:        resourceNetboxIpamASNExists,
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

func resourceNetboxIpamASNCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)

	asn := int64(d.Get("asn").(int))
	rirID := int32(d.Get("rir_id").(int))
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritableASNRequestWithDefaults()
	newResource.SetAsn(asn)
	newResource.SetRir(rirID)
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	if tenantID := int32(d.Get("tenant_id").(int)); tenantID != 0 {
		newResource.SetTenant(tenantID)
	}

	resourceCreated, response, err := client.IpamAPI.IpamAsnsCreate(ctx).WritableASNRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resourceCreated.GetId() == 0 {
		return diag.FromErr(errors.New("Bug Netbox - TO BE FIXED"))
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxIpamASNRead(ctx, d, m)
}

func resourceNetboxIpamASNRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.IpamAPI.IpamAsnsRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("asn", resource.GetAsn()); err != nil {
		return util.GenerateErrorMessage(nil, err)
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

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("provider_count", resource.GetProviderCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("rir_id", resource.GetRir().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("site_count", resource.GetSiteCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tenant_id", resource.GetTenant().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamASNUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	resource := netbox.NewWritableASNRequestWithDefaults()
	//resource.SetAsn(int64(d.Get("asn").(int)))
	//resource.SetRir(int32(d.Get("rir_id").(int)))

	if d.HasChange("asn") {
		resource.SetAsn(int64(d.Get("asn").(int)))
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			resource.SetDescription(description.(string))
		} else {
			resource.SetDescription("")
		}
	}

	if d.HasChange("rir_id") {
		rirID := int32(d.Get("rir_id").(int))
		resource.SetRir(rirID)
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		tenantID := int32(d.Get("tenant_id").(int))
		if tenantID != 0 {
			resource.SetTenant(tenantID)
		} else {
			resource.SetTenantNil()
		}
	}

	if _, response, err := client.IpamAPI.IpamAsnsUpdate(ctx, int32(resourceID)).WritableASNRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamASNRead(ctx, d, m)
}

func resourceNetboxIpamASNDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamASNExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	if response, err := client.IpamAPI.IpamAsnsDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamASNExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamAsnsRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
