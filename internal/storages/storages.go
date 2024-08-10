package storages

// Run init when package is imported
func init() {
	RegisterStorage("void", newVoidStorage)
	RegisterStorage("sqlite", newSqliteStorage)
	// TODO: register all storage generators here
}
