package storages

import "github.com/hop-/goi/internal/core"

// Run init when package is imported
func init() {
	core.RegisterStorage("sqlite", newSqliteStorage)
	// TODO: register all storage generators here
}
