package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDAuditStatus() *schema.Resource {
	return &schema.Resource{
		Description: "Tenant audit status",
		Read:        dataSourceIBMAppIDAuditStatusRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The auditing status of the tenant.",
			},
		},
	}
}

func dataSourceIBMAppIDAuditStatusRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	auditStatus, resp, err := appIDClient.GetAuditStatusWithContext(context.TODO(), &appid.GetAuditStatusOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error getting AppID audit status: %s\n%s", err, resp)
	}

	d.Set("is_active", *auditStatus.IsActive)
	d.SetId(tenantID)

	return nil
}
