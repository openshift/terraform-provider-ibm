package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDLanguages() *schema.Resource {
	return &schema.Resource{
		Description: "User localization configuration",
		Read:        dataSourceIBMAppIDLanguagesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"languages": {
				Description: "The list of languages that can be used to customize email templates for Cloud Directory",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDLanguagesRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	langs, resp, err := appIDClient.GetLocalizationWithContext(context.TODO(), &appid.GetLocalizationOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error getting AppID languages: %s\n%s", err, resp)
	}

	d.Set("languages", langs.Languages)
	d.SetId(tenantID)

	return nil
}
