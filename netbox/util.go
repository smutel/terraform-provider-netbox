package netbox

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	netboxclient "github.com/smutel/go-netbox/v3/netbox/client"
	"github.com/smutel/go-netbox/v3/netbox/client/ipam"
	"github.com/smutel/go-netbox/v3/netbox/client/virtualization"
	"github.com/smutel/go-netbox/v3/netbox/models"
)

// Type of vm interface in Netbox
const vMInterfaceType string = "virtualization.vminterface"


func expandToInt64Slice(v []interface{}) []int64 {
	s := make([]int64, len(v))
	for i, val := range v {
		if strVal, ok := val.(int64); ok {
			s[i] = strVal
		}
	}

	return s
}


func getVMIDForInterface(m interface{}, objectID int64) (int64, error) {
	client := m.(*netboxclient.NetBoxAPI)

	objectIDStr := fmt.Sprintf("%d", objectID)
	paramsInterface := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&objectIDStr)
	interfaces, err := client.Virtualization.VirtualizationInterfacesList(
		paramsInterface, nil)

	if err != nil {
		return 0, err
	}

	for _, i := range interfaces.Payload.Results {
		if i.ID == objectID {
			if i.VirtualMachine != nil {
				return i.VirtualMachine.ID, nil
			}
		}
	}
	return 0, fmt.Errorf("virtual machine not found")
}

func isprimary(m interface{}, objectID *int64, ipID int64, ip4 bool) (bool, error) {
	client := m.(*netboxclient.NetBoxAPI)

	if objectID == nil {
		return false, nil
	}

	objectIDStr := strconv.FormatInt(*objectID, 10)
	paramsInterface := virtualization.NewVirtualizationInterfacesListParams().WithID(
		&objectIDStr)
	interfaces, err := client.Virtualization.VirtualizationInterfacesList(
		paramsInterface, nil)

	if err != nil {
		return false, err
	}

	for _, i := range interfaces.Payload.Results {
		if i.ID == *objectID {
			if i.VirtualMachine != nil {
				vmIDStr := fmt.Sprintf("%d", i.VirtualMachine.ID)
				paramsVM := virtualization.NewVirtualizationVirtualMachinesListParams().WithID(&vmIDStr)
				vms, err := client.Virtualization.VirtualizationVirtualMachinesList(
					paramsVM, nil)

				if err != nil {
					return false, err
				}

				if *vms.Payload.Count != 1 {
					return false, fmt.Errorf("Cannot set an interface as primary when " +
						"the interface is not attached to a VM.")
				}

				if ip4 && vms.Payload.Results[0].PrimaryIp4 != nil {
					return vms.Payload.Results[0].PrimaryIp4.ID == ipID, nil
				} else if !ip4 && vms.Payload.Results[0].PrimaryIp6 != nil {
					return vms.Payload.Results[0].PrimaryIp6.ID == ipID, nil
				} else {
					return false, nil
				}
			} else {
				return false, fmt.Errorf("Cannot set an interface as primary when the " +
					"interface is not attached to a VM.")
			}
		}
	}

	return false, nil
}

func updatePrimaryStatus(m interface{}, vmid, ipid int64, primary bool) error {
	client := m.(*netboxclient.NetBoxAPI)

	emptyFields := make(map[string]interface{})
	dropFields := []string{
		"created",
		"last_updated",
		"name",
		"cluster",
		"tags",
	}

	params := &models.WritableVirtualMachineWithConfigContext{}
	if primary {
		params.PrimaryIp4 = &ipid
	} else {
		params.PrimaryIp4 = nil
		emptyFields["primary_ip4"] = nil
	}
	vm := virtualization.NewVirtualizationVirtualMachinesPartialUpdateParams().WithData(params)
	vm.SetID(vmid)
	_, err := client.Virtualization.VirtualizationVirtualMachinesPartialUpdate(
		vm, nil, newRequestModifierOperation(emptyFields, dropFields))
	if err != nil {
		return err
	}
	return nil
}

// Convert URL in content_type
func convertURIContentType(uri strfmt.URI) string {
	uriSplit := strings.Split(uri.String(), "/")
	firstLevel := uriSplit[len(uriSplit)-4]
	secondLevel := uriSplit[len(uriSplit)-3]

	slRegexp := regexp.MustCompile(`s$`)
	secondLevel = strings.ReplaceAll(secondLevel, "-", "")
	secondLevel = slRegexp.ReplaceAllString(secondLevel, "")

	contentType := firstLevel + "." + secondLevel
	return contentType
}

func getNewAvailableIPForIPRange(client *netboxclient.NetBoxAPI, id int64) (*models.IPAddress, error) {
	params := ipam.NewIpamIPRangesAvailableIpsCreateParams().WithID(id)
	params.Data = []*models.WritableAvailableIP{
		{},
	}
	list, err := client.Ipam.IpamIPRangesAvailableIpsCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

func getNewAvailableIPForPrefix(client *netboxclient.NetBoxAPI, id int64) (*models.IPAddress, error) {
	params := ipam.NewIpamPrefixesAvailableIpsCreateParams().WithID(id)
	params.Data = []*models.WritableAvailableIP{
		{},
	}
	list, err := client.Ipam.IpamPrefixesAvailableIpsCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

func getNewAvailablePrefix(client *netboxclient.NetBoxAPI, id int64, length int64) (*models.Prefix, error) {
	params := ipam.NewIpamPrefixesAvailablePrefixesCreateParams().WithID(id)
	params.Data = []*models.PrefixLength{
		{PrefixLength: &length},
	}
	list, err := client.Ipam.IpamPrefixesAvailablePrefixesCreate(params, nil)
	if err != nil {
		return nil, err
	}
	return list.Payload[0], nil
}

type requestModifier struct {
	origwriter      runtime.ClientRequestWriter
	overwritefields map[string]interface{}
	dropfields      []string
}

func (o requestModifier) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {
	err := o.origwriter.WriteToRequest(r, reg)
	if err != nil {
		return err
	}

	jsonString, err := json.Marshal(r.GetBodyParam())
	if err != nil {
		return err
	}

	var objmap map[string]interface{}
	err = json.Unmarshal(jsonString, &objmap)
	if err != nil {
		return err
	}
	for _, k := range o.dropfields {
		delete(objmap, k)
	}
	for k, v := range o.overwritefields {
		objmap[k] = v
	}

	err = r.SetBodyParam(objmap)
	return err
}

func newRequestModifierOperation(empty map[string]interface{}, drop []string) func(*runtime.ClientOperation) {
	return func(op *runtime.ClientOperation) {
		if len(empty) > 0 || len(drop) > 0 {
			tmp := requestModifier{
				origwriter:      op.Params,
				overwritefields: empty,
				dropfields:      drop,
			}
			op.Params = tmp
		}
	}
}
