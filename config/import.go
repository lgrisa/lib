package config

import (
	"embed"
	_ "embed"
)

//go:embed holiday/*
var HolidayFs embed.FS
