package influxdbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func ResourceTaskByFlux() *schema.Resource {
	return &schema.Resource{
		Create: resourceTaskByFluxCreate,
		Read:   resourceTaskRead,
		Update: resourceTaskUpdate,
		Delete: resourceTaskDelete,
		Schema: map[string]*schema.Schema{
			"flux": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTaskByFluxCreate(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(influxdb2.Client)

	// https://github.com/influxdata/influxdb-client-go/blob/master/domain/types.gen.go
	flux := d.Get("flux").(string)
	orgId := d.Get("org_id").(string)
	result, err := influx.TasksAPI().CreateTaskByFlux(context.Background(), flux, orgId)
	if err != nil {
		return fmt.Errorf("error creating task: %v", err)
	}
	d.SetId(result.Id)

	// START: doing update after creation because status was not set during creation
	resourceTaskUpdate(d, meta)
	// END: doing update after creation because status was not set during creation

	return resourceTaskRead(d, meta)
}

func resourceTaskRead(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(influxdb2.Client)
	result, err := influx.TasksAPI().GetTaskByID(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error getting task: %v", err)
	}

	d.Set("name", result.Name)
	d.Set("org_id", result.OrgID)
	d.Set("flux", result.Flux)
	d.Set("every", result.Every)
	d.Set("offset", result.Offset)
	d.Set("status", result.Status)
	d.Set("created_at", result.CreatedAt.String())
	d.Set("updated_at", result.UpdatedAt.String())

	return nil
}

func resourceTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(influxdb2.Client)

	var status domain.TaskStatusType
	if d.Get("status").(string) == "active" {
		status = domain.TaskStatusTypeActive
	} else {
		status = domain.TaskStatusTypeInactive
	}

	updateTask := &domain.Task{
		Id:     d.Id(),
		Flux:   d.Get("flux").(string),
		OrgID:  d.Get("org_id").(string),
		Status: &status,
	}
	_, err := influx.TasksAPI().UpdateTask(context.Background(), updateTask)

	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	return resourceTaskRead(d, meta)
}

func resourceTaskDelete(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(influxdb2.Client)
	err := influx.TasksAPI().DeleteTaskWithID(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}
	d.SetId("")
	return nil
}
