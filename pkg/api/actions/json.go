package actions

import (
	"encoding/json"
)

func (a *action) toJSON(value any) (string, error) {
	var bytes []byte
	var err error

	if a.initializer.Config().PrettyPrint && !a.initializer.Config().RawPrint {
		spacing := ""
		for i := 0; i < a.initializer.Config().PrettySpacing; i++ {
			spacing += " "
		}

		bytes, err = json.MarshalIndent(value, "", spacing)
	} else {
		bytes, err = json.Marshal(value)
	}

	return string(bytes), err
}
