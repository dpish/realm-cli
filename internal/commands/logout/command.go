package logout

import (
	"github.com/10gen/realm-cli/internal/cli"
	"github.com/10gen/realm-cli/internal/cli/user"
	"github.com/10gen/realm-cli/internal/terminal"
)

// Command is the `logout` command
type Command struct{}

// Handler is the command handler
func (cmd *Command) Handler(profile *user.Profile, ui terminal.UI, clients cli.Clients) error {
	user := profile.Credentials()
	user.PrivateAPIKey = "" // ensures subsequent `login` commands prompt for password

	profile.SetCredentials(user)
	profile.ClearSession()

	if err := profile.Save(); err != nil {
		return err
	}

	ui.Print(terminal.NewTextLog("Successfully logged out"))
	return nil
}