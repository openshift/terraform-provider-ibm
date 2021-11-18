package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppIDRedirectURLs() *schema.Resource {
	return &schema.Resource{
		Description: "Redirect URIs that can be used as callbacks of App ID authentication flow",
		Create:      resourceIBMAppIDRedirectURLsCreate,
		Read:        resourceIBMAppIDRedirectURLsRead,
		Update:      resourceIBMAppIDRedirectURLsUpdate,
		Delete:      resourceIBMAppIDRedirectURLsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service `tenantId`",
			},
			"urls": {
				Description: "A list of redirect URLs",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDRedirectURLsRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	urls, resp, err := appIDClient.GetRedirectUrisWithContext(context.TODO(), &appid.GetRedirectUrisOptions{
		TenantID: &tenantID,
	})
	if err != nil {
		return fmt.Errorf("Error loading AppID Cloud Directory redirect urls: %s\n%s", err, resp)
	}

	if err := d.Set("urls", urls.RedirectUris); err != nil {
		return fmt.Errorf("Error setting AppID Cloud Directory redirect urls: %s", err)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDRedirectURLsCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	urls := d.Get("urls")

	redirectURLs := expandStringList(urls.([]interface{}))
	resp, err := appIDClient.UpdateRedirectUrisWithContext(context.TODO(), &appid.UpdateRedirectUrisOptions{
		TenantID: &tenantID,
		RedirectUrisArray: &appid.RedirectURIConfig{
			RedirectUris: redirectURLs,
		},
	})

	if err != nil {
		return fmt.Errorf("Error updating AppID Cloud Directory redirect URLs: %s\n%s", err, resp)
	}

	d.SetId(tenantID)
	return resourceIBMAppIDRedirectURLsRead(d, meta)
}

func resourceIBMAppIDRedirectURLsUpdate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	urls := d.Get("urls")

	redirectURLs := expandStringList(urls.([]interface{}))
	resp, err := appIDClient.UpdateRedirectUrisWithContext(context.TODO(), &appid.UpdateRedirectUrisOptions{
		TenantID: &tenantID,
		RedirectUrisArray: &appid.RedirectURIConfig{
			RedirectUris: redirectURLs,
		},
	})

	if err != nil {
		return fmt.Errorf("Error updating AppID Cloud Directory redirect URLs: %s\n%s", err, resp)
	}

	return resourceIBMAppIDRedirectURLsRead(d, meta)
}

func resourceIBMAppIDRedirectURLsDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	resp, err := appIDClient.UpdateRedirectUrisWithContext(context.TODO(), &appid.UpdateRedirectUrisOptions{
		TenantID: &tenantID,
		RedirectUrisArray: &appid.RedirectURIConfig{
			RedirectUris: []string{},
		},
	})

	if err != nil {
		return fmt.Errorf("Error resetting AppID Cloud Directory redirect URLs: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
