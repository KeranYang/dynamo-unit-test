package version_snapshot

// ID represents the unique identifier for a single version snapshot.
type ID struct {
	// AssetID is the unique identifier for the asset.
	// Required.
	AssetID string `json:"assetId"`
	// ResourceName is the name of the resource.
	// Required.
	ResourceName string `json:"resourceName"`
	// ResourceType is the type of the resource.
	// Required.
	ResourceType string `json:"resourceType"`
	// Location is the location of the resource.
	// Required.
	Location string `json:"location"`
}

type VersionSnapshot struct {
	ID
	// SnapshotSpec is the specification of the version snapshot.
	SnapshotSpec []byte `json:"snapshotSpec"`
}
