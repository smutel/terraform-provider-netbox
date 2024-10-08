package dcim

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxDcimRack() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a rack within Netbox.",
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
				Description:  "A unique tag used to identify this rack.",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   util.TrimString,
				Description: "Comments for this rack.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this rack.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this rack was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"desc_units": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "True if rack units are numbered top-to-bottom.",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of devices associated to this rack.",
			},
			"facility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
				Description:  "Local facility ID or description.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this rack was last updated.",
			},
			"location_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the location for this rack.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this rack.",
			},
			"outer_depth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Outer depth of this rack.",
			},
			"outer_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Outer unit among mm or in of this rack.",
				ValidateFunc: validation.StringInSlice([]string{"mm", "in"}, false),
			},
			"outer_width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Outer width of this rack.",
			},
			"power_feed_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The power feed count of this rack.",
			},
			"role_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the role associated to this rack.",
			},
			"serial": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
				Description:  "The serial number of this rack.",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the site where this rack is attached.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				ValidateFunc: validation.StringInSlice([]string{"reserved", "available", "planned", "active", "deprecated"}, false),
				Description:  "The status among reserved, available, planned, active or deprecated (active by default) of this rack.",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The tenant of this rack.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"2-post-frame", "4-post-frame", "4-post-cabinet", "wall-frame", "wall-frame-vertical", "wall-cabinet", "wall-cabinet-vertical"}, false),
				Description:  "The type among 2-post-frame, 4-post-frame, 4-post-cabinet, wall-frame, wall-frame-vertical, wall-cabinet or wall-cabinet-vertical (active by default) of this rack.",
			},
			"height": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Height in rack units of this rack.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this rack.",
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
				Description: "The type among 10, 19, 21 or 23 (inches) of this rack.",
			},
		},
	}
}

func resourceNetboxDcimRackCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	assetTag := d.Get("asset_tag").(string)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	descendingUnits := d.Get("desc_units").(bool)
	facilityID := d.Get("facility").(string)
	locationID := int32(d.Get("location_id").(int))
	name := d.Get("name").(string)
	outerDepth := int32(d.Get("outer_depth").(int))
	outerUnit := d.Get("outer_unit").(string)
	outerWidth := int32(d.Get("outer_width").(int))
	roleID := int32(d.Get("role_id").(int))
	siteID := int32(d.Get("site_id").(int))
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	rackType := d.Get("type").(string)
	tenantID := int32(d.Get("tenant_id").(int))
	height := int32(d.Get("height").(int))
	width := int32(d.Get("width").(int))

	newResource := netbox.NewWritableRackRequestWithDefaults()
	newResource.SetAssetTag(assetTag)
	newResource.SetComments(d.Get("comments").(string))
	newResource.SetCustomFields(customFields)
	newResource.SetDescUnits(descendingUnits)
	newResource.SetFacilityId(facilityID)
	newResource.SetName(name)
	newResource.SetOuterUnit(netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
	newResource.SetSerial(d.Get("serial").(string))
	newResource.SetStatus(netbox.PatchedWritableRackRequestStatus(status))
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	newResource.SetType(netbox.PatchedWritableRackRequestType(rackType))
	newResource.SetUHeight(height)
	newResource.SetWidth(netbox.PatchedWritableRackRequestWidth(width))

	if siteID != 0 {
		b, err := brief.GetBriefSiteRequestFromID(client, ctx, siteID)
		if err != nil {
			return err
		}
		newResource.SetSite(*b)
	}

	if locationID != 0 {
		b, err := brief.GetBriefLocationRequestFromID(client, ctx, locationID)
		if err != nil {
			return err
		}
		newResource.SetLocation(*b)
	}

	if outerDepth != 0 {
		newResource.SetOuterDepth(outerDepth)
	}

	if outerWidth != 0 {
		newResource.SetOuterWidth(outerWidth)
	}

	if roleID != 0 {
		b, err := brief.GetBriefRackRoleRequestFromID(client, ctx, roleID)
		if err != nil {
			return err
		}
		newResource.SetRole(*b)
	}

	if tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(client, ctx, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	_, response, err := client.DcimAPI.DcimRacksCreate(ctx).WritableRackRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if resourceID, err := util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	} else {
		d.SetId(fmt.Sprintf("%d", resourceID))
	}

	return resourceNetboxDcimRackRead(ctx, d, m)
}

func resourceNetboxDcimRackRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.DcimAPI.DcimRacksRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("asset_tag", resource.GetAssetTag()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("comments", resource.GetComments()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(resource.GetUrl())); err != nil {
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

	if err = d.Set("desc_units", resource.GetDescUnits()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("device_count", resource.DeviceCount); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("facility", resource.GetFacilityId()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("location_id", resource.GetLocation().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("outer_depth", resource.GetOuterDepth()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("outer_unit", resource.GetOuterUnit().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("outer_width", resource.GetOuterWidth()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("power_feed_count", resource.GetPowerfeedCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("role_id", resource.GetRole().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("serial", resource.GetSerial()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("site_id", resource.GetSite().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("status", resource.GetStatus().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tenant_id", resource.GetTenant().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("type", resource.GetType().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("height", resource.GetUHeight()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("width", resource.GetWidth().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimRackUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableRackRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))
	b, errDiag := brief.GetBriefSiteRequestFromID(client, ctx, int32(d.Get("site_id").(int)))
	if errDiag != nil {
		return errDiag
	}
	resource.SetSite(*b)

	if d.HasChange("asset_tag") {
		assetTag := d.Get("asset_tag").(string)
		resource.SetAssetTag(assetTag)
	}

	if d.HasChange("comments") {
		comments := d.Get("comments").(string)
		resource.SetComments(comments)
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("desc_units") {
		descendingUnits := d.Get("desc_units").(bool)
		resource.SetDescUnits(descendingUnits)
	}

	if d.HasChange("facility") {
		if facility, exist := d.GetOk("facility"); exist {
			resource.SetFacilityId(facility.(string))
		} else {
			resource.SetFacilityIdNil()
		}
	}

	if d.HasChange("location_id") {
		if locationID, exist := d.GetOk("location_id"); exist {
			b, err := brief.GetBriefLocationRequestFromID(client, ctx, int32(locationID.(int)))
			if err != nil {
				return err
			}
			resource.SetLocation(*b)
		} else {
			resource.SetLocationNil()
		}
	}

	if d.HasChange("outer_unit") {
		outerUnit := d.Get("outer_unit").(string)
		resource.SetOuterUnit(netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
	}

	if d.HasChange("outer_depth") {
		if outerDepth, exist := d.GetOk("outer_depth"); exist {
			resource.SetOuterDepth(int32(outerDepth.(int)))

			outerUnit := d.Get("outer_unit").(string)
			resource.SetOuterUnit(netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
		} else {
			resource.SetOuterDepthNil()
		}
	}

	if d.HasChange("outer_width") {
		if outerWidth, exist := d.GetOk("outer_width"); exist {
			resource.SetOuterWidth(int32(outerWidth.(int)))

			outerUnit := d.Get("outer_unit").(string)
			resource.SetOuterUnit(netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
		} else {
			resource.SetOuterWidthNil()
		}
	}

	if d.HasChange("role_id") {
		if roleID, exist := d.GetOk("role_id"); exist {
			b, err := brief.GetBriefRackRoleRequestFromID(client, ctx, int32(roleID.(int)))
			if err != nil {
				return err
			}
			resource.SetRole(*b)
		} else {
			resource.SetRoleNil()
		}
	}

	if d.HasChange("serial") {
		serial := d.Get("serial").(string)
		resource.SetSerial(serial)
	}

	if d.HasChange("status") {
		resource.SetStatus(netbox.PatchedWritableRackRequestStatus(d.Get("status").(string)))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		if tenantID, exist := d.GetOk("tenant_id"); exist {
			b, err := brief.GetBriefTenantRequestFromID(client, ctx, int32(tenantID.(int)))
			if err != nil {
				return err
			}
			resource.SetTenant(*b)
		} else {
			resource.SetTenantNil()
		}
	}

	if d.HasChange("type") {
		rackType := d.Get("type").(string)
		resource.SetType(netbox.PatchedWritableRackRequestType(rackType))
	}

	if d.HasChange("height") {
		height := int32(d.Get("height").(int))
		resource.SetUHeight(height)
	}

	if d.HasChange("width") {
		width := int64(d.Get("width").(int))
		resource.SetWidth(netbox.PatchedWritableRackRequestWidth(width))
	}

	if _, response, err := client.DcimAPI.DcimRacksUpdate(ctx, int32(resourceID)).WritableRackRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimRackRead(ctx, d, m)
}

func resourceNetboxDcimRackDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimRackExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int64"))
	}

	if response, err := client.DcimAPI.DcimRacksDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimRackExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimRacksRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
