package dcim

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/dcim"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxDcimRack() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a rack (dcim module) within Netbox.",
		CreateContext: resourceNetboxDcimRackCreate,
		ReadContext:   resourceNetboxDcimRackRead,
		UpdateContext: resourceNetboxDcimRackUpdate,
		DeleteContext: resourceNetboxDcimRackDelete,
		Exists:        resourceNetboxDcimRackExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"asset_tag": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
				Description:  "A unique tag used to identify this rack (dcim module).",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   util.TrimString,
				Description: "Comments for this rack (dcim module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this rack (dcim module) was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"desc_units": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True if rack (dcim module) units are numbered top-to-bottom.",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of devices associated to this rack (dcim module).",
			},
			"facility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
				Description:  "Local facility ID or description (dcim module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this rack was last updated (dcim module).",
			},
			"location_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the location for this rack (dcim module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this rack (dcim module).",
			},
			"outer_depth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Outer depth of this rack (dcim module).",
			},
			"outer_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Outer unit among mm or in of this rack (dcim module).",
				ValidateFunc: validation.StringInSlice([]string{"mm", "in"}, false),
			},
			"outer_width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Outer width of this rack (dcim module).",
			},
			"power_feed_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The power feed count of this rack (dcim module).",
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the role associated to this rack (dcim module).",
			},
			"serial": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
				Description:  "The serial number of this rack (dcim module).",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the site where this rack (dcim module) is attached.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				ValidateFunc: validation.StringInSlice([]string{"reserved", "available", "planned", "active", "deprecated"}, false),
				Description:  "The status among reserved, available, planned, active or deprecated (active by default) of this rack (dcim module).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The tenant of this rack (dcim module).",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"2-post-frame", "4-post-frame", "4-post-cabinet", "wall-frame", "wall-frame-vertical", "wall-cabinet", "wall-cabinet-vertical"}, false),
				Description:  "The type among 2-post-frame, 4-post-frame, 4-post-cabinet, wall-frame, wall-frame-vertical, wall-cabinet or wall-cabinet-vertical (active by default) of this rack (dcim module).",
			},
			"height": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Height in rack units of this rack (dcim module).",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this rack (dcim module).",
			},
			"width": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(int)
					if v != 10 && v != 19 && v != 21 && v != 23 {
						errs = append(errs, fmt.Errorf("%q must be 10/19/21 or 23, got: %d", key, v))
					}
					return
				},
				Description: "The type among 10, 19, 21 or 23 (inches) of this rack (dcim module).",
			},
		},
	}
}

var rackRequiredFields = []string{
	"created",
	"height",
	"last_updated",
	"name",
	"site",
	"tags",
	"width",
	"desc_units",
}

func resourceNetboxDcimRackCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	assetTag := d.Get("asset_tag").(string)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	descendingUnits := d.Get("desc_units").(bool)
	facilityID := d.Get("facility").(string)
	locationID := int64(d.Get("location_id").(int))
	name := d.Get("name").(string)
	outerDepth := int64(d.Get("outer_depth").(int))
	outerUnit := d.Get("outer_unit").(string)
	outerWidth := int64(d.Get("outer_width").(int))
	roleID := int64(d.Get("role_id").(int))
	siteID := int64(d.Get("site_id").(int))
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))
	height := int64(d.Get("height").(int))
	width := int64(d.Get("width").(int))

	newResource := &models.WritableRack{
		AssetTag:     &assetTag,
		Comments:     d.Get("comments").(string),
		CustomFields: customFields,
		DescUnits:    descendingUnits,
		FacilityID:   &facilityID,
		Name:         &name,
		OuterUnit:    outerUnit,
		Serial:       d.Get("serial").(string),
		Site:         &siteID,
		Status:       d.Get("status").(string),
		Tags:         tag.ConvertTagsToNestedTags(tags),
		Type:         d.Get("type").(string),
		UHeight:      height,
		Width:        width,
	}

	if locationID != 0 {
		newResource.Location = &locationID
	}
	if outerDepth != 0 {
		newResource.OuterDepth = &outerDepth
	}
	if outerWidth != 0 {
		newResource.OuterWidth = &outerWidth
	}
	if roleID != 0 {
		newResource.Role = &roleID
	}
	if tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	resource := dcim.NewDcimRacksCreateParams().WithData(newResource)

	resourceCreated, err := client.Dcim.DcimRacksCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxDcimRackRead(ctx, d, m)
}

func resourceNetboxDcimRackRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := dcim.NewDcimRacksListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimRacksList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]
	if err = d.Set("asset_tag", resource.AssetTag); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comments", resource.Comments); err != nil {
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
	if err = d.Set("desc_units", resource.DescUnits); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("device_count", resource.DeviceCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("facility", resource.FacilityID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("location_id", util.GetNestedLocationID(resource.Location)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("outer_depth", resource.OuterDepth); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("outer_unit", util.GetRackOuterUnit(resource.OuterUnit)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("outer_width", resource.OuterWidth); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("power_feed_count", resource.PowerfeedCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("role_id", util.GetNestedRackRoleID(resource.Role)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("serial", resource.Serial); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("site_id", util.GetNestedSiteID(resource.Site)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("status", util.GetRackStatusValue(resource.Status)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tenant_id", util.GetNestedTenantID(resource.Tenant)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type", util.GetRackTypeValue(resource.Type)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("height", resource.UHeight); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("width", resource.Width.Value); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimRackUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableRack{}

	modifiedFields := map[string]interface{}{}
	if d.HasChange("asset_tag") {
		assetTag := d.Get("asset_tag").(string)
		params.AssetTag = &assetTag
		modifiedFields["asset_tag"] = assetTag
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
	if d.HasChange("desc_units") {
		descendingUnits := d.Get("desc_units").(bool)
		params.DescUnits = descendingUnits
		modifiedFields["desc_units"] = descendingUnits
	}
	if d.HasChange("facility") {
		facility := d.Get("facility").(string)
		params.FacilityID = &facility
		modifiedFields["facility"] = facility
	}
	if d.HasChange("location_id") {
		locationID := int64(d.Get("location_id").(int))
		params.Location = &locationID
		modifiedFields["location"] = locationID
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("outer_depth") {
		outerDepth := int64(d.Get("outer_depth").(int))
		params.OuterDepth = &outerDepth
		modifiedFields["outer_depth"] = outerDepth
	}
	if d.HasChange("outer_unit") {
		outerUnit := d.Get("outer_unit").(string)
		params.OuterUnit = outerUnit
		modifiedFields["outer_unit"] = outerUnit
	}
	if d.HasChange("outer_width") {
		outerWidth := int64(d.Get("outer_width").(int))
		params.OuterWidth = &outerWidth
		modifiedFields["outer_width"] = outerWidth
	}
	if d.HasChange("role_id") {
		roleID := int64(d.Get("role_id").(int))
		params.Role = &roleID
		modifiedFields["role"] = roleID
	}
	if d.HasChange("serial") {
		serial := d.Get("serial").(string)
		params.Serial = serial
		modifiedFields["serial"] = serial
	}
	if d.HasChange("site_id") {
		siteID := int64(d.Get("site_id").(int))
		params.Site = &siteID
		modifiedFields["site"] = siteID
	}
	if d.HasChange("status") {
		params.Status = d.Get("status").(string)
	}
	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}
	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		params.Tenant = &tenantID
		modifiedFields["tenant"] = tenantID
	}
	if d.HasChange("type") {
		RackType := d.Get("type").(string)
		params.Type = RackType
		modifiedFields["type"] = RackType
	}
	if d.HasChange("height") {
		height := int64(d.Get("height").(int))
		params.UHeight = height
		modifiedFields["height"] = height
	}
	if d.HasChange("width") {
		width := int64(d.Get("width").(int))
		params.Width = width
		modifiedFields["width"] = width
	}

	resource := dcim.NewDcimRacksPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Dcim.DcimRacksPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, rackRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxDcimRackRead(ctx, d, m)
}

func resourceNetboxDcimRackDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxDcimRackExists(d, m)
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

	resource := dcim.NewDcimRacksDeleteParams().WithID(id)
	if _, err := client.Dcim.DcimRacksDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimRackExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := dcim.NewDcimRacksListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimRacksList(params, nil)
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
