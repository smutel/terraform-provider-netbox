package ipam

import (
	"context"
	"regexp"
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

func ResourceNetboxIpamIPAddresses() *schema.Resource {
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
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"address", "prefix", "ip_range"},
				ValidateFunc: validation.IsCIDR,
				Description:  "The IP address (with mask) used for this IP address (ipam module). Required if both prefix and ip_range are not set.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this aggregate was created.",
			},
			"content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content type of this IP address (ipam module).",
			},
			"custom_field": &customfield.CustomFieldSchema,
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      nil,
				ValidateFunc: validation.StringLenBetween(1, 200),
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
			"family": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Family of IP address (IPv4 or IPv6).",
			},
			"ip_range": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Description: "The ip-range id for automatic IP assignment. Required if both prefix and address are not set.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when this aggregate was last updated.",
			},
			"nat_inside_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the NAT inside of this IP address (ipam module).",
			},
			// "nat_outside": {
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "The IDs of the NAT outside of this IP address (ipam module).",
			// 	Elem:        schema.TypeInt,
			// },
			"object_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The ID of the object where this resource is attached to.",
				RequiredWith: []string{"object_type"},
			},
			"object_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ValidateFunc: validation.StringInSlice([]string{
					vMInterfaceType, deviceInterfaceType, fhrpgroupType}, false),
				Description: "The object type among virtualization.vminterface, dcim.interface or ipam.fhrpgroup (empty by default).",
			},
			"prefix": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Description: "The prefix id for automatic IP assignment. Required if both address and ip_range are not set.",
			},
			"primary_ip4": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Deprecated:  "Use new netbox_virtualization_primary_ip resource instead",
				Description: "Set this resource as primary IPv4 (false by default).",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return d.GetRawConfig().GetAttr("primary_ip4").IsNull()
				},
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
				ValidateFunc: validation.StringInSlice([]string{"active",
					"reserved", "deprecated", "dhcp", "slaac"}, false),
				Description: "The status among of this IP address (ipam module) active, reserved, deprecated, dhcp, slaac (active by default).",
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
			"vrf_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the vrf attached to this IP address (ipam module).",
			},
		},
	}
}

var ipAddressRequiredFields = []string{
	"address",
	"created",
	"last_updated",
	"tags",
}

func resourceNetboxIpamIPAddressesCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)

	var address string
	var addressid *int64
	if stateaddress, ok := d.GetOk("address"); ok {
		address = stateaddress.(string)
	} else if prefixid, ok := d.GetOk("prefix"); ok {
		ip, err := getNewAvailableIPForPrefix(client, int64(prefixid.(int)))
		if err != nil {
			return diag.FromErr(err)
		}
		address = *ip.Address
		addressid = &ip.ID
		if err := d.Set("address", address); err != nil {
			return diag.FromErr(err)
		}
	} else if rangeid, ok := d.GetOk("ip_range"); ok {
		ip, err := getNewAvailableIPForIPRange(client, int64(rangeid.(int)))
		if err != nil {
			return diag.FromErr(err)
		}
		address = *ip.Address
		addressid = &ip.ID
		if err := d.Set("address", address); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("exactly one of (address, ip_range, prefix) must be specified")
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	dnsName := d.Get("dns_name").(string)
	role := d.Get("role").(string)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := &models.WritableIPAddress{
		Address:      &address,
		CustomFields: &customFields,
		Description:  description,
		DNSName:      dnsName,
		Role:         role,
		Status:       status,
		Tags:         tag.ConvertTagsToNestedTags(tags),
	}

	if natInsideID := int64(d.Get("nat_inside_id").(int)); natInsideID != 0 {
		newResource.NatInside = &natInsideID
	}

	objectID := int64(0)
	objectType := ""
	if d.Get("object_id").(int) != 0 {
		objectID = int64(d.Get("object_id").(int))
		objectType = d.Get("object_type").(string)
		newResource.AssignedObjectID = &objectID
		newResource.AssignedObjectType = &objectType
	}

	if tenantID := int64(d.Get("tenant_id").(int)); tenantID != 0 {
		newResource.Tenant = &tenantID
	}

	if vrfID := int64(d.Get("vrf_id").(int)); vrfID != 0 {
		newResource.Vrf = &vrfID
	}

	if addressid == nil {
		resource := ipam.NewIpamIPAddressesCreateParams().WithData(newResource)

		resourceCreated, err := client.Ipam.IpamIPAddressesCreate(resource, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		addressid = &resourceCreated.Payload.ID
	} else {
		resource := ipam.NewIpamIPAddressesUpdateParams().WithID(*addressid).WithData(newResource)

		_, err := client.Ipam.IpamIPAddressesUpdate(resource, nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.FormatInt(*addressid, 10))
	if primaryIP := d.Get("primary_ip4").(bool); primaryIP {
		if err := setPrimaryIP(m, *addressid, objectID, objectType, true); err != nil {
			return diag.FromErr(err)
		}
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

	if len(resources.Payload.Results) != 1 {
		d.SetId("")
		return nil
	}

	resource := resources.Payload.Results[0]

	if err = d.Set("content_type", util.ConvertURIContentType(resource.URL)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("address", resource.Address); err != nil {
		return diag.FromErr(err)
	}

	objectType, objectID := util.GetIPAddressAssignedObject(resource)
	if err = d.Set("object_id", objectID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("object_type", objectType); err != nil {
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
	if err = d.Set("dns_name", resource.DNSName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("family", util.GetIPAddressFamilyLabel(resource.Family)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", resource.LastUpdated.String()); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("nat_inside_id", util.GetNestedIPAddressID(resource.NatInside)); err != nil {
		return diag.FromErr(err)
	}
	// if err = d.Set("nat_outside", util.ConvertNestedIPsToIPs(resource.NatOutside)); err != nil {
	// 	return diag.FromErr(err)
	// }

	isPrimary, err := isprimary(m, resource.AssignedObjectID, resource.ID, (*resource.Family.Value == 4))
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("primary_ip4", isPrimary); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("role", util.GetIPAddressRoleValue(resource.Role)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("status", util.GetIPAddressStatusValue(resource.Status)); err != nil {
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
	if err = d.Set("vrf_id", util.GetNestedVrfID(resource.Vrf)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetboxIpamIPAddressesUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableIPAddress{}
	modifiedFields := map[string]interface{}{}

	// Required parameters
	if d.HasChange("address") {
		address := d.Get("address").(string)
		params.Address = &address
	}
	if d.HasChange("object_id") || d.HasChange("object_type") {
		objectID := int64(d.Get("object_id").(int))
		objectType := d.Get("object_type").(string)
		params.AssignedObjectID = &objectID
		params.AssignedObjectType = &objectType
		modifiedFields["assigned_object_id"] = objectID
		modifiedFields["assigned_object_type"] = objectType
	}
	if d.HasChange("custom_field") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_field")
		customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(stateCustomFields.(*schema.Set).List(), resourceCustomFields.(*schema.Set).List())
		params.CustomFields = &customFields
	}

	if d.HasChange("description") {
		params.Description = d.Get("description").(string)
		modifiedFields["description"] = params.Description
	}
	if d.HasChange("dns_name") {
		params.DNSName = d.Get("dns_name").(string)
		modifiedFields["dns_name"] = params.DNSName
	}
	if d.HasChange("nat_inside_id") {
		natInsideID := int64(d.Get("nat_inside_id").(int))
		params.NatInside = &natInsideID
		modifiedFields["nat_inside"] = natInsideID
	}
	if d.HasChange("role") {
		role := d.Get("role").(string)
		params.Role = role
		modifiedFields["role"] = role
	}
	if d.HasChange("status") {
		status := d.Get("status").(string)
		params.Status = status
		modifiedFields["status"] = status
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
	if d.HasChange("vrf_id") {
		vrfID := int64(d.Get("vrf_id").(int))
		params.Vrf = &vrfID
		modifiedFields["vrf"] = vrfID
	}

	resource := ipam.NewIpamIPAddressesPartialUpdateParams().WithData(
		params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamIPAddressesPartialUpdate(resource, nil, requestmodifier.NewNetboxRequestModifier(modifiedFields, ipAddressRequiredFields))
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.GetRawConfig().GetAttr("primary_ip4").IsNull() {
		objectChanged := d.HasChange("object_id") || d.HasChange("object_type")
		primaryIP4 := d.Get("primary_ip4").(bool)

		if (objectChanged && primaryIP4) ||
			(!objectChanged && d.HasChange("primary_ip4")) ||
			(d.HasChange("primary_ip4") && primaryIP4) {
			objectID := int64(d.Get("object_id").(int))
			objectType := d.Get("object_type").(string)

			err = setPrimaryIP(client, resourceID, objectID, objectType, primaryIP4)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

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
