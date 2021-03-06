package organization

import (
	"github.com/cloudfoundry/cli/cf/api/organizations"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	. "github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/flags"
)

type ListOrgs struct {
	ui      terminal.UI
	config  core_config.Reader
	orgRepo organizations.OrganizationRepository
}

func (cmd ListOrgs) Execute(fc flags.FlagContext) {
	cmd.ui.Say(T("Getting orgs as {{.Username}}...\n",
		map[string]interface{}{"Username": terminal.EntityNameColor(cmd.config.Username())}))

	noOrgs := true
	table := cmd.ui.Table([]string{T("name")})

	orgs, apiErr := cmd.orgRepo.ListOrgs()
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
	}
	for _, org := range orgs {
		table.Add(org.Name)
		noOrgs = false
	}

	table.Print()

	if apiErr != nil {
		cmd.ui.Failed(T("Failed fetching orgs.\n{{.ApiErr}}",
			map[string]interface{}{"ApiErr": apiErr}))
		return
	}

	if noOrgs {
		cmd.ui.Say(T("No orgs found"))
	}
}

func init() {
	command_registry.Register(&ListOrgs{})
}

func (cmd *ListOrgs) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "orgs",
		ShortName:   "o",
		Description: T("List all orgs"),
		Usage:       "CF_NAME orgs",
	}
}

func (cmd *ListOrgs) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed("Incorrect Usage. No argument required\n\n" + command_registry.Commands.CommandUsage("orgs"))
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd *ListOrgs) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.orgRepo = deps.RepoLocator.GetOrganizationRepository()
	return cmd
}
