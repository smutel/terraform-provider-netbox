package dcim

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

func ResourceNetboxDcimLocation() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a location (dcim module) within Netbox.",
		CreateContext: resourceNetboxDcimLocationCreate,
		ReadContext:   resourceNetboxDcimLocationRead,
		UpdateContext: resourceNetboxDcimLocationUpdate,
		DeleteContext: resourceNetboxDcimLocationDelete,
		Exists:        resourceNetboxDcimLocationExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this location (dcim module) was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"depth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Depth of this location (dcim module).",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
				Description:  "Description of this location (dcim module).",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of devices in this location (dcim module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this location was last updated (dcim module).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this location (dcim module).",
			},
			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the parent for this location (dcim module).",
			},
			"rack_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of racks in this location (dcim module).",
			},
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The slug of this site (dcim module).",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The site where this location (dcim module) is.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				ValidateFunc: validation.StringInSlice([]string{"planned", "staging", "active", "decommissioning", "retired"}, false),
				Description:  "The status among planned, staging, active, decommissioning or retired (active by default) of this location (dcim module).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The tenant of this location (dcim module).",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this location (dcim module).",
			},
		},
	}
}

func resourceNetboxDcimLocationCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	name := d.Get("name").(string)
	parentID := int32(d.Get("parent_id").(int))
	siteID := int32(d.Get("site_id").(int))
	slug := d.Get("slug").(string)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int32(d.Get("tenant_id").(int))

	newResource := netbox.NewWritableLocationRequestWithDefaults()
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetName(name)
	newResource.SetSite(siteID)
	newResource.SetSlug(slug)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	s, err := netbox.NewLocationStatusValueFromValue(status)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetStatus(*s)

	if parentID != 0 {
		newResource.SetParent(parentID)
	}

	if tenantID != 0 {
		newResource.SetTenant(tenantID)
	}

	resourceCreated, response, err := client.DcimAPI.DcimLocationsCreate(ctx).WritableLocationRequest(*newResource).Execute()
	if response.StatusCode != 201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	// NETBOX BUG - TO BE FIXED
	if resourceCreated.GetId() == 0 {
		return diag.FromErr(errors.New("Bug Netbox - TO BE FIXED"))
	}

	d.SetId(fmt.Sprintf("%d", resourceCreated.GetId()))

	return resourceNetboxDcimLocationRead(ctx, d, m)
}

func resourceNetboxDcimLocationRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.DcimAPI.DcimLocationsRetrieve(ctx, int32(resourceID)).Execute()

	if response.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
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

	if err = d.Set("depth", resource.GetDepth()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("parent_id", resource.GetParent()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("rack_count", resource.GetRackCount()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("site_id", resource.GetSite()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("slug", resource.GetSlug()); err != nil {
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

	if err = d.Set("url", resource.GetUrl()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxDcimLocationUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	resource := netbox.NewWritableLocationRequestWithDefaults()

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

	if d.HasChange("name") {
		resource.SetName(d.Get("name").(string))
	}

	if d.HasChange("parent_id") {
		parentID := int32(d.Get("parent_id").(int))
		if parentID != 0 {
			resource.SetParent(parentID)
		} else {
			resource.SetParentNil()
		}
	}

	if d.HasChange("site_id") {
		siteID := int32(d.Get("site_id").(int))
		resource.SetSite(siteID)
	}

	if d.HasChange("slug") {
		resource.SetSlug(d.Get("slug").(string))
	}

	if d.HasChange("status") {
		s, err := netbox.NewLocationStatusValueFromValue(d.Get("status").(string))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetStatus(*s)
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

	if _, response, err := client.DcimAPI.DcimLocationsUpdate(ctx, int32(resourceID)).WritableLocationRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxDcimLocationRead(ctx, d, m)
}

func resourceNetboxDcimLocationDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxDcimLocationExists(d, m)
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

	if response, err := client.DcimAPI.DcimLocationsDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxDcimLocationExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.DcimAPI.DcimLocationsRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
