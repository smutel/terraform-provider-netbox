package dcim

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/dcim"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v4/netbox/internal/util"
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
				Description:  "The status of this site. Alowed values: \"active\" (default), \"planned\", \"staging\", \"decommisioning\", \"retired\".",
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

var siteRequiredFields = []string{
	"created",
	"last_updated",
	"asns",
	"name",
	"slug",
	"tags",
}

func resourceNetboxDcimSiteCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	asns := d.Get("asns").(*schema.Set).List()
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	groupID := int64(d.Get("group_id").(int))
	latitude := d.Get("latitude").(float64)
	longitude := d.Get("longitude").(float64)
	name := d.Get("name").(string)
	regionID := int64(d.Get("region_id").(int))
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))

	newResource := &models.WritableSite{
		Asns:            util.ToListofInts(asns),
		Comments:        d.Get("comments").(string),
		CustomFields:    customFields,
		Facility:        d.Get("facility").(string),
		Description:     d.Get("description").(string),
		Latitude:        &latitude,
		Longitude:       &longitude,
		Name:            &name,
		PhysicalAddress: d.Get("physical_address").(string),
		ShippingAddress: d.Get("shipping_address").(string),
		Slug:            &slug,
		Status:          d.Get("status").(string),
		Tags:            tag.ConvertTagsToNestedTags(tags),
		TimeZone:        d.Get("time_zone").(string),
	}

	if groupID != 0 {
		newResource.Group = &groupID
	}
	if regionID != 0 {
		newResource.Region = &regionID
	}
	if tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	resource := dcim.NewDcimSitesCreateParams().WithData(newResource)

	resourceCreated, err := client.Dcim.DcimSitesCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxDcimSiteRead(ctx, d, m)
}

func resourceNetboxDcimSiteRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := dcim.NewDcimSitesListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimSitesList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]
	if err = d.Set("asns", util.ConvertNestedASNsToASNs(resource.Asns)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("circuit_count", resource.CircuitCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comments", resource.Comments); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
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
	if err = d.Set("device_count", resource.DeviceCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("facility", resource.Facility); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("group_id", util.GetNestedSiteGroupID(resource.Group)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("latitude", resource.Latitude); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("longitude", resource.Longitude); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("physical_address", resource.PhysicalAddress); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("prefix_count", resource.PrefixCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("rack_count", resource.RackCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("region_id", util.GetNestedRegionID(resource.Region)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("shipping_address", resource.ShippingAddress); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("slug", resource.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("status", resource.Status.Value); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tenant_id", util.GetNestedTenantID(resource.Tenant)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("time_zone", resource.TimeZone); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("virtualmachine_count", resource.VirtualmachineCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vlan_count", resource.VlanCount); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimSiteUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableSite{}

	modifiedFields := map[string]interface{}{}
	if d.HasChange("asns") {
		params.Asns = util.ToListofInts(d.Get("asns").(*schema.Set).List())
	}
	if d.HasChange("comments") {
		comments := d.Get("comments").(string)
		params.Comments = comments
		modifiedFields["comments"] = comments
	}
	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		params.Description = description
		modifiedFields["description"] = description
	}
	if d.HasChange("facility") {
		facility := d.Get("facility").(string)
		params.Facility = facility
		modifiedFields["facility"] = facility
	}
	if d.HasChange("group_id") {
		groupID := int64(d.Get("group_id").(int))
		params.Group = &groupID
		modifiedFields["group"] = groupID
	}
	if d.HasChange("latitude") {
		latitude := d.Get("latitude").(float64)
		params.Latitude = &latitude
		modifiedFields["latitude"] = latitude
	}
	if d.HasChange("longitude") {
		longitude := d.Get("longitude").(float64)
		params.Longitude = &longitude
		modifiedFields["longitude"] = longitude
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("physical_address") {
		physicalAddress := d.Get("physical_address").(string)
		params.PhysicalAddress = physicalAddress
		modifiedFields["physical_address"] = physicalAddress
	}
	if d.HasChange("region_id") {
		regionID := int64(d.Get("region_id").(int))
		params.Region = &regionID
		modifiedFields["region"] = regionID
	}
	if d.HasChange("status") {
		params.Status = d.Get("status").(string)
	}
	if d.HasChange("slug") {
		slug := d.Get("slug").(string)
		params.Slug = &slug
	}
	if d.HasChange("shipping_address") {
		shippingAddress := d.Get("shipping_address").(string)
		params.ShippingAddress = shippingAddress
		modifiedFields["shipping_address"] = shippingAddress
	}
	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		params.Tenant = &tenantID
		modifiedFields["tenant"] = tenantID
	}
	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}
	if d.HasChange("time_zone") {
		timeZone := d.Get("time_zone").(string)
		params.TimeZone = timeZone
		modifiedFields["time_zone"] = timeZone
	}

	resource := dcim.NewDcimSitesPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Dcim.DcimSitesPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, siteRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxDcimSiteRead(ctx, d, m)
}

func resourceNetboxDcimSiteDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxDcimSiteExists(d, m)
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

	resource := dcim.NewDcimSitesDeleteParams().WithID(id)
	if _, err := client.Dcim.DcimSitesDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimSiteExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := dcim.NewDcimSitesListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimSitesList(params, nil)
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
