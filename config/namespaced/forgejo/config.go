package forgejo

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures all forgejo resources (namespaced / Crossplane v2).
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("forgejo_organization", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "Organization"
	})
	p.AddResourceConfigurator("forgejo_repository", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "Repository"
		r.References["owner"] = config.Reference{
			TerraformName: "forgejo_organization",
		}
	})
	p.AddResourceConfigurator("forgejo_team", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "Team"
		r.References["organization"] = config.Reference{
			TerraformName: "forgejo_organization",
		}
	})
	p.AddResourceConfigurator("forgejo_team_member", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "TeamMember"
		r.References["team_id"] = config.Reference{
			TerraformName: "forgejo_team",
		}
		r.References["user"] = config.Reference{
			TerraformName: "forgejo_user",
		}
	})
	p.AddResourceConfigurator("forgejo_user", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "User"
	})
	p.AddResourceConfigurator("forgejo_personal_access_token", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "PersonalAccessToken"
		r.References["user"] = config.Reference{
			TerraformName: "forgejo_user",
		}
	})
	p.AddResourceConfigurator("forgejo_collaborator", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "Collaborator"
		r.References["repository_id"] = config.Reference{
			TerraformName: "forgejo_repository",
		}
		r.References["user"] = config.Reference{
			TerraformName: "forgejo_user",
		}
	})
	p.AddResourceConfigurator("forgejo_deploy_key", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "DeployKey"
		r.References["repository_id"] = config.Reference{
			TerraformName: "forgejo_repository",
		}
	})
	p.AddResourceConfigurator("forgejo_gpg_key", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "GpgKey"
	})
	p.AddResourceConfigurator("forgejo_ssh_key", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "SshKey"
	})
	p.AddResourceConfigurator("forgejo_branch_protection", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "BranchProtection"
		r.References["repository_id"] = config.Reference{
			TerraformName: "forgejo_repository",
		}
	})
	p.AddResourceConfigurator("forgejo_organization_action_secret", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "OrganizationActionSecret"
	})
	p.AddResourceConfigurator("forgejo_organization_action_variable", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "OrganizationActionVariable"
	})
	p.AddResourceConfigurator("forgejo_repository_action_secret", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "RepositoryActionSecret"
		r.References["repository_id"] = config.Reference{
			TerraformName: "forgejo_repository",
		}
	})
	p.AddResourceConfigurator("forgejo_repository_action_variable", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "RepositoryActionVariable"
		r.References["repository_id"] = config.Reference{
			TerraformName: "forgejo_repository",
		}
	})
	p.AddResourceConfigurator("forgejo_repository_webhook", func(r *config.Resource) {
		r.ShortGroup = "forgejo"
		r.Kind = "RepositoryWebhook"
		r.References["repository_id"] = config.Reference{
			TerraformName: "forgejo_repository",
		}
	})
}
