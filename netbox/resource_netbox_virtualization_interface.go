package netbox

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	netboxclient "github.com/smutel/go-netbox/netbox/client"
	"github.com/smutel/go-netbox/netbox/client/virtualization"
	"github.com/smutel/go-netbox/netbox/models"
)

func resourceNetboxVirtualizationInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVirtualizationInterfaceCreate,
		Read:   resourceNetboxVirtualizationInterfaceRead,
		Update: resourceNetboxVirtualizationInterfaceUpdate,
		Delete: resourceNetboxVirtualizationInterfaceDelete,
		Exists: resourceNetboxVirtualizationInterfaceExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      " ",
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^([A-Z0-9]{2}:){5}[A-Z0-9]{2}$"),
					"Must be like AA:AA:AA:AA:AA"),
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"access", "tagged",
					"tagged-all"}, false),
				ForceNew: true,
			},
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65536),
				ForceNew:     true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"tagged_vlans": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"untagged_vlan": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"virtualmachine_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceNetboxVirtualizationInterfaceCreate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	macAddress := d.Get("mac_address").(string)
	mode := d.Get("mode").(string)
	mtu := int64(d.Get("mtu").(int))
	name := d.Get("name").(string)
	taggedVlans := d.Get("tagged_vlans").(*schema.Set).List()
	tags := d.Get("tag").(*schema.Set).List()
	untaggedVlan := int64(d.Get("untagged_vlan").(int))
	virtualmachineID := int64(d.Get("virtualmachine_id").(int))

	newResource := &models.WritableVMInterface{
		Description:    description,
		Enabled:        enabled,
		Mode:           mode,
		Name:           &name,
		TaggedVlans:    expandToInt64Slice(taggedVlans),
		Tags:           convertTagsToNestedTags(tags),
		VirtualMachine: &virtualmachineID,
	}

	if macAddress != "" {
		newResource.MacAddress = &macAddress
	}

	if mtu != 0 {
		newResource.Mtu = &mtu
	}

	if untaggedVlan != 0 {
		newResource.UntaggedVlan = &untaggedVlan
	}

	resource := virtualization.NewVirtualizationInterfacesCreateParams().WithData(newResource)

	resourceCreated, err := client.Virtualization.VirtualizationInterfacesCreate(
		resource, nil)

	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(resourceCreated.Payload.ID, 10))

	return resourceNetboxVirtualizationInterfaceRead(d, m)
}

func resourceNetboxVirtualizationInterfaceRead(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceID := d.Id()
	params := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationInterfacesList(
		params, nil)

	if err != nil {
		return err
	}

	for _, resource := range resources.Payload.Results {
		if strconv.FormatInt(resource.ID, 10) == d.Id() {
			var description string

			if resource.Description == "" {
				description = " "
			} else {
				description = resource.Description
			}

			if err = d.Set("description", description); err != nil {
				return err
			}

			if err = d.Set("enabled", resource.Enabled); err != nil {
				return err
			}

			if err = d.Set("mac_address", resource.MacAddress); err != nil {
				return err
			}

			if resource.Mode == nil {
				if err = d.Set("mode", ""); err != nil {
					return err
				}
			} else {
				if err = d.Set("mode", resource.Mode.Value); err != nil {
					return err
				}
			}

			if err = d.Set("mtu", resource.Mtu); err != nil {
				return err
			}

			if err = d.Set("name", resource.Name); err != nil {
				return err
			}

			if err = d.Set("tagged_vlans", resource.TaggedVlans); err != nil {
				return err
			}

			if err = d.Set("tag", convertNestedTagsToTags(
				resource.Tags)); err != nil {
				return err
			}

			if err = d.Set("untagged_vlan", resource.UntaggedVlan); err != nil {
				return err
			}

			if resource.VirtualMachine == nil {
				if err = d.Set("virtualmachine_id", 0); err != nil {
					return err
				}
			} else {
				if err = d.Set("virtualmachine_id",
					resource.VirtualMachine.ID); err != nil {
					return err
				}
			}

			if err = d.Set("type", VMInterfaceType); err != nil {
				return err
			}

			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceNetboxVirtualizationInterfaceUpdate(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)
	params := &models.WritableVMInterface{}

	// Required parameters
	name := d.Get("name").(string)
	params.Name = &name
	virtualMachineID := int64(d.Get("virtualmachine_id").(int))
	params.VirtualMachine = &virtualMachineID
	taggedVlans := d.Get("tagged_vlans").(*schema.Set).List()
	params.TaggedVlans = expandToInt64Slice(taggedVlans)

	if d.HasChange("description") {
		description := d.Get("description").(string)
		params.Description = description
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		params.Enabled = enabled
	}

	if d.HasChange("mac_address") {
		macAddress := d.Get("mac_address").(string)
		params.MacAddress = &macAddress
	}

	if d.HasChange("mode") {
		mode := d.Get("mode").(string)
		params.Mode = mode
	}

	if d.HasChange("mtu") {
		mtu := int64(d.Get("mtu").(int))
		params.Mtu = &mtu
	}

	tags := d.Get("tag").(*schema.Set).List()
	params.Tags = convertTagsToNestedTags(tags)

	if d.HasChange("untagged_vlan") {
		untaggedVlan := int64(d.Get("untagged_vlan").(int))
		params.UntaggedVlan = &untaggedVlan
	}

	resource := virtualization.NewVirtualizationInterfacesPartialUpdateParams().WithData(params)

	resourceID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Unable to convert ID into int64")
	}

	resource.SetID(resourceID)

	_, err = client.Virtualization.VirtualizationInterfacesPartialUpdate(
		resource, nil)
	if err != nil {
		return err
	}

	return resourceNetboxVirtualizationInterfaceRead(d, m)
}

func resourceNetboxVirtualizationInterfaceDelete(d *schema.ResourceData,
	m interface{}) error {
	client := m.(*netboxclient.NetBoxAPI)

	resourceExists, err := resourceNetboxVirtualizationInterfaceExists(d, m)
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

	p := virtualization.NewVirtualizationInterfacesDeleteParams().WithID(id)
	if _, err := client.Virtualization.VirtualizationInterfacesDelete(
		p, nil); err != nil {
		return err
	}

	return nil
}

func resourceNetboxVirtualizationInterfaceExists(d *schema.ResourceData,
	m interface{}) (b bool,
	e error) {
	client := m.(*netboxclient.NetBoxAPI)
	resourceExist := false

	resourceID := d.Id()
	params := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&resourceID)
	resources, err := client.Virtualization.VirtualizationInterfacesList(
		params, nil)
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
