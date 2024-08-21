package storages

// Run init when package is imported
func init() {
	RegisterStorage("void", newVoidStorage)
	RegisterStorage("sqlite", newSqliteStorage)

	// Register all storage generators here
}
