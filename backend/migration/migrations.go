//go:build !ignore

package migration

import "embed"

//go:embed *.sql
var MigrationsFS embed.FS
