package orchestrator

import (
	"fmt"
	"time"

	"github.com/looplab/fsm"
	"github.com/pkg/errors"
)

type backupWorkflow struct {
	Backuper
	*fsm.FSM

	backupErrors   Error
	deploymentName string
	events         []fsm.EventDesc
	deployment     Deployment
	artifact       Artifact
}

const (
	StateReady            = "ready"
	StateDeploymentExists = "deployment-exists"
	StateHasBackupScript  = "is-backupable"
	StateArtifactCreated  = "artifact-created"
	StateLocked           = "locked"
	StateBackedup         = "backed-up"
	StateUnlocked         = "unlocked"
	StateDrained          = "drained"
	StateFinished         = "finished"
)
const (
	EventCheckDeployment          = "check-deployment-exists"
	EventCheckHasBackupScript     = "check-is-backupable"
	EventCreateEmptyLocalArtifact = "create-artifact"
	EventPrebackupLock            = "pre-backup-lock"
	EventBackup                   = "backup"
	EventPostBackupUnlock         = "post-backup-unlock"
	EventDrain                    = "drain"
	EventCleanup                  = "cleanup"
)

func newBackupCheckWorkflow(backuper Backuper, deploymentName string) *backupWorkflow {
	bw := &backupWorkflow{
		Backuper:       backuper,
		deployment:     nil,
		deploymentName: deploymentName,
		events: fsm.Events{
			{Name: EventCheckDeployment, Src: []string{StateReady}, Dst: StateDeploymentExists},
			{Name: EventCheckHasBackupScript, Src: []string{StateDeploymentExists}, Dst: StateHasBackupScript},
			{Name: EventCleanup, Src: []string{StateDeploymentExists, StateHasBackupScript, StateArtifactCreated, StateUnlocked, StateDrained}, Dst: StateFinished},
		},
	}

	bw.FSM = fsm.NewFSM(
		StateReady,
		bw.events,
		fsm.Callbacks{
			beforeEvent(EventCheckDeployment):      bw.checkDeployment,
			beforeEvent(EventCheckHasBackupScript): bw.checkHasBackupScript,
			EventCleanup:                           bw.cleanup,
		},
	)

	return bw
}

func newBackupWorkflow(backuper Backuper, deploymentName string) *backupWorkflow {
	bw := &backupWorkflow{
		Backuper:       backuper,
		deployment:     nil,
		deploymentName: deploymentName,
		events: fsm.Events{
			{Name: EventCheckDeployment, Src: []string{StateReady}, Dst: StateDeploymentExists},
			{Name: EventCheckHasBackupScript, Src: []string{StateDeploymentExists}, Dst: StateHasBackupScript},
			{Name: EventCreateEmptyLocalArtifact, Src: []string{StateHasBackupScript}, Dst: StateArtifactCreated},
			{Name: EventPrebackupLock, Src: []string{StateArtifactCreated}, Dst: StateLocked},
			{Name: EventBackup, Src: []string{StateLocked}, Dst: StateBackedup},
			{Name: EventPostBackupUnlock, Src: []string{StateBackedup, StateArtifactCreated}, Dst: StateUnlocked},
			{Name: EventDrain, Src: []string{StateUnlocked}, Dst: StateDrained},
			{Name: EventCleanup, Src: []string{StateDeploymentExists, StateHasBackupScript, StateArtifactCreated, StateUnlocked, StateDrained}, Dst: StateFinished},
		},
	}

	bw.FSM = fsm.NewFSM(
		StateReady,
		bw.events,
		fsm.Callbacks{
			beforeEvent(EventCheckDeployment):          bw.checkDeployment,
			beforeEvent(EventCheckHasBackupScript):     bw.checkHasBackupScript,
			beforeEvent(EventCreateEmptyLocalArtifact): bw.createEmptyLocalArtifact,
			beforeEvent(EventPrebackupLock):            bw.prebackupLock,
			beforeEvent(EventBackup):                   bw.backup,
			beforeEvent(EventPostBackupUnlock):         bw.postBackupUnlock,
			beforeEvent(EventDrain):                    bw.drain,
			EventCleanup:                               bw.cleanup,
		},
	)

	return bw
}

func (bw *backupWorkflow) Run() Error {
	for _, e := range bw.events {
		if bw.Can(e.Name) {
			bw.Event(e.Name) //TODO: err
		}
	}
	return bw.backupErrors
}

func (bw *backupWorkflow) checkDeployment(e *fsm.Event) {
	bw.Logger.Info("", "Running pre-checks for backup of %s...\n", bw.deploymentName)

	exists := bw.ArtifactManager.Exists(bw.deploymentName)
	if exists {
		bw.backupErrors = append(bw.backupErrors, errors.Errorf("artifact %s already exists", bw.deploymentName))
		e.Cancel()
		return
	}

	deployment, err := bw.DeploymentManager.Find(bw.deploymentName)
	if err != nil {
		bw.backupErrors = append(bw.backupErrors, err)
		e.Cancel()
		return
	}

	bw.deployment = deployment
}

func (bw *backupWorkflow) checkHasBackupScript(e *fsm.Event) {
	if !bw.deployment.HasBackupScript() {
		bw.backupErrors = append(bw.backupErrors, errors.Errorf("Deployment '%s' has no backup scripts", bw.deploymentName))
		e.Cancel()
		return
	}

	if bw.deployment.ArtifactDirExists() {
		bw.backupErrors = append(bw.backupErrors, errors.Errorf("Deployment '%s' - %s already exists", bw.deploymentName, ArtifactDirectory))
		e.Cancel()
		return
	}

	if !bw.deployment.HasUniqueCustomBackupNames() {
		bw.backupErrors = append(bw.backupErrors, errors.Errorf("Multiple jobs in deployment '%s' specified the same backup name", bw.deploymentName))
		e.Cancel()
	}

	if err := bw.deployment.CustomArtifactNamesMatch(); err != nil {
		bw.backupErrors = append(bw.backupErrors, err)
		e.Cancel()
	}
}

func (bw *backupWorkflow) cleanup(e *fsm.Event) {
	if err := bw.deployment.Cleanup(); err != nil {
		bw.backupErrors = append(bw.backupErrors, NewCleanupError(fmt.Sprintf("Deployment '%s' failed while cleaning up with error: %v", bw.deploymentName, err)))
	}
}

func (bw *backupWorkflow) createEmptyLocalArtifact(e *fsm.Event) {
	bw.Logger.Info("", "Starting backup of %s...\n", bw.deploymentName)
	var err error
	bw.artifact, err = bw.ArtifactManager.Create(bw.deploymentName, bw.Logger)
	if err != nil {
		bw.backupErrors = append(bw.backupErrors, err)
		e.Cancel()
		return
	}

	err = bw.DeploymentManager.SaveManifest(bw.deploymentName, bw.artifact)
	if err != nil {
		bw.backupErrors = append(bw.backupErrors, err)
		e.Cancel()
		return
	}

}

func (bw *backupWorkflow) prebackupLock(e *fsm.Event) {
	err := bw.deployment.PreBackupLock()

	if err != nil {
		bw.backupErrors = append(bw.backupErrors, NewLockError(err.Error()))
		e.Cancel()
	}
}

func (bw *backupWorkflow) backup(e *fsm.Event) {
	err := bw.deployment.Backup()

	if err != nil {
		bw.backupErrors = append(bw.backupErrors, NewBackupError(err.Error()))
	}
}

func (bw *backupWorkflow) postBackupUnlock(e *fsm.Event) {
	err := bw.deployment.PostBackupUnlock()

	if err != nil {
		bw.backupErrors = append(bw.backupErrors, NewPostBackupUnlockError(err.Error()))
	}
}

func (bw *backupWorkflow) drain(e *fsm.Event) {
	if bw.backupErrors.IsFatal() { // TODO: how do we remove this?
		e.Cancel()
		return
	}
	err := bw.deployment.CopyRemoteBackupToLocal(bw.artifact)

	if err != nil {
		bw.backupErrors = append(bw.backupErrors, err)
		return
	}

	bw.Logger.Info("", "Backup created of %s on %v\n", bw.deploymentName, time.Now())
}

func beforeEvent(eventName string) string {
	return "before_" + eventName
}
