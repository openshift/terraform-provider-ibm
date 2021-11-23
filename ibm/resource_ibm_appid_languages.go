package ibm

import (
	"context"
	"fmt"
	"log"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppIDLanguages() *schema.Resource {
	return &schema.Resource{
		Description: "User localization configuration",
		Create:      resourceIBMAppIDLanguagesCreate,
		Read:        resourceIBMAppIDLanguagesRead,
		Delete:      resourceIBMAppIDLanguagesDelete,
		Update:      resourceIBMAppIDLanguagesCreate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"languages": {
				Description: "The list of languages that can be used to customize email templates for Cloud Directory",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDLanguagesRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	langs, resp, err := appIDClient.GetLocalizationWithContext(context.TODO(), &appid.GetLocalizationOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing language configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID languages: %s\n%s", err, resp)
	}

	d.Set("languages", langs.Languages)
	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDLanguagesCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	languages := expandStringList(d.Get("languages").([]interface{}))

	input := &appid.UpdateLocalizationOptions{
		TenantID:  &tenantID,
		Languages: languages,
	}

	resp, err := appIDClient.UpdateLocalizationWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error updating AppID languages: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDLanguagesRead(d, meta)
}

func resourceIBMAppIDLanguagesDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.UpdateLocalizationOptions{
		TenantID:  &tenantID,
		Languages: []string{"en"}, // AppID default
	}

	resp, err := appIDClient.UpdateLocalizationWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error resetting AppID languages: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
