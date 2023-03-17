package influxdbv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func ResourceTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTaskCreate,
		Read:   resourceTaskRead,
		Update: resourceTaskUpdate,
		Delete: resourceTaskDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flux": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"every": {
				Type:     schema.TypeString,
				Required: true,
			},
			"offset": {
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

func resourceTaskCreate(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(influxdb2.Client)

	// https://github.com/influxdata/influxdb-client-go/blob/master/domain/types.gen.go
	desc := d.Get("description").(string)
	every := d.Get("every").(string)
	offset := d.Get("offset").(string)
	// status := domain.TaskStatusTypeInactive
	var status domain.TaskStatusType
	if d.Get("status").(string) == "active" {
		status = domain.TaskStatusTypeActive
	} else {
		status = domain.TaskStatusTypeInactive
	}

	newTask := &domain.Task{
		Name: d.Get("name").(string),
		// Description: d.Get("description").(string),
		Description: &desc,
		Flux:        d.Get("flux").(string),
		OrgID:       d.Get("org_id").(string),
		// Every:       d.Get("every").(*string),
		Every: &every,
		// Offset:      d.Get("offset").(*string),
		Offset: &offset,
		// Status: d.Get("status").(*domain.TaskStatusType),
		Status: &status,
	}
	result, err := influx.TasksAPI().CreateTask(context.Background(), newTask)
	if err != nil {
		return fmt.Errorf("error creating task: %v", err)
	}
	d.SetId(result.Id)
	return resourceTaskRead(d, meta)
}

func resourceTaskRead(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(influxdb2.Client)
	result, err := influx.TasksAPI().GetTaskByID(context.Background(), d.Get("Id").(string))
	if err != nil {
		return fmt.Errorf("error getting task: %v", err)
	}

	d.Set("name", result.Name)
	d.Set("description", result.Description)
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

	updateTask := &domain.Task{
		Id:          d.Get("Id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(*string),
		Flux:        d.Get("flux").(string),
		OrgID:       d.Get("org_id").(string),
		Every:       d.Get("every").(*string),
		Offset:      d.Get("offset").(*string),
		Status:      d.Get("status").(*domain.TaskStatusType),
	}
	_, err := influx.TasksAPI().UpdateTask(context.Background(), updateTask)

	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	return resourceTaskRead(d, meta)
}

func resourceTaskDelete(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(influxdb2.Client)
	err := influx.TasksAPI().DeleteTaskWithID(context.Background(), d.Get("Id").(string))
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}
	d.SetId("")
	return nil
}
