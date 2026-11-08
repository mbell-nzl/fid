package export

// ResourceLimits are the ceilings applied to an export's uncompressed members.
// Backup restore uses the same ceilings before reconstructing an archive.
type ResourceLimits struct {
	NonArtifactBytes int64
	ArtifactBytes    int64
	ExpandedBytes    int64
	Entries          int
}

// ArchiveResourceLimits returns the existing export/import resource ceilings.
func ArchiveResourceLimits() ResourceLimits {
	return ResourceLimits{NonArtifactBytes: maxArchiveNonArtifactBytes, ArtifactBytes: maxArchiveArtifactBytes,
		ExpandedBytes: maxArchiveExpandedBytes, Entries: maxArchiveEntries}
}
