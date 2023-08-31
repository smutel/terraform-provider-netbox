package ipam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxIpamVrf() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a vrf (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamVrfCreate,
		ReadContext:   resourceNetboxIpamVrfRead,
		UpdateContext: resourceNetboxIpamVrfUpdate,
		DeleteContext: resourceNetboxIpamVrfDelete,
		// Exists:        resourceNetboxIpamVrfExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this VRF (ipam module).",
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
				Description: "Comments for this VRF (ipam module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this VRF (ipam module).",
			},
			"enforce_unique": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Prevent duplicate prefixes/IP addresses within this VRF (ipam module)",
			},
			"export_targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Array of ID of exported vrf targets attached to this VRF (ipam module).",
			},
			"import_targets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Array of ID of imported vrf targets attached to this VRF (ipam module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this VRF was created.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this VRF (ipam module).",
			},
			"rd": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 21),
				Description:  "The Route Distinguisher (RFC 4364) of this VRF (ipam module).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the tenant where this VRF (ipam module) is attached.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this VRF (ipam module).",
			},
		},
	}
}

var vrfRequiredFields = []string{
	"created",
	"enforce_unique",
	"last_updated",
	"name",
	"tags",
}

func resourceNetboxIpamVrfCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)

	name := d.Get("name").(string)
	rd := d.Get("rd").(string)
	tags := d.Get("tag").(*schema.Set).List()
	exportTargets := d.Get("export_targets").([]interface{})
	exportTargetsID64 := []int64{}
	importTargets := d.Get("import_targets").([]interface{})
	importTargetsID64 := []int64{}

	for _, id := range exportTargets {
		exportTargetsID64 = append(exportTargetsID64, int64(id.(int)))
	}

	for _, id := range importTargets {
		importTargetsID64 = append(importTargetsID64, int64(id.(int)))
	}

	newResource := &models.WritableVRF{
		Comments:      d.Get("comments").(string),
		CustomFields:  &customFields,
		Description:   d.Get("description").(string),
		EnforceUnique: d.Get("enforce_unique").(bool),
		ExportTargets: exportTargetsID64,
		ImportTargets: importTargetsID64,
		Name:          &name,
		Rd:            &rd,
		Tags:          tag.ConvertTagsToNestedTags(tags),
	}

	if tenantID := int64(d.Get("tenant_id").(int)); tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	resource := ipam.NewIpamVrfsCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamVrfsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxIpamVrfRead(ctx, d, m)
}

func resourceNetboxIpamVrfRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamVrfsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVrfsList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]

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

	if err = d.Set("enforce_unique", resource.EnforceUnique); err != nil {
		return diag.FromErr(err)
	}

	exportTargetsObject := resource.ExportTargets
	exportTargetsInt := []int64{}
	for _, ip := range exportTargetsObject {
		exportTargetsInt = append(exportTargetsInt, ip.ID)
	}
	if err = d.Set("export_targets", exportTargetsInt); err != nil {
		return diag.FromErr(err)
	}

	importTargetsObject := resource.ImportTargets
	importTargetsInt := []int64{}
	for _, ip := range importTargetsObject {
		importTargetsInt = append(importTargetsInt, ip.ID)
	}
	if err = d.Set("import_targets", importTargetsInt); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("rd", resource.Rd); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("tenant_id", util.GetNestedTenantID(resource.Tenant)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("url", resource.URL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamVrfUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modiefiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableVRF{}

	if d.HasChange("comments") {
		comments := d.Get("comments")
		if comments != "" {
			params.Comments = comments.(string)
		} else {
			modiefiedFields["comments"] = ""
		}
	}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
		modiefiedFields["description"] = params.Description
	}

	if d.HasChange("enforce_unique") {
		params.EnforceUnique = d.Get("enforce_unique").(bool)
		modiefiedFields["enforce_unique"] = params.EnforceUnique
	}

	exportTargets := d.Get("export_targets").([]interface{})
	exportTargetsID64 := []int64{}
	for _, id := range exportTargets {
		exportTargetsID64 = append(exportTargetsID64, int64(id.(int)))
	}

	params.ExportTargets = exportTargetsID64

	importTargets := d.Get("import_targets").([]interface{})
	importTargetsID64 := []int64{}
	for _, id := range importTargets {
		importTargetsID64 = append(importTargetsID64, int64(id.(int)))
	}

	params.ImportTargets = importTargetsID64

	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}

	if d.HasChange("rd") {
		rd := d.Get("rd").(string)
		params.Rd = &rd
	}

	if d.HasChange("tag") {
		tags := d.Get("tag").(*schema.Set).List()
		params.Tags = tag.ConvertTagsToNestedTags(tags)
	}

	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		if tenantID != 0 {
			params.Tenant = &tenantID
		} else {
			modiefiedFields["tenant"] = nil
		}
	}

	resource := ipam.NewIpamVrfsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamVrfsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modiefiedFields, vrfRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamVrfRead(ctx, d, m)
}

func resourceNetboxIpamVrfDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamVrfExists(d, m)
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

	resource := ipam.NewIpamVrfsDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamVrfsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamVrfExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamVrfsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamVrfsList(params, nil)
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
