// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	branchprotection "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/branchprotection"
	collaborator "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/collaborator"
	deploykey "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/deploykey"
	gpgkey "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/gpgkey"
	organization "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/organization"
	organizationactionsecret "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/organizationactionsecret"
	organizationactionvariable "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/organizationactionvariable"
	personalaccesstoken "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/personalaccesstoken"
	repository "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/repository"
	repositoryactionsecret "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/repositoryactionsecret"
	repositoryactionvariable "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/repositoryactionvariable"
	repositorywebhook "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/repositorywebhook"
	sshkey "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/sshkey"
	team "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/team"
	teammember "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/teammember"
	user "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/forgejo/user"
	providerconfig "github.com/crossplane-contrib/provider-upjet-forgejo/internal/controller/namespaced/providerconfig"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		branchprotection.Setup,
		collaborator.Setup,
		deploykey.Setup,
		gpgkey.Setup,
		organization.Setup,
		organizationactionsecret.Setup,
		organizationactionvariable.Setup,
		personalaccesstoken.Setup,
		repository.Setup,
		repositoryactionsecret.Setup,
		repositoryactionvariable.Setup,
		repositorywebhook.Setup,
		sshkey.Setup,
		team.Setup,
		teammember.Setup,
		user.Setup,
		providerconfig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		branchprotection.SetupGated,
		collaborator.SetupGated,
		deploykey.SetupGated,
		gpgkey.SetupGated,
		organization.SetupGated,
		organizationactionsecret.SetupGated,
		organizationactionvariable.SetupGated,
		personalaccesstoken.SetupGated,
		repository.SetupGated,
		repositoryactionsecret.SetupGated,
		repositoryactionvariable.SetupGated,
		repositorywebhook.SetupGated,
		sshkey.SetupGated,
		team.SetupGated,
		teammember.SetupGated,
		user.SetupGated,
		providerconfig.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		branchprotection.SetupWebhookWithManager,
		collaborator.SetupWebhookWithManager,
		deploykey.SetupWebhookWithManager,
		gpgkey.SetupWebhookWithManager,
		organization.SetupWebhookWithManager,
		organizationactionsecret.SetupWebhookWithManager,
		organizationactionvariable.SetupWebhookWithManager,
		personalaccesstoken.SetupWebhookWithManager,
		repository.SetupWebhookWithManager,
		repositoryactionsecret.SetupWebhookWithManager,
		repositoryactionvariable.SetupWebhookWithManager,
		repositorywebhook.SetupWebhookWithManager,
		sshkey.SetupWebhookWithManager,
		team.SetupWebhookWithManager,
		teammember.SetupWebhookWithManager,
		user.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
