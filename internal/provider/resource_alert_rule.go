package provider

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	mykibana "github.com/qonto/terraform-provider-kibana/pkg/kibana_api"
)

func resourceAlertRule() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider AlertRule.",

		CreateContext: resourceAlertRuleCreate,
		ReadContext:   resourceAlertRuleRead,
		UpdateContext: resourceAlertRuleUpdate,
		DeleteContext: resourceAlertRuleDelete,

		Schema: map[string]*schema.Schema{
			"space_id": {
				Description: "An identifier for the space. If space_id is not provided in the URL, the default space is used.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "A name to reference and search.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tags": {
				Description: "A list of keywords to reference and search.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rule_type_id": {
				Description: "The ID of the rule type that you want to call when the rule is scheduled to run.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"schedule": {
				Description: "The schedule specifying when this rule should be run, using one of the available schedule formats.",
				Type:        schema.TypeMap,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"throttle": {
				Description: "How often this rule should fire the same actions. This will prevent the rule from sending out the same notification over and over. For example, if a rule with a schedule of 1 minute stays in a triggered state for 90 minutes, setting a throttle of 10m or 1h will prevent it from sending 90 notifications during this period.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"notify_when": {
				Description: "The condition for throttling the notification: onActionGroupChange, onActiveAlert, or onThrottleInterval.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Indicates if you want to run the rule on an interval basis after it is created.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"consumer": {
				Description: "The name of the application that owns the rule. This name has to match the Kibana Feature name, as that dictates the required RBAC privileges.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"params": {
				Description:      "The parameters to pass to the rule type executor params value. This will also validate against the rule type params validator, if defined.",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: rawJsonEqual,
			},
			"actions": {
				Description: "An array of the following action objects.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the connector saved object to execute.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"group": {
							Description: "Grouping actions is recommended for escalations for different types of alerts. If you don’t need this, set this value to default.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"params": {
							Description:      "The map to the params that the connector type will receive. ` params` are handled as Mustache templates and passed a default set of context.",
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: rawJsonEqual,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAlertRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	var alertId string
	var diags diag.Diagnostics
	client := meta.(mykibana.KibanaAPI)
	alert := mykibana.Alert{}
	alert.Name = d.Get("name").(string)
	tags := d.Get("tags").([]interface{})
	for _, tag := range tags {
		alert.Tags = append(alert.Tags, tag.(string))
	}
	alert.RuleTypeId = d.Get("rule_type_id").(string)
	schedule := d.Get("schedule").(map[string]interface{})
	alert.Schedule = make(map[string]string)
	for key, val := range schedule {
		alert.Schedule[key] = val.(string)
	}
	alert.Throttle = d.Get("throttle").(string)
	alert.NotifyWhen = d.Get("notify_when").(string)
	// alert.Enabled = d.Get("enabled").(bool)
	alert.Consumer = d.Get("consumer").(string)
	params := d.Get("params").(string)
	alert.Params = json.RawMessage([]byte(params))
	actionsInterface := d.Get("actions").([]interface{})
	actionsList := make([]map[string]interface{}, 0, len(actionsInterface))
	for _, action := range actionsInterface {
		actionsList = append(actionsList, action.(map[string]interface{}))
	}
	alert.Actions, err = deflateActions(actionsList)
	if err != nil {
		return diag.FromErr(err)
	}
	alertId, err = client.CreateAlertRule(alert)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(alertId)
	// Hacky. for some reason setting enabled to false at rule creation doesn't work, need fix.
	if d.Get("enabled") != nil && d.Get("enabled") == false {
		client.DisableRule(alertId)
	}
	resourceAlertRuleRead(ctx, d, meta)
	return diags
}

func resourceAlertRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(mykibana.KibanaAPI)
	alertId := d.Id()
	alert, err := client.ReadAlertRule(alertId)
	if err != nil {
		return diag.FromErr(err)
	}
	flattenedActions, err := flattenActions(alert.Actions)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("actions", flattenedActions)
	paramsBytes, err := json.Marshal(alert.Params)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("params", string(paramsBytes))
	d.Set("consumer", alert.Consumer)
	d.Set("enabled", alert.Enabled)
	d.Set("name", alert.Name)
	d.Set("notify_when", alert.NotifyWhen)
	d.Set("rule_type_id", alert.RuleTypeId)
	d.Set("schedule", alert.Schedule)
	d.Set("tags", alert.Tags)
	d.Set("throttle", alert.Throttle)

	return diags
}

func resourceAlertRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	var diags diag.Diagnostics
	client := meta.(mykibana.KibanaAPI)
	alertId := d.Id()
	alert := mykibana.Alert{}
	alert.Name = d.Get("name").(string)
	alert.Tags = []string{}
	tags := d.Get("tags").([]interface{})
	for _, tag := range tags {
		alert.Tags = append(alert.Tags, tag.(string))
	}
	schedule := d.Get("schedule").(map[string]interface{})
	alert.Schedule = make(map[string]string)
	for key, val := range schedule {
		alert.Schedule[key] = val.(string)
	}
	alert.Throttle = d.Get("throttle").(string)
	alert.NotifyWhen = d.Get("notify_when").(string)
	params := d.Get("params").(string)
	alert.Params = json.RawMessage([]byte(params))
	actionsInterface := d.Get("actions").([]interface{})
	actionsList := make([]map[string]interface{}, 0, len(actionsInterface))
	for _, action := range actionsInterface {
		actionsList = append(actionsList, action.(map[string]interface{}))
	}
	alert.Actions, err = deflateActions(actionsList)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("enabled") {
		if enabled := d.Get("enabled").(bool); enabled {
			err = client.EnableRule(alertId)
		} else {
			err = client.DisableRule(alertId)
		}
		if err != nil {
			return diag.FromErr(err)
		}
		// resourceAlertRuleRead(ctx, d, meta)
	}
	err = client.UpdateAlertRule(alertId, alert)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceAlertRuleRead(ctx, d, meta)
	return diags
}

func resourceAlertRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(mykibana.KibanaAPI)
	alertId := d.Id()
	err := client.DeleteAlertRule(alertId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func deflateActions(actionArray []map[string]interface{}) ([]mykibana.Action, error) {
	actions := []mykibana.Action{}
	for _, flatAction := range actionArray {
		var action mykibana.Action
		id := flatAction["id"].(string)
		action.Id = id
		group := flatAction["group"].(string)
		action.Group = group
		params := flatAction["params"].(string)
		action.Params = json.RawMessage([]byte(params))
		actions = append(actions, action)
	}
	return actions, nil
}

func flattenActions(actions []mykibana.Action) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0, len(actions))
	for _, a := range actions {
		action := make(map[string]interface{})
		action["id"] = a.Id
		action["group"] = a.Group
		paramsBytes, err := json.Marshal(a.Params)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to marshal Action")
		}
		action["params"] = string(paramsBytes)
		res = append(res, action)
	}
	return res, nil
}

func rawJsonEqual(k, oldValue, newValue string, d *schema.ResourceData) bool {
	var oldInterface, newInterface interface{}
	if err := json.Unmarshal([]byte(oldValue), &oldInterface); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newValue), &newInterface); err != nil {
		return false
	}
	return reflect.DeepEqual(oldInterface, newInterface)
}
