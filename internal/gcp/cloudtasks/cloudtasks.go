package cloudtasks

import (
	"context"

	"google.golang.org/api/cloudtasks/v2"
	"google.golang.org/api/option"
	"metrio.net/fougere-lite/internal/utils"
)

type Client struct {
	cloudtasksService *cloudtasks.Service
}

func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	cloudtasksService, err := cloudtasks.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		cloudtasksService: cloudtasksService,
	}, nil
}

func (c *Client) create(queue TaskQueue) error {
	utils.Logger.Infof("[%s] creating queue", queue.Name)
	spec := c.createStorageSpec(queue)
	parent := "projects/" + queue.ProjectId + "/locations/" + queue.Region
	_, err := c.cloudtasksService.Projects.Locations.Queues.Create(parent, spec).Do()
	if err != nil {
		utils.Logger.Errorf("[%s] error creating queue: %s", queue.Name, err)
		return err
	}

	return nil
}

func (c *Client) update(queue TaskQueue) error {
	spec := c.createStorageSpec(queue)
	utils.Logger.Infof("[%s] updating bucket", spec.Name)
	_, err := c.cloudtasksService.Projects.Locations.Queues.Patch(spec.Name, spec).Do()
	if err != nil {
		utils.Logger.Errorf("[%s] error updating queue: %s", spec.Name, err)
		return err
	}
	return nil
}

func (c *Client) createStorageSpec(queue TaskQueue) *cloudtasks.Queue {
	return &cloudtasks.Queue{
		// AppEngineRoutingOverride: &cloudtasks.AppEngineRouting{},
		Name: "projects/" + queue.ProjectId + "/locations/" + queue.Region + "/queues/" + queue.Name,
		// PurgeTime:                "",
		// RateLimits:               &cloudtasks.RateLimits{},
		// RetryConfig:              &cloudtasks.RetryConfig{},
		// StackdriverLoggingConfig: &cloudtasks.StackdriverLoggingConfig{},
		// State:                    "",
		// ServerResponse:           googleapi.ServerResponse{},
		// ForceSendFields:          []string{},
		// NullFields:               []string{},
	}
}
