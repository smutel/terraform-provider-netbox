package ipam

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	netbox "github.com/netbox-community/go-netbox/v4"
	"github.com/smutel/terraform-provider-netbox/v7/netbox/internal/customfield"
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

func resourceNetboxIpamIPAddressesCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	var address string
	if stateaddress, ok := d.GetOk("address"); ok {
		address = stateaddress.(string)
	} else if prefixid, ok := d.GetOk("prefix"); ok {
		ip, errDiag := getNewAvailableIPForPrefix(client, ctx, int32(prefixid.(int)))
		if errDiag != nil {
			return errDiag
		}
		address = ip.GetAddress()
		if err := d.Set("address", address); err != nil {
			return diag.FromErr(err)
		}
	} else if rangeid, ok := d.GetOk("ip_range"); ok {
		ip, errDiag := getNewAvailableIPForIPRange(client, ctx, int32(rangeid.(int)))
		if errDiag != nil {
			return errDiag
		}
		address = ip.GetAddress()
		if err := d.Set("address", address); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return util.GenerateErrorMessage(nil, errors.New("exactly one of "+
			"(address, ip_range, prefix) must be specified"))
	}

	resourceCustomFields := d.Get("custom_field").(*schema.Set).List()
	customFields := customfield.ConvertCustomFieldsFromTerraformToAPI(nil, resourceCustomFields)
	description := d.Get("description").(string)
	dnsName := d.Get("dns_name").(string)
	role := d.Get("role").(string)
	status := d.Get("status").(string)
	tags := d.Get("tag").(*schema.Set).List()

	newResource := netbox.NewWritableIPAddressRequestWithDefaults()
	newResource.SetAddress(address)
	newResource.SetCustomFields(customFields)
	newResource.SetDescription(description)
	newResource.SetDnsName(dnsName)
	newResource.SetTags(tag.ConvertTagsToNestedTagRequest(tags))

	s, err := netbox.NewPatchedWritableIPAddressRequestStatusFromValue(status)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetStatus(*s)

	r, err := netbox.NewPatchedWritableIPAddressRequestRoleFromValue(role)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}
	newResource.SetRole(*r)

	if natInsideID := int32(d.Get("nat_inside_id").(int)); natInsideID != 0 {
		newResource.SetNatInside(natInsideID)
	}

	objectID := int32(0)
	objectID64 := int64(0)
	objectType := ""
	if d.Get("object_id").(int) != 0 {
		objectID = int32(d.Get("object_id").(int))
		objectID64 = int64(d.Get("object_id").(int))
		objectType = d.Get("object_type").(string)
		newResource.SetAssignedObjectId(objectID64)
		newResource.SetAssignedObjectType(objectType)
	}

	if tenantID := int32(d.Get("tenant_id").(int)); tenantID != 0 {
		newResource.SetTenant(tenantID)
	}

	if vrfID := int32(d.Get("vrf_id").(int)); vrfID != 0 {
		newResource.SetVrf(vrfID)
	}

	addressid := int32(0)

	requestAddress := []string{address}
	resource, response, err := client.IpamAPI.IpamIpAddressesList(ctx).Address(requestAddress).Execute()
	if response.StatusCode == 200 {
		addressid = resource.Results[0].GetId()
	}

	if addressid == 0 {
		resourceCreated, response, err := client.IpamAPI.IpamIpAddressesCreate(ctx).WritableIPAddressRequest(*newResource).Execute()
		if response.StatusCode != 201 && err != nil {
			return util.GenerateErrorMessage(response, err)
		}
		addressid = resourceCreated.GetId()
	} else {
		if _, response, err := client.IpamAPI.IpamIpAddressesUpdate(ctx, int32(addressid)).WritableIPAddressRequest(*newResource).Execute(); err != nil {
			return util.GenerateErrorMessage(response, err)
		}
	}

	d.SetId(fmt.Sprintf("%d", addressid))
	if primaryIP := d.Get("primary_ip4").(bool); primaryIP {
		if err := setPrimaryIP(ctx, client, addressid, objectID, objectType, true); err != nil {
			return err
		}
	}

	return resourceNetboxIpamIPAddressesRead(ctx, d, m)
}

func resourceNetboxIpamIPAddressesRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, _ := strconv.Atoi(d.Id())
	resource, response, err := client.IpamAPI.IpamIpAddressesRetrieve(ctx, int32(resourceID)).Execute()

	if err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if err = d.Set("content_type", util.ConvertURLContentType(resource.GetUrl())); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("address", resource.GetAddress()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("object_id", resource.GetAssignedObjectId()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("object_type", resource.GetAssignedObjectType()); err != nil {
		return diag.FromErr(err)
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

	if err = d.Set("dns_name", resource.GetDnsName()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("family", resource.GetFamily().Value); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("last_updated", resource.GetLastUpdated().String()); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("nat_inside_id", resource.GetNatInside().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	// if err = d.Set("nat_outside", resource.GetNatOutside().Id); err != nil {
	// return util.GenerateErrorMessage(nil, err)
	// }

	isPrimary, errDiag := isprimary(ctx, client, resource.GetAssignedObjectId(), resource.GetId(), (*resource.GetFamily().Value == 4))
	if err != nil {
		return errDiag
	}

	if err = d.Set("primary_ip4", isPrimary); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if err = d.Set("role", resource.GetRole().Value); err != nil {
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

	if err = d.Set("vrf_id", resource.GetVrf().Id); err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	return nil
}

func resourceNetboxIpamIPAddressesUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}
	resource := netbox.NewWritableIPAddressRequestWithDefaults()

	// Required parameters
	if d.HasChange("address") {
		resource.SetAddress(d.Get("address").(string))
	}

	if d.HasChange("object_id") || d.HasChange("object_type") {
		objectID := int64(d.Get("object_id").(int))
		objectType := d.Get("object_type").(string)
		resource.SetAssignedObjectId(objectID)
		resource.SetAssignedObjectType(objectType)
	}

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

	if d.HasChange("dns_name") {
		resource.SetDnsName(d.Get("dns_name").(string))
	}

	if d.HasChange("nat_inside_id") {
		natInsideID := int32(d.Get("nat_inside_id").(int))
		if natInsideID != 0 {
			resource.SetNatInside(natInsideID)
		} else {
			resource.SetNatInsideNil()
		}
	}

	if d.HasChange("role") {
		role := d.Get("role").(string)
		r, err := netbox.NewPatchedWritableIPAddressRequestRoleFromValue(role)
		if err != nil {
			return util.GenerateErrorMessage(nil, err)
		}
		resource.SetRole(*r)
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		s, err := netbox.NewPatchedWritableIPAddressRequestStatusFromValue(status)
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

	if d.HasChange("vrf_id") {
		vrfID := int32(d.Get("vrf_id").(int))
		if vrfID != 0 {
			resource.SetVrf(vrfID)
		} else {
			resource.SetVrfNil()
		}
	}

	if _, response, err := client.IpamAPI.IpamIpAddressesUpdate(ctx, int32(resourceID)).WritableIPAddressRequest(*resource).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	if !d.GetRawConfig().GetAttr("primary_ip4").IsNull() {
		objectChanged := d.HasChange("object_id") || d.HasChange("object_type")
		primaryIP4 := d.Get("primary_ip4").(bool)

		if (objectChanged && primaryIP4) ||
			(!objectChanged && d.HasChange("primary_ip4")) ||
			(d.HasChange("primary_ip4") && primaryIP4) {
			objectID := int32(d.Get("object_id").(int))
			objectType := d.Get("object_type").(string)

			errDiag := setPrimaryIP(ctx, client, int32(resourceID), objectID, objectType, primaryIP4)
			if errDiag != nil {
				return errDiag
			}
		}
	}

	return resourceNetboxIpamIPAddressesRead(ctx, d, m)
}

func resourceNetboxIpamIPAddressesDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	client := m.(*netbox.APIClient)

	resourceExists, err := resourceNetboxIpamIPAddressesExists(d, m)
	if err != nil {
		return util.GenerateErrorMessage(nil, err)
	}

	if !resourceExists {
		return nil
	}

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return util.GenerateErrorMessage(nil, errors.New("Unable to convert ID into int"))
	}

	if response, err := client.IpamAPI.IpamIpAddressesDestroy(ctx, int32(resourceID)).Execute(); err != nil {
		return util.GenerateErrorMessage(response, err)
	}

	return nil
}

func resourceNetboxIpamIPAddressesExists(d *schema.ResourceData,
	m interface{}) (b bool, e error) {
	client := m.(*netbox.APIClient)

	resourceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, err
	}

	_, http, err := client.IpamAPI.IpamIpAddressesRetrieve(nil, int32(resourceID)).Execute()
	if err != nil && http.StatusCode == 404 {
		return false, nil
	} else if err == nil && http.StatusCode == 200 {
		return true, nil
	} else {
		return false, err
	}
}
