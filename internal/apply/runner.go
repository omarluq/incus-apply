package apply

import (
	"github.com/abiosoft/incus-apply/internal/config"
	"github.com/abiosoft/incus-apply/internal/incus"
	"github.com/abiosoft/incus-apply/internal/resource"
)

// upsertAction is the planned action for a single resource.
type upsertAction int

const (
	upsertSkip    upsertAction = iota // no changes detected
	upsertCreate                      // resource does not exist
	upsertUpdate                      // resource exists and has changes
	upsertReplace                     // resource must be deleted and recreated
)

// upsertPlan captures the create/update/skip decisions made during the diff
// phase so the execution phase can act on them without recomputing.
type upsertPlan struct {
	res    *config.Resource
	action upsertAction
}

// deletePlan captures the delete/skip decisions made during the diff phase
// so the execution phase can act on them without re-checking existence.
type deletePlan struct {
	res  *config.Resource
	skip bool // true when the resource was not found during diff
}

// runner executes planned operations against the Incus client.
type runner struct {
	opts    *Options
	client  incus.Client
	result  result
	printer printer
}

func (r *runner) printSummary() {
	r.printer.Print(r.opts.Quiet, r.result)
}

// upsert handles create-or-update logic for a single resource.
// Returns an error only if FailFast is enabled and an error occurs.
func (r *runner) upsert(plan upsertPlan) error {
	switch plan.action {
	case upsertCreate:
		return r.create(plan.res, formatResourceID(plan.res))
	case upsertUpdate:
		return r.update(plan.res, formatResourceID(plan.res))
	case upsertReplace:
		return r.replace(plan.res, formatResourceID(plan.res))
	default: // upsertSkip
		r.result.unchanged++
		return nil
	}
}

func (r *runner) replace(res *config.Resource, resourceID string) error {
	deleteResult := r.client.Delete(res)
	if deleteResult.Error != nil {
		return r.result.recordError(r.opts.FailFast, resourceID, "replace delete failed", deleteResult.Error)
	}
	createResult := r.client.Create(res)
	if createResult.Error != nil {
		return r.result.recordError(r.opts.FailFast, resourceID, "replace create failed", createResult.Error)
	}
	printColored(r.opts.Quiet, colorYellow, "! %s replaced", resourceID)
	r.result.replaced++

	if resource.Type(res.Type) == resource.TypeInstance && r.opts.Launch {
		startResult := r.client.Start(res)
		if startResult.Error != nil {
			return r.result.recordError(r.opts.FailFast, resourceID, "start failed", startResult.Error)
		}
	}

	return nil
}

// update applies an update for a resource that was already confirmed to have changes.
func (r *runner) update(res *config.Resource, resourceID string) error {
	result := r.client.Update(res)
	if result.Error != nil {
		return r.result.recordError(r.opts.FailFast, resourceID, "update failed", result.Error)
	}
	printColored(r.opts.Quiet, colorYellow, "~ %s updated", resourceID)
	r.result.updated++
	return nil
}

// create handles creating a new resource.
func (r *runner) create(res *config.Resource, resourceID string) error {
	result := r.client.Create(res)
	if result.Error != nil {
		return r.result.recordError(r.opts.FailFast, resourceID, "create failed", result.Error)
	}
	printColored(r.opts.Quiet, colorGreen, "+ %s created", resourceID)
	r.result.created++

	// Start newly created instances unless --launch=false
	if resource.Type(res.Type) == resource.TypeInstance && r.opts.Launch {
		startResult := r.client.Start(res)
		if startResult.Error != nil {
			return r.result.recordError(r.opts.FailFast, resourceID, "start failed", startResult.Error)
		} else {
			printInfo(r.opts.Quiet, "  └─ started")
		}
	}
	return nil
}

// delete handles deletion of a single resource based on the pre-computed plan.
func (r *runner) delete(p deletePlan) error {
	resourceID := formatResourceID(p.res)

	if p.skip {
		printColored(r.opts.Quiet, colorDim, "= %s (not found)", resourceID)
		r.result.skipped++
		return nil
	}

	result := r.client.Delete(p.res)
	if result.Error != nil {
		return r.result.recordError(r.opts.FailFast, resourceID, "delete failed", result.Error)
	}
	printColored(r.opts.Quiet, colorRed, "- %s deleted", resourceID)
	r.result.deleted++
	return nil
}
