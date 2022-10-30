package requestmodifier

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type netboxRequestModifier struct {
	origwriter     runtime.ClientRequestWriter
	modifiedFields map[string]interface{}
	requiredFields []string
}

func (o netboxRequestModifier) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {
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

	for _, k := range o.requiredFields {
		if val, ok := objmap[k]; ok && (val == nil || val == "0001-01-01T00:00:00.000Z") {
			delete(objmap, k)
		}
	}

	for k, v := range o.modifiedFields {
		switch v.(type) {
		case int64:
			if v == int64(0) {
				objmap[k] = nil
			}
		case float64:
			if v == float64(0) {
				objmap[k] = nil
			}
		case bool:
			if v == false {
				objmap[k] = false
			}
		case string:
			if v == "" {
				objmap[k] = ""
			}
		case nil:
			objmap[k] = nil
		default:
			return fmt.Errorf("unknown type for '%s'", k)
		}
	}

	err = r.SetBodyParam(objmap)
	return err
}

func NewNetboxRequestModifier(modifiedFields map[string]interface{}, requiredFields []string) func(*runtime.ClientOperation) {
	return func(op *runtime.ClientOperation) {
		if len(modifiedFields) > 0 || len(requiredFields) > 0 {
			tmp := netboxRequestModifier{
				origwriter:     op.Params,
				modifiedFields: modifiedFields,
				requiredFields: requiredFields,
			}
			op.Params = tmp
		}
	}
}
