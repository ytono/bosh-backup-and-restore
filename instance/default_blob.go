package instance

import (
	"fmt"
	"io"

	"github.com/pivotal-cf/pcf-backup-and-restore/orchestrator"
)

func NewDefaultBlob(instance orchestrator.Instance, sshConn SSHConnection, logger Logger) *DefaultBlob {
	return &DefaultBlob{
		Instance:      instance,
		SSHConnection: sshConn,
		Logger:        logger,
	}
}

type DefaultBlob struct {
	orchestrator.Instance
	SSHConnection
	Logger
}

func (d *DefaultBlob) StreamFromRemote(writer io.Writer) error {
	d.Logger.Debug("", "Streaming backup from instance %s/%s", d.Name(), d.ID())
	stderr, exitCode, err := d.Stream("sudo tar -C /var/vcap/store/backup -zc .", writer)

	d.Logger.Debug("", "Stderr: %s", string(stderr))

	if err != nil {
		d.Logger.Debug("", "Error running instance backup scripts. Exit code %d, error %s", exitCode, err.Error())
	}

	if exitCode != 0 {
		return fmt.Errorf("Instance backup scripts returned %d. Error: %s", exitCode, stderr)
	}

	return err
}

func (d *DefaultBlob) StreamBackupToRemote(reader io.Reader) error {
	stdout, stderr, exitCode, err := d.logAndRun("sudo mkdir -p /var/vcap/store/backup/", "create backup directory on remote")

	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("Creating backup directory on the remote returned %d. Error: %s", exitCode, stderr)
	}

	d.Logger.Debug("", "Streaming backup to instance %s/%s", d.Instance.Name(), d.Instance.ID())
	stdout, stderr, exitCode, err = d.StreamStdin("sudo sh -c 'tar -C /var/vcap/store/backup -zx'", reader)

	d.Logger.Debug("", "Stdout: %s", string(stdout))
	d.Logger.Debug("", "Stderr: %s", string(stderr))

	if err != nil {
		d.Logger.Debug("", "Error streaming backup to remote instance. Exit code %d, error %s", exitCode, err.Error())
	}

	if exitCode != 0 {
		return fmt.Errorf("Streaming backup to remote returned %d. Error: %s", exitCode, stderr)
	}

	return err
}

func (d *DefaultBlob) logAndRun(cmd, label string) ([]byte, []byte, int, error) {
	d.Logger.Debug("", "Running %s on %s/%s", label, d.Name(), d.ID())

	stdout, stderr, exitCode, err := d.Run(cmd)
	d.Logger.Debug("", "Stdout: %s", string(stdout))
	d.Logger.Debug("", "Stderr: %s", string(stderr))

	if err != nil {
		d.Logger.Debug("", "Error running %s on instance %s/%s. Exit code %d, error: %s", label, d.Name(), d.ID(), exitCode, err.Error())
	}

	return stdout, stderr, exitCode, err
}

func (d *DefaultBlob) Delete() error {
	return nil
}
