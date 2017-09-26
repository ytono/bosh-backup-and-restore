package factory

import (
	"time"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/backup"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orderer"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/standalone"
)

func BuildDirectorBackuper(host,
	username,
	privateKeyPath string,
	hasDebug bool) *orchestrator.Backuper {

	logger := BuildLogger(hasDebug)
	deploymentManager := standalone.NewDeploymentManager(logger,
		host,
		username,
		privateKeyPath,
		instance.NewJobFinder(logger),
		ssh.NewConnection,
	)

	return orchestrator.NewBackuper(backup.BackupDirectoryManager{}, logger, deploymentManager, orderer.NewDirectorLockOrderer(), time.Now)
}
