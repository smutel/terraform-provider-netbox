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

func ResourceNetboxIpamRouteTargets() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a Route Targets (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamRouteTargetsCreate,
		ReadContext:   resourceNetboxIpamRouteTargetsRead,
		UpdateContext: resourceNetboxIpamRouteTargetsUpdate,
		DeleteContext: resourceNetboxIpamRouteTargetsDelete,
		// Exists:        resourceNetboxIpamRouteTargetsExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this Route Targets was created.",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				StateFunc:   util.TrimString,
				Description: "Comments for this Route Targets (ipam module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this Route Targets (ipam module).",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this Route Targets was created.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The name of this Route Targets (ipam module).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "ID of the tenant where this Route Targets (ipam module) is attached.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this Route Targets (ipam module).",
			},
		},
	}
}

var routeTargetsRequiredFields = []string{
	"created",
	"last_updated",
	"name",
	"tags",
}

func resourceNetboxIpamRouteTargetsCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)

	comments := d.Get("comments").(string)
	name := d.Get("name").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableRouteTarget{
		Comments:     comments,
		CustomFields: &customFields,
		Description:  d.Get("description").(string),
		Name:         &name,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	if tenantID := int64(d.Get("tenant_id").(int)); tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	resource := ipam.NewIpamRouteTargetsCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamRouteTargetsCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxIpamRouteTargetsRead(ctx, d, m)
}

func resourceNetboxIpamRouteTargetsRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamRouteTargetsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamRouteTargetsList(params, nil)
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

	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", resource.Name); err != nil {
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

func resourceNetboxIpamRouteTargetsUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	modiefiedFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableRouteTarget{}

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

	if d.HasChange("name") {
		name := d.Get("name").(string)
		params.Name = &name
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

	resource := ipam.NewIpamRouteTargetsPartialUpdateParams().WithData(params)

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamRouteTargetsPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modiefiedFields, routeTargetsRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamRouteTargetsRead(ctx, d, m)
}

func resourceNetboxIpamRouteTargetsDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamRouteTargetsExists(d, m)
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

	resource := ipam.NewIpamRouteTargetsDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamRouteTargetsDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamRouteTargetsExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamRouteTargetsListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamRouteTargetsList(params, nil)
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
