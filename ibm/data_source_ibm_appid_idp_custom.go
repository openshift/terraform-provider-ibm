package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDIDPCustom() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAppIDIDPCustomRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_key": {
				Description: "This is the public key used to validate your signed JWT. It is required to be a PEM in the RS256 or greater format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDIDPCustomRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	config, resp, err := appIDClient.GetCustomIDPWithContext(context.TODO(), &appid.GetCustomIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error loading AppID custom IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *config.IsActive)

	if config.Config != nil && config.Config.PublicKey != nil {
		if err := d.Set("public_key", *config.Config.PublicKey); err != nil {
			return fmt.Errorf("failed setting config: %s", err)
		}
	}

	d.SetId(tenantID)

	return nil
}
