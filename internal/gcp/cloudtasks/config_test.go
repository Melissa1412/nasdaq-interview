package cloudtasks

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var validTaskConfig = []byte(`
cloudTasks:
  queue1:
    region: us-central1
    projectId: some-project
    minBackoff: 1s
    maxBackoff: 10s
    maxConcurrentDispatches: 1000
    maxDispatchesPerSecond: 500.0`)

var invalidConfig = []byte(`
cloudTasks:
  some-queue:
    region: 
      - should_not_be_an_array`)

var _ = Describe("config", func() {
	BeforeEach(func() {
		viper.Reset()
		viper.SetConfigType("yaml")
	})
	Describe("GetTaskConfig", func() {
		It("should successfully parse a task queue config", func() {
			err := viper.ReadConfig(bytes.NewBuffer(validTaskConfig))
			Expect(err).ToNot(HaveOccurred())
			taskConfig, err := GetTaskConfig(viper.GetViper(), "some-client")
			Expect(err).To(BeNil())
			Expect(len(taskConfig.TaskQueues)).To(Equal(1))
			queue := taskConfig.TaskQueues["queue1"]
			Expect(queue.Region).To(Equal("us-central1"))
			Expect(queue.ProjectId).To(Equal("some-project"))
			Expect(queue.MinBackoff).To(Equal("1s"))
			Expect(queue.MaxBackoff).To(Equal("10s"))
			Expect(queue.MaxConcurrentDispatches).To(Equal(1000))
			Expect(queue.MaxDispatchesPerSecond).To(Equal(500.0))
		})
		It("returns an error if cannot parse the config", func() {
			err := viper.ReadConfig(bytes.NewBuffer(invalidConfig))
			Expect(err).ToNot(HaveOccurred())
			_, err = GetTaskConfig(viper.GetViper(), "some-client")
			Expect(err).NotTo(BeNil())
		})
	})
	Context("validates storage buckets", func() {
		It("should not detect error", func() {
			config := &Config{
				TaskQueues: map[string]TaskQueue{
					"foooo": {
						Region:                  "us-central1",
						ProjectId:               "mock-project",
						Name:                    "foooo",
						MinBackoff:              "1s",
						MaxBackoff:              "100s",
						MaxConcurrentDispatches: 100,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should detect empty name storage buckets", func() {
			config := &Config{
				TaskQueues: map[string]TaskQueue{
					"foooo": {
						Region:                  "us-central1",
						ProjectId:               "mock-project",
						MinBackoff:              "1s",
						MaxBackoff:              "100s",
						MaxConcurrentDispatches: 100,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})
		It("should detect an empty region", func() {
			config := &Config{
				TaskQueues: map[string]TaskQueue{
					"foooo": {
						ProjectId:               "mock-project",
						Name:                    "foooo",
						MinBackoff:              "1s",
						MaxBackoff:              "100s",
						MaxConcurrentDispatches: 100,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})
		It("should detect a missing project id", func() {
			config := &Config{
				TaskQueues: map[string]TaskQueue{
					"foooo": {
						Region:                  "us-central1",
						Name:                    "foooo",
						MinBackoff:              "1s",
						MaxBackoff:              "100s",
						MaxConcurrentDispatches: 100,
						MaxDispatchesPerSecond:  100.0,
					},
				},
			}
			err := ValidateConfig(config)
			Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		})
		// It("should detect a missing min backoff time", func() {
		// 	config := &Config{
		// 		TaskQueues: map[string]TaskQueue{
		// 			"foooo": {
		// 				Region:                  "us-central1",
		// 				ProjectId:               "mock-project",
		// 				Name:                    "foooo",
		// 				MaxBackoff:              "100s",
		// 				MaxConcurrentDispatches: 100,
		// 				MaxDispatchesPerSecond:  100.0,
		// 			},
		// 		},
		// 	}
		// 	err := ValidateConfig(config)
		// 	Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		// })
		// It("should detect a missing max backoff time", func() {
		// 	config := &Config{
		// 		TaskQueues: map[string]TaskQueue{
		// 			"foooo": {
		// 				Region:                  "us-central1",
		// 				ProjectId:               "mock-project",
		// 				Name:                    "foooo",
		// 				MinBackoff:              "1s",
		// 				MaxConcurrentDispatches: 100,
		// 				MaxDispatchesPerSecond:  100.0,
		// 			},
		// 		},
		// 	}
		// 	err := ValidateConfig(config)
		// 	Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		// })
		// It("should detect missing max concurrent dispatches field", func() {
		// 	config := &Config{
		// 		TaskQueues: map[string]TaskQueue{
		// 			"foooo": {
		// 				Region:                  "us-central1",
		// 				ProjectId:               "mock-project",
		// 				Name:                    "foooo",
		// 				MinBackoff:              "1s",
		// 				MaxBackoff:              "100s",
		// 				MaxDispatchesPerSecond:  100.0,
		// 			},
		// 		},
		// 	}
		// 	err := ValidateConfig(config)
		// 	Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		// })
		// It("should detect missing max dispatches per second field", func() {
		// 	config := &Config{
		// 		TaskQueues: map[string]TaskQueue{
		// 			"foooo": {
		// 				Region:                  "us-central1",
		// 				ProjectId:               "mock-project",
		// 				Name:                    "foooo",
		// 				MinBackoff:              "1s",
		// 				MaxBackoff:              "100s",
		// 				MaxConcurrentDispatches: 100,
		// 			},
		// 		},
		// 	}
		// 	err := ValidateConfig(config)
		// 	Expect(err).Should(MatchError(ContainSubstring("validate failed on the required rule")))
		// })
	})
})
