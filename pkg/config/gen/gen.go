package main

import (
	cfg "github.com/conductorone/baton-panther/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("panther", cfg.Config)
}
