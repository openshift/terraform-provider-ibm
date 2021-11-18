// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/atrackerv1"
)

func resourceIBMAtrackerRoute() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMAtrackerRouteCreate,
		Read:     resourceIBMAtrackerRouteRead,
		Update:   resourceIBMAtrackerRouteUpdate,
		Delete:   resourceIBMAtrackerRouteDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_atracker_route", "name"),
				Description:  "The name of the route. The name must be 1000 characters or less and cannot include any special characters other than `(space) - . _ :`.",
			},
			"receive_global_events": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates whether or not all global events should be forwarded to this region.",
			},
			"rules": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Routing rules that will be evaluated in their order of the array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_ids": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The target ID List. Only 1 target id is supported.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the route resource.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the route.",
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the route creation time.",
			},
			"updated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the route last updated time.",
			},
		},
	}
}

func resourceIBMAtrackerRouteValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 0)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 -._:]+$`,
			MinValueLength:             1,
			MaxValueLength:             1000,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_atracker_route", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMAtrackerRouteCreate(d *schema.ResourceData, meta interface{}) error {
	atrackerClient, err := meta.(ClientSession).AtrackerV1()
	if err != nil {
		return err
	}

	createRouteOptions := &atrackerv1.CreateRouteOptions{}

	createRouteOptions.SetName(d.Get("name").(string))
	createRouteOptions.SetReceiveGlobalEvents(d.Get("receive_global_events").(bool))
	var rules []atrackerv1.Rule
	for _, e := range d.Get("rules").([]interface{}) {
		value := e.(map[string]interface{})
		rulesItem := resourceIBMAtrackerRouteMapToRule(value)
		rules = append(rules, rulesItem)
	}
	createRouteOptions.SetRules(rules)

	route, response, err := atrackerClient.CreateRouteWithContext(context.TODO(), createRouteOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateRouteWithContext failed %s\n%s", err, response)
		return fmt.Errorf("CreateRouteWithContext failed %s\n%s", err, response)
	}

	d.SetId(*route.ID)

	return resourceIBMAtrackerRouteRead(d, meta)
}

func resourceIBMAtrackerRouteMapToRule(ruleMap map[string]interface{}) atrackerv1.Rule {
	rule := atrackerv1.Rule{}

	targetIds := []string{}
	for _, targetIdsItem := range ruleMap["target_ids"].([]interface{}) {
		targetIds = append(targetIds, targetIdsItem.(string))
	}
	rule.TargetIds = targetIds

	return rule
}

func resourceIBMAtrackerRouteRead(d *schema.ResourceData, meta interface{}) error {
	atrackerClient, err := meta.(ClientSession).AtrackerV1()
	if err != nil {
		return err
	}

	getRouteOptions := &atrackerv1.GetRouteOptions{}

	getRouteOptions.SetID(d.Id())

	route, response, err := atrackerClient.GetRouteWithContext(context.TODO(), getRouteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetRouteWithContext failed %s\n%s", err, response)
		return fmt.Errorf("GetRouteWithContext failed %s\n%s", err, response)
	}

	if err = d.Set("name", route.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("receive_global_events", route.ReceiveGlobalEvents); err != nil {
		return fmt.Errorf("Error setting receive_global_events: %s", err)
	}
	rules := []map[string]interface{}{}
	for _, rulesItem := range route.Rules {
		rulesItemMap := resourceIBMAtrackerRouteRuleToMap(rulesItem)
		rules = append(rules, rulesItemMap)
	}
	if err = d.Set("rules", rules); err != nil {
		return fmt.Errorf("Error setting rules: %s", err)
	}
	if err = d.Set("crn", route.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("version", intValue(route.Version)); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}
	if err = d.Set("created", dateTimeToString(route.Created)); err != nil {
		return fmt.Errorf("Error setting created: %s", err)
	}
	if err = d.Set("updated", dateTimeToString(route.Updated)); err != nil {
		return fmt.Errorf("Error setting updated: %s", err)
	}

	return nil
}

func resourceIBMAtrackerRouteRuleToMap(rule atrackerv1.Rule) map[string]interface{} {
	ruleMap := map[string]interface{}{}

	ruleMap["target_ids"] = rule.TargetIds

	return ruleMap
}

func resourceIBMAtrackerRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	atrackerClient, err := meta.(ClientSession).AtrackerV1()
	if err != nil {
		return err
	}

	replaceRouteOptions := &atrackerv1.ReplaceRouteOptions{}

	replaceRouteOptions.SetID(d.Id())
	replaceRouteOptions.SetName(d.Get("name").(string))
	replaceRouteOptions.SetReceiveGlobalEvents(d.Get("receive_global_events").(bool))
	var rules []atrackerv1.Rule
	for _, e := range d.Get("rules").([]interface{}) {
		value := e.(map[string]interface{})
		rulesItem := resourceIBMAtrackerRouteMapToRule(value)
		rules = append(rules, rulesItem)
	}
	replaceRouteOptions.SetRules(rules)

	_, response, err := atrackerClient.ReplaceRouteWithContext(context.TODO(), replaceRouteOptions)
	if err != nil {
		log.Printf("[DEBUG] ReplaceRouteWithContext failed %s\n%s", err, response)
		return fmt.Errorf("ReplaceRouteWithContext failed %s\n%s", err, response)
	}

	return resourceIBMAtrackerRouteRead(d, meta)
}

func resourceIBMAtrackerRouteDelete(d *schema.ResourceData, meta interface{}) error {
	atrackerClient, err := meta.(ClientSession).AtrackerV1()
	if err != nil {
		return err
	}

	deleteRouteOptions := &atrackerv1.DeleteRouteOptions{}

	deleteRouteOptions.SetID(d.Id())

	response, err := atrackerClient.DeleteRouteWithContext(context.TODO(), deleteRouteOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteRouteWithContext failed %s\n%s", err, response)
		return fmt.Errorf("DeleteRouteWithContext failed %s\n%s", err, response)
	}

	d.SetId("")

	return nil
}
