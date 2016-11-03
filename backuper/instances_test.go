package backuper_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcf-backup-and-restore/backuper"
	"github.com/pivotal-cf/pcf-backup-and-restore/backuper/fakes"
)

var _ = Describe("Instances", func() {
	Context("AllBackupable", func() {
		var (
			instance1           *fakes.FakeInstance
			instance2           *fakes.FakeInstance
			instance3           *fakes.FakeInstance
			instances           backuper.Instances
			backupableInstances backuper.Instances
			backupableError     error
		)
		BeforeEach(func() {
			instance1 = new(fakes.FakeInstance)
			instance2 = new(fakes.FakeInstance)
			instance3 = new(fakes.FakeInstance)
		})
		JustBeforeEach(func() {
			backupableInstances, backupableError = instances.AllBackupable()
		})
		Context("Single instance, not backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(false, nil)
				instances = backuper.Instances{instance1}
			})
			It("Checks that instance is backupable", func() {
				Expect(instance1.IsBackupableCallCount()).To(Equal(1))
			})
			It("returns no instances", func() {
				Expect(backupableInstances).To(BeEmpty())
			})
			It("doesn't fail", func() {
				Expect(backupableError).NotTo(HaveOccurred())
			})
		})
		Context("Single instance, backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(true, nil)
				instances = backuper.Instances{instance1}
			})
			It("Checks that instance is backupable", func() {
				Expect(instance1.IsBackupableCallCount()).To(Equal(1))
			})
			It("returns the instance", func() {
				Expect(backupableInstances).To(ConsistOf(instance1))
			})
			It("doesn't fail", func() {
				Expect(backupableError).NotTo(HaveOccurred())
			})
		})
		Context("Multiple instances, one backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(false, nil)
				instance2.IsBackupableReturns(false, nil)
				instance3.IsBackupableReturns(true, nil)
				instances = backuper.Instances{instance1, instance2, instance3}
			})
			It("Checks that instance is backupable", func() {
				Expect(instance1.IsBackupableCallCount()).To(Equal(1))
				Expect(instance2.IsBackupableCallCount()).To(Equal(1))
				Expect(instance3.IsBackupableCallCount()).To(Equal(1))
			})
			It("returns true", func() {
				Expect(backupableInstances).To(ConsistOf(instance3))
			})
			It("dosent fail", func() {
				Expect(backupableError).NotTo(HaveOccurred())
			})
		})

		Context("Multiple instances, all backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(true, nil)
				instance2.IsBackupableReturns(true, nil)
				instance3.IsBackupableReturns(true, nil)
				instances = backuper.Instances{instance1, instance2, instance3}
			})
			It("Checks all instances", func() {
				Expect(instance1.IsBackupableCallCount()).To(Equal(1))
				Expect(instance2.IsBackupableCallCount()).To(Equal(1))
				Expect(instance3.IsBackupableCallCount()).To(Equal(1))
			})
			It("returns true", func() {
				Expect(backupableInstances).To(ConsistOf(instance1, instance2, instance3))
			})
			It("dosent fail", func() {
				Expect(backupableError).NotTo(HaveOccurred())
			})
		})
		Context("Multiple instances, none backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(false, nil)
				instance2.IsBackupableReturns(false, nil)
				instance3.IsBackupableReturns(false, nil)
				instances = backuper.Instances{instance1, instance2, instance3}
			})
			It("Checks all instances", func() {
				Expect(instance1.IsBackupableCallCount()).To(Equal(1))
				Expect(instance2.IsBackupableCallCount()).To(Equal(1))
				Expect(instance3.IsBackupableCallCount()).To(Equal(1))
			})
			It("returns false", func() {
				Expect(backupableInstances).To(BeEmpty())
			})
			It("dosent fail", func() {
				Expect(backupableError).NotTo(HaveOccurred())
			})
		})
	})

	Context("IsEmpty", func() {
		It("is true when no instnace", func() {
			Expect(backuper.Instances{}.IsEmpty()).To(BeTrue())
		})

		It("is false if instances", func() {
			Expect(backuper.Instances{new(fakes.FakeInstance)}.IsEmpty()).To(BeFalse())
		})
	})

	Context("Backup", func() {
		var (
			instance1   *fakes.FakeInstance
			instance2   *fakes.FakeInstance
			instance3   *fakes.FakeInstance
			instances   backuper.Instances
			backupError error
		)
		BeforeEach(func() {
			instance1 = new(fakes.FakeInstance)
			instance2 = new(fakes.FakeInstance)
			instance3 = new(fakes.FakeInstance)
		})

		JustBeforeEach(func() {
			backupError = instances.Backup()
		})

		Context("Single instance, backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(true, nil)
				instances = backuper.Instances{instance1}
			})
			It("does not fail", func() {
				Expect(backupError).NotTo(HaveOccurred())
			})
			It("backs up the instance", func() {
				Expect(instance1.BackupCallCount()).To(Equal(1))
			})
		})

		Context("Multiple instances, all backupable", func() {
			BeforeEach(func() {
				instance1.IsBackupableReturns(true, nil)
				instance2.IsBackupableReturns(true, nil)
				instances = backuper.Instances{instance1, instance2}
			})
			It("does not fail", func() {
				Expect(backupError).NotTo(HaveOccurred())
			})
			It("backs up the only the backupable instance", func() {
				Expect(instance1.BackupCallCount()).To(Equal(1))
				Expect(instance2.BackupCallCount()).To(Equal(1))
			})
		})
	})

	Context("Cleanup", func() {
		var (
			instance1    *fakes.FakeInstance
			instance2    *fakes.FakeInstance
			instance3    *fakes.FakeInstance
			instances    backuper.Instances
			cleanupError error
		)
		BeforeEach(func() {
			instance1 = new(fakes.FakeInstance)
			instance2 = new(fakes.FakeInstance)
			instance3 = new(fakes.FakeInstance)
		})
		JustBeforeEach(func() {
			cleanupError = instances.Cleanup()
		})
		Context("single instance", func() {
			BeforeEach(func() {
				instance1.CleanupReturns(nil)
				instances = backuper.Instances{instance1}
			})
			It("calls cleanup", func() {
				Expect(instance1.CleanupCallCount()).To(Equal(1))
			})
			It("dosen't fail", func() {
				Expect(cleanupError).NotTo(HaveOccurred())
			})
		})
		Context("multiple instances", func() {
			BeforeEach(func() {
				instance1.CleanupReturns(nil)
				instance2.CleanupReturns(nil)
				instance3.CleanupReturns(nil)
				instances = backuper.Instances{instance1, instance2, instance3}
			})
			It("calls cleanup on all", func() {
				Expect(instance1.CleanupCallCount()).To(Equal(1))
				Expect(instance2.CleanupCallCount()).To(Equal(1))
				Expect(instance3.CleanupCallCount()).To(Equal(1))
			})
			It("dosen't fail", func() {
				Expect(cleanupError).NotTo(HaveOccurred())
			})
		})

		Context("faliure, single instance", func() {
			var actualError = fmt.Errorf("So Wrong!")
			BeforeEach(func() {
				instance1.CleanupReturns(actualError)
				instances = backuper.Instances{instance1}
			})
			It("calls cleanup", func() {
				Expect(instance1.CleanupCallCount()).To(Equal(1))
			})
			It("fails", func() {
				Expect(cleanupError).To(MatchError(actualError))
			})
		})

		Context("faliure, multiple instance", func() {
			var actualError = fmt.Errorf("the test is rigged!")

			BeforeEach(func() {
				instance1.CleanupReturns(nil)
				instance2.CleanupReturns(actualError)
				instance3.CleanupReturns(nil)
				instances = backuper.Instances{instance1, instance2, instance3}
			})
			It("calls cleanup, till instnace fails", func() {
				Expect(instance1.CleanupCallCount()).To(Equal(1))
				Expect(instance2.CleanupCallCount()).To(Equal(1))
			})
			It("does not call cleanup after that", func() {
				Expect(instance3.CleanupCallCount()).To(Equal(0))
			})
			It("fails", func() {
				Expect(cleanupError).To(MatchError(actualError))
			})
		})
	})
})
