package oks

import (
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/services/oks/okshelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestRequiresReplace(t *testing.T) {
	t.Parallel()

	configMap := func(name, namespace string) string {
		manifest := "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: " + name + "\n"
		if namespace != "" {
			manifest += "  namespace: " + namespace + "\n"
		}
		return manifest + "data:\n  key: value\n"
	}

	cases := map[string]struct {
		plan, state string
		// planNamespaced sets "default" on the planned object, as GetResourceInterfaceFromManifest does for a namespaced kind
		planNamespaced bool
		expected       bool
	}{
		"namespace_omitted_in_both": {
			plan:           configMap("app", ""),
			state:          configMap("app", ""),
			planNamespaced: true,
			expected:       false,
		},
		"namespace_default_added": {
			plan:     configMap("app", "default"),
			state:    configMap("app", ""),
			expected: false,
		},
		"namespace_default_removed": {
			plan:           configMap("app", ""),
			state:          configMap("app", "default"),
			planNamespaced: true,
			expected:       false,
		},
		"namespace_changed": {
			plan:     configMap("app", "other"),
			state:    configMap("app", ""),
			expected: true,
		},
		"name_changed": {
			plan:     configMap("app-2", "default"),
			state:    configMap("app", "default"),
			expected: true,
		},
		"kind_changed": {
			plan:     "apiVersion: v1\nkind: Secret\nmetadata:\n  name: app\n",
			state:    configMap("app", ""),
			expected: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			planObj, err := okshelpers.FromYAML(tc.plan)
			require.NoError(t, err)
			if tc.planNamespaced {
				planObj.SetNamespace(okshelpers.NamespaceOrDefault(planObj.GetNamespace()))
			}
			stateObj, err := okshelpers.FromYAML(tc.state)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, manifestRequiresReplace(planObj, stateObj))
		})
	}
}
