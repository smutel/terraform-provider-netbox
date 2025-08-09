// Copyright (c)
// SPDX-License-Identifier: MIT

package dcim

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
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
				ValidateFunc: validation.StringLenBetween(0, util.Const50),
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
				ValidateFunc: validation.StringLenBetween(0, util.Const50),
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
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name of this rack.",
			},
			"outer_depth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Outer depth of this rack.",
			},
			"outer_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Outer unit among mm or in of this rack.",
				ValidateFunc: validation.StringInSlice([]string{"mm", "in"},
					false),
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
				ValidateFunc: validation.StringLenBetween(0, util.Const50),
				Description:  "The serial number of this rack.",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the site where this rack is attached.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"reserved",
					"available", "planned", "active", "deprecated"}, false),
				Description: "The status among reserved, available, planned, " +
					"active or deprecated (active by default) of this rack.",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The tenant of this rack.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"2-post-frame",
					"4-post-frame", "4-post-cabinet", "wall-frame",
					"wall-frame-vertical", "wall-cabinet",
					"wall-cabinet-vertical"}, false),
				Description: "The type among 2-post-frame, 4-post-frame, " +
					"4-post-cabinet, wall-frame, wall-frame-vertical, " +
					"wall-cabinet or wall-cabinet-vertical (active by " +
					"default) of this rack.",
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
				ValidateFunc: func(val any, key string) (warns []string,
					errs []error) {
					v := val.(int)
					if v != util.Const10 && v != util.Const19 &&
						v != util.Const21 && v != util.Const23 {
						errs = append(errs,
							fmt.Errorf("%q must be 10/19/21 or 23, got: %d",
								key, v))
					}
					return warns, errs
				},
				Description: "The type among 10, 19, 21 or 23 (inches) " +
					"of this rack.",
			},
		},
	}
}

func resourceNetboxDcimRackCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	assetTag := d.Get("asset_tag").(string)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	descendingUnits := d.Get("desc_units").(bool)
	facilityID := d.Get("facility").(string)
	locationID := d.Get("location_id").(int)
	name := d.Get("name").(string)
	outerDepth := d.Get("outer_depth").(int)
	outerUnit := d.Get("outer_unit").(string)
	outerWidth := d.Get("outer_width").(int)
	roleID := d.Get("role_id").(int)
	siteID := d.Get("site_id").(int)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	rackType := d.Get("type").(string)
	tenantID := d.Get("tenant_id").(int)
	height := d.Get("height").(int)
	width := d.Get("width").(int)
	newResource := netbox.NewWritableRackRequestWithDefaults()
	newResource.SetAssetTag(assetTag)
	newResource.SetComments(d.Get("comments").(string))
	newResource.SetCustomFields(customFields)
	newResource.SetDescUnits(descendingUnits)
	newResource.SetFacilityId(facilityID)
	newResource.SetName(name)
	newResource.SetOuterUnit(
		netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
	newResource.SetSerial(d.Get("serial").(string))
	newResource.SetStatus(netbox.PatchedWritableRackRequestStatus(status))
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	newResource.SetType(netbox.PatchedWritableRackRequestType(rackType))

	height32, err := safecast.ToInt32(height)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetUHeight(height32)

	width32, err := safecast.ToInt32(width)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetWidth(netbox.PatchedWritableRackRequestWidth(width32))

	if siteID != 0 {
		b, err := brief.GetBriefSiteRequestFromID(ctx, client, siteID)
		if err != nil {
			return err
		}
		newResource.SetSite(*b)
	}

	if locationID != 0 {
		b, err := brief.GetBriefLocationRequestFromID(ctx, client, locationID)
		if err != nil {
			return err
		}
		newResource.SetLocation(*b)
	}

	if outerDepth != 0 {
		outerDepth32, err := safecast.ToInt32(outerDepth)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetOuterDepth(outerDepth32)
	}

	if outerWidth != 0 {
		outerWidth32, err := safecast.ToInt32(outerWidth)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		newResource.SetOuterWidth(outerWidth32)
	}

	if roleID != 0 {
		b, err := brief.GetBriefRackRoleRequestFromID(ctx, client, roleID)
		if err != nil {
			return err
		}
		newResource.SetRole(*b)
	}

	if tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	_, response, err := client.DcimAPI.DcimRacksCreate(
		ctx).WritableRackRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxDcimRackRead(ctx, d, m)
}

func resourceNetboxDcimRackRead(ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.DcimAPI.DcimRacksRetrieve(ctx,
		int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
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

	if err = d.Set("content_type", util.ConvertURLContentType(
		resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("created", resource.GetCreated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.UpdateCustomFieldsFromAPI(resourceCustomFields,
		resource.GetCustomFields())

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

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
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

	if err = d.Set("power_feed_count",
		resource.GetPowerfeedCount()); err != nil {
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

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
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

//nolint:gocyclo
func resourceNetboxDcimRackUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewWritableRackRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))
	b, errDiag := brief.GetBriefSiteRequestFromID(ctx, client,
		d.Get("site_id").(int))
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
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(),
			resourceCustomFields.(*schema.Set).List())
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
			b, err := brief.GetBriefLocationRequestFromID(ctx, client,
				locationID.(int))
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
		resource.SetOuterUnit(
			netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
	}

	if d.HasChange("outer_depth") {
		if outerDepth, exist := d.GetOk("outer_depth"); exist {
			outerDepth32, err := safecast.ToInt32(outerDepth.(int))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			resource.SetOuterDepth(outerDepth32)

			outerUnit := d.Get("outer_unit").(string)
			resource.SetOuterUnit(
				netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
		} else {
			resource.SetOuterDepthNil()
		}
	}

	if d.HasChange("outer_width") {
		if outerWidth, exist := d.GetOk("outer_width"); exist {
			outerWidth32, err := safecast.ToInt32(outerWidth.(int))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			resource.SetOuterWidth(outerWidth32)

			outerUnit := d.Get("outer_unit").(string)
			resource.SetOuterUnit(
				netbox.PatchedWritableRackRequestOuterUnit(outerUnit))
		} else {
			resource.SetOuterWidthNil()
		}
	}

	if d.HasChange("role_id") {
		if roleID, exist := d.GetOk("role_id"); exist {
			b, err := brief.GetBriefRackRoleRequestFromID(ctx, client,
				roleID.(int))
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
		resource.SetStatus(netbox.PatchedWritableRackRequestStatus(
			d.Get("status").(string)))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		if tenantID, exist := d.GetOk("tenant_id"); exist {
			b, err := brief.GetBriefTenantRequestFromID(ctx, client,
				tenantID.(int))
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
		height := d.Get("height").(int)
		height32, err := safecast.ToInt32(height)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetUHeight(height32)
	}

	if d.HasChange("width") {
		width32, err := safecast.ToInt32(d.Get("width").(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetWidth(netbox.PatchedWritableRackRequestWidth(width32))
	}

	if _, response, err := client.DcimAPI.DcimRacksUpdate(ctx,
		int32(resourceID)).WritableRackRequest(
		*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimRackRead(ctx, d, m)
}

func resourceNetboxDcimRackDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimRackExists(d, m)
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

	if response, err := client.DcimAPI.DcimRacksDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimRackExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimRacksRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
