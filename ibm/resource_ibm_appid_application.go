package ibm

import (
	"context"
	"fmt"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"strings"
)

func resourceIBMAppIDApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMAppIDApplicationCreate,
		Read:   resourceIBMAppIDApplicationRead,
		Delete: resourceIBMAppIDApplicationDelete,
		Update: resourceIBMAppIDApplicationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:  "The application name to be registered. Application name cannot exceed 50 characters.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"type": {
				Description:  "The type of application to be registered. Allowed types are `regularwebapp` and `singlepageapp`, default is `regularwebapp`.",
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "regularwebapp",
				ValidateFunc: validation.StringInSlice([]string{"regularwebapp", "singlepageapp"}, false),
			},
			"secret": {
				Description: "The `secret` is a secret known only to the application and the authorization server",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"oauth_server_url": {
				Description: "Base URL for common OAuth endpoints, like `/authorization`, `/token` and `/publickeys`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"profiles_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discovery_endpoint": {
				Description: "This URL returns OAuth Authorization Server Metadata",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceIBMAppIDApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	appName := d.Get("name").(string)
	appType := d.Get("type").(string)

	input := &appid.RegisterApplicationOptions{
		TenantID: &tenantID,
		Name:     &appName,
		Type:     &appType,
	}

	app, resp, err := appIDClient.RegisterApplicationWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error creating AppID application: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *app.ClientID))

	return resourceIBMAppIDApplicationRead(d, meta)
}

func resourceIBMAppIDApplicationRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return fmt.Errorf("Incorrect ID %s: ID should be a combination of tenantID/clientID", d.Id())
	}

	tenantID := idParts[0]
	clientID := idParts[1]

	app, resp, err := appIDClient.GetApplicationWithContext(context.TODO(), &appid.GetApplicationOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID application '%s' is not found, removing from state", clientID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID application: %s\n%s", err, resp)
	}

	if app.Name != nil {
		d.Set("name", *app.Name)
	}

	if app.Secret != nil {
		d.Set("secret", *app.Secret)
	}

	if app.OAuthServerURL != nil {
		d.Set("oauth_server_url", *app.OAuthServerURL)
	}

	if app.ProfilesURL != nil {
		d.Set("profiles_url", *app.ProfilesURL)
	}

	if app.DiscoveryEndpoint != nil {
		d.Set("discovery_endpoint", *app.DiscoveryEndpoint)
	}

	if app.Type != nil {
		d.Set("type", *app.Type)
	}

	d.Set("tenant_id", tenantID)
	d.Set("client_id", clientID)

	return nil
}

func resourceIBMAppIDApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") {
		appIDClient, err := meta.(ClientSession).AppIDAPI()

		if err != nil {
			return err
		}

		tenantID := d.Get("tenant_id").(string)
		appName := d.Get("name").(string)
		clientID := d.Get("client_id").(string)

		_, resp, err := appIDClient.UpdateApplicationWithContext(context.TODO(), &appid.UpdateApplicationOptions{
			TenantID: &tenantID,
			Name:     &appName,
			ClientID: &clientID,
		})

		if err != nil {
			return fmt.Errorf("Error updating AppID application: %s\n%s", err, resp)
		}
	}

	return resourceIBMAppIDApplicationRead(d, meta)
}

func resourceIBMAppIDApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	resp, err := appIDClient.DeleteApplicationWithContext(context.TODO(), &appid.DeleteApplicationOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		return fmt.Errorf("Error deleting AppID application: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
