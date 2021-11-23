package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppIDThemeText() *schema.Resource {
	return &schema.Resource{
		Description: "Update theme texts of the App ID login widget",
		Create:      resourceIBMAppIDThemeTextCreate,
		Read:        resourceIBMAppIDThemeTextRead,
		Update:      resourceIBMAppIDThemeTextUpdate,
		Delete:      resourceIBMAppIDThemeTextDelete,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AppID instance GUID",
			},
			"tab_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"footnote": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIBMAppIDThemeTextRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	text, resp, err := appIDClient.GetThemeTextWithContext(context.TODO(), &appid.GetThemeTextOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing AppID theme text configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID theme text: %s\n%s", err, resp)
	}

	if text.TabTitle != nil {
		d.Set("tab_title", *text.TabTitle)
	}

	if text.Footnote != nil {
		d.Set("footnote", *text.Footnote)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDThemeTextCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeTextOptions{
		TenantID: &tenantID,
		TabTitle: helpers.String(d.Get("tab_title").(string)),
		Footnote: helpers.String(d.Get("footnote").(string)),
	}

	resp, err := appIDClient.PostThemeTextWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error setting AppID theme text: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDThemeTextRead(d, meta)
}

func resourceIBMAppIDThemeTextUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceIBMAppIDThemeTextCreate(d, meta)
}

func resourceIBMAppIDThemeTextDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeTextOptions{
		TenantID: &tenantID,
		TabTitle: helpers.String("Login"),
		Footnote: helpers.String("Powered by App ID"),
	}

	resp, err := appIDClient.PostThemeTextWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error resetting AppID theme text: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
