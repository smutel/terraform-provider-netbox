package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxIpamIPAddresses() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamIPAddressesCreate,
		Read:   resourceNetboxIpamIPAddressesRead,
		Update: resourceNetboxIpamIPAddressesUpdate,
		Delete: resourceNetboxIpamIPAddressesDelete,
		Exists: resourceNetboxIpamIPAddressesExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.IsCIDR,
			},
			"custom_fields": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				// terraform default behavior sees a difference between null and an empty string
				// therefore we override the default, because null in terraform results in empty string in netbox
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// function is called for each member of map
					// including additional call on the amount of entries
					// we ignore the count, because the actual state always returns the amount of existing custom_fields and all are optional in terraform
					if k == CustomFieldsRegex {
						return true
					}
					return old == new
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      " ",
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"dns_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  " ",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9_.]{1,255}$"),
					"Must be like ^[-a-zA-Z0-9_.]{1,255}$"),
			},
			"nat_inside_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nat_outside_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"object_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"object_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  VMInterfaceType,
				ValidateFunc: validation.StringInSlice([]string{
					VMInterfaceType, "dcim.interface"}, false),
			},
			"primary_ip4": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ValidateFunc: validation.StringInSlice([]string{"loopback",
					"secondary", "anycast", "vip", "vrrp", "hsrp", "glbp", "carp"},
					false),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active",
				ValidateFunc: validation.StringInSlice([]string{"container", "active",
					"reserved", "deprecated", "dhcp"}, false),
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
			"tenant_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vrf_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNetboxIpamIPAddressesCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	address := d.Get("address").(string)
	resourceCustomFields := d.Get("custom_fields").(map[string]interface{})
	customFields := convertCustomFieldsFromTerraformToAPICreate(resourceCustomFields)
	description := d.Get("description").(string)
	dnsName := d.Get("dns_name").(string)
	natInsideID := int64(d.Get("nat_inside_id").(int))
	natOutsideID := int64(d.Get("nat_outside_id").(int))
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

	if natOutsideID != 0 {
		newResource.NatOutside = &natOutsideID
	}

	var info InfosForPrimary
	if primaryIP4 && objectID != 0 {
		if objectType == VMInterfaceType {
			var err error
			info, err = getInfoForPrimary(m, objectID)
			if err != nil {
				return err
			}
		}
	}

	if objectID != 0 {
		newResource.AssignedObjectID = &objectID
		newResource.AssignedObjectType = objectType
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
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	err = updatePrimaryStatus(client, info, resourceCreated.Payload.ID)
	if err != nil {
		return err
	}

	return resourceNetboxIpamIPAddressesRead(d, m)
}

func resourceNetboxIpamIPAddressesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	idNum, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid resource id: %w", err)
	}

	ipResp, err := client.Ipam.IpamIPAddressesRead(
		ipam.NewIpamIPAddressesReadParams().WithID(idNum), nil,
	)
	if err != nil {
		return err
	}
	resource := ipResp.Payload

	if err = d.Set("address", *resource.Address); err != nil {
		return err
	}

	customFields := convertCustomFieldsFromAPIToTerraform(resource.CustomFields)
	if err = d.Set("custom_fields", customFields); err != nil {
		return err
	}

	if err = d.Set("description", hackEmptyString(resource.Description)); err != nil {
		return err
	}

	if err = d.Set("dns_name", hackEmptyString(resource.DNSName)); err != nil {
		return err
	}

	var natInside *int64
	if resource.NatInside != nil {
		natInside = &resource.NatInside.ID
	}
	if err = d.Set("nat_inside_id", natInside); err != nil {
		return err
	}

	var natOutside *int64
	if resource.NatOutside != nil {
		natOutside = &resource.NatOutside.ID
	}
	if err = d.Set("nat_outside_id", natOutside); err != nil {
		return err
	}

	if err = d.Set("object_id", resource.AssignedObjectID); err != nil {
		return err
	}

	var primaryIP4 bool
	if *resource.AssignedObjectID != 0 {
		if resource.AssignedObjectType == VMInterfaceType {
			info, err := getInfoForPrimary(m, *resource.AssignedObjectID)
			if err != nil {
				return err
			}

			primaryIP4 = info.vmPrimaryIP4ID == resource.ID
		}
	}
	if err = d.Set("primary_ip4", primaryIP4); err != nil {
		return err
	}

	objectType := resource.AssignedObjectType
	if objectType == "" {
		objectType = VMInterfaceType
	}
	if err = d.Set("object_type", objectType); err != nil {
		return err
	}

	var role *string
	if resource.Role != nil {
		role = resource.Role.Value
	}
	if err = d.Set("role", role); err != nil {
		return err
	}

	var status *string
	if resource.Status != nil {
		status = resource.Status.Value
	}
	if err = d.Set("status", status); err != nil {
		return err
	}

	if err = d.Set("tag", convertNestedTagsToTags(resource.Tags)); err != nil {
		return err
	}

	var tenantID *int64
	if resource.Tenant != nil {
		tenantID = &resource.Tenant.ID
	}
	if err = d.Set("tenant_id", tenantID); err != nil {
		return err
	}

	var vrfID *int64
	if resource.Vrf != nil {
		vrfID = &resource.Vrf.ID
	}
	if err = d.Set("vrf_id", vrfID); err != nil {
		return err
	}

	return nil
}

func resourceNetboxIpamIPAddressesUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableIPAddress{}
	// primary_ip4 := false

	// Required parameters
	address := d.Get("address").(string)
	params.Address = &address

	if d.HasChange("custom_fields") {
		stateCustomFields, resourceCustomFields := d.GetChange("custom_fields")
		customFields := convertCustomFieldsFromTerraformToAPIUpdate(stateCustomFields, resourceCustomFields)
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
		params.DNSName = d.Get("dns_name").(string)
	}

	if d.HasChange("nat_inside_id") {
		natInsideID := int64(d.Get("nat_inside_id").(int))
		if natInsideID != 0 {
			params.NatInside = &natInsideID
		}
	}

	if d.HasChange("nat_outside_id") {
		natInsideID := int64(d.Get("nat_outside_id").(int))
		if natInsideID != 0 {
			params.NatInside = &natInsideID
		}
	}

	if d.HasChange("object_id") || d.HasChange("object_type") {
		// primary_ip4 = true
		objectID := int64(d.Get("object_id").(int))
		params.AssignedObjectID = &objectID

		var objectType string
		if params.AssignedObjectType == "" {
			objectType = VMInterfaceType
		} else {
			objectType = d.Get("object_type").(string)
		}
		params.AssignedObjectType = objectType
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
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Ipam.IpamIPAddressesPartialUpdate(resource, nil)
	if err != nil {
		return err
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
	 *           return err
	 *         }
	 *       }
	 *     }
	 *
	 *     var ipID int64
	 *     ipID = 0
	 *     if isPrimary {
	 *       ipID, err = strconv.ParseInt(d.Id(), 10, 64)
	 *       if err != nil {
	 *         return fmt.Errorf("Unable to convert ID into int64")
	 *       }
	 *     }
	 *     err = updatePrimaryStatus(client, info, ipID)
	 *     if err != nil {
	 *       return err
	 *     }
	 *   }
	 */

	return resourceNetboxIpamIPAddressesRead(d, m)
}

func resourceNetboxIpamIPAddressesDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxIpamIPAddressesExists(d, m)
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

	resource := ipam.NewIpamIPAddressesDeleteParams().WithID(id)
	if _, err := client.Ipam.IpamIPAddressesDelete(resource, nil); err != nil {
		return err
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
