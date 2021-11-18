package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDThemeText() *schema.Resource {
	return &schema.Resource{
		Description: "The theme texts of the App ID login widget",
		Read:        dataSourceIBMAppIDThemeTextRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"tab_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"footnote": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDThemeTextRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	text, resp, err := appIDClient.GetThemeTextWithContext(context.TODO(), &appid.GetThemeTextOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error getting AppID theme text: %s\n%s", err, resp)
	}

	if text.TabTitle != nil {
		d.Set("tab_title", *text.TabTitle)
	}

	if text.Footnote != nil {
		d.Set("footnote", *text.Footnote)
	}

	d.SetId(tenantID)

	return nil
}
