package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const defaultHeaderColor = "#EEF2F5" // AppID default

func resourceIBMAppIDThemeColor() *schema.Resource {
	return &schema.Resource{
		Description: "Colors of the App ID login widget",
		Create:      resourceIBMAppIDThemeColorCreate,
		Update:      resourceIBMAppIDThemeColorUpdate,
		Read:        resourceIBMAppIDThemeColorRead,
		Delete:      resourceIBMAppIDThemeColorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AppID instance GUID",
			},
			"header_color": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDThemeColorRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	colors, resp, err := appIDClient.GetThemeColorWithContext(context.TODO(), &appid.GetThemeColorOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing AppID theme color configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID theme colors: %s\n%s", err, resp)
	}

	if colors.HeaderColor != nil {
		d.Set("header_color", *colors.HeaderColor)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDThemeColorCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeColorOptions{
		TenantID:    &tenantID,
		HeaderColor: helpers.String(d.Get("header_color").(string)),
	}

	resp, err := appIDClient.PostThemeColorWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error setting AppID theme color: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDThemeColorRead(d, meta)
}

func resourceIBMAppIDThemeColorUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceIBMAppIDThemeColorCreate(d, meta)
}

func resourceIBMAppIDThemeColorDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeColorOptions{
		TenantID:    &tenantID,
		HeaderColor: helpers.String(defaultHeaderColor),
	}

	resp, err := appIDClient.PostThemeColorWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error resetting AppID theme color: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
