// This file was generated by counterfeiter
package fakes

import (
	"io"
	"sync"

	"github.com/pivotal-cf/pcf-backup-and-restore/backuper"
)

type FakeArtifact struct {
	CreateFileStub        func(backuper.Instance) (io.WriteCloser, error)
	createFileMutex       sync.RWMutex
	createFileArgsForCall []struct {
		arg1 backuper.Instance
	}
	createFileReturns struct {
		result1 io.WriteCloser
		result2 error
	}
	ReadFileStub        func(backuper.Instance) (io.ReadCloser, error)
	readFileMutex       sync.RWMutex
	readFileArgsForCall []struct {
		arg1 backuper.Instance
	}
	readFileReturns struct {
		result1 io.ReadCloser
		result2 error
	}
	AddChecksumStub        func(backuper.Instance, map[string]string) error
	addChecksumMutex       sync.RWMutex
	addChecksumArgsForCall []struct {
		arg1 backuper.Instance
		arg2 map[string]string
	}
	addChecksumReturns struct {
		result1 error
	}
	CalculateChecksumStub        func(backuper.Instance) (map[string]string, error)
	calculateChecksumMutex       sync.RWMutex
	calculateChecksumArgsForCall []struct {
		arg1 backuper.Instance
	}
	calculateChecksumReturns struct {
		result1 map[string]string
		result2 error
	}
	DeploymentMatchesStub        func(string, []backuper.Instance) (bool, error)
	deploymentMatchesMutex       sync.RWMutex
	deploymentMatchesArgsForCall []struct {
		arg1 string
		arg2 []backuper.Instance
	}
	deploymentMatchesReturns struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeArtifact) CreateFile(arg1 backuper.Instance) (io.WriteCloser, error) {
	fake.createFileMutex.Lock()
	fake.createFileArgsForCall = append(fake.createFileArgsForCall, struct {
		arg1 backuper.Instance
	}{arg1})
	fake.recordInvocation("CreateFile", []interface{}{arg1})
	fake.createFileMutex.Unlock()
	if fake.CreateFileStub != nil {
		return fake.CreateFileStub(arg1)
	} else {
		return fake.createFileReturns.result1, fake.createFileReturns.result2
	}
}

func (fake *FakeArtifact) CreateFileCallCount() int {
	fake.createFileMutex.RLock()
	defer fake.createFileMutex.RUnlock()
	return len(fake.createFileArgsForCall)
}

func (fake *FakeArtifact) CreateFileArgsForCall(i int) backuper.Instance {
	fake.createFileMutex.RLock()
	defer fake.createFileMutex.RUnlock()
	return fake.createFileArgsForCall[i].arg1
}

func (fake *FakeArtifact) CreateFileReturns(result1 io.WriteCloser, result2 error) {
	fake.CreateFileStub = nil
	fake.createFileReturns = struct {
		result1 io.WriteCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifact) ReadFile(arg1 backuper.Instance) (io.ReadCloser, error) {
	fake.readFileMutex.Lock()
	fake.readFileArgsForCall = append(fake.readFileArgsForCall, struct {
		arg1 backuper.Instance
	}{arg1})
	fake.recordInvocation("ReadFile", []interface{}{arg1})
	fake.readFileMutex.Unlock()
	if fake.ReadFileStub != nil {
		return fake.ReadFileStub(arg1)
	} else {
		return fake.readFileReturns.result1, fake.readFileReturns.result2
	}
}

func (fake *FakeArtifact) ReadFileCallCount() int {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	return len(fake.readFileArgsForCall)
}

func (fake *FakeArtifact) ReadFileArgsForCall(i int) backuper.Instance {
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	return fake.readFileArgsForCall[i].arg1
}

func (fake *FakeArtifact) ReadFileReturns(result1 io.ReadCloser, result2 error) {
	fake.ReadFileStub = nil
	fake.readFileReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifact) AddChecksum(arg1 backuper.Instance, arg2 map[string]string) error {
	fake.addChecksumMutex.Lock()
	fake.addChecksumArgsForCall = append(fake.addChecksumArgsForCall, struct {
		arg1 backuper.Instance
		arg2 map[string]string
	}{arg1, arg2})
	fake.recordInvocation("AddChecksum", []interface{}{arg1, arg2})
	fake.addChecksumMutex.Unlock()
	if fake.AddChecksumStub != nil {
		return fake.AddChecksumStub(arg1, arg2)
	} else {
		return fake.addChecksumReturns.result1
	}
}

func (fake *FakeArtifact) AddChecksumCallCount() int {
	fake.addChecksumMutex.RLock()
	defer fake.addChecksumMutex.RUnlock()
	return len(fake.addChecksumArgsForCall)
}

func (fake *FakeArtifact) AddChecksumArgsForCall(i int) (backuper.Instance, map[string]string) {
	fake.addChecksumMutex.RLock()
	defer fake.addChecksumMutex.RUnlock()
	return fake.addChecksumArgsForCall[i].arg1, fake.addChecksumArgsForCall[i].arg2
}

func (fake *FakeArtifact) AddChecksumReturns(result1 error) {
	fake.AddChecksumStub = nil
	fake.addChecksumReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifact) CalculateChecksum(arg1 backuper.Instance) (map[string]string, error) {
	fake.calculateChecksumMutex.Lock()
	fake.calculateChecksumArgsForCall = append(fake.calculateChecksumArgsForCall, struct {
		arg1 backuper.Instance
	}{arg1})
	fake.recordInvocation("CalculateChecksum", []interface{}{arg1})
	fake.calculateChecksumMutex.Unlock()
	if fake.CalculateChecksumStub != nil {
		return fake.CalculateChecksumStub(arg1)
	} else {
		return fake.calculateChecksumReturns.result1, fake.calculateChecksumReturns.result2
	}
}

func (fake *FakeArtifact) CalculateChecksumCallCount() int {
	fake.calculateChecksumMutex.RLock()
	defer fake.calculateChecksumMutex.RUnlock()
	return len(fake.calculateChecksumArgsForCall)
}

func (fake *FakeArtifact) CalculateChecksumArgsForCall(i int) backuper.Instance {
	fake.calculateChecksumMutex.RLock()
	defer fake.calculateChecksumMutex.RUnlock()
	return fake.calculateChecksumArgsForCall[i].arg1
}

func (fake *FakeArtifact) CalculateChecksumReturns(result1 map[string]string, result2 error) {
	fake.CalculateChecksumStub = nil
	fake.calculateChecksumReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifact) DeploymentMatches(arg1 string, arg2 []backuper.Instance) (bool, error) {
	var arg2Copy []backuper.Instance
	if arg2 != nil {
		arg2Copy = make([]backuper.Instance, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.deploymentMatchesMutex.Lock()
	fake.deploymentMatchesArgsForCall = append(fake.deploymentMatchesArgsForCall, struct {
		arg1 string
		arg2 []backuper.Instance
	}{arg1, arg2Copy})
	fake.recordInvocation("DeploymentMatches", []interface{}{arg1, arg2Copy})
	fake.deploymentMatchesMutex.Unlock()
	if fake.DeploymentMatchesStub != nil {
		return fake.DeploymentMatchesStub(arg1, arg2)
	} else {
		return fake.deploymentMatchesReturns.result1, fake.deploymentMatchesReturns.result2
	}
}

func (fake *FakeArtifact) DeploymentMatchesCallCount() int {
	fake.deploymentMatchesMutex.RLock()
	defer fake.deploymentMatchesMutex.RUnlock()
	return len(fake.deploymentMatchesArgsForCall)
}

func (fake *FakeArtifact) DeploymentMatchesArgsForCall(i int) (string, []backuper.Instance) {
	fake.deploymentMatchesMutex.RLock()
	defer fake.deploymentMatchesMutex.RUnlock()
	return fake.deploymentMatchesArgsForCall[i].arg1, fake.deploymentMatchesArgsForCall[i].arg2
}

func (fake *FakeArtifact) DeploymentMatchesReturns(result1 bool, result2 error) {
	fake.DeploymentMatchesStub = nil
	fake.deploymentMatchesReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifact) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createFileMutex.RLock()
	defer fake.createFileMutex.RUnlock()
	fake.readFileMutex.RLock()
	defer fake.readFileMutex.RUnlock()
	fake.addChecksumMutex.RLock()
	defer fake.addChecksumMutex.RUnlock()
	fake.calculateChecksumMutex.RLock()
	defer fake.calculateChecksumMutex.RUnlock()
	fake.deploymentMatchesMutex.RLock()
	defer fake.deploymentMatchesMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeArtifact) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ backuper.Artifact = new(FakeArtifact)
