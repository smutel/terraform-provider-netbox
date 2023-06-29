package virtualization

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/virtualization"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/util"
)

func ResourceNetboxVirtualizationCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a tag (extra module) within Netbox.",
		CreateContext: resourceNetboxVirtualizationClusterCreate,
		ReadContext:   resourceNetboxVirtualizationClusterRead,
		UpdateContext: resourceNetboxVirtualizationClusterUpdate,
		DeleteContext: resourceNetboxVirtualizationClusterDelete,
		Exists:        resourceNetboxVirtualizationClusterExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this cluster (virtualization module).",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   util.TrimString,
				Description: "Comments for this cluster (virtualization module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this cluster was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of devices in this cluster.",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The cluster group of this cluster.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this cluster was last updated.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this cluster (virtualization module).",
			},
			"site_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The site of this cluster.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"offline", "active",
					"planned", "staging", "decommissioning"}, false),
				Description: "The status among offline, active, planned, staging or decommissioning (active by default).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the tenant where this cluster is attached.",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Default:     nil,
				Description: "Type of this cluster.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this cluster (virtualization module).",
			},
			"virtualmachine_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of virtual machines in this cluster.",
			},
		},
	}
}

var clusterRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"type",
	"tags",
}

func resourceNetboxVirtualizationClusterCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	groupID := int64(d.Get("group_id").(int))
	name := d.Get("name").(string)
	typeID := int64(d.Get("type_id").(int))
	siteID := int64(d.Get("site_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableCluster{
		Comments:     d.Get("comments").(string),
		CustomFields: customFields,
		Name:         &name,
		Tags:         tag.ConvertTagsToNestedTags(tags),
		Type:         &typeID,
		Status:       d.Get("status").(string),
	}

	if groupID != 0 {
		newResource.Group = &groupID
	}

	if siteID != 0 {
		newResource.Site = &siteID
	}
	if tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	resource := virtualization.NewVirtualizationClustersCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationClustersCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationClusterRead(ctx, d, m)
}

func resourceNetboxVirtualizationClusterRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationClustersListParams().WithID(&resourceID)
	resources, err := client.Virtualization.VirtualizationClustersList(params, nil)
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
	if err = d.Set("device_count", resource.DeviceCount); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("group_id", util.GetNestedClusterGroupID(resource.Group)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resource.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("site_id", util.GetNestedSiteID(resource.Site)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("status", util.GetClusterStatusValue(resource.Status)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tenant_id", util.GetNestedTenantID(resource.Tenant)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("type_id", resource.Type.ID); err != nil {
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

func resourceNetboxVirtualizationClusterUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modifiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableCluster{}

	if d.HasChange("comments") {
		params.Comments = d.Get("comments").(string)
		modifiedFields["comments"] = params.Comments
	}
	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("group_id") {
		groupID := int64(d.Get("group_id").(int))
		params.Group = &groupID
		modifiedFields["group"] = groupID
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
	}
	if d.HasChange("status") {
		params.Status = d.Get("status").(string)
	}
	if d.HasChange("site_id") {
		siteID := int64(d.Get("site_id").(int))
		params.Site = &siteID
		modifiedFields["site"] = siteID
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
	if d.HasChange("type_id") {
		typeID := int64(d.Get("type_id").(int))
		params.Type = &typeID
	}

	resource := virtualization.NewVirtualizationClustersPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationClustersPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, clusterRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxVirtualizationClusterRead(ctx, d, m)
}

func resourceNetboxVirtualizationClusterDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationClusterExists(d, m)
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

	resource := virtualization.NewVirtualizationClustersDeleteParams().WithID(id)
	if _, err := client.Virtualization.VirtualizationClustersDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxVirtualizationClusterExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := virtualization.NewVirtualizationClustersListParams().WithID(&resourceID)
	resources, err := client.Virtualization.VirtualizationClustersList(params, nil)
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
