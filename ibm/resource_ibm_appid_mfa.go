package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppIDMFA() *schema.Resource {
	return &schema.Resource{
		Read:   resourceIBMAppIDMFARead,
		Create: resourceIBMAppIDMFACreate,
		Update: resourceIBMAppIDMFACreate,
		Delete: resourceIBMAppIDMFADelete,
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
				Description: "`true` if MFA is active",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
	}
}

func resourceIBMAppIDMFARead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	mfa, resp, err := appIDClient.GetMFAConfigWithContext(context.TODO(), &appid.GetMFAConfigOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing AppID MFA configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID MFA configuration: %s\n%s", err, resp)
	}

	if mfa.IsActive != nil {
		d.Set("is_active", *mfa.IsActive)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDMFACreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	input := &appid.UpdateMFAConfigOptions{
		TenantID: &tenantID,
		IsActive: &isActive,
	}

	_, resp, err := appIDClient.UpdateMFAConfigWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error updating AppID MFA configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDMFARead(d, meta)
}

func resourceIBMAppIDMFADelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.UpdateMFAConfigOptions{
		TenantID: &tenantID,
		IsActive: helpers.Bool(false),
	}

	_, resp, err := appIDClient.UpdateMFAConfigWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error resetting AppID MFA configuration: %s\n%s", err, resp)
	}

	d.SetId("")
	return nil
}
