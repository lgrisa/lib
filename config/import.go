package config

import (
	"embed"
	_ "embed"
)

//go:embed holiday/*
var HolidayFs embed.FS

//go:embed bin/GeoLite2-Country.mmdb
var IpDbFs []byte
