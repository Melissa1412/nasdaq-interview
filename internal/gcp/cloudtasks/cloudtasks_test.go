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

	BeforeEach(func() {
		taskConfig = TaskQueue{
			Name:       "queue1",
			Region:     "northamerica-northeast1",
			ProjectId:  "projet-123",
			ClientName: "banane",
		}
	})
	Describe("create queue", func() {
		It("successfully creates the queue", func() {
			mockServerCalls := make(chan utils.MockServerCall, 1)
			mockServerCalls <- utils.MockServerCall{
				UrlMatchFunc: func(url string) bool {
					return strings.HasPrefix(url, "/v2/projects/projet-123/locations/northamerica-northeast1")
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
					return strings.HasPrefix(url, "/v2/projects/projet-123/locations/northamerica-northeast1/queues/queue1")
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
