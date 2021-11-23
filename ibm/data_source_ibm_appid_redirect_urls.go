package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDRedirectURLs() *schema.Resource {
	return &schema.Resource{
		Description: "Redirect URIs that can be used as callbacks of App ID authentication flow",
		Read:        dataSourceIBMAppIDRedirectURLsRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"urls": {
				Description: "A list of redirect URLs",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDRedirectURLsRead(d *schema.ResourceData, meta interface{}) error {
	appidClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	urls, resp, err := appidClient.GetRedirectUrisWithContext(context.TODO(), &appid.GetRedirectUrisOptions{
		TenantID: &tenantID,
	})
	if err != nil {
		return fmt.Errorf("Error loading Cloud Directory AppID redirect urls: %s\n%s", err, resp)
	}

	if err := d.Set("urls", urls.RedirectUris); err != nil {
		return fmt.Errorf("Error setting Cloud Directory AppID redirect URLs: %s", err)
	}

	d.SetId(tenantID)

	return nil
}
