package netbox

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxIpamAggregate() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamAggregateCreate,
		Read:   resourceNetboxIpamAggregateRead,
		Update: resourceNetboxIpamAggregateUpdate,
		Delete: resourceNetboxIpamAggregateDelete,
		Exists: resourceNetboxIpamAggregateExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"custom_field": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{"text", "integer", "boolean",
								"date", "url", "selection", "multiple"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			"prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 256),
			},
			"rir_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetboxIpamAggregateCreate(d *schema.ResourceData,
	m interface{}) error {
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
			return err
		}

		dateAddedFmt := strfmt.Date(dateAddedTime)
		newResource.DateAdded = &dateAddedFmt
	}

	resource := ipam.NewIpamAggregatesCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamAggregatesCreate(resource, nil)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxIpamAggregateRead(d, m)
}

func resourceNetboxIpamAggregateRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamAggregatesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamAggregatesList(params, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
				return err
			}

			var dateAdded string
			if resource.DateAdded == nil {
				dateAdded = ""
			} else {
				dateAdded = resource.DateAdded.String()
			}

			if err = d.Set("date_added", dateAdded); err != nil {
				return err
			}

			var description interface{}
			if resource.Description == "" {
				description = nil
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return err
			}

			if err = d.Set("prefix", resource.Prefix); err != nil {
				return err
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
				return err
			}

			if resource.Rir == nil {
				if err = d.Set("rir_id", nil); err != nil {
					return err
				}
			} else {
				if err = d.Set("rir_id", resource.Rir.ID); err != nil {
					return err
				}
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxIpamAggregateUpdate(d *schema.ResourceData,
	m interface{}) error {
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
				return err
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
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamAggregatesPartialUpdate(resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxIpamAggregateRead(d, m)
}

func resourceNetboxIpamAggregateDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamAggregateExists(d, m)
	if err != nil {
		return err
	}

	if !resourceExists {
		return nil
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource := ipam.NewIpamAggregatesDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamAggregatesDelete(resource, nil); err != nil {
		return err
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
