package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDIDPFacebook() *schema.Resource {
	return &schema.Resource{
		Description: "Returns the Facebook identity provider configuration.",
		Read:        dataSourceIBMAppIDIDPFacebookRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Description: "`true` if Facebook IDP configuration is active",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"config": {
				Description: "Facebook IDP configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Facebook application id",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"application_secret": {
							Description: "Facebook application secret",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"redirect_url": {
				Description: "Paste the URI into the Valid OAuth redirect URIs field in the Facebook Login section of the Facebook Developers Portal",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDIDPFacebookRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	fb, resp, err := appIDClient.GetFacebookIDPWithContext(context.TODO(), &appid.GetFacebookIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error loading AppID Facebook IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *fb.IsActive)

	if fb.RedirectURL != nil {
		d.Set("redirect_url", *fb.RedirectURL)
	}

	if fb.Config != nil {
		if err := d.Set("config", flattenIBMAppIDFacebookIDPConfig(fb.Config)); err != nil {
			return fmt.Errorf("Failed setting AppID Facebook IDP config: %s", err)
		}
	}

	d.SetId(tenantID)

	return nil
}

func flattenIBMAppIDFacebookIDPConfig(config *appid.FacebookConfigParamsConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	mConfig := map[string]interface{}{}
	mConfig["application_id"] = *config.IDPID
	mConfig["application_secret"] = *config.Secret

	return []interface{}{mConfig}
}
