package ipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ccoveille/go-safecast"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/smutel/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/brief"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v8/netbox/internal/util"
)

func ResourceNetboxIpamVrf() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a vrf within Netbox.",
		CreateContext: resourceNetboxIpamVrfCreate,
		ReadContext:   resourceNetboxIpamVrfRead,
		UpdateContext: resourceNetboxIpamVrfUpdate,
		DeleteContext: resourceNetboxIpamVrfDelete,
		// Exists:				resourceNetboxIpamVrfExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this VRF.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this VRF was created.",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				StateFunc:   util.TrimString,
				Description: "Comments for this VRF.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, util.Const200),
				Description:  "The description of this VRF.",
			},
			"enforce_unique": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Prevent duplicate prefixes/IP addresses " +
					"within this VRF",
			},
			"export_targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
					Description: "One of the ID of exported vrf targets " +
						"attached to this VRF.",
				},
				Description: "Array of ID of exported vrf targets attached " +
					"to this VRF.",
			},
			"import_targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
					Description: "One of the ID of imported vrf targets " +
						"attached to this VRF.",
				},
				Description: "Array of ID of imported vrf targets " +
					"attached to this VRF.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this VRF was created.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const100),
				Description:  "The name of this VRF.",
			},
			"rd": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, util.Const21),
				Description:  "The Route Distinguisher (RFC 4364) of this VRF.",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the tenant where this VRF is attached.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this VRF.",
			},
		},
	}
}

func resourceNetboxIpamVrfCreate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
		nil, resourceCustomFields)

	tags := d.Get("tag").(*schema.Set).List()
	exportTargets := d.Get("export_targets").([]any)
	exportTargetsID := []int32{}
	importTargets := d.Get("import_targets").([]any)
	importTargetsID := []int32{}

	for _, id := range exportTargets {
		id32, err := safecast.ToInt32(id.(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		exportTargetsID = append(exportTargetsID, id32)
	}

	for _, id := range importTargets {
		id32, err := safecast.ToInt32(id.(int))
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		importTargetsID = append(importTargetsID, id32)
	}

	newResource := netbox.NewVRFRequestWithDefaults()
	newResource.SetComments(d.Get("comments").(string))
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(d.Get("description").(string))
	newResource.SetEnforceUnique(d.Get("enforce_unique").(bool))
	newResource.SetExportTargets(exportTargetsID)
	newResource.SetImportTargets(importTargetsID)
	newResource.SetName(d.Get("name").(string))
	newResource.SetRd(d.Get("rd").(string))
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
		b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
		if err != nil {
			return err
		}
		newResource.SetTenant(*b)
	}

	_, response, err := client.IpamAPI.IpamVrfsCreate(ctx).VRFRequest(
		*newResource).Execute()
	if response.StatusCode != util.Const201 && err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	var resourceID int32
	if resourceID, err = util.UnmarshalID(response.Body); resourceID == 0 {
		return util.GenerateErrorMessage(response, err)
	}

	d.SetId(fmt.Sprintf("%d", resourceID))
	return resourceNetboxIpamVrfRead(ctx, d, m)
}

func resourceNetboxIpamVrfRead(ctx context.Context, d *schema.ResourceData,
	m any) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	resource, response, err := client.IpamAPI.IpamVrfsRetrieve(
		ctx, int32(resourceID)).Execute()

	if response.StatusCode == util.Const404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("comments", resource.GetComments()); err != nil {
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

	if err = d.Set("enforce_unique",
		resource.GetEnforceUnique()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	exportTargets := resource.GetExportTargets()
	exportTargetsID := []int32{}
	importTargets := resource.GetImportTargets()
	importTargetsID := []int32{}

	for _, exportTarget := range exportTargets {
		exportTargetsID = append(exportTargetsID, exportTarget.GetId())
	}

	for _, importTarget := range importTargets {
		importTargetsID = append(importTargetsID, importTarget.GetId())
	}

	if err = d.Set("export_targets", exportTargetsID); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("import_targets", importTargetsID); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated",
		resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("name", resource.GetName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("rd", resource.GetRd()); err != nil {
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

func resourceNetboxIpamVrfUpdate(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return util.GenerateErrorMessage(nil,
			errors.New("Unable to convert ID into int64"))
	}
	resource := netbox.NewVRFRequestWithDefaults()

	// Required fields
	resource.SetName(d.Get("name").(string))

	if d.HasChange("comments") {
		resource.SetComments(d.Get("comments").(string))
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(
			stateCustomFields.(*schema.Set).List(),
			resourceCustomFields.(*schema.Set).List())
		resource.SetCustomFields(customFields)
	}

	if d.HasChange("description") {
		resource.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("enforce_unique") {
		resource.SetEnforceUnique(d.Get("enforce_unique").(bool))
	}

	if d.HasChange("export_targets") {
		exportTargets := d.Get("export_targets").([]any)
		exportTargetsID := []int32{}

		for _, id := range exportTargets {
			id32, err := safecast.ToInt32(id.(int))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			exportTargetsID = append(exportTargetsID, id32)
		}

		resource.SetExportTargets(exportTargetsID)
	}

	if d.HasChange("import_targets") {
		importTargets := d.Get("import_targets").([]any)
		importTargetsID := []int32{}

		for _, id := range importTargets {
			id32, err := safecast.ToInt32(id.(int))
			if err != nil {
				return util.GenerateErrorMessage(nil, err)
			}
			importTargetsID = append(importTargetsID, id32)
		}

		resource.SetImportTargets(importTargetsID)
	}

	if d.HasChange("rd") {
		resource.SetRd(d.Get("rd").(string))
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		resource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))
	}

	if d.HasChange("tenant_id") {
		if tenantID := d.Get("tenant_id").(int); tenantID != 0 {
			b, err := brief.GetBriefTenantRequestFromID(ctx, client, tenantID)
			if err != nil {
				return err
			}
			resource.SetTenant(*b)
		} else {
			resource.SetTenantNil()
		}
	}

	if _, response, err := client.IpamAPI.IpamVrfsUpdate(ctx,
		int32(resourceID)).VRFRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return resourceNetboxIpamVrfRead(ctx, d, m)
}

func resourceNetboxIpamVrfDelete(ctx context.Context,
	d *schema.ResourceData, m any) diag.Diagnostics {

	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamVrfExists(d, m)
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

	if response, err := client.IpamAPI.IpamVrfsDestroy(
		ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamVrfExists(d *schema.ResourceData,
	m any) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.ParseInt(d.Id(), util.Const10, util.Const32)
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamVrfsRetrieve(nil,
		int32(resourceID)).Execute()
	if err != nil && http.StatusCode == util.Const404 {
		return false, nil
	} else if err == nil && http.StatusCode == util.Const200 {
		return true, nil
	}

	return false, err
}
