package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"octopus"
	"time"
)

const (
	resourceKeyEnvironmentName        = "name"
	resourceKeyEnvironmentDescription = "description"
	resourceCreateTimeoutEnvironment  = 30 * time.Minute
	resourceUpdateTimeoutEnvironment  = 10 * time.Minute
	resourceDeleteTimeoutEnvironment  = 15 * time.Minute
)

const computedPropertyDescription = "<computed>"

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			resourceKeyEnvironmentName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment name.",
			},
			resourceKeyEnvironmentDescription: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The environment description.",
			},
		},
	}
}

// Create an environment resource.
func resourceEnvironmentCreate(data *schema.ResourceData, provider interface{}) error {
	name := data.Get(resourceKeyEnvironmentName).(string)
	description := data.Get(resourceKeyEnvironmentDescription).(string)

	log.Printf("Create environment named '%s'.", name)

	client := provider.(*octopus.Client)

	environment, err := client.CreateEnvironment(name, description, 0)
	if err != nil {
		return err
	}

	data.SetId(environment.ID)

	return nil
}

// Read an environment resource.
func resourceEnvironmentRead(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	name := data.Get(resourceKeyEnvironmentName).(string)

	log.Printf("Read environment '%s' (name = '%s').", id, name)

	client := provider.(*octopus.Client)
	environment, err := client.GetEnvironment(id)
	if err != nil {
		return err
	}

	if environment == nil {
		// Environment has been deleted.
		data.SetId("")

		return nil
	}

	data.Set(resourceKeyEnvironmentName, environment.Name)
	data.Set(resourceKeyEnvironmentDescription, environment.Description)

	return nil
}

// Update an environment resource.
func resourceEnvironmentUpdate(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()

	log.Printf("Update environment '%s'.", id)

	if !(data.HasChange(resourceKeyEnvironmentName) || data.HasChange(resourceKeyEnvironmentDescription)) {
		return nil // Nothing to do.
	}

	client := provider.(*octopus.Client)
	environment, err := client.GetEnvironment(id)
	if err != nil {
		return err
	}
	if environment != nil {
		// Environment has been deleted.
		data.SetId("")

		return nil
	}

	if data.HasChange(resourceKeyEnvironmentName) {
		environment.Name = data.Get(resourceKeyEnvironmentName).(string)
	}

	if data.HasChange(resourceKeyEnvironmentDescription) {
		environment.Description = data.Get(resourceKeyEnvironmentDescription).(string)
	}

	_, err = client.UpdateEnvironment(environment)

	return err
}

// Delete an environment resource.
func resourceEnvironmentDelete(data *schema.ResourceData, provider interface{}) error {
	id := data.Id()
	name := data.Get(resourceKeyEnvironmentName).(string)

	log.Printf("Delete Environment '%s' (name = '%s').", id, name)

	client := provider.(*octopus.Client)

	return client.DeleteEnvironment(id)
}