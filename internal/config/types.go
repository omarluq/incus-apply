// Package config handles parsing and validation of incus-apply configuration files.
package config

// Base contains fields common to all Incus resource types.
type Base struct {
	Type        string                    `yaml:"type" json:"type"`                                   // Resource type: instance, profile, network, etc.
	Name        string                    `yaml:"name" json:"name"`                                   // Resource name (unique within type)
	Project     string                    `yaml:"project,omitempty" json:"project,omitempty"`         // --project flag (can be overridden by CLI)
	Config      map[string]string         `yaml:"config,omitempty" json:"config,omitempty"`           // Key-value config options
	Devices     map[string]map[string]any `yaml:"devices,omitempty" json:"devices,omitempty"`         // Device configurations
	Description string                    `yaml:"description,omitempty" json:"description,omitempty"` // Resource description
	SourceFile  string                    `yaml:"-" json:"-"`                                         // Path to source file (set during parsing)
}

// Resource represents a single resource configuration from a .incus.yaml file.
// It embeds BaseConfig for common fields and includes type-specific fields
// used based on the resource Type.
type Resource struct {
	Base `yaml:",inline"`

	PreviewRedactPrefixes []string `yaml:"-" json:"-"`

	// --- Instance-Specific Fields ---

	Image string `yaml:"image,omitempty" json:"image,omitempty"` // Image for instances (e.g., images:debian/12)
	VM    bool   `yaml:"vm,omitempty" json:"vm,omitempty"`       // Create VM instead of container
	Empty bool   `yaml:"empty,omitempty" json:"empty,omitempty"` // Create empty instance (no image)

	// --- Instance Create Flags ---

	Profiles []string `yaml:"profiles,omitempty" json:"profiles,omitempty"` // --profile flags
	Storage  string   `yaml:"storage,omitempty" json:"storage,omitempty"`   // --storage flag
	Network  string   `yaml:"network,omitempty" json:"network,omitempty"`   // --network flag
	Target   string   `yaml:"target,omitempty" json:"target,omitempty"`     // --target flag (cluster member)

	// --- Storage Volume/Bucket Specific ---

	Pool string `yaml:"pool,omitempty" json:"pool,omitempty"` // Pool name for storage volumes/buckets

	// --- Network Specific ---

	NetworkType string `yaml:"networkType,omitempty" json:"networkType,omitempty"` // Network type (bridge, macvlan, etc.)

	// --- Storage Pool Specific ---

	Driver string `yaml:"driver,omitempty" json:"driver,omitempty"` // Storage driver (dir, zfs, btrfs, etc.)
	Source string `yaml:"source,omitempty" json:"source,omitempty"` // Source path/device for storage pool

	// --- Network ACL Specific ---

	Ingress []map[string]any `yaml:"ingress,omitempty" json:"ingress,omitempty"` // Ingress firewall rules
	Egress  []map[string]any `yaml:"egress,omitempty" json:"egress,omitempty"`   // Egress firewall rules
}

// Instance represents configuration specific to Incus instances.
type Instance struct {
	Base `yaml:",inline"`

	Image    string   `yaml:"image,omitempty" json:"image,omitempty"`
	VM       bool     `yaml:"vm,omitempty" json:"vm,omitempty"`
	Empty    bool     `yaml:"empty,omitempty" json:"empty,omitempty"`
	Profiles []string `yaml:"profiles,omitempty" json:"profiles,omitempty"`
	Storage  string   `yaml:"storage,omitempty" json:"storage,omitempty"`
	Network  string   `yaml:"network,omitempty" json:"network,omitempty"`
	Target   string   `yaml:"target,omitempty" json:"target,omitempty"`
}

// StoragePool represents configuration specific to storage pools.
type StoragePool struct {
	Base `yaml:",inline"`

	Driver string `yaml:"driver,omitempty" json:"driver,omitempty"`
	Source string `yaml:"source,omitempty" json:"source,omitempty"`
}

// StorageResource represents configuration for storage volumes and buckets.
type StorageResource struct {
	Base `yaml:",inline"`

	Pool string `yaml:"pool,omitempty" json:"pool,omitempty"`
}

// Network represents configuration specific to networks.
type Network struct {
	Base `yaml:",inline"`

	NetworkType string `yaml:"networkType,omitempty" json:"networkType,omitempty"`
}

// NetworkACL represents configuration specific to network ACLs.
type NetworkACL struct {
	Base `yaml:",inline"`

	Ingress []map[string]any `yaml:"ingress,omitempty" json:"ingress,omitempty"`
	Egress  []map[string]any `yaml:"egress,omitempty" json:"egress,omitempty"`
}

// Stdin represents configuration data passed to incus commands via stdin.
// Incus edit commands accept YAML on stdin to modify resource configuration.
type Stdin struct {
	Config      map[string]string         `yaml:"config,omitempty"`
	Devices     map[string]map[string]any `yaml:"devices,omitempty"`
	Description string                    `yaml:"description,omitempty"`
	Profiles    []string                  `yaml:"profiles,omitempty"`
	Ingress     []map[string]any          `yaml:"ingress,omitempty"` // Network ACL ingress rules
	Egress      []map[string]any          `yaml:"egress,omitempty"`  // Network ACL egress rules
}

// Validate checks if required fields are present in the resource configuration.
func (r Resource) Validate() error {
	if r.Type == "" {
		return &ValidationError{Field: "type", Message: "type is required"}
	}
	if r.Type != "vars" && r.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	return nil
}

// Vars represents a `type: vars` document that declares variables
// for interpolation in resource configs within the same file (or globally).
type Vars struct {
	Vars   map[string]string `yaml:"vars,omitempty"`
	Files  []string          `yaml:"files,omitempty"` // .env files to load
	Global bool              `yaml:"global,omitempty"`

	SourceFile string `yaml:"-"`
}

// ValidationError represents a field-level configuration validation error.
type ValidationError struct {
	Field   string // The field that failed validation
	Message string // Description of the validation failure
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
