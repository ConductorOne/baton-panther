package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	Token = field.StringField(
		"token",
		field.WithDisplayName("API key"),
		field.WithDescription("API token used to authenticate to the Panther API."),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)

	URL = field.StringField(
		"url",
		field.WithDisplayName("API URL"),
		field.WithPlaceholder("Your Panther API URL"),
		field.WithDescription("API url of your panther account."),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	Token,
	URL,
},
	field.WithConnectorDisplayName("Panther"),
	field.WithIconUrl("/static/app-icons/panther.svg"),
	field.WithHelpUrl("/docs/baton/panther"),
)
