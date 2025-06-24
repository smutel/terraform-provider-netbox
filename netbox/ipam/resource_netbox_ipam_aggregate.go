package ipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxIpamAggregate() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an aggregate within Netbox.",
		CreateContext: resourceNetboxIpamAggregateCreate,
		ReadContext:   resourceNetboxIpamAggregateRead,
		UpdateContext: resourceNetboxIpamAggregateUpdate,
		DeleteContext: resourceNetboxIpamAggregateDelete,
		Exists:        resourceNetboxIpamAggregateExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this aggregate.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this aggregate was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"date_added": {
				Type:     schema.TypeString,
				Optional: true,
				//nolint:revive
				ValidateFunc: func(val any, key string) (
					warns []string, errs []error) {
					v := val.(string)
					_, err := time.Parse("2006-01-02", v)

					if err != nil {
						errs = append(errs,
							fmt.Errorf("date_added in not in the good format "+
								"YYYY-MM-DD"))
					}
					return warns, errs
				},
				Description: "Date when this aggregate was added. " +
					"Format *YYYY-MM-DD*.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The description of this aggregate.",
			},
			"family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP family of this aggregate.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this aggregate was last updated.",
			},
			"prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, util.Const256),
				Description:  "The network prefix of this aggregate.",
			},
			"rir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The RIR id linked to this aggregate.",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this object is attached.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this tag.",
			},
		},
	}
}

func resourceNetboxIpamAggregateCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
		nil, resourceCustomFields)
	dateAdded := d.Get("date_added").(string)
	description := d.Get("description").(string)
	prefix := d.Get("prefix").(string)
	rirID := d.Get("rir_id").(int)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritableAggregateRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetPrefix(prefix)
	b, err := brief.GetBriefRIRRequestFromID(ctx, client, rirID)
	if err != nil {
		return err
	}
	newResource.SetRir(*b)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	if dateAdded != "" {
		newResource.SetDateAdded(dateAdded)
	}

	_, response, errDiag :=
		client.IpamAPI.IpamAggregatesCreate(
			ctx).WritableAggregateRequest(
			*newResource).Execute()
	if response.StatusCode != util.Const201 && errDiag != nil {
		return util.GenerateErrorMessage(response, errDiag)
	}

	var resourceID int32
	if resourceID, errDiag = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, errDiag)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamAggregateRead(ctx, d, m)
}

func resourceNetboxIpamAggregateRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err :=
		client.IpamAPI.IpamAggregatesRetrieve(
			ctx, int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
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

	if err = d.Set("date_added", resource.GetDateAdded()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("description", resource.GetDescription()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("family", resource.GetFamily().Label); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("prefix", resource.GetPrefix()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("rir_id", resource.GetRir().Id); err != nil {
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

func resourceNetboxIpamAggregateUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableAggregateRequestWithDefaults()
	rirID := d.Get("rir_id").(int)
	b, errDiag := brief.GetBriefRIRRequestFromID(ctx, client, rirID)
	if errDiag != nil {
		return errDiag
	}
	resource.SetRir(*b)
	resource.SetPrefix(d.Get("prefix").(string))

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(),
			resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("date_added") {
		if dateAdded, exist := d.GetOk("date_added"); exist {
			resource.SetDateAdded(dateAdded.(string))
		} else {
			resource.SetDateAddedNil()
		}
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

	if _, response, err := client.IpamAPI.IpamAggregatesUpdate(ctx,
		int32(resourceID)).WritableAggregateRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamAggregateRead(ctx, d, m)
}

func resourceNetboxIpamAggregateDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamAggregateExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int"))
	}

	if response, err := client.IpamAPI.IpamAggregatesDestroy(
		ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamAggregateExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamAggregatesRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
