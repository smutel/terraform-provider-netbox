package dcim

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v3"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxDcimSite() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a site (dcim module) within Netbox.",
		CreateContext: resourceNetboxDcimSiteCreate,
		ReadContext:   resourceNetboxDcimSiteRead,
		UpdateContext: resourceNetboxDcimSiteUpdate,
		DeleteContext: resourceNetboxDcimSiteDelete,
		Exists:        resourceNetboxDcimSiteExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"asns": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "ASNs",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"circuit_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of circuits associated to this site (dcim module).",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   util.TrimString,
				Description: "Comments for this site (dcim module).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this site (dcim module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this site was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
				Description:  "The description of this site (dcim module).",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of devices associated to this site (dcim module).",
			},
			"facility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
				Description:  "Local facility ID or description.",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The site group for this site (dcim module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this site was last updated.",
			},
			"latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "GPS coordinate (latitude).",
			},
			"longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "GPS coordinate (longitude)",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this site (dcim module).",
			},
			"physical_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
				StateFunc:    util.TrimString,
				Description:  "The physical address of this site (dcim module).",
			},
			"prefix_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of prefixes associated to this site (dcim module).",
			},
			"rack_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of racks associated to this site (dcim module).",
			},
			"region_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The description of this site (dcim module).",
			},
			"shipping_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
				StateFunc:    util.TrimString,
				Description:  "The shipping address of this site (dcim module).",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this site (dcim module).",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				ValidateFunc: validation.StringInSlice([]string{"planned", "staging", "active", "decommisioning", "retired"}, false),
				Description:  "The status of this site. Allowed values: \"active\" (default), \"planned\", \"staging\", \"decommisioning\", \"retired\".",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The tenant of this site (dcim module).",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Timezone this site is in.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this site (dcim module).",
			},
			"virtualmachine_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of virtual machines associated to this site (dcim module).",
			},
			"vlan_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of vlans associated to this site (dcim module).",
			},
		},
	}
}

func resourceNetboxDcimSiteCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	asns := d.Get("asns").(*schema.Set).List()
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	groupID := int32(d.Get("group_id").(int))
	latitude := d.Get("latitude").(float64)
	longitude := d.Get("longitude").(float64)
	name := d.Get("name").(string)
	regionID := int32(d.Get("region_id").(int))
	slug := d.Get("slug").(string)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int32(d.Get("tenant_id").(int))

	newResource := netbox.NewWritableSiteRequestWithDefaults()
	newResource.SetAsns(util.ToListofInts(asns))
	newResource.SetComments(d.Get("comments").(string))
	newResource.SetCustomFields(customFields)
	newResource.SetFacility(d.Get("facility").(string))
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetLatitude(latitude)
	newResource.SetLongitude(longitude)
	newResource.SetName(name)
	newResource.SetPhysicalAddress(d.Get("physical_address").(string))
	newResource.SetShippingAddress(d.Get("shipping_address").(string))
	newResource.SetSlug(slug)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	s, err := netbox.NewLocationStatusValueFromValue(status)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetStatus(*s)
	if groupID != 0 {
		newResource.SetGroup(groupID)
	}

	if regionID != 0 {
		newResource.SetRegion(regionID)
	}

	if tenantID != 0 {
		newResource.SetTenant(tenantID)
	}

	if timezone := d.Get("time_zone").(string); timezone != "" {
		newResource.SetTimeZone(timezone)
	}

	resourceCreated, response, err := client.DcimAPI.DcimSitesCreate(ctx).WritableSiteRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resourceCreated.GetId() == 0 {
		return util.GenerateErrorMessage(response, errors.New("Bug Netbox - TO BE FIXED"))
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxDcimSiteRead(ctx, d, m)
}

func resourceNetboxDcimSiteRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.DcimAPI.DcimSitesRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("asns", resource.GetAsns()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("circuit_count", resource.GetCircuitCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("comments", resource.GetComments()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", util.ConvertURIContentType(strfmt.URI(resource.GetUrl()))); err != nil {
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

	if err = d.Set("device_count", resource.GetDeviceCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("facility", resource.GetFacility()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("group_id", resource.GetGroup().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("latitude", resource.GetLatitude()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("longitude", resource.GetLongitude()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("physical_address", resource.GetPhysicalAddress()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("prefix_count", resource.GetPrefixCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("rack_count", resource.GetRackCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("region_id", resource.GetRegion().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("shipping_address", resource.GetShippingAddress()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("status", resource.GetStatus().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.GetTags())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tenant_id", resource.GetTenant().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("time_zone", resource.GetTimeZone()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("virtualmachine_count", resource.GetVirtualmachineCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("vlan_count", resource.GetVlanCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimSiteUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	resource := netbox.NewWritableSiteRequestWithDefaults()

	if d.HasChange("asns") {
		resource.SetAsns(util.ToListofInts(d.Get("asns").(*schema.Set).List()))
	}

	if d.HasChange("comments") {
		resource.SetComments(d.Get("comments").(string))
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("facility") {
		resource.SetFacility(d.Get("facility").(string))
	}

	if d.HasChange("group_id") {
		if groupID, exist := d.GetOk("group_id"); exist {
			resource.SetGroup(int32(groupID.(int)))
		} else {
			resource.SetGroupNil()
		}
	}

	if d.HasChange("latitude") {
		resource.SetLatitude(d.Get("latitude").(float64))
	}

	if d.HasChange("longitude") {
		resource.SetLongitude(d.Get("longitude").(float64))
	}

	if d.HasChange("name") {
		resource.SetName(d.Get("name").(string))
	}

	if d.HasChange("physical_address") {
		resource.SetPhysicalAddress(d.Get("physical_address").(string))
	}

	if d.HasChange("region_id") {
		if regionID, exist := d.GetOk("region_id"); exist {
			resource.SetRegion(int32(regionID.(int)))
		} else {
			resource.SetRegionNil()
		}
	}

	if d.HasChange("status") {
		s, err := netbox.NewLocationStatusValueFromValue(d.Get("status").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetStatus(*s)
	}

	if d.HasChange("slug") {
		resource.SetSlug(d.Get("slug").(string))
	}

	if d.HasChange("shipping_address") {
		resource.SetShippingAddress(d.Get("shipping_address").(string))
	}

	if d.HasChange("tenant_id") {
		tenantID := int32(d.Get("tenant_id").(int))
		if tenantID != 0 {
			resource.SetTenant(tenantID)
		} else {
			resource.SetTenantNil()
		}
	}

	if d.HasChange("tenant_id") {
		tenantID := int32(d.Get("tenant_id").(int))
		if tenantID != 0 {
			resource.SetTenant(tenantID)
		} else {
			resource.SetTenantNil()
		}
	}

	if d.HasChange("time_zone") {
		TimeZone := d.Get("time_zone").(string)
		if TimeZone != "" {
			resource.SetTimeZone(TimeZone)
		} else {
			resource.SetTimeZoneNil()
		}
	}

	if _, response, err := client.DcimAPI.DcimSitesUpdate(ctx, int32(resourceID)).WritableSiteRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimSiteRead(ctx, d, m)
}

func resourceNetboxDcimSiteDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimSiteExists(d, m)
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

	if response, err := client.DcimAPI.DcimSitesDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimSiteExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimSitesRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
