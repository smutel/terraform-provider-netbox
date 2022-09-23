package netbox

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
)

func resourceNetboxIpamAggregate() *schema.Resource {
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
			"custom_field": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the existing custom field.",
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"text", "integer", "boolean",
								"date", "url", "selection", "multiple"}, false),
							Description: "Type of the existing custom field (text, integer, boolean, url, selection, multiple).",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the existing custom field.",
						},
					},
				},
				Description: "Existing custom fields to associate to this aggregate (ipam module).",
			},
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
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the existing tag.",
						},
						"slug": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Slug of the existing tag.",
						},
					},
				},
				Description: "Existing tag to associate to this aggregate (ipam module).",
			},
		},
	}
}

func resourceNetboxIpamAggregateCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	dateAdded := d.Get("date_added").(string)
	description := d.Get("description").(string)
	prefix := d.Get("prefix").(string)
	rirID := int64(d.Get("rir_id").(int))
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableAggregate{
		CustomFields: &customFields,
		Description:  description,
		Prefix:       &prefix,
		Tags:         convertTagsToNestedTags(tags),
	}

	if rirID != 0 {
		newResource.Rir = &rirID
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
			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

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

			var description interface{}
			if resource.Description == "" {
				description = nil
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("prefix", resource.Prefix); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return diag.FromErr(err)
			}

			if resource.Rir == nil {
				if err = d.Set("rir_id", nil); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("rir_id", resource.Rir.ID); err != nil {
					return diag.FromErr(err)
				}
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
	params := &models.WritableAggregate{}

	// Required parameters
	prefix := d.Get("prefix").(string)
	params.Prefix = &prefix

	rirID := int64(d.Get("rir_id").(int))
	params.Rir = &rirID

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
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
		}
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			params.Description = " "
		}
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	resource := ipam.NewIpamAggregatesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamAggregatesPartialUpdate(resource, nil)
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
