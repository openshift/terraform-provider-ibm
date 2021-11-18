package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDPasswordRegex() *schema.Resource {
	return &schema.Resource{
		Description: "The regular expression used by App ID for password strength validation",
		Read:        dataSourceIBMAppIDPasswordRegexRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"base64_encoded_regex": {
				Description: "The regex expression rule for acceptable password encoded in base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"error_message": {
				Description: "Custom error message",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"regex": {
				Description: "The escaped regex expression rule for acceptable password",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDPasswordRegexRead(d *schema.ResourceData, meta interface{}) error {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	pw, resp, err := appIDClient.GetCloudDirectoryPasswordRegexWithContext(context.TODO(), &appid.GetCloudDirectoryPasswordRegexOptions{
		TenantID: &tenantID,
	})

	if err != nil {
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

	d.SetId(tenantID)

	return nil
}
