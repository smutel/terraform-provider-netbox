package requestmodifier

import (
	"encoding/json"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

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

func NewRequestModifierOperation(empty map[string]interface{}, drop []string) func(*runtime.ClientOperation) {
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
