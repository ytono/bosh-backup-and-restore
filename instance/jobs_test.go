package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/pcf-backup-and-restore/instance"
)

var _ = Describe("Jobs", func() {
	var jobs instance.Jobs
	var scripts instance.BackupAndRestoreScripts
	var artifactNames map[string]string

	BeforeEach(func() {
		artifactNames = map[string]string{}
	})

	JustBeforeEach(func() {
		jobs = instance.NewJobs(scripts, artifactNames)
	})

	Describe("NewJobs", func() {
		Context("when there are two jobs each with a backup script", func() {
			BeforeEach(func() {
				scripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/foo/bin/p-backup",
					"/var/vcap/jobs/bar/bin/p-backup",
				}
			})
			It("groups scripts to create jobs", func() {
				Expect(jobs).To(ConsistOf(
					instance.NewJob(instance.BackupAndRestoreScripts{"/var/vcap/jobs/foo/bin/p-backup"}, ""),
					instance.NewJob(instance.BackupAndRestoreScripts{"/var/vcap/jobs/bar/bin/p-backup"}, ""),
				))
			})
		})

		Context("when there is one job with a backup script", func() {
			BeforeEach(func() {
				scripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/foo/bin/p-backup",
				}
			})
			It("groups scripts to create jobs", func() {
				Expect(jobs).To(ConsistOf(
					instance.NewJob(instance.BackupAndRestoreScripts{"/var/vcap/jobs/foo/bin/p-backup"}, ""),
				))
			})
		})

		Context("when there is one job with a backup script and an blob name", func() {
			BeforeEach(func() {
				scripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/foo/bin/p-backup",
				}
				artifactNames = map[string]string{
					"foo": "a-bosh-backup",
				}
			})

			It("creates a job with the correct blob name", func() {
				Expect(jobs).To(ConsistOf(
					instance.NewJob(
						instance.BackupAndRestoreScripts{"/var/vcap/jobs/foo/bin/p-backup"},
						"a-bosh-backup",
					),
				))
			})
		})

		Context("when there are two jobs, both with backup scripts and unique metadata names", func() {
			BeforeEach(func() {
				scripts = instance.BackupAndRestoreScripts{
					"/var/vcap/jobs/foo/bin/p-backup",
					"/var/vcap/jobs/bar/bin/p-backup",
				}
				artifactNames = map[string]string{
					"foo": "a-bosh-backup",
					"bar": "another-backup",
				}
			})

			It("creates two jobs with the correct blob names", func() {
				Expect(jobs).To(ConsistOf(
					instance.NewJob(
						instance.BackupAndRestoreScripts{"/var/vcap/jobs/foo/bin/p-backup"},
						"a-bosh-backup",
					),
					instance.NewJob(
						instance.BackupAndRestoreScripts{"/var/vcap/jobs/bar/bin/p-backup"},
						"another-backup",
					),
				))
			})
		})

	})

	Context("contains jobs with backup script", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/foo/bin/p-backup",
				"/var/vcap/jobs/bar/bin/p-restore",
			}
		})

		Describe("Backupable", func() {
			It("returns the backupable job", func() {
				Expect(jobs.Backupable()).To(ConsistOf(
					instance.NewJob(instance.BackupAndRestoreScripts{"/var/vcap/jobs/foo/bin/p-backup"}, ""),
				))
			})
		})

		Describe("AnyAreBackupable", func() {
			It("returns true", func() {
				Expect(jobs.AnyAreBackupable()).To(BeTrue())
			})
		})
	})

	Context("contains no jobs with backup script", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/bar/bin/p-restore",
			}
		})

		Describe("Backupable", func() {
			It("returns empty", func() {
				Expect(jobs.Backupable()).To(BeEmpty())
			})
		})

		Describe("AnyAreBackupable", func() {
			It("returns false", func() {
				Expect(jobs.AnyAreBackupable()).To(BeFalse())
			})
		})
	})

	Context("contains jobs with pre-backup-lock scripts", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/foo/bin/p-pre-backup-lock",
				"/var/vcap/jobs/foo/bin/p-backup",
				"/var/vcap/jobs/bar/bin/p-restore",
			}
		})

		Describe("PreBackupable", func() {
			It("returns the lockable job", func() {
				Expect(jobs.PreBackupable()).To(ConsistOf(instance.NewJob(
					instance.BackupAndRestoreScripts{
						"/var/vcap/jobs/foo/bin/p-pre-backup-lock",
						"/var/vcap/jobs/foo/bin/p-backup",
					}, ""),
				))
			})
		})

		Describe("AnyArePreBackupable", func() {
			It("returns true", func() {
				Expect(jobs.AnyArePreBackupable()).To(BeTrue())
			})
		})
	})
	Context("contains no jobs with pre-backup-lock scripts", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/bar/bin/p-restore",
			}
		})
		Describe("PreBackupable", func() {
			It("returns empty", func() {
				Expect(jobs.PreBackupable()).To(BeEmpty())
			})
		})

		Describe("AnyArePreBackupable", func() {
			It("returns false", func() {
				Expect(jobs.AnyArePreBackupable()).To(BeFalse())
			})
		})
	})

	Context("contains jobs with post-backup-lock scripts", func() {

		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/foo/bin/p-backup",
				"/var/vcap/jobs/foo/bin/p-post-backup-unlock",
				"/var/vcap/jobs/bar/bin/p-restore",
			}
		})

		Describe("PostBackupable", func() {
			It("returns the unlockable job", func() {
				Expect(jobs.PostBackupable()).To(ConsistOf(instance.NewJob(
					instance.BackupAndRestoreScripts{
						"/var/vcap/jobs/foo/bin/p-post-backup-unlock",
						"/var/vcap/jobs/foo/bin/p-backup",
					}, ""),
				))
			})
		})

		Describe("AnyArePostBackupable", func() {
			It("returns true", func() {
				Expect(jobs.AnyArePostBackupable()).To(BeTrue())
			})
		})
	})
	Context("contains no jobs with backup script", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/bar/bin/p-restore",
			}
		})

		Describe("PostBackupable", func() {
			It("returns empty", func() {
				Expect(jobs.PostBackupable()).To(BeEmpty())
			})
		})
	})

	Context("contains jobs with restore scripts", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/foo/bin/p-backup",
				"/var/vcap/jobs/foo/bin/p-post-backup-unlock",
				"/var/vcap/jobs/bar/bin/p-restore",
			}
		})

		Describe("Restorable", func() {
			It("returns the unlockable job", func() {
				Expect(jobs.Restorable()).To(ConsistOf(instance.NewJob(
					instance.BackupAndRestoreScripts{"/var/vcap/jobs/bar/bin/p-restore"}, ""),
				))
			})
		})

		Describe("AnyAreRestorable", func() {
			It("returns true", func() {
				Expect(jobs.AnyAreRestorable()).To(BeTrue())
			})
		})
	})

	Context("contains no jobs with backup script", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/bar/bin/p-backup",
			}
		})

		It("returns empty", func() {
			Expect(jobs.Restorable()).To(BeEmpty())
		})
	})

	Context("contains no jobs with named blobs", func() {
		Describe("WithNamedBlobs", func() {
			It("returns empty", func() {
				Expect(jobs.WithNamedBlobs()).To(BeEmpty())
			})
		})

		Describe("NamedBlobs", func() {
			It("returns empty", func() {
				Expect(jobs.NamedBlobs()).To(BeEmpty())
			})
		})
	})

	Context("contains jobs with a named blob", func() {
		BeforeEach(func() {
			artifactNames = map[string]string{
				"bar": "my-cool-blob",
			}
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/bar/bin/p-backup",
				"/var/vcap/jobs/bar/bin/p-restore",
				"/var/vcap/jobs/foo/bin/p-backup",
				"/var/vcap/jobs/baz/bin/p-restore",
			}
		})

		Describe("NamedBlobs", func() {
			It("returns a list of blob names", func() {
				Expect(jobs.NamedBlobs()).To(ConsistOf("my-cool-blob"))
			})
		})

		Describe("WithNamedBlobs", func() {
			It("returns jobs with named blobs", func() {
				Expect(jobs.WithNamedBlobs()).To(ConsistOf(instance.NewJob(
					instance.BackupAndRestoreScripts{
						"/var/vcap/jobs/bar/bin/p-backup",
						"/var/vcap/jobs/bar/bin/p-restore",
					}, "my-cool-blob"),
				))
			})
		})
	})

	Context("contains jobs with multiple named blobs", func() {
		BeforeEach(func() {
			scripts = instance.BackupAndRestoreScripts{
				"/var/vcap/jobs/foo/bin/p-backup",
				"/var/vcap/jobs/bar/bin/p-backup",
			}
			artifactNames = map[string]string{
				"foo": "a-bosh-backup",
				"bar": "another-backup",
			}
		})

		Describe("NamedBlobs", func() {
			It("returns a list of blob names", func() {
				Expect(jobs.NamedBlobs()).To(ConsistOf("a-bosh-backup", "another-backup"))
			})
		})
	})
})