package command

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/cli/flags"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/factory"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type DirectorRestoreCommand struct {
}

func NewDirectorRestoreCommand() DirectorRestoreCommand {
	return DirectorRestoreCommand{}
}

func (cmd DirectorRestoreCommand) Cli() cli.Command {
	return cli.Command{
		Name:    "restore",
		Aliases: []string{"r"},
		Usage:   "Restore a deployment from backup",
		Action:  cmd.Action,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "artifact-path",
				Usage: "Path to the artifact to restore",
			},
		},
	}
}

func (cmd DirectorRestoreCommand) Action(c *cli.Context) error {
	trapSigint(false)

	if err := flags.Validate([]string{"artifact-path"}, c); err != nil {
		return err
	}

	directorName := ExtractNameFromAddress(c.Parent().String("host"))
	artifactPath := c.String("artifact-path")

	restorer := factory.BuildDirectorRestorer(
		c.Parent().String("host"),
		c.Parent().String("username"),
		c.Parent().String("private-key-path"),
		c.GlobalBool("debug"),
	)

	restoreErr := restorer.Restore(directorName, artifactPath)
	return processError(restoreErr)
}

func processError(err orchestrator.Error) error {
	errorCode, errorMessage, errorWithStackTrace := orchestrator.ProcessError(err)
	if err := writeStackTrace(errorWithStackTrace); err != nil {
		return errors.Wrap(err, err.Error())
	}

	return cli.NewExitError(errorMessage, errorCode)
}
