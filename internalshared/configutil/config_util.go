// +build !enterprise

package configutil

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

type EntSharedConfig struct{}

type UnusedKeyMap map[string][]token.Pos

type ConfigError struct {
	Problem  string
	Position token.Pos
}

func (c *ConfigError) String() string {
	return fmt.Sprintf("%s at %s", c.Problem, c.Position.String())
}

type ValidatableConfig interface {
	Validate() []ConfigError
}

func (ec *EntSharedConfig) ParseConfig(list *ast.ObjectList) error {
	return nil
}

func ParseEntropy(result *SharedConfig, list *ast.ObjectList, blockName string) error {
	return nil
}

// Creates the ConfigErrors for unused fields, which occur in various structs
func ValidateUnusedFields(unusedKeyPositions UnusedKeyMap, sourceFilePath string) []ConfigError {
	if unusedKeyPositions == nil {
		return nil
	}
	var errors []ConfigError
	for field, positions := range unusedKeyPositions {
		problem := fmt.Sprintf("unknown field %s found in configuration", field)
		for _, pos := range positions {
			if pos.Filename == "" && sourceFilePath != "" {
				pos.Filename = sourceFilePath
			}
			errors = append(errors, ConfigError{
				Problem:  problem,
				Position: pos,
			})
		}
	}
	return errors
}

// UnusedFieldDifference returns all the keys in map a that are not present in map b, and also not present in foundKeys.
func UnusedFieldDifference(a, b UnusedKeyMap, foundKeys []string) UnusedKeyMap {
	if a == nil {
		return nil
	}
	if b == nil {
		return a
	}
	res := make(UnusedKeyMap)
	for k, v := range a {
		if _, ok := b[k]; !ok && !strutil.StrListContainsCaseInsensitive(foundKeys, govalidator.UnderscoreToCamelCase(k)) {
			res[k] = v
		}
	}
	return res
}
