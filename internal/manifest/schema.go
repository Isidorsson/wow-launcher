package manifest

// File describes one downloadable artifact + its destination in the WoW install.
//
// Path is relative to the WoW install root, using forward slashes.
// Example: "Data/patch-4.MPQ" or "Data/enUS/patch-enUS-4.MPQ".
type File struct {
	Path     string `json:"path"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	SHA256   string `json:"sha256"`
	Required bool   `json:"required"`
	Label    string `json:"label,omitempty"` // optional, shown for optional packs
}

// Manifest is the server-side description of patch state.
//
// Versioned filenames on the CDN let manifests be immutable per-version
// while patch URLs stay infinitely cacheable.
type Manifest struct {
	SchemaVersion      int    `json:"schema_version"`       // bump on breaking changes
	ServerID           string `json:"server_id"`            // must match config.toml entry
	ClientVersion      string `json:"client_version"`       // informational, e.g. "3.3.5a-12340"
	Realmlist          string `json:"realmlist"`            // value written to realmlist.wtf
	Locale             string `json:"locale"`               // enUS, enGB, deDE, etc.
	MinLauncherVersion string `json:"min_launcher_version"` // semver; below = refuse launch
	Files              []File `json:"files"`
}

// Signed wraps a manifest with an Ed25519 signature over the canonical JSON
// of the inner Manifest. PubkeyHex must match the launcher's embedded pubkey.
type Signed struct {
	Manifest     Manifest `json:"manifest"`
	SignatureHex string   `json:"signature"`
	PubkeyHex    string   `json:"pubkey"`
}
