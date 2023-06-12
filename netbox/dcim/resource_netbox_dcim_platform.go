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
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxDcimPlatform() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a platform (dcim module) within Netbox.",
		CreateContext: resourceNetboxDcimPlatformCreate,
		ReadContext:   resourceNetboxDcimPlatformRead,
		UpdateContext: resourceNetboxDcimPlatformUpdate,
		DeleteContext: resourceNetboxDcimPlatformDelete,
		Exists:        resourceNetboxDcimPlatformExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this platform (dcim module).",
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
				ValidateFunc: validation.StringLenBetween(0, 100),
				Description:  "The description of this platform (dcim module).",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of devices this platform (dcim module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this platform was last updated.",
			},
			"manufacturer_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The manufacturer of this platform (dcim module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this platform (dcim module).",
			},
			"napalm_args": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Argument for the napalm driver.",
			},
			"napalm_driver": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
				Description:  "The napalm driver",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this platform (dcim module).",
			},
			"tag": &tag.TagSchema,
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this platform (dcim module).",
			},
			"virtualmachine_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of virtual machines of this platform (dcim module).",
			},
		},
	}
}

var platformRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"slug",
	"tags",
}

func resourceNetboxDcimPlatformCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	manufacturerID := int64(d.Get("manufacturer_id").(int))
	name := d.Get("name").(string)
	napalmArgs := d.Get("napalm_args").(string)
	slug := d.Get("slug").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritablePlatform{
		CustomFields: customFields,
		Description:  d.Get("description").(string),
		Name:         &name,
		NapalmArgs:   &napalmArgs,
		NapalmDriver: d.Get("napalm_driver").(string),
		Slug:         &slug,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}
	if manufacturerID != 0 {
		newResource.Manufacturer = &manufacturerID
	}

	resource := dcim.NewDcimPlatformsCreateParams().WithData(newResource)

	resourceCreated, err := client.Dcim.DcimPlatformsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxDcimPlatformRead(ctx, d, m)
}

func resourceNetboxDcimPlatformRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := dcim.NewDcimPlatformsListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimPlatformsList(params, nil)
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
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("manufacturer_id", util.GetNestedManufacturerID(resource.Manufacturer)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("napalm_args", resource.NapalmArgs); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("napalm_driver", resource.NapalmDriver); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("slug", resource.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("virtualmachine_count", resource.VirtualmachineCount); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimPlatformUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modifiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritablePlatform{}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		params.Description = description
		modifiedFields["description"] = description
	}
	if d.HasChange("manufacturer_id") {
		manufacturerID := int64(d.Get("manufacturer_id").(int))
		params.Manufacturer = &manufacturerID
		modifiedFields["manufacturer"] = manufacturerID
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("napalm_args") {
		napalmArgs := d.Get("napalm_args").(string)
		params.NapalmArgs = &napalmArgs
		modifiedFields["napalm_args"] = napalmArgs
	}
	if d.HasChange("napalm_driver") {
		napalmDriver := d.Get("napalm_driver").(string)
		params.NapalmDriver = napalmDriver
		modifiedFields["napalm_driver"] = napalmDriver
	}
	if d.HasChange("slug") {
		slug := d.Get("slug").(string)
		params.Slug = &slug
	}
	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}

	resource := dcim.NewDcimPlatformsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Dcim.DcimPlatformsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, platformRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxDcimPlatformRead(ctx, d, m)
}

func resourceNetboxDcimPlatformDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxDcimPlatformExists(d, m)
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

	resource := dcim.NewDcimPlatformsDeleteParams().WithID(id)
	if _, err := client.Dcim.DcimPlatformsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxDcimPlatformExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := dcim.NewDcimPlatformsListParams().WithID(&resourceID)
	resources, err := client.Dcim.DcimPlatformsList(params, nil)
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
