package apply

import "github.com/abiosoft/incus-apply/internal/incus"

// Action represents the type of change for a resource.
type Action string

const (
	ActionCreate   Action = "create"
	ActionUpdate   Action = "update"
	ActionReplace  Action = "replace"
	ActionUnchange Action = "unchanged"
	ActionDelete   Action = "delete"
	ActionNotFound Action = "not found"
)

// OutputItem represents a single resource in the output.
type OutputItem struct {
	ResourceID string             `json:"resource_id"`
	Changes    []incus.DiffChange `json:"changes,omitempty"`
	Note       string             `json:"note,omitempty"`
}

// OutputGroup represents a group of resources under a common action.
type OutputGroup struct {
	Action Action       `json:"action"`
	Items  []OutputItem `json:"items"`
}

// Output holds the complete preview output.
type Output struct {
	FileCount     int           `json:"file_count,omitempty"`
	ResourceCount int           `json:"resource_count,omitempty"`
	Groups        []OutputGroup `json:"groups"`
	Summary       string        `json:"summary,omitempty"`
}

// AddGroup appends a group to the output if it has items.
func (o *Output) AddGroup(action Action, items []OutputItem) {
	if len(items) > 0 {
		o.Groups = append(o.Groups, OutputGroup{Action: action, Items: items})
	}
}
