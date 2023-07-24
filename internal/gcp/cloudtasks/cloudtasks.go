package cloudtasks

import (
	"context"
	"net/http"

	"google.golang.org/api/cloudtasks/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"metrio.net/fougere-lite/internal/common"
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

func (c *Client) Create(config *Config) error {
	createChannel := make(chan common.Response, len(config.TaskQueues))
	for _, queue := range config.TaskQueues {
		go func(resp chan common.Response, queue TaskQueue) {
			name := "projects/" + queue.ProjectId + "/locations/" + queue.Region + "/queues/" + queue.Name
			_, err := c.get(name)
			if err != nil {
				if e, ok := err.(*googleapi.Error); ok && e.Code == http.StatusNotFound {
					utils.Logger.Debug("[%s] queue not found", name)

					if err := c.create(queue); err != nil {
						resp <- common.Response{Err: err}
						return
					}
				} else {
					utils.Logger.Errorf("[%s] error getting queue: %s", name, err)
					resp <- common.Response{Err: err}
					return
				}
			} else {
				if err := c.update(queue); err != nil {
					resp <- common.Response{Err: err}
					return
				}
			}
			resp <- common.Response{}
		}(createChannel, queue)
	}
	for range config.TaskQueues {
		resp := <-createChannel
		if resp.Err != nil {
			return resp.Err
		}
	}
	return nil
}

func (c *Client) get(name string) (*cloudtasks.Queue, error) {
	utils.Logger.Debug("[%s] getting queue", name)
	queue, err := c.cloudtasksService.Projects.Locations.Queues.Get(name).Do()
	if err != nil {
		return nil, err
	}
	return queue, nil
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
	utils.Logger.Infof("[%s] updating queue", spec.Name)
	_, err := c.cloudtasksService.Projects.Locations.Queues.Patch(spec.Name, spec).Do()
	if err != nil {
		utils.Logger.Errorf("[%s] error updating queue: %s", spec.Name, err)
		return err
	}
	return nil
}

func (c *Client) createStorageSpec(queue TaskQueue) *cloudtasks.Queue {
	return &cloudtasks.Queue{
		Name: "projects/" + queue.ProjectId + "/locations/" + queue.Region + "/queues/" + queue.Name,
		RateLimits: &cloudtasks.RateLimits{
			MaxDispatchesPerSecond:  queue.MaxDispatchesPerSecond,
			MaxConcurrentDispatches: int64(queue.MaxConcurrentDispatches),
		},
		RetryConfig: &cloudtasks.RetryConfig{
			MinBackoff: queue.MinBackoff,
			MaxBackoff: queue.MaxBackoff,
		},
	}
}
