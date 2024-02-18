package main

import (
	"github.com/nats-io/stan.go"
	"github.com/xeipuuv/gojsonschema"
)

func validateData(msg *stan.Msg) bool {
	loader := gojsonschema.NewStringLoader(config.ValidationSchema)
	documentLoader := gojsonschema.NewStringLoader(string(msg.Data))

	result, err := gojsonschema.Validate(loader, documentLoader)
	if err != nil {
		return false
	}

	if result.Valid() {
		return true
	} else {
		return false
	}
}
