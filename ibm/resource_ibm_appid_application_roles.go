package ibm

import (
	"context"
	"fmt"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func resourceIBMAppIDApplicationRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMAppIDApplicationRolesCreate,
		Read:   resourceIBMAppIDApplicationRolesRead,
		Delete: resourceIBMAppIDApplicationRolesDelete,
		Update: resourceIBMAppIDApplicationRolesUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"roles": {
				Description: "A list of role IDs for roles that you want to be assigned to the application (this is different from AppID role access)",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDApplicationRolesCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	roles := expandStringList(d.Get("roles").([]interface{}))

	roleOpts := &appid.PutApplicationsRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: roles,
		},
	}

	_, resp, err := appIDClient.PutApplicationsRolesWithContext(context.TODO(), roleOpts)

	if err != nil {
		return fmt.Errorf("Error setting application roles: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))

	return resourceIBMAppIDApplicationRolesRead(d, meta)
}

func resourceIBMAppIDApplicationRolesRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return fmt.Errorf("Incorrect ID %s: ID should be a combination of tenantID/clientID", d.Id())
	}

	tenantID := idParts[0]
	clientID := idParts[1]

	roles, resp, err := appIDClient.GetApplicationRolesWithContext(context.TODO(), &appid.GetApplicationRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID application '%s' is not found, removing roles from state", clientID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID application roles: %s\n%s", err, resp)
	}

	var appRoles []interface{}

	if roles.Roles != nil {
		for _, v := range roles.Roles {
			appRoles = append(appRoles, *v.ID)
		}
	}

	if err := d.Set("roles", appRoles); err != nil {
		return fmt.Errorf("Error setting application roles: %s", err)
	}

	d.Set("tenant_id", tenantID)
	d.Set("client_id", clientID)

	return nil
}

func resourceIBMAppIDApplicationRolesUpdate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	roles := expandStringList(d.Get("roles").([]interface{}))

	roleOpts := &appid.PutApplicationsRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: roles,
		},
	}

	_, resp, err := appIDClient.PutApplicationsRolesWithContext(context.TODO(), roleOpts)

	if err != nil {
		return fmt.Errorf("Error updating application roles: %s\n%s", err, resp)
	}

	return resourceIBMAppIDApplicationRolesRead(d, meta)
}

func resourceIBMAppIDApplicationRolesDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	roleOpts := &appid.PutApplicationsRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: []string{},
		},
	}

	_, resp, err := appIDClient.PutApplicationsRolesWithContext(context.TODO(), roleOpts)

	if err != nil {
		return fmt.Errorf("Error clearing application roles: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
