package ibm

import (
	"context"
	"fmt"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func resourceIBMAppIDUserRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Manage AppID user roles",
		Read:        resourceIBMAppIDUserRolesRead,
		Create:      resourceIBMAppIDUserRolesCreate,
		Delete:      resourceIBMAppIDUserRolesDelete,
		Update:      resourceIBMAppIDUserRolesUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AppID instance GUID",
			},
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user's identifier ('subject' in identity token)",
			},
			"role_ids": {
				Description: "A set of AppID role IDs that should be assigned to the user",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIBMAppIDUserRolesRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return fmt.Errorf("Incorrect ID %s: ID should be a combination of tenantID/subject", id)
	}

	tenantID := idParts[0]
	subject := idParts[1]

	d.Set("tenant_id", tenantID)
	d.Set("subject", subject)

	roles, resp, err := appIDClient.GetUserRolesWithContext(context.TODO(), &appid.GetUserRolesOptions{
		TenantID: &tenantID,
		ID:       &subject,
	})

	if err != nil {
		log.Printf("[DEBUG] Error getting AppID user roles: %s\n%s", err, resp)
		return fmt.Errorf("Error getting AppID user roles: %s", err)
	}

	if roles.Roles != nil {
		if err := d.Set("role_ids", flattenAppIDUserRoleIDs(roles.Roles)); err != nil {
			return fmt.Errorf("Error setting AppID user role_ids: %s", err)
		}
	}

	return nil
}

func resourceIBMAppIDUserRolesCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	subject := d.Get("subject").(string)
	roleIds := d.Get("role_ids").(*schema.Set)

	input := &appid.UpdateUserRolesOptions{
		TenantID: &tenantID,
		ID:       &subject,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: expandStringList(roleIds.List()),
		},
	}

	_, resp, err := appIDClient.UpdateUserRolesWithContext(context.TODO(), input)

	if err != nil {
		log.Printf("[DEBUG] Error updating AppID user roles: %s\n%s", err, resp)
		return fmt.Errorf("Error updating AppID user roles: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, subject))
	return resourceIBMAppIDUserRolesRead(d, meta)
}

func resourceIBMAppIDUserRolesDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	subject := d.Get("subject").(string)

	input := &appid.UpdateUserRolesOptions{
		TenantID: &tenantID,
		ID:       &subject,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: []string{},
		},
	}

	_, resp, err := appIDClient.UpdateUserRolesWithContext(context.TODO(), input)

	if err != nil {
		log.Printf("[DEBUG] Error deleting AppID user roles: %s\n%s", err, resp)
		return fmt.Errorf("Error deleting AppID user roles: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMAppIDUserRolesUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceIBMAppIDUserRolesCreate(d, meta)
}

func flattenAppIDUserRoleIDs(r []appid.GetUserRolesResponseRolesItem) []string {
	result := make([]string, len(r))
	for i, role := range r {
		result[i] = *role.ID
	}
	return result
}
