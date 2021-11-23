package ibm

import (
	"context"
	"fmt"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func resourceIBMAppIDApplicationScopes() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMAppIDApplicationScopesCreate,
		Read:   resourceIBMAppIDApplicationScopesRead,
		Delete: resourceIBMAppIDApplicationScopesDelete,
		Update: resourceIBMAppIDApplicationScopesUpdate,
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
				Required:    true,
				ForceNew:    true,
			},
			"scopes": {
				Description: "A `scope` is a runtime action in your application that you register with IBM Cloud App ID to create an access permission",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDApplicationScopesCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	scopes := expandStringList(d.Get("scopes").([]interface{}))

	scopeOpts := &appid.PutApplicationsScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Scopes:   scopes,
	}

	_, resp, err := appIDClient.PutApplicationsScopesWithContext(context.TODO(), scopeOpts)

	if err != nil {
		return fmt.Errorf("Error setting application scopes: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))

	return resourceIBMAppIDApplicationScopesRead(d, meta)
}

func resourceIBMAppIDApplicationScopesRead(d *schema.ResourceData, meta interface{}) error {
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

	scopes, resp, err := appIDClient.GetApplicationScopesWithContext(context.TODO(), &appid.GetApplicationScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID application '%s' is not found, removing scopes from state", clientID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID application scopes: %s\n%s", err, resp)
	}

	if err := d.Set("scopes", scopes.Scopes); err != nil {
		return fmt.Errorf("Error setting application scopes: %s", err)
	}

	d.Set("tenant_id", tenantID)
	d.Set("client_id", clientID)

	return nil
}

func resourceIBMAppIDApplicationScopesUpdate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	scopes := expandStringList(d.Get("scopes").([]interface{}))

	scopeOpts := &appid.PutApplicationsScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Scopes:   scopes,
	}

	_, resp, err := appIDClient.PutApplicationsScopesWithContext(context.TODO(), scopeOpts)

	if err != nil {
		return fmt.Errorf("Error updating application scopes: %s\n%s", err, resp)
	}

	return resourceIBMAppIDApplicationScopesRead(d, meta)
}

func resourceIBMAppIDApplicationScopesDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	scopeOpts := &appid.PutApplicationsScopesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Scopes:   []string{},
	}

	_, resp, err := appIDClient.PutApplicationsScopesWithContext(context.TODO(), scopeOpts)

	if err != nil {
		return fmt.Errorf("Error clearing application scopes: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
