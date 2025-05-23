package version_snapshot

type VersionSnapshotManager interface {
	// GetVersionSnapshot retrieves a version snapshot from the database.
	GetVersionSnapshot(id ID) (*VersionSnapshot, error)
	// WriteVersionSnapshot writes a version snapshot to the database.
	WriteVersionSnapshot(snapshot *VersionSnapshot) error
}
