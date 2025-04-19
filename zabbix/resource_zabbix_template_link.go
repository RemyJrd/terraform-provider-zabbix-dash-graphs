package zabbix

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceZabbixTemplateLink returns a *schema.Resource that corresponds to a Zabbix template link
func resourceZabbixTemplateLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixTemplateLinkCreate,
		Read:   resourceZabbixTemplateLinkRead,
		Update: resourceZabbixTemplateLinkUpdate,
		Delete: resourceZabbixTemplateLinkDelete,
		Exists: resourceZabbixTemplateLinkExists,
		Schema: map[string]*schema.Schema{
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"clear_templates": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"item": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     schemaTemplateItem(),
			},
			"trigger": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     schemaTemplateTrigger(),
			},
			"lld_rule": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     schemaTemplateLLDRule(),
			},
		},
	}
}

func schemaTemplateItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"trigger_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func schemaTemplateLLDRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lld_rule_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// Simple placeholder implementations to satisfy the interface

func resourceZabbixTemplateLinkCreate(d *schema.ResourceData, meta interface{}) error {
	_ = meta.(*ZabbixGraphAPI) // Placeholder for future implementation
	template_id := d.Get("template_id").(string)
	d.SetId(template_id)
	return resourceZabbixTemplateLinkRead(d, meta)
}

func resourceZabbixTemplateLinkRead(d *schema.ResourceData, meta interface{}) error {
	_ = meta.(*ZabbixGraphAPI) // Placeholder for future implementation
	return nil
}

func resourceZabbixTemplateLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	_ = meta.(*ZabbixGraphAPI) // Placeholder for future implementation
	return resourceZabbixTemplateLinkRead(d, meta)
}

func resourceZabbixTemplateLinkDelete(d *schema.ResourceData, meta interface{}) error {
	_ = meta.(*ZabbixGraphAPI) // Placeholder for future implementation
	return nil
}

func resourceZabbixTemplateLinkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	_ = meta.(*ZabbixGraphAPI) // Placeholder for future implementation
	return true, nil
}
