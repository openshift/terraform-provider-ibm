package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppIDIDPFacebook() *schema.Resource {
	return &schema.Resource{
		Description: "Update Facebook identity provider configuration.",
		Create:      resourceIBMAppIDIDPFacebookCreate,
		Read:        resourceIBMAppIDIDPFacebookRead,
		Delete:      resourceIBMAppIDIDPFacebookDelete,
		Update:      resourceIBMAppIDIDPFacebookUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"is_active": {
				Description: "`true` if Facebook IDP configuration is active",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"config": {
				Description: "Facebook IDP configuration",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Facebook application id",
							Type:        schema.TypeString,
							Required:    true,
						},
						"application_secret": {
							Description: "Facebook application secret",
							Type:        schema.TypeString,
							Required:    true,
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

func resourceIBMAppIDIDPFacebookRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	fb, resp, err := appIDClient.GetFacebookIDPWithContext(context.TODO(), &appid.GetFacebookIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Facebook IDP configuration from state", tenantID)
			d.SetId("")
			return nil
		}

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

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDIDPFacebookCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	config := &appid.SetFacebookIDPOptions{
		TenantID: &tenantID,
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: &isActive,
		},
	}

	if isActive {
		config.IDP.Config = expandAppIDFBIDPConfig(d.Get("config").([]interface{}))
	}

	_, resp, err := appIDClient.SetFacebookIDPWithContext(context.TODO(), config)

	if err != nil {
		return fmt.Errorf("Error applying AppID Facebook IDP configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDIDPFacebookRead(d, meta)
}

func resourceIBMAppIDIDPFacebookDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	config := appIDFacebookIDPConfigDefaults(tenantID)

	_, resp, err := appIDClient.SetFacebookIDPWithContext(context.TODO(), config)

	if err != nil {
		return fmt.Errorf("Error resetting AppID Facebook IDP configuration: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppIDIDPFacebookUpdate(d *schema.ResourceData, m interface{}) error {
	// since this is configuration we can reuse create method
	return resourceIBMAppIDIDPFacebookCreate(d, m)
}

func expandAppIDFBIDPConfig(cfg []interface{}) *appid.FacebookGoogleConfigParamsConfig {
	config := &appid.FacebookGoogleConfigParamsConfig{}

	if len(cfg) == 0 || cfg[0] == nil {
		return nil
	}

	mCfg := cfg[0].(map[string]interface{})

	config.IDPID = helpers.String(mCfg["application_id"].(string))
	config.Secret = helpers.String(mCfg["application_secret"].(string))

	return config
}

func appIDFacebookIDPConfigDefaults(tenantID string) *appid.SetFacebookIDPOptions {
	return &appid.SetFacebookIDPOptions{
		TenantID: &tenantID,
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: helpers.Bool(false),
		},
	}
}
