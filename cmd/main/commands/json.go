package commands

import (
	"encoding/json"
	"github.com/jhight/firestore-cli/pkg/config"
)

func toJSON(cfg config.Config, value any) (string, error) {
	var bytes []byte
	var err error

	if cfg.PrettyPrint && !cfg.Raw {
		spacing := ""
		for i := 0; i < cfg.PrettySpacing; i++ {
			spacing += " "
		}

		bytes, err = json.MarshalIndent(value, "", spacing)
	} else {
		bytes, err = json.Marshal(value)
	}

	return string(bytes), err
}
