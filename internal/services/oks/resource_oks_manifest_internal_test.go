package oks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseManifestImportID(t *testing.T) {
	t.Parallel()

	const clusterID = "11111111-2222-3333-4444-555555555555"

	valid := map[string]struct {
		id       string
		expected manifestImportID
	}{
		"cluster_scoped_grouped_api": {
			id: "cluster_id=" + clusterID + ",apiVersion=oks.dev/v1beta2,kind=NodePool,name=pool-1",
			expected: manifestImportID{
				ClusterId:  clusterID,
				APIVersion: "oks.dev/v1beta2",
				Kind:       "NodePool",
				Name:       "pool-1",
			},
		},
		"namespaced_core_api_any_order": {
			id: "name=coredns-custom,namespace=kube-system,kind=ConfigMap,apiVersion=v1,cluster_id=" + clusterID,
			expected: manifestImportID{
				ClusterId:  clusterID,
				APIVersion: "v1",
				Kind:       "ConfigMap",
				Name:       "coredns-custom",
				Namespace:  "kube-system",
			},
		},
	}

	for name, tc := range valid {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			importID, err := parseManifestImportID(tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, importID)
		})
	}

	invalid := map[string]struct {
		id    string
		error string
	}{
		"empty": {
			id:    "",
			error: "is not a key=value pair",
		},
		"slash_separated": {
			id:    clusterID + "/oks.dev/v1beta2/NodePool/pool-1",
			error: "is not a key=value pair",
		},
		"empty_value": {
			id:    "cluster_id=" + clusterID + ",apiVersion=v1,kind=ConfigMap,name=",
			error: "is not a key=value pair",
		},
		"missing_name": {
			id:    "cluster_id=" + clusterID + ",apiVersion=v1,kind=ConfigMap",
			error: `missing key "name"`,
		},
		"missing_cluster_id": {
			id:    "apiVersion=v1,kind=ConfigMap,name=example",
			error: `missing key "cluster_id"`,
		},
		"unknown_key": {
			id:    "cluster_id=" + clusterID + ",apiVersion=v1,kind=ConfigMap,name=example,group=core",
			error: `unknown key "group"`,
		},
		"duplicate_key": {
			id:    "cluster_id=" + clusterID + ",apiVersion=v1,kind=ConfigMap,name=example,name=other",
			error: `duplicate key "name"`,
		},
	}

	for name, tc := range invalid {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := parseManifestImportID(tc.id)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.error)
			assert.Contains(t, err.Error(), manifestImportIDFormat)
		})
	}
}

func TestManifestPatchOptions(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		forceConflicts types.Bool
		expectForce    bool
	}{
		"null":  {forceConflicts: types.BoolNull(), expectForce: false},
		"false": {forceConflicts: types.BoolValue(false), expectForce: false},
		"true":  {forceConflicts: types.BoolValue(true), expectForce: true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			patchOptions := manifestPatchOptions(tc.forceConflicts)
			assert.Equal(t, manifestFieldManager, patchOptions.FieldManager)
			if tc.expectForce {
				require.NotNil(t, patchOptions.Force)
				assert.True(t, *patchOptions.Force)
			} else {
				assert.Nil(t, patchOptions.Force)
			}
		})
	}
}
