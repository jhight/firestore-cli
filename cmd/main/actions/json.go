package actions

import (
	"encoding/json"
)

func (a *Action) toJSON(value any) (string, error) {
	var bytes []byte
	var err error

	if a.cfg.PrettyPrint && !a.cfg.Raw {
		spacing := ""
		for i := 0; i < a.cfg.PrettySpacing; i++ {
			spacing += " "
		}

		bytes, err = json.MarshalIndent(value, "", spacing)
	} else {
		bytes, err = json.Marshal(value)
	}

	return string(bytes), err
}
