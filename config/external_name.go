package config

import (
	"fmt"
	"math"
	"strconv"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs — all svalabs/forgejo v1.6.0 resources (full schema).
//
// Most resources use IdentifierFromProvider: TF computed `id` (string after our
// svalabs patch) is the Crossplane external-name.
//
// forgejo_user must NOT use ParameterAsIdentifier("login"): that helper still
// reads tfstate["id"] via IDAsExternalName, so after Create external-name becomes
// the numeric id and the next plan rewrites login → ForceNew destroy
// (lifecycle.prevent_destroy). Keep login in forProvider; external-name = id.
//
// Collaborator has no TF `id`; DeployKey/Webhook use key_id / webhook_id.
var ExternalNameConfigs = map[string]config.ExternalName{
	"forgejo_organization":                 config.IdentifierFromProvider,
	"forgejo_repository":                   config.IdentifierFromProvider,
	"forgejo_team":                         config.IdentifierFromProvider,
	"forgejo_team_member":                  config.IdentifierFromProvider,
	"forgejo_user":                         config.IdentifierFromProvider,
	"forgejo_personal_access_token":        config.IdentifierFromProvider,
	"forgejo_collaborator":                 collaboratorExternalName(),
	"forgejo_deploy_key":                   altIDExternalName("key_id"),
	"forgejo_gpg_key":                      config.IdentifierFromProvider,
	"forgejo_ssh_key":                      config.IdentifierFromProvider,
	"forgejo_branch_protection":            config.IdentifierFromProvider,
	"forgejo_organization_action_secret":   config.IdentifierFromProvider,
	"forgejo_organization_action_variable": config.IdentifierFromProvider,
	"forgejo_repository_action_secret":     config.IdentifierFromProvider,
	"forgejo_repository_action_variable":   config.IdentifierFromProvider,
	"forgejo_repository_webhook":           altIDExternalName("webhook_id"),
}

func collaboratorExternalName() config.ExternalName {
	return config.ExternalName{
		SetIdentifierArgumentFn: config.NopSetIdentifierArgument,
		GetExternalNameFn: func(tfstate map[string]any) (string, error) {
			user, _ := tfstate["user"].(string)
			if user == "" {
				return "", errors.New("cannot find user in tfstate")
			}
			rid, err := numericAttrString(tfstate, "repository_id")
			if err != nil {
				return "", err
			}
			return rid + "/" + user, nil
		},
		GetIDFn:                config.ExternalNameAsID,
		DisableNameInitializer: true,
	}
}

func altIDExternalName(attr string) config.ExternalName {
	return config.ExternalName{
		// Do not write key_id/webhook_id into TF *config*: Framework marks
		// them Computed-only → "Invalid Configuration for Read-Only Attribute".
		// After Create they live in TF state; GetExternalNameFn / GetIDFn
		// keep Crossplane external-name in sync (same pattern as Collaborator).
		SetIdentifierArgumentFn: config.NopSetIdentifierArgument,
		GetExternalNameFn: func(tfstate map[string]any) (string, error) {
			return numericAttrString(tfstate, attr)
		},
		GetIDFn:                config.ExternalNameAsID,
		DisableNameInitializer: true,
	}
}

func numericAttrString(tfstate map[string]any, attr string) (string, error) {
	v, ok := tfstate[attr]
	if !ok || v == nil {
		return "", errors.Errorf("cannot find %s in tfstate", attr)
	}
	switch n := v.(type) {
	case string:
		if n == "" {
			return "", errors.Errorf("%s is empty", attr)
		}
		return n, nil
	case float64:
		if math.Trunc(n) != n {
			return "", errors.Errorf("%s is not an integer: %v", attr, n)
		}
		return strconv.FormatInt(int64(n), 10), nil
	case int64:
		return strconv.FormatInt(n, 10), nil
	case int:
		return strconv.Itoa(n), nil
	default:
		return "", errors.Errorf("cannot work with %s type %T: %v", attr, v, fmt.Sprint(v))
	}
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
