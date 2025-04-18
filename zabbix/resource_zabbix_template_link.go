package zabbix

import (
	"log"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceZabbixTemplateLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTemplateLinkCreate,
		Read:   resourceZabbixTemplateLinkRead,
		Exists: resourceZabbixTemplateLinkExists,
		Update: resourceZabbixTemplateLinkUpdate,
		Delete: resourceZabbixTemplateLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"item": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplateItem(),
				Optional: true,
			},
			"trigger": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplateTrigger(),
				Optional: true,
			},
			"lld_rule": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schemaTemplatelldRule(),
				Optional: true,
			},
			"graph": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"graph_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func schemaTemplateItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"local": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"item_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func schemaTemplateTrigger() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"local": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"trigger_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func schemaTemplatelldRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"local": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lld_rule_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceZabbixTemplateLinkCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	template_id := d.Get("template_id").(string)

	items, item_set := d.GetOk("item")
	triggers, trigger_set := d.GetOk("trigger")
	graphs, graph_set := d.GetOk("graph")

	// First, we get the state of the template (items, triggers and graphs)
	params := map[string]interface{}{
		"hostids":          template_id,
		"preservekeys":     true,
		"selectItems":      []string{"itemid"},
		"selectTriggers":   []string{"triggerid"},
		"selectGraphs":     []string{"graphid"},
		"output":           []string{"name"},
		"filter":           map[string]interface{}{},
		"templated_hosts":  true,
		"with_items":       true,
		"with_triggers":    true,
		"with_graphs":      true,
		"monitored_hosts":  true,
		"with_monitored_items": true,
		"with_monitored_triggers": true,
		"with_monitored_graphs": true,
		"with_simple_graph_items": true,
	}

	templates, err := api.TemplatesGet(params)
	if err != nil {
		return err
	}

	if len(templates) != 1 {
		return fmt.Errorf("Expected exactly one template with id %s", template_id)
	}

	template := templates[0]
	template_items := template.Items
	template_triggers := template.Triggers
	template_graphs := template.Graphs

	log.Printf("[DEBUG] Items for template %v : %v", template.Name, template_items)

	// Next, we get the list of items, triggers and graphs to keep:
	var items_to_keep []zabbix.Item
	var items_to_keep_id []string
	var triggers_to_keep []zabbix.Trigger
	var triggers_to_keep_id []string
	var graphs_to_keep []zabbix.Graph
	var graphs_to_keep_id []string

	if item_set {
		for _, i := range items.(*schema.Set).List() {
			item_id := i.(map[string]interface{})["item_id"].(string)
			log.Printf("[DEBUG] item_id: %s", item_id)
			items_to_keep_id = append(items_to_keep_id, item_id)
		}
	}

	if trigger_set {
		for _, i := range triggers.(*schema.Set).List() {
			trigger_id := i.(map[string]interface{})["trigger_id"].(string)
			log.Printf("[DEBUG] trigger_id: %s", trigger_id)
			triggers_to_keep_id = append(triggers_to_keep_id, trigger_id)
		}
	}

	if graph_set {
		for _, i := range graphs.(*schema.Set).List() {
			graph_id := i.(map[string]interface{})["graph_id"].(string)
			log.Printf("[DEBUG] graph_id: %s", graph_id)
			graphs_to_keep_id = append(graphs_to_keep_id, graph_id)
		}
	}

	// Now we identify what to delete
	var items_to_delete []zabbix.Item
	var triggers_to_delete []zabbix.Trigger
	var graphs_to_delete []zabbix.Graph

	for _, template_item := range template_items {
		found := false
		for _, item_id := range items_to_keep_id {
			if item_id == template_item.ItemID {
				found = true
				items_to_keep = append(items_to_keep, template_item)
				break
			}
		}
		if !found {
			log.Printf("[DEBUG] Item to delete %s", template_item.ItemID)
			items_to_delete = append(items_to_delete, template_item)
		} else {
			log.Printf("[DEBUG] Item to keep %s", template_item.ItemID)
		}
	}

	for _, template_trigger := range template_triggers {
		found := false
		for _, trigger_id := range triggers_to_keep_id {
			if trigger_id == template_trigger.TriggerID {
				found = true
				triggers_to_keep = append(triggers_to_keep, template_trigger)
				break
			}
		}
		if !found {
			log.Printf("[DEBUG] Trigger to delete %s", template_trigger.TriggerID)
			triggers_to_delete = append(triggers_to_delete, template_trigger)
		} else {
			log.Printf("[DEBUG] Trigger to keep %s", template_trigger.TriggerID)
		}
	}

	for _, template_graph := range template_graphs {
		found := false
		for _, graph_id := range graphs_to_keep_id {
			if graph_id == template_graph.GraphID {
				found = true
				graphs_to_keep = append(graphs_to_keep, template_graph)
				break
			}
		}
		if !found {
			log.Printf("[DEBUG] Graph to delete %s", template_graph.GraphID)
			graphs_to_delete = append(graphs_to_delete, template_graph)
		} else {
			log.Printf("[DEBUG] Graph to keep %s", template_graph.GraphID)
		}
	}

	// Delete what we have to delete
	var item_ids_to_delete []string
	for _, item := range items_to_delete {
		item_ids_to_delete = append(item_ids_to_delete, item.ItemID)
	}
	if len(item_ids_to_delete) > 0 {
		log.Printf("[DEBUG] Deleting %d items", len(item_ids_to_delete))
		err = api.ItemsDelete(item_ids_to_delete)
		if err != nil {
			return err
		}
	}

	var trigger_ids_to_delete []string
	for _, trigger := range triggers_to_delete {
		trigger_ids_to_delete = append(trigger_ids_to_delete, trigger.TriggerID)
	}
	if len(trigger_ids_to_delete) > 0 {
		log.Printf("[DEBUG] Expected to delete %d trigger", len(trigger_ids_to_delete))
		err = api.TriggersDelete(trigger_ids_to_delete)
		if err != nil {
			return err
		}
	}

	var graph_ids_to_delete []string
	for _, graph := range graphs_to_delete {
		graph_ids_to_delete = append(graph_ids_to_delete, graph.GraphID)
	}
	if len(graph_ids_to_delete) > 0 {
		log.Printf("[DEBUG] Expected to delete %d graphs", len(graph_ids_to_delete))
		err = api.GraphsDelete(graph_ids_to_delete)
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Template link created")

	d.SetId(template_id)

	return nil
}

func resourceZabbixTemplateLinkRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	log.Printf("[DEBUG] Reading template link")

	template_id := d.Id()

	params := map[string]interface{}{
		"output":         "extend",
		"templateids":    template_id,
		"selectItems":    []string{"itemid"},
		"selectTriggers": []string{"triggerid"},
		"selectGraphs":   []string{"graphid"},
	}

	template, err := api.TemplatesGet(params)
	if err != nil {
		return err
	}

	if len(template) != 1 {
		return fmt.Errorf("Expected 1 template and got %d templates", len(template))
	}

	t := template[0]

	// Check if items are set
	if i, ok := d.GetOk("item"); ok {
		var newItems []map[string]string
		for _, u := range i.(*schema.Set).List() {
			newItems = append(newItems, map[string]string{
				"item_id": u.(map[string]interface{})["item_id"].(string),
			})
		}
		d.Set("item", newItems)
	}

	// Check if triggers are set
	if i, ok := d.GetOk("trigger"); ok {
		var newTriggers []map[string]string
		for _, u := range i.(*schema.Set).List() {
			newTriggers = append(newTriggers, map[string]string{
				"trigger_id": u.(map[string]interface{})["trigger_id"].(string),
			})
		}
		d.Set("trigger", newTriggers)
	}

	// Check if graphs are set
	if i, ok := d.GetOk("graph"); ok {
		var newGraphs []map[string]string
		for _, u := range i.(*schema.Set).List() {
			newGraphs = append(newGraphs, map[string]string{
				"graph_id": u.(map[string]interface{})["graph_id"].(string),
			})
		}
		d.Set("graph", newGraphs)
	}

	d.Set("template_id", template_id)
	return nil
}

func resourceZabbixTemplateLinkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func resourceZabbixTemplateLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbix.API)

	err := updateZabbixTemplateItems(d, api)
	if err != nil {
		return err
	}
	err = updateZabbixTemplateTriggers(d, api)
	if err != nil {
		return err
	}
	err = updateZabbixTemplateDiscoveryRules(d, api)
	if err != nil {
		return err
	}
	return resourceZabbixTemplateLinkRead(d, meta)
}

func resourceZabbixTemplateLinkDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func getTerraformTemplateItems(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	items, err := api.ItemsGet(params)
	if err != nil {
		return nil, err
	}

	itemsTerraform := make([]interface{}, len(items))
	for i, item := range items {
		var itemTerraform = make(map[string]interface{})

		itemTerraform["local"] = true
		itemTerraform["item_id"] = item.ItemID
		itemsTerraform[i] = itemTerraform
	}
	return itemsTerraform, nil
}

func getTerraformTemplateTriggers(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	triggers, err := api.TriggersGet(params)
	if err != nil {
		return nil, err
	}

	triggersTerraform := make([]interface{}, len(triggers))
	for i, trigger := range triggers {
		var triggerTerraform = make(map[string]interface{})

		triggerTerraform["local"] = true
		triggerTerraform["trigger_id"] = trigger.TriggerID
		triggersTerraform[i] = triggerTerraform
	}
	return triggersTerraform, nil
}

func getTerraformTemplateLLDRules(d *schema.ResourceData, api *zabbix.API) ([]interface{}, error) {
	params := zabbix.Params{
		"output": "extend",
		"templateids": []string{
			d.Get("template_id").(string),
		},
		"inherited": false,
	}
	lldRules, err := api.DiscoveryRulesGet(params)
	if err != nil {
		return nil, err
	}

	lldRulesTerraform := make([]interface{}, len(lldRules))
	for i, lldRule := range lldRules {
		var lldRuleTerraform = make(map[string]interface{})

		lldRuleTerraform["local"] = true
		lldRuleTerraform["lld_rule_id"] = lldRule.ItemID
		lldRulesTerraform[i] = lldRuleTerraform
	}
	return lldRulesTerraform, nil
}

func updateZabbixTemplateItems(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("item") {
		oldV, newV := d.GetChange("item")
		oldItems := oldV.(*schema.Set).List()
		newItems := newV.(*schema.Set).List()
		var deletedItems []string
		templatedItems, err := api.ItemsGet(zabbix.Params{
			"templateids": []string{
				d.Get("template_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Found templated item %#v", templatedItems)
		for _, oldItem := range oldItems {
			oldItemValue := oldItem.(map[string]interface{})
			exist := false

			if oldItemValue["local"] == true {
				continue
			}

			for _, newItem := range newItems {
				newItemValue := newItem.(map[string]interface{})
				if newItemValue["item_id"].(string) == oldItemValue["item_id"].(string) {
					exist = true
				}
			}

			if !exist {
				templated := false

				for _, templatedItem := range templatedItems {
					if templatedItem.ItemID == oldItemValue["item_id"].(string) {
						templated = true
						break
					}
				}
				if !templated {
					deletedItems = append(deletedItems, oldItemValue["item_id"].(string))
				}
			}
		}
		if len(deletedItems) > 0 {
			log.Printf("[DEBUG] template link will delete item with ids : %#v", deletedItems)
			_, err := api.ItemsDeleteIDs(deletedItems)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateZabbixTemplateTriggers(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("trigger") {
		oldV, newV := d.GetChange("trigger")
		oldTriggers := oldV.(*schema.Set).List()
		newTriggers := newV.(*schema.Set).List()
		var deletedTriggers []string
		templatedTriggers, err := api.TriggersGet(zabbix.Params{
			"output": "extend",
			"templateids": []string{
				d.Get("template_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] found templated trigger %#v", templatedTriggers)
		for _, oldTrigger := range oldTriggers {
			oldTriggerValue := oldTrigger.(map[string]interface{})
			exist := false

			if oldTriggerValue["local"] == true {
				continue
			}

			for _, newTrigger := range newTriggers {
				newTriggerValue := newTrigger.(map[string]interface{})
				if oldTriggerValue["trigger_id"].(string) == newTriggerValue["trigger_id"].(string) {
					exist = true
				}
			}

			if !exist {
				templated := false

				for _, templatedTrigger := range templatedTriggers {
					if templatedTrigger.TriggerID == oldTriggerValue["trigger_id"].(string) {
						templated = true
						break
					}
				}
				if !templated {
					deletedTriggers = append(deletedTriggers, oldTriggerValue["trigger_id"].(string))
				}
			}
		}
		if len(deletedTriggers) > 0 {
			log.Printf("[DEBUG] template link will delete trigger with ids : %#v", deletedTriggers)
			_, err := api.TriggersDeleteIDs(deletedTriggers)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateZabbixTemplateDiscoveryRules(d *schema.ResourceData, api *zabbix.API) error {
	if d.HasChange("lld_rule") {
		oldV, newV := d.GetChange("lld_rule")
		oldlldRules := oldV.(*schema.Set).List()
		newlldRules := newV.(*schema.Set).List()
		var deletedlldRules []string
		templatedlldRules, err := api.DiscoveryRulesGet(zabbix.Params{
			"output": "extend",
			"templateids": []string{
				d.Get("template_id").(string),
			},
			"inherited": true,
		})

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] found templated lldRule %#v", templatedlldRules)
		for _, oldlldRule := range oldlldRules {
			oldlldRuleValue := oldlldRule.(map[string]interface{})
			exist := false

			if oldlldRuleValue["local"] == true {
				continue
			}

			for _, newlldRule := range newlldRules {
				newlldRuleValue := newlldRule.(map[string]interface{})
				if oldlldRuleValue["lld_rule_id"].(string) == newlldRuleValue["lld_rule_id"].(string) {
					exist = true
				}
			}

			if !exist {
				templated := false

				for _, templatedlldRule := range templatedlldRules {
					if templatedlldRule.ItemID == oldlldRuleValue["lld_rule_id"].(string) {
						templated = true
						break
					}
				}
				if !templated {
					deletedlldRules = append(deletedlldRules, oldlldRuleValue["lld_rule_id"].(string))
				}
			}
		}
		if len(deletedlldRules) > 0 {
			log.Printf("[DEBUG] template link will delete lldRule with ids : %#v", deletedlldRules)
			_, err := api.DiscoveryRulesDeletesIDs(deletedlldRules)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
