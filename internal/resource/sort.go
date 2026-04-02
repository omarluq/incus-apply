package resource

import (
	"sort"

	"github.com/abiosoft/incus-apply/internal/config"
)

// defaultPriority is assigned to unknown resource types.
// Higher value means processed later during apply, earlier during delete.
const defaultPriority = 100

// SortForApply returns a copy of resources sorted by priority for creation.
// Resources with lower priority values are created first to satisfy dependencies
// (e.g., projects/networks before instances that use them).
func SortForApply(resources []*config.Resource) []*config.Resource {
	return sortResources(resources, func(a, b int) bool { return a < b })
}

// SortForDelete returns a copy of resources sorted by priority for deletion.
// Resources with higher priority values are deleted first (reverse of creation order)
// to respect dependency relationships.
func SortForDelete(resources []*config.Resource) []*config.Resource {
	return sortResources(resources, func(a, b int) bool { return a > b })
}

// sortResources creates a sorted copy using the provided comparison function.
// Uses stable sort to preserve relative order of same-priority resources.
func sortResources(resources []*config.Resource, less func(a, b int) bool) []*config.Resource {
	sorted := make([]*config.Resource, len(resources))
	copy(sorted, resources)

	sort.SliceStable(sorted, func(i, j int) bool {
		pi := getPriority(sorted[i].Type)
		pj := getPriority(sorted[j].Type)
		return less(pi, pj)
	})

	return sorted
}

// getPriority returns the priority for a resource type.
// Unknown types get defaultPriority, ensuring they are processed last during
// apply and first during delete (safest behavior for unknown resources).
func getPriority(t string) int {
	if meta, ok := GetTypeMeta(t); ok {
		return meta.Priority
	}
	return defaultPriority
}
