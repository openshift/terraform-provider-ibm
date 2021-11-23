package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDMFA() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAppIDMFARead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Description: "`true` if MFA is active",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDMFARead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	mfa, resp, err := appIDClient.GetMFAConfigWithContext(context.TODO(), &appid.GetMFAConfigOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error getting IBM AppID MFA configuration: %s\n%s", err, resp)
	}

	if mfa.IsActive != nil {
		d.Set("is_active", *mfa.IsActive)
	}

	d.SetId(tenantID)

	return nil
}
