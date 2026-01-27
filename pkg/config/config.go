package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	Token = field.StringField(
		"token",
		field.WithDescription("API token used to authenticate to the Panther API."),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)

	URL = field.StringField(
		"url",
		field.WithDescription("API url of your panther account."),
		field.WithRequired(true),
	)
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	Token,
	URL,
})
