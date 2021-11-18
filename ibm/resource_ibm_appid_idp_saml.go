package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceIBMAppIDIDPSAML() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMAppIDIDPSAMLCreate,
		Read:   resourceIBMAppIDIDPSAMLRead,
		Delete: resourceIBMAppIDIDPSAMLDelete,
		Update: resourceIBMAppIDIDPSAMLUpdate,
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
				Description: "SAML IDP activation",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"config": {
				Description: "SAML IDP configuration",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_id": {
							Description: "Unique name for an Identity Provider",
							Type:        schema.TypeString,
							Required:    true,
						},
						"sign_in_url": {
							Description: "SAML SSO url",
							Type:        schema.TypeString,
							Required:    true,
						},
						"certificates": {
							Description: "List of certificates, primary and optional secondary",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MaxItems: 2,
							Required: true,
						},
						"display_name": {
							Description: "Provider name",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"encrypt_response": {
							Description: "`true` if SAML responses should be encrypted",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"sign_request": {
							Description: "`true` if SAML requests should be signed",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"include_scoping": {
							Description: "`true` if scopes are included",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"authn_context": {
							Description: "SAML authNContext configuration",
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"class": {
										Description: "List of `authnContext` classes",
										Type:        schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional: true,
									},
									"comparison": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"exact", "maximum", "minimum", "better"}, false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceIBMAppIDIDPSAMLRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	saml, resp, err := appIDClient.GetSAMLIDPWithContext(context.TODO(), &appid.GetSAMLIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing IDP SAML from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading AppID SAML IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *saml.IsActive)

	if saml.Config != nil {
		if err := d.Set("config", flattenAppIDIDPSAMLConfig(saml.Config)); err != nil {
			return fmt.Errorf("failed setting config: %s", err)
		}
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDIDPSAMLCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	config := &appid.SetSAMLIDPOptions{
		TenantID: &tenantID,
		IsActive: &isActive,
	}

	if isActive {
		if cfg, ok := d.GetOk("config"); ok {
			config.Config = expandAppIDIDPSAMLConfig(cfg.([]interface{}))
		}
	}

	_, resp, err := appIDClient.SetSAMLIDPWithContext(context.TODO(), config)

	if err != nil {
		return fmt.Errorf("Error applying SAML IDP configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDIDPSAMLRead(d, meta)
}

func expandAppIDIDPSAMLAuthNContext(ctx []interface{}) *appid.SAMLConfigParamsAuthnContext {
	authNContext := &appid.SAMLConfigParamsAuthnContext{}

	if len(ctx) == 0 || ctx[0] == nil {
		return nil
	}

	mContext := ctx[0].(map[string]interface{})

	if comparison, ok := mContext["comparison"]; ok {
		authNContext.Comparison = helpers.String(comparison.(string))
	}

	if class, ok := mContext["class"]; ok {
		authNContext.Class = expandStringList(class.([]interface{}))
	}

	return authNContext
}

func expandAppIDIDPSAMLConfig(cfg []interface{}) *appid.SAMLConfigParams {
	config := &appid.SAMLConfigParams{}

	if len(cfg) == 0 || cfg[0] == nil {
		return nil
	}

	mCfg := cfg[0].(map[string]interface{})

	config.EntityID = helpers.String(mCfg["entity_id"].(string))
	config.SignInURL = helpers.String(mCfg["sign_in_url"].(string))

	if dispName, ok := mCfg["display_name"]; ok {
		config.DisplayName = helpers.String(dispName.(string))
	}

	if encResponse, ok := mCfg["encrypt_response"]; ok {
		config.EncryptResponse = helpers.Bool(encResponse.(bool))
	}

	if signRequest, ok := mCfg["sign_request"]; ok {
		config.SignRequest = helpers.Bool(signRequest.(bool))
	}

	if includeScoping, ok := mCfg["include_scoping"]; ok {
		config.IncludeScoping = helpers.Bool(includeScoping.(bool))
	}

	if certificates, ok := mCfg["certificates"]; ok {
		config.Certificates = []string{}

		for _, cert := range certificates.([]interface{}) {
			if cert != nil {
				config.Certificates = append(config.Certificates, cert.(string))
			}
		}
	}

	if ctx, ok := mCfg["authn_context"]; ok {
		config.AuthnContext = expandAppIDIDPSAMLAuthNContext(ctx.([]interface{}))
	}

	return config
}

func appIDIDPSAMLConfigDefaults(tenantID string) *appid.SetSAMLIDPOptions {
	return &appid.SetSAMLIDPOptions{
		IsActive: helpers.Bool(false),
		TenantID: helpers.String(tenantID),
	}
}

func resourceIBMAppIDIDPSAMLDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	config := appIDIDPSAMLConfigDefaults(tenantID)

	_, resp, err := appIDClient.SetSAMLIDPWithContext(context.TODO(), config)

	if err != nil {
		return fmt.Errorf("Error resetting SAML IDP configuration: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppIDIDPSAMLUpdate(d *schema.ResourceData, m interface{}) error {
	// since this is configuration we can reuse create method
	return resourceIBMAppIDIDPSAMLCreate(d, m)
}
