package ibm

import (
	"context"
	"fmt"

	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAppIDIDPSAMLMetadata() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve SAML metadata",
		Read:        dataSourceIBMAppIDIDPSAMLMetadataRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"metadata": {
				Type:        schema.TypeString,
				Description: "SAML Metadata",
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDIDPSAMLMetadataRead(d *schema.ResourceData, meta interface{}) error {
	appidClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return err
	}

	tenantID := d.Get("tenant_id").(string)

	metadata, resp, err := appidClient.GetSAMLMetadataWithContext(context.TODO(), &appid.GetSAMLMetadataOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return fmt.Errorf("Error loading AppID SAML metadata: %s\n%s", err, resp)
	}

	if err := d.Set("metadata", metadata); err != nil {
		return fmt.Errorf("Error setting AppID SAML metadata: %s", err)
	}

	d.SetId(tenantID)

	return nil
}
