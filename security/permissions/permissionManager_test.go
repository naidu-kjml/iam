package permissions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCreateRules(t *testing.T) {
	var expected []Rule

	for _, rule := range []string{"CS", "CS Systems", "Engineering"} {
		r := &Rule{
			Type:  "team",
			Value: rule,
		}

		expected = append(expected, *r)
	}

	rules, err := createRules([]string{"team:CS", "team:CS Systems", "team:Engineering"})
	assert.NoError(t, err)
	assert.Equal(t, expected, rules)

	rules, err = createRules([]string{"team:CS:invalid"})
	assert.Nil(t, rules)
	assert.Error(t, err)
}

func TestCreatePermissions(t *testing.T) {
	expected := []Permission{
		{
			Name: "comments:write",
			Rules: []Rule{
				{
					Type:  "team",
					Value: "CS",
				},
			},
		},
	}

	permissions, err := createPermissions(map[string][]string{"comments:write": {"team:CS"}})
	assert.NoError(t, err)
	assert.Equal(t, expected, permissions)
}

func TestFlattenDocumentNode(t *testing.T) {
	test := `---
foo:
  bar:
    - first
    - second
foo-bar:
  - single`
	expected := map[string][]string{
		"foo:bar": {"first", "second"},
		"foo-bar": {"single"},
	}

	var node yaml.Node
	err := yaml.Unmarshal([]byte(test), &node)
	require.NoError(t, err)

	result, _ := flattenDocumentNode(&node)
	assert.Equal(t, expected, result)
}

func TestGetApplicablePermissions(t *testing.T) {
	type TestCase struct {
		Membership []string
		Expected   []string
	}

	permissions, _ := createPermissions(map[string][]string{
		"key1:read":  {"team:Engineering/Content + CS Systems"},
		"key1:write": {"team:Engineering"},
		"key2:read":  {"team:CS"},
		"key2:write": {"team:Systems"},
	},
	)

	tests := []TestCase{
		{
			Membership: []string{"Engineering", "Engineering/Content + cs systems"},
			Expected:   []string{"key1:read", "key1:write"},
		},
		{
			Membership: []string{""},
			Expected:   []string(nil),
		},
		{
			Membership: []string{"CS"},
			Expected:   []string{"key2:read"},
		},
	}

	for _, test := range tests {
		result := getApplicablePermissions(permissions, test.Membership)
		assert.ElementsMatch(t, test.Expected, result)
	}
}
