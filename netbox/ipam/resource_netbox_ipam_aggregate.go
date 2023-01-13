package ipam

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/models"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/customfield"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/requestmodifier"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/tag"
	"github.com/smutel/terraform-provider-netbox/v6/netbox/internal/util"
)

func ResourceNetboxIpamAggregate() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an aggregate (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamAggregateCreate,
		ReadContext:   resourceNetboxIpamAggregateRead,
		UpdateContext: resourceNetboxIpamAggregateUpdate,
		DeleteContext: resourceNetboxIpamAggregateDelete,
		Exists:        resourceNetboxIpamAggregateExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this aggregate (ipam module).",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this aggregate was created.",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"date_added": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := time.Parse("2006-01-02", v)

					if err != nil {
						errs = append(errs, fmt.Errorf("date_added in not in the good format YYYY-MM-DD"))
					}
					return
				},
				Description: "Date when this aggregate was added. Format *YYYY-MM-DD*.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "The description of this aggregate (ipam module).",
			},
			"family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP family of this aggregate.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this aggregate was last updated.",
			},
			"prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 256),
				Description:  "The network prefix of this aggregate (ipam module).",
			},
			"rir_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The RIR id linked to this aggregate (ipam module).",
			},
			"tag": &tag.TagSchema,
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this object is attached.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link to this tag (extra module).",
			},
		},
	}
}

func resourceNetboxIpamAggregateCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	dateAdded := d.Get("date_added").(string)
	description := d.Get("description").(string)
	prefix := d.Get("prefix").(string)
	rirID := int64(d.Get("rir_id").(int))
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableAggregate{
		CustomFields: &customFields,
		Description:  description,
		Prefix:       &prefix,
		Rir:          &rirID,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	if tenantID := int64(d.Get("tenant_id").(int)); tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	if dateAdded != "" {
		dateAddedTime, err := time.Parse("2006-01-02", dateAdded)
		if err != nil {
			return diag.FromErr(err)
		}

		dateAddedFmt := strfmt.Date(dateAddedTime)
		newResource.DateAdded = &dateAddedFmt
	}

	resource := ipam.NewIpamAggregatesCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamAggregatesCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxIpamAggregateRead(ctx, d, m)
}

func resourceNetboxIpamAggregateRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamAggregatesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamAggregatesList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
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

			var dateAdded string
			if resource.DateAdded == nil {
				dateAdded = ""
			} else {
				dateAdded = resource.DateAdded.String()
			}

			if err = d.Set("date_added", dateAdded); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("description", resource.Description); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("family", resource.Family.Label); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("prefix", resource.Prefix); err != nil {
				return diag.FromErr(err)
			}

			var rirID *int64
			rirID = nil
			if resource.Rir != nil {
				rirID = &resource.Rir.ID
			}
			if err = d.Set("rir_id", rirID); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("tag", tag.ConvertNestedTagsToTags(resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			var tenantID *int64
			tenantID = nil
			if resource.Tenant != nil {
				tenantID = &resource.Tenant.ID
			}
			if err = d.Set("tenant_id", tenantID); err != nil {
				return diag.FromErr(err)
			}
			if err = d.Set("url", resource.URL); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamAggregateUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	dropFields := []string{
		"created",
		"last_updated",
	}
	emptyFields := make(map[string]interface{})

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}
	params := &models.WritableAggregate{}

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}
	if d.HasChange("date_added") {
		dateAdded := d.Get("date_added").(string)

		if dateAdded != "" {
			dateAddedTime, err := time.Parse("2006-01-02", dateAdded)
			if err != nil {
				return diag.FromErr(err)
			}

			dateAddedFmt := strfmt.Date(dateAddedTime)
			params.DateAdded = &dateAddedFmt
		} else {
			emptyFields["date_added"] = nil
		}
	}
	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			emptyFields["description"] = ""
		}
	}
	if d.HasChange("prefix") {
		prefix := d.Get("prefix").(string)
		params.Prefix = &prefix
	} else {
		dropFields = append(dropFields, "prefix")
	}
	if d.HasChange("rir_id") {
		rirID := int64(d.Get("rir_id").(int))
		params.Rir = &rirID
	} else {
		dropFields = append(dropFields, "rir")
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = tag.ConvertTagsToNestedTags(tags)

	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		if tenantID != 0 {
			params.Tenant = &tenantID
		} else {
			emptyFields["tenant"] = nil
		}
	}

	resource := ipam.NewIpamAggregatesPartialUpdateParams().WithData(params)
	resource.SetID(resourceID)

	_, err = client.Ipam.IpamAggregatesPartialUpdate(resource, nil, requestmodifier.NewRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamAggregateRead(ctx, d, m)
}

func resourceNetboxIpamAggregateDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamAggregateExists(d, m)
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

	resource := ipam.NewIpamAggregatesDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamAggregatesDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamAggregateExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamAggregatesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamAggregatesList(params, nil)
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
