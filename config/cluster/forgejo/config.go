package forgejo

import "github.com/crossplane/upjet/v2/pkg/config"

const (
	shortGroup = "forgejo"

	tfOrganization = "forgejo_organization"
	tfRepository   = "forgejo_repository"
	tfTeam         = "forgejo_team"
	tfTeamMember   = "forgejo_team_member"
	tfUser         = "forgejo_user"
	tfPAT          = "forgejo_personal_access_token"
	tfCollaborator = "forgejo_collaborator"
	tfDeployKey    = "forgejo_deploy_key"
	tfGpgKey       = "forgejo_gpg_key"
	tfSshKey       = "forgejo_ssh_key"
	tfBranchProt   = "forgejo_branch_protection"
	tfOrgSecret    = "forgejo_organization_action_secret"
	tfOrgVariable  = "forgejo_organization_action_variable"
	tfRepoSecret   = "forgejo_repository_action_secret"
	tfRepoVariable = "forgejo_repository_action_variable"
	tfRepoWebhook  = "forgejo_repository_webhook"
)

// Configure configures all forgejo resources (cluster-scoped).
func Configure(p *config.Provider) {
	p.AddResourceConfigurator(tfOrganization, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "Organization"
	})
	p.AddResourceConfigurator(tfRepository, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "Repository"
		r.References["owner"] = config.Reference{
			TerraformName: tfOrganization,
		}
	})
	p.AddResourceConfigurator(tfTeam, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "Team"
		r.References["organization"] = config.Reference{
			TerraformName: tfOrganization,
		}
	})
	p.AddResourceConfigurator(tfTeamMember, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "TeamMember"
		r.References["team_id"] = config.Reference{
			TerraformName: tfTeam,
		}
		r.References["user"] = config.Reference{
			TerraformName: tfUser,
		}
	})
	p.AddResourceConfigurator(tfUser, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "User"
	})
	p.AddResourceConfigurator(tfPAT, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "PersonalAccessToken"
		r.References["user"] = config.Reference{
			TerraformName: tfUser,
		}
	})
	p.AddResourceConfigurator(tfCollaborator, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "Collaborator"
		r.References["repository_id"] = config.Reference{
			TerraformName: tfRepository,
		}
		r.References["user"] = config.Reference{
			TerraformName: tfUser,
		}
	})
	p.AddResourceConfigurator(tfDeployKey, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "DeployKey"
		r.References["repository_id"] = config.Reference{
			TerraformName: tfRepository,
		}
	})
	p.AddResourceConfigurator(tfGpgKey, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "GpgKey"
	})
	p.AddResourceConfigurator(tfSshKey, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "SshKey"
	})
	p.AddResourceConfigurator(tfBranchProt, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "BranchProtection"
		r.References["repository_id"] = config.Reference{
			TerraformName: tfRepository,
		}
	})
	p.AddResourceConfigurator(tfOrgSecret, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "OrganizationActionSecret"
	})
	p.AddResourceConfigurator(tfOrgVariable, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "OrganizationActionVariable"
	})
	p.AddResourceConfigurator(tfRepoSecret, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "RepositoryActionSecret"
		r.References["repository_id"] = config.Reference{
			TerraformName: tfRepository,
		}
	})
	p.AddResourceConfigurator(tfRepoVariable, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "RepositoryActionVariable"
		r.References["repository_id"] = config.Reference{
			TerraformName: tfRepository,
		}
	})
	p.AddResourceConfigurator(tfRepoWebhook, func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "RepositoryWebhook"
		r.References["repository_id"] = config.Reference{
			TerraformName: tfRepository,
		}
	})
}
