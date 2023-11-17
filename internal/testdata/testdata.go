package testdata

import "embed"

//go:embed files/*.sql
var TestFS embed.FS
