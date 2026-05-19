package db

import "embed"

//go:embed migrations/*.sql
var Migracoes embed.FS
