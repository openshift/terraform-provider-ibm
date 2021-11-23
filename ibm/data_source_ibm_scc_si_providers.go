// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/findingsv1"
)

func dataSourceIBMSccSiProviders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSccSiProvidersRead,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"providers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The providers requested.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the provider in the form '{account_id}/providers/{provider_id}'.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the provider.",
						},
					},
				},
			},
			"limit": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: ValidatePageSize,
				Description:  "The number of elements returned in the current instance. The default is 200.",
			},
			"skip": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The offset is the index of the item from which you want to start returning data from. The default is 0.",
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of providers available.",
			},
		},
	}
}

func dataSourceIBMSccSiProvidersRead(d *schema.ResourceData, meta interface{}) error {
	findingsClient, err := meta.(ClientSession).FindingsV1()
	if err != nil {
		return err
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := d.Get("account_id").(string)
	log.Println(fmt.Sprintf("[DEBUG] using specified AccountID %s", accountID))
	if accountID == "" {
		accountID = userDetails.userAccount
		log.Println(fmt.Sprintf("[DEBUG] AccountID not spedified, using %s", accountID))
	}
	findingsClient.AccountID = &accountID

	listProvidersOptions := &findingsv1.ListProvidersOptions{}

	if skip, ok := d.GetOk("skip"); ok {
		listProvidersOptions.SetSkip(int64(skip.(int)))
	}
	if limit, ok := d.GetOk("limit"); ok {
		listProvidersOptions.SetLimit(int64(limit.(int)))
	}

	apiProviders, totalCount, err := collectAllProviders(findingsClient, context.TODO(), listProvidersOptions)
	if err != nil {
		log.Printf("[DEBUG] ListProvidersWithContext failed %s", err)
		return fmt.Errorf("ListProvidersWithContext failed %s", err)
	}

	d.SetId(dataSourceIBMSccSiProvidersID(d))

	if apiProviders != nil {
		err = d.Set("providers", dataSourceAPIListProvidersResponseFlattenProviders(apiProviders))
		if err != nil {
			return fmt.Errorf("Error setting providers %s", err)
		}
	}
	if err = d.Set("total_count", intValue(totalCount)); err != nil {
		return fmt.Errorf("Error setting total_count: %s", err)
	}

	return nil
}

// dataSourceIBMSccSiProviderID returns a reasonable ID for the list.
func dataSourceIBMSccSiProvidersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func collectAllProviders(findingsClient *findingsv1.FindingsV1, ctx context.Context, options *findingsv1.ListProvidersOptions) ([]findingsv1.APIProvider, *int64, error) {
	finalList := []findingsv1.APIProvider{}
	totalCount := core.Int64Ptr(0)

	for {
		apiListProvidersResponse, response, err := findingsClient.ListProvidersWithContext(context.TODO(), options)
		if err != nil {
			return nil, core.Int64Ptr(0), fmt.Errorf("%s\n%s", err, response)
		}

		totalCount = apiListProvidersResponse.TotalCount

		finalList = append(finalList, apiListProvidersResponse.Providers...)

		// if user has specified some limit, then stop once finalList has length equal to the limit specified
		if options.Limit != nil && int64(len(finalList)) == *options.Limit {
			break
		}

		// if skip is specified, then stop once the finalList has length equal to difference of totalCount and skip
		if options.Skip != nil && int64(len(finalList)) == *apiListProvidersResponse.TotalCount-*options.Skip {
			break
		}

		// if user has not specified the limit, then stop once the finalList has length equal to totalCount
		if options.Limit == nil && int64(len(finalList)) == *apiListProvidersResponse.TotalCount {
			break
		}
	}

	return finalList, totalCount, nil
}

func dataSourceAPIListProvidersResponseFlattenProviders(result []findingsv1.APIProvider) (providers []map[string]interface{}) {
	for _, providersItem := range result {
		providers = append(providers, dataSourceAPIListProvidersResponseProvidersToMap(providersItem))
	}

	return providers
}

func dataSourceAPIListProvidersResponseProvidersToMap(providersItem findingsv1.APIProvider) (providersMap map[string]interface{}) {
	providersMap = map[string]interface{}{}

	if providersItem.Name != nil {
		providersMap["name"] = providersItem.Name
	}
	if providersItem.ID != nil {
		providersMap["id"] = providersItem.ID
	}

	return providersMap
}
