package security

import (
	"fmt"
	"os"
	"strings"

	"gitlab.skypicker.com/platform/security/iam/shared"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// PermissionManager is used to map user groups to permissions
type PermissionManager interface {
	GetUserPermissions(string, []string) ([]string, error)
}

// YamlPermissionManager is used to map okta user groups to permissions defined in a local yaml configuration
type YamlPermissionManager struct{}

// NewYamlPermissionManager creates an permission manager
func NewYamlPermissionManager() *YamlPermissionManager {
	return &YamlPermissionManager{}
}

// Rule specifies a group of users
type Rule struct {
	Type  string
	Value string
}

// Permission indicates an action, and who can perform it
type Permission struct {
	Name  string
	Rules []Rule
}

func createRules(rules []string) ([]Rule, error) {
	var out []Rule

	for _, rule := range rules {
		split := strings.Split(rule, ":")
		if len(split) != 2 {
			// Keep going, return valid permissions and skipped rules in err?
			return nil, errors.New("failed to parse yaml permission rules")
		}
		r := &Rule{
			Type:  split[0],
			Value: split[1],
		}

		out = append(out, *r)
	}

	return out, nil
}

func createPermissions(flatmap map[string][]string) ([]Permission, error) {
	var out []Permission

	for name, rules := range flatmap {
		rules, err := createRules(rules)
		if err != nil {
			// Keep going, return valid permissions and skipped permissions in err?
			return nil, err
		}

		p := Permission{
			Name:  name,
			Rules: rules,
		}

		out = append(out, p)
	}

	return out, nil
}

func readFile(filename string) (yaml.Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return yaml.Node{}, errors.Wrap(err, "error opening permissions file")
	}

	var contents yaml.Node
	yd := yaml.NewDecoder(file)
	if err = yd.Decode(&contents); err != nil {
		return yaml.Node{}, errors.Wrap(err, "error reading permissions file")
	}

	return contents, nil
}

func buildPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + ":" + key
}

func flattenDocumentNode(node *yaml.Node) (map[string][]string, error) {
	if node.Kind != yaml.DocumentNode {
		return nil, fmt.Errorf("yaml: line: %d column: %d: unexpected configuration", node.Line, node.Column)
	}

	out := make(map[string][]string)
	for _, content := range node.Content {
		if err := flattenNode(content, out, ""); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func flattenNode(node *yaml.Node, out map[string][]string, prefix string) error {
	switch node.Kind {
	case yaml.SequenceNode:
		var value []string
		if err := node.Decode(&value); err != nil {
			return fmt.Errorf("yaml: line: %d column: %d: failed to parse rule: %s", node.Line, node.Column, err.Error())
		}
		out[prefix] = value
	case yaml.MappingNode:
		var mapping map[string]yaml.Node
		if err := node.Decode(&mapping); err != nil {
			return fmt.Errorf("yaml: line: %d column: %d: failed to parse action: %s", node.Line, node.Column, err.Error())
		}

		// indexing because of gocritic - rangeValCopy
		for key := range mapping {
			value := mapping[key]
			if err := flattenNode(&value, out, buildPrefix(prefix, key)); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("yaml: line: %d column: %d: Kind '%d' isn't expected", node.Line, node.Column, node.Kind)
	}

	return nil
}

func getServicePermissions(service string) ([]Permission, error) {
	err := checkServiceName(service)
	if err != nil {
		return nil, err
	}

	filename, err := shared.JoinURL("permissions/", strings.ToLower(service)+".yaml")
	if err != nil {
		return nil, err
	}

	contents, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	// Flatten contents
	servicePermissions, err := flattenDocumentNode(&contents)
	if err != nil {
		return nil, err
	}

	permissions, err := createPermissions(servicePermissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetUserPermissions returns only permissions (associated with a service) that the given user has.
func (p YamlPermissionManager) GetUserPermissions(service string, teamMembership []string) ([]string, error) {
	allPermissions, err := getServicePermissions(service)
	if err != nil {
		return nil, err
	}
	return getApplicablePermissions(allPermissions, teamMembership), nil
}

func getApplicablePermissions(allPermissions []Permission, teamMembership []string) []string {
	var permissions []string

	for _, permission := range allPermissions {
		for _, rule := range permission.Rules {
			if rule.Type == "team" {
				if shared.StringInSlice(strings.TrimSpace(rule.Value), teamMembership) {
					permissions = append(permissions, permission.Name)
					// Skip other rules if permission is already given
					break
				}
			}
		}
	}

	return permissions
}
