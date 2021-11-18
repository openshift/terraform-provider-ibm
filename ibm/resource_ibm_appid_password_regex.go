package ibm

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMAppIDPasswordRegex() *schema.Resource {
	return &schema.Resource{
		Description: "The regular expression used by App ID for password strength validation",
		Create:      resourceIBMAppIDPasswordRegexCreate,
		Read:        resourceIBMAppIDPasswordRegexRead,
		Delete:      resourceIBMAppIDPasswordRegexDelete,
		Update:      resourceIBMAppIDPasswordRegexUpdate,
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
			"base64_encoded_regex": {
				Description: "The regex expression rule for acceptable password encoded in base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"error_message": {
				Description: "Custom error message",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"regex": {
				Description: "The escaped regex expression rule for acceptable password",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceIBMAppIDPasswordRegexRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Id()

	pw, resp, err := appIDClient.GetCloudDirectoryPasswordRegexWithContext(context.TODO(), &appid.GetCloudDirectoryPasswordRegexOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Password Regex from state", tenantID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading AppID Cloud Directory password regex: %s\n%s", err, resp)
	}

	if pw.Base64EncodedRegex != nil {
		d.Set("base64_encoded_regex", *pw.Base64EncodedRegex)
	}

	if pw.Regex != nil {
		d.Set("regex", *pw.Regex)
	}

	if pw.ErrorMessage != nil {
		d.Set("error_message", *pw.ErrorMessage)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDPasswordRegexCreate(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)
	regex := d.Get("regex").(string)

	input := &appid.SetCloudDirectoryPasswordRegexOptions{
		TenantID:           &tenantID,
		Base64EncodedRegex: helpers.String(b64.StdEncoding.EncodeToString([]byte(regex))),
	}

	if msg, ok := d.GetOk("error_message"); ok {
		input.ErrorMessage = helpers.String(msg.(string))
	}

	_, resp, err := appIDClient.SetCloudDirectoryPasswordRegexWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error setting AppID Cloud Directory password regex: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDPasswordRegexRead(d, meta)
}

func resourceIBMAppIDPasswordRegexUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceIBMAppIDPasswordRegexCreate(d, meta)
}

func resourceIBMAppIDPasswordRegexDelete(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.SetCloudDirectoryPasswordRegexOptions{
		TenantID:           &tenantID,
		Base64EncodedRegex: helpers.String(""),
	}

	_, resp, err := appIDClient.SetCloudDirectoryPasswordRegexWithContext(context.TODO(), input)

	if err != nil {
		return fmt.Errorf("Error resetting AppID Cloud Directory password regex: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
