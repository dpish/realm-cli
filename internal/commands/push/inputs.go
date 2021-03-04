package push

import (
	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/cloud/realm"
	"github.com/10gen/realm-cli/internal/local"
	"github.com/10gen/realm-cli/internal/terminal"
	"github.com/10gen/realm-cli/internal/utils/flags"
)

const (
	flagAppDirectory      = "app-dir"
	flagAppDirectoryShort = "p"
	flagAppDirectoryUsage = "provide the path to a Realm app containing the changes to push"

	flagProject      = "project"
	flagProjectUsage = "the MongoDB cloud project id"

	flagTo      = "to"
	flagToShort = "t"
	flagToUsage = "choose a Realm app to push changes towards"

	flagDryRun      = "dry-run"
	flagDryRunShort = "x"
	flagDryRunUsage = "include to run without pushing any changes to the Realm server"

	flagIncludeDependencies      = "include-dependencies"
	flagIncludeDependenciesShort = "d"
	flagIncludeDependenciesUsage = "include to push Realm app dependencies changes as well"

	flagIncludeHosting      = "include-hosting"
	flagIncludeHostingShort = "s"
	flagIncludeHostingUsage = "include to push Realm app hosting changes as well"

	flagResetCDNCache      = "reset-cdn-cache"
	flagResetCDNCacheShort = "c"
	flagResetCDNCacheUsage = "include to reset the Realm app hosting CDN cache"
)

type to struct {
	GroupID string
	AppID   string
}

type inputs struct {
	AppDirectory        string
	Project             string
	To                  string
	IncludeDependencies bool
	IncludeHosting      bool
	ResetCDNCache       bool
	DryRun              bool
}

func (i *inputs) Resolve(profile *cli.Profile, ui terminal.UI) error {
	wd := i.AppDirectory
	if wd == "" {
		wd = profile.WorkingDirectory
	}

	app, appErr := local.LoadAppConfig(wd)
	if appErr != nil {
		return appErr
	}

	if i.AppDirectory == "" {
		if app.RootDir == "" {
			return errProjectNotFound{}
		}
		i.AppDirectory = app.RootDir
	}

	if i.To == "" {
		i.To = app.Option()
	}

	return nil
}

func (i inputs) resolveTo(ui terminal.UI, client realm.Client) (to, error) {
	t := to{GroupID: i.Project}

	if i.To == "" {
		return t, nil
	}

	app, err := cli.ResolveApp(ui, client, realm.AppFilter{GroupID: i.Project, App: i.To})
	if err != nil {
		if _, ok := err.(cli.ErrAppNotFound); !ok {
			return to{}, err
		}
	}

	t.GroupID = app.GroupID
	t.AppID = app.ID
	return t, nil
}

func (i inputs) args(omitDryRun bool) []flags.Arg {
	args := make([]flags.Arg, 0, 7)
	if i.Project != "" {
		args = append(args, flags.Arg{flagProject, i.Project})
	}
	if i.AppDirectory != "" {
		args = append(args, flags.Arg{flagAppDirectory, i.AppDirectory})
	}
	if i.To != "" {
		args = append(args, flags.Arg{flagTo, i.To})
	}
	if i.IncludeDependencies {
		args = append(args, flags.Arg{Name: flagIncludeDependencies})
	}
	if i.IncludeHosting {
		args = append(args, flags.Arg{Name: flagIncludeHosting})
	}
	if i.ResetCDNCache {
		args = append(args, flags.Arg{Name: flagResetCDNCache})
	}
	if i.DryRun && !omitDryRun {
		args = append(args, flags.Arg{Name: flagDryRun})
	}
	return args
}
