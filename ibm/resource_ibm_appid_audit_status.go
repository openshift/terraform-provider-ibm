package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppIDAuditStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMAppIDAuditStatusCreate,
		Read:   resourceIBMAppIDAuditStatusRead,
		Delete: resourceIBMAppIDAuditStatusDelete,
		Update: resourceIBMAppIDAuditStatusUpdate,
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
			"is_active": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The auditing status of the tenant.",
			},
		},
	}
}

func resourceIBMAppIDAuditStatusRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	auditStatus, resp, err := appIDClient.GetAuditStatusWithContext(context.TODO(), &appid.GetAuditStatusOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing audit status configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting AppID audit status: %s\n%s", err, resp)
	}

	d.Set("is_active", *auditStatus.IsActive)
	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDAuditStatusCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	resp, err := appIDClient.SetAuditStatusWithContext(context.TODO(), &appid.SetAuditStatusOptions{
		TenantID: &tenantID,
		IsActive: &isActive,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing audit status from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error setting AppID audit status: %s\n%s", err, resp)
	}

	d.SetId(tenantID)
	return resourceIBMAppIDAuditStatusRead(d, meta)
}

func resourceIBMAppIDAuditStatusDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	resp, err := appIDClient.SetAuditStatusWithContext(context.TODO(), &appid.SetAuditStatusOptions{
		TenantID: &tenantID,
		IsActive: helpers.Bool(false),
	})

	if err != nil {
		return fmt.Errorf("Error resetting AppID audit status: %s\n%s", err, resp)
	}

	d.SetId("")
	return nil
}

func resourceIBMAppIDAuditStatusUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceIBMAppIDAuditStatusCreate(d, m)
}
