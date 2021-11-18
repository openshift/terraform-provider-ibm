package ibm

import (
	"context"
	"fmt"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDApplicationRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAppIDApplicationRolesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Required:    true,
			},
			"roles": {
				Description: "Defined roles for an application that is registered with an App ID instance",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Application role ID",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Application role name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDApplicationRolesRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	roles, resp, err := appIDClient.GetApplicationRolesWithContext(context.TODO(), &appid.GetApplicationRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		return fmt.Errorf("Error getting AppID application roles: %s\n%s", err, resp)
	}

	if err := d.Set("roles", flattenAppIDApplicationRoles(roles.Roles)); err != nil {
		return fmt.Errorf("Error setting AppID application roles: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))
	return nil
}

func flattenAppIDApplicationRoles(r []appid.GetUserRolesResponseRolesItem) []interface{} {
	var result []interface{}

	if r == nil {
		return result
	}

	for _, v := range r {
		role := map[string]interface{}{
			"id": *v.ID,
		}

		if v.Name != nil {
			role["name"] = *v.Name
		}

		result = append(result, role)
	}

	return result
}
