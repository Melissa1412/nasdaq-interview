// ©Copyright 2022 Metrio
package cloudtasks

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
	"metrio.net/fougere-lite/internal/utils"
)

// Helper method to create client
func getMockedClient(url string) *Client {
	client, err := NewClient(context.Background(), option.WithoutAuthentication(), option.WithEndpoint(url))
	if err != nil {
		Fail(err.Error())
	}
	return client
}

var _ = Describe("Storage client", func() {
	var taskConfig TaskQueue
	var parent string
	var queueName string

	BeforeEach(func() {
		taskConfig = TaskQueue{
			Name:       "queue1",
			Region:     "northamerica-northeast1",
			ProjectId:  "projet-123",
			ClientName: "banane",
		}
		parent = "projects/" + taskConfig.ProjectId + "/locations/" + taskConfig.Region
		queueName = parent + "/queues/" + taskConfig.Name
	})

	Describe("create storage spec", func() {
		It("succesfully creates storage spec", func() {
			mockServerCalls := make(chan utils.MockServerCall, 0)
			mockServer := utils.NewMockServer(mockServerCalls)
			defer mockServer.Close()

			client := getMockedClient(mockServer.URL)

			task := client.createStorageSpec(taskConfig)
			Expect(task.Name).To(Equal(queueName))
		})
	})
	Describe("create queue", func() {
		It("successfully creates the queue", func() {
			mockServerCalls := make(chan utils.MockServerCall, 1)
			mockServerCalls <- utils.MockServerCall{
				UrlMatchFunc: func(url string) bool {
					return strings.HasPrefix(url, "/v2/"+parent)
				},
				Method: "post",
			}
			mockServer := utils.NewMockServer(mockServerCalls)
			defer mockServer.Close()

			client := getMockedClient(mockServer.URL)

			err := client.create(taskConfig)
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Describe("update bucket", func() {
		It("successfully updates the bucket", func() {
			mockServerCalls := make(chan utils.MockServerCall, 1)
			mockServerCalls <- utils.MockServerCall{
				UrlMatchFunc: func(url string) bool {
					return strings.HasPrefix(url, "/v2/"+queueName)
				},
				Method: "patch",
			}
			mockServer := utils.NewMockServer(mockServerCalls)
			defer mockServer.Close()

			client := getMockedClient(mockServer.URL)

			err := client.update(taskConfig)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
