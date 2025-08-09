// Copyright (c)
// SPDX-License-Identifier: MIT

package ipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxIpamASN() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a ASN within Netbox.",
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
				Description:  "The ASN number of this ASN.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this ASN.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this ASN was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The description of this ASN.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last date when this ASN was updated.",
			},
			"provider_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of providers for this ASN.",
			},
			"rir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The RIR for this ASN.",
			},
			"site_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of sites for this ASN.",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The rir for this ASN.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this ASN.",
			},
		},
	}
}

func resourceNetboxIpamASNCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
		nil, resourceCustomFields)

	asn := int64(d.Get("asn").(int))
	rirID := d.Get("rir_id").(int)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewASNRequestWithDefaults()
	newResource.SetAsn(asn)
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	b, err := brief.GetBriefRIRRequestFromID(ctx, client, rirID)
	if err != nil {
		return err
	}
	newResource.SetRir(*b)

	if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	_, response, errDiag := client.IpamAPI.IpamAsnsCreate(
		ctx).ASNRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && errDiag != nil {
		return util.GenerateErrorMessage(response, errDiag)
	}

	var resourceID int32
	if resourceID, errDiag = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, errDiag)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamASNRead(ctx, d, m)
}

func resourceNetboxIpamASNRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.IpamAPI.IpamAsnsRetrieve(
		ctx, int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("asn", resource.GetAsn()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type",
		util.ConvertURLContentType(resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(
		resourceCustomFields, resource.GetCustomFields())

	if err = d.Set("custom_field", customFields); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("provider_count",
		resource.GetProviderCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("rir_id", resource.GetRir().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("site_count", resource.GetSiteCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
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

func resourceNetboxIpamASNUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewASNRequestWithDefaults()

	// Required fields
	resource.SetAsn(int64(d.Get("asn").(int)))
	rirID := d.Get("rir_id").(int)
	b, errDiag := brief.GetBriefRIRRequestFromID(ctx, client, rirID)
	if errDiag != nil {
		return errDiag
	}
	resource.SetRir(*b)

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields :=
			customfield.ConvertCustomFieldsFromTerraformToAPI(
				stateCustomFields.(*schema.Set).List(),
				resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			resource.SetDescription(description.(string))
		} else {
			resource.SetDescription("")
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		tenantID := d.Get("tenant_id").(int)
		if tenantID != 0 {
			b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
			if err != nil {
				return err
			}
			resource.SetTenant(*b)
		} else {
			resource.SetTenantNil()
		}
	}

	if _, response, err := client.IpamAPI.IpamAsnsUpdate(
		ctx, int32(resourceID)).ASNRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamASNRead(ctx, d, m)
}

func resourceNetboxIpamASNDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamASNExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}

	if response, err := client.IpamAPI.IpamAsnsDestroy(
		ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamASNExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamAsnsRetrieve(
		nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
