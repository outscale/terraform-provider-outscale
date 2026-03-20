package testacc

import (
	"context"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runPlanCheck(ctx context.Context, resourceAddress string, allowedAttributes []string, plan *tfjson.Plan) error {
	check := ExpectEmptyPlanExcept(resourceAddress, allowedAttributes...)
	req := plancheck.CheckPlanRequest{Plan: plan}
	resp := &plancheck.CheckPlanResponse{}
	check.CheckPlan(ctx, req, resp)
	return resp.Error
}

func TestExpectEmptyPlanExcept(t *testing.T) {
	t.Run("empty plan passes", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{
					Address: "outscale_vm.vm",
					Change:  &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionNoop}},
				},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		assert.NoError(t, err)
	})

	t.Run("allowed attribute change passes", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{
					Address: "outscale_vm.vm",
					Change: &tfjson.Change{
						Actions: []tfjson.Action{tfjson.ActionUpdate},
						Before:  map[string]any{"id": "vm-123", "state": "running"},
						After:   map[string]any{"id": "vm-123", "state": "stopped"},
					},
				},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		assert.NoError(t, err)
	})

	t.Run("non-allowed attribute change fails", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{
					Address: "outscale_vm.vm",
					Change: &tfjson.Change{
						Actions: []tfjson.Action{tfjson.ActionUpdate},
						Before:  map[string]any{"id": "vm-123", "vm_type": "tinav7.c1r1p1"},
						After:   map[string]any{"id": "vm-123", "vm_type": "tinav6.c1r1p1"},
					},
				},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "vm_type")
	})

	t.Run("change on different resource fails", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{
					Address: "outscale_nic.nic",
					Change: &tfjson.Change{
						Actions: []tfjson.Action{tfjson.ActionUpdate},
						Before:  map[string]any{"id": "nic-123"},
						After:   map[string]any{"id": "nic-456"},
					},
				},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "outscale_nic.nic")
	})

	t.Run("multiple allowed attributes change passes", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{
					Address: "outscale_vm.vm",
					Change: &tfjson.Change{
						Actions: []tfjson.Action{tfjson.ActionUpdate},
						Before:  map[string]any{"state": "running", "public_ip": "0.0.0.0"},
						After:   map[string]any{"state": "stopped", "public_ip": "1.1.1.1"},
					},
				},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state", "public_ip"}, plan)
		assert.NoError(t, err)
	})

	t.Run("mixed resources with allowed changes only on target passes", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{
					Address: "outscale_vm.vm",
					Change: &tfjson.Change{
						Actions: []tfjson.Action{tfjson.ActionUpdate},
						Before:  map[string]any{"state": "running"},
						After:   map[string]any{"state": "stopped"},
					},
				},
				{
					Address: "outscale_nic.nic",
					Change:  &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionNoop}},
				},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		assert.NoError(t, err)
	})

	t.Run("nil change is ignored", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{Address: "outscale_vm.vm", Change: nil},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		assert.NoError(t, err)
	})

	t.Run("allowed attribute with no actual change passes", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{
					Address: "outscale_vm.vm",
					Change: &tfjson.Change{
						Actions: []tfjson.Action{tfjson.ActionNoop},
						Before:  map[string]any{"state": "running"},
						After:   map[string]any{"state": "running"},
					},
				},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		assert.NoError(t, err)
	})

	t.Run("all resources unchanged passes", func(t *testing.T) {
		plan := &tfjson.Plan{
			ResourceChanges: []*tfjson.ResourceChange{
				{Address: "outscale_vm.vm", Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionNoop}}},
				{Address: "outscale_nic.nic", Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionNoop}}},
				{Address: "outscale_net.net", Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionNoop}}},
			},
		}
		err := runPlanCheck(t.Context(), "outscale_vm.vm", []string{"state"}, plan)
		assert.NoError(t, err)
	})
}
