package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDThemeColor() *schema.Resource {
	return &schema.Resource{
		Description: "Colors of the App ID login widget",
		Read:        dataSourceIBMAppIDThemeColorRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"header_color": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDThemeColorRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	colors, resp, err := appIDClient.GetThemeColorWithContext(context.TODO(), &appid.GetThemeColorOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error getting AppID theme colors: %s\n%s", err, resp)
	}

	if colors.HeaderColor != nil {
		d.Set("header_color", *colors.HeaderColor)
	}

	d.SetId(tenantID)

	return nil
}
