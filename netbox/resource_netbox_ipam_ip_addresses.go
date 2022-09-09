package netbox

import (
	"context"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/ipam"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxIpamIPAddresses() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an IP address (ipam module) within Netbox.",
		CreateContext: resourceNetboxIpamIPAddressesCreate,
		ReadContext:   resourceNetboxIpamIPAddressesRead,
		UpdateContext: resourceNetboxIpamIPAddressesUpdate,
		DeleteContext: resourceNetboxIpamIPAddressesDelete,
		Exists:        resourceNetboxIpamIPAddressesExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "The IP address (with mask) used for this IP address (ipam module).",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this IP address (ipam module).",
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
				Description: "Existing custom fields to associate to this IP address (ipam module).",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 100),
				Description:  "The description of this IP address (ipam module).",
			},
			"dns_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_.]{1,255}$"),
					"Must be like ^[-a-zA-Z0-9_.]{1,255}$"),
				Description: "The DNS name of this IP address (ipam module).",
			},
			"nat_inside_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the NAT inside of this IP address (ipam module).",
			},
			"object_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the object where this resource is attached to.",
			},
			"object_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ValidateFunc: validation.StringInSlice([]string{
					VMInterfaceType, "dcim.interface"}, false),
				Description: "The object type among virtualization.vminterface or dcim.interface (empty by default).",
			},
			"primary_ip4": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Set this resource as primary IPv4 (false by default).",
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ValidateFunc: validation.StringInSlice([]string{"loopback",
					"secondary", "anycast", "vip", "vrrp", "hsrp", "glbp", "carp"},
					false),
				Description: "The role among loopback, secondary, anycast, vip, vrrp, hsrp, glbp, carp of this IP address (ipam module).",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"container", "active",
					"reserved", "deprecated", "dhcp"}, false),
				Description: "The status among of this IP address (ipam module) container, active, reserved, deprecated (active by default).",
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
				Description: "Existing tag to associate to this IP address (ipam module).",
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the tenant where this object is attached.",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vrf attached to this IP address (ipam module).",
			},
		},
	}
}

func resourceNetboxIpamIPAddressesCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	address := d.Get("address").(string)
	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := convertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	dnsName := d.Get("dns_name").(string)
	natInsideID := int64(d.Get("nat_inside_id").(int))
	objectID := int64(d.Get("object_id").(int))
	objectType := d.Get("object_type").(string)
	primaryIP4 := d.Get("primary_ip4").(bool)
	role := d.Get("role").(string)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()
	tenantID := int64(d.Get("tenant_id").(int))
	vrfID := int64(d.Get("vrf_id").(int))

	newResource := &models.WritableIPAddress{
		Address:      &address,
		CustomFields: &customFields,
		Description:  description,
		DNSName:      dnsName,
		Role:         role,
		Status:       status,
		Tags:         convertTagsToNestedTags(tags),
	}

	if natInsideID != 0 {
		newResource.NatInside = &natInsideID
	}

	var info InfosForPrimary
	if primaryIP4 && objectID != 0 {
		if objectType == VMInterfaceType {
			var err error
			info, err = getInfoForPrimary(m, objectID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if objectID != 0 {
		newResource.AssignedObjectID = &objectID
		newResource.AssignedObjectType = &objectType
	}

	if tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	if vrfID != 0 {
		newResource.Vrf = &vrfID
	}

	resource := ipam.NewIpamIPAddressesCreateParams().WithData(newResource)

	resourceCreated, err := client.Ipam.IpamIPAddressesCreate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	err = updatePrimaryStatus(client, info, resourceCreated.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetboxIpamIPAddressesRead(ctx, d, m)
}

func resourceNetboxIpamIPAddressesRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := ipam.NewIpamIPAddressesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamIPAddressesList(params, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			if err = d.Set("address", resource.Address); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("content_type", convertURIContentType(resource.URL)); err != nil {
				return diag.FromErr(err)
			}

			resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
			customFields := updateCustomFieldsFromAPI(resourceCustomFields, resource.CustomFields)

			if err = d.Set("custom_field", customFields); err != nil {
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

			var dnsName interface{}
			if resource.DNSName == "" {
				dnsName = nil
			} else {
				dnsName = resource.DNSName
			}

			if err = d.Set("dns_name", dnsName); err != nil {
				return diag.FromErr(err)
			}

			var natInsideID *int64
			natInsideID = nil
			if resource.NatInside != nil {
				natInsideID = &resource.NatInside.ID
			}

			if err = d.Set("nat_inside_id", natInsideID); err != nil {
				return diag.FromErr(err)
			}

			var assignedObjectID *int64
			assignedObjectID = nil
			if resource.AssignedObjectID != nil {
				assignedObjectID = resource.AssignedObjectID

				var info InfosForPrimary
				if *resource.AssignedObjectID != 0 {
					if *resource.AssignedObjectType == VMInterfaceType {
						var err error
						info, err = getInfoForPrimary(m, *resource.AssignedObjectID)
						if err != nil {
							return diag.FromErr(err)
						}

						if info.vmPrimaryIP4ID == resource.ID {
							if err = d.Set("primary_ip4", true); err != nil {
								return diag.FromErr(err)
							}
						} else {
							if err = d.Set("primary_ip4", false); err != nil {
								return diag.FromErr(err)
							}
						}
					}
				}
			} else if err = d.Set("primary_ip4", false); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("object_id", assignedObjectID); err != nil {
				return diag.FromErr(err)
			}

			objectType := resource.AssignedObjectType
			if objectType != nil {
				*objectType = VMInterfaceType

				if err = d.Set("object_type", *objectType); err != nil {
					return diag.FromErr(err)
				}
			} else {
				if err = d.Set("object_type", nil); err != nil {
					return diag.FromErr(err)
				}
			}

			var roleValue *string
			roleValue = nil
			if resource.Role != nil {
				roleValue = resource.Role.Value
			}
			if err = d.Set("role", roleValue); err != nil {
				return diag.FromErr(err)
			}

			var resourceStatus *string
			resourceStatus = nil
			if resource.Status != nil {
				resourceStatus = resource.Status.Value
			}
			if err = d.Set("status", resourceStatus); err != nil {
				return diag.FromErr(err)
			}

			if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
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

			var vrfID *int64
			vrfID = nil
			if resource.Vrf != nil {
				vrfID = &resource.Vrf.ID
			}
			if err = d.Set("vrf_id", vrfID); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceNetboxIpamIPAddressesUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableIPAddress{}
	// primary_ip4 := false

	// Required parameters
	address := d.Get("address").(string)
	params.Address = &address

	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := convertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		if description, exist := d.GetOk("description"); exist {
			params.Description = description.(string)
		} else {
			params.Description = " "
		}
	}

	if d.HasChange("dns_name") {
		if dnsName, exist := d.GetOk("dns_name"); exist {
			params.DNSName = dnsName.(string)
		} else {
			params.DNSName = " "
		}
	}

	if d.HasChange("nat_inside_id") {
		natInsideID := int64(d.Get("nat_inside_id").(int))
		if natInsideID != 0 {
			params.NatInside = &natInsideID
		}
	}

	if d.HasChange("object_id") || d.HasChange("object_type") {
		// primary_ip4 = true
		objectID := int64(d.Get("object_id").(int))
		params.AssignedObjectID = &objectID

		var objectType string
		if *params.AssignedObjectType == "" {
			objectType = VMInterfaceType
		} else {
			objectType = d.Get("object_type").(string)
		}
		*params.AssignedObjectType = objectType
	}

	if d.HasChange("role") {
		role := d.Get("role").(string)
		params.Role = role
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		params.Status = status
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	if d.HasChange("tenant_id") {
		tenantID := int64(d.Get("tenant_id").(int))
		if tenantID != 0 {
			params.Tenant = &tenantID
		}
	}

	if d.HasChange("vrf_id") {
		vrfID := int64(d.Get("vrf_id").(int))
		if vrfID != 0 {
			params.Vrf = &vrfID
		}
	}

	resource := ipam.NewIpamIPAddressesPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamIPAddressesPartialUpdate(resource, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	/*
	 *   if primary_ip4 || d.HasChange("primary_ip4") {
	 *     var info InfosForPrimary
	 *     objectID := int64(d.Get("object_id").(int))
	 *     objectType := d.Get("object_type").(string)
	 *     isPrimary := d.Get("primary_ip4").(bool)
	 *     if objectID != 0 {
	 *       if objectType == VMInterfaceType {
	 *         var err error
	 *         info, err = getInfoForPrimary(m, objectID)
	 *         if err != nil {
	 *           return diag.FromErr(err)
	 *         }
	 *       }
	 *     }
	 *
	 *     var ipID int64
	 *     ipID = 0
	 *     if isPrimary {
	 *       ipID, err = strconv.ParseInt(d.Id(), 10, 64)
	 *       if err != nil {
	 *         return diag.Errorf("Unable to convert ID into int64")
	 *       }
	 *     }
	 *     err = updatePrimaryStatus(client, info, ipID)
	 *     if err != nil {
	 *       return diag.FromErr(err)
	 *     }
	 *   }
	 */

	return resourceNetboxIpamIPAddressesRead(ctx, d, m)
}

func resourceNetboxIpamIPAddressesDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamIPAddressesExists(d, m)
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

	resource := ipam.NewIpamIPAddressesDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamIPAddressesDelete(resource, nil); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamIPAddressesExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := ipam.NewIpamIPAddressesListParams().WithID(&resourceID)
	resources, err := client.Ipam.IpamIPAddressesList(params, nil)
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
