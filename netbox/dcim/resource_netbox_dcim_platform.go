package dcim

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

func ResourceNetboxDcimPlatform() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a platform within Netbox.",
		CreateContext: resourceNetboxDcimPlatformCreate,
		ReadContext:   resourceNetboxDcimPlatformRead,
		UpdateContext: resourceNetboxDcimPlatformUpdate,
		DeleteContext: resourceNetboxDcimPlatformDelete,
		Exists:        resourceNetboxDcimPlatformExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"config_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The config template used for this platform.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this platform.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this platform was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, util.Const100),
				Description:  "The description of this platform.",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of devices this platform.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this platform was last updated.",
			},
			"manufacturer_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The manufacturer of this platform.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name of this platform.",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The slug of this platform.",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this platform.",
			},
			"virtualmachine_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of virtual machines of this platform.",
			},
		},
	}
}

func resourceNetboxDcimPlatformCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	configTemplateID := d.Get("config_template_id").(int)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil,
		resourceCustomFields)
	description := d.Get("description").(string)
	manufacturerID := d.Get("manufacturer_id").(int)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewPlatformRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetName(name)
	newResource.SetSlug(slug)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	if manufacturerID != 0 {
		b, err := brief.GetBriefManufacturerRequestFromID(ctx, client,
			manufacturerID)
		if err != nil {
			return err
		}
		newResource.SetManufacturer(*b)
	}

	if configTemplateID != 0 {
		b, err := brief.GetBriefConfigTemplateRequestFromID(ctx, client,
			manufacturerID)
		if err != nil {
			return err
		}
		newResource.SetConfigTemplate(*b)
	}

	_, response, err := client.DcimAPI.DcimPlatformsCreate(
		ctx).PlatformRequest(*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxDcimPlatformRead(ctx, d, m)
}

func resourceNetboxDcimPlatformRead(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.DcimAPI.DcimPlatformsRetrieve(ctx,
		int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("config_template_id",
		resource.GetConfigTemplate().Id); err != nil {
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

	if err = d.Set("device_count", resource.GetDeviceCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("manufacturer_id",
		resource.GetManufacturer().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("tag",
		tag.ConvertNestedTagRequestToTags(resource.Tags)); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("virtualmachine_count",
		resource.GetVirtualmachineCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimPlatformUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewPlatformRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))
	resource.SetSlug(d.Get("slug").(string))

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
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

	if d.HasChange("manufacturer_id") {
		if manufacturerID, exist := d.GetOk("manufacturer_id"); exist {
			//nolint:revive
			if manufacturer, err := brief.GetBriefManufacturerRequestFromID(ctx,
				client, manufacturerID.(int)); err == nil {
				resource.SetManufacturer(*manufacturer)
			} else {
				return err
			}
		} else {
			resource.SetManufacturerNil()
		}
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if _, response, err := client.DcimAPI.DcimPlatformsUpdate(ctx,
		int32(resourceID)).PlatformRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimPlatformRead(ctx, d, m)
}

func resourceNetboxDcimPlatformDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimPlatformExists(d, m)
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

	if response, err := client.DcimAPI.DcimPlatformsDestroy(ctx,
		int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimPlatformExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimPlatformsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
