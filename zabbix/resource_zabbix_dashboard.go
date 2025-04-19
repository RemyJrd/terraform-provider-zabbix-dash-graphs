package zabbix

import (
	"fmt"
	"log"
	"strings"

	zabbixapi "github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Note: The Dashboard, DashboardPage, DashboardWidget structs are now defined in types.go

// ResourceZabbixDashboard creates the Zabbix dashboard resource
func resourceZabbixDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixDashboardCreate,
		Read:   resourceZabbixDashboardRead,
		Update: resourceZabbixDashboardUpdate,
		Delete: resourceZabbixDashboardDelete,
		Exists: resourceZabbixDashboardExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"userid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "User ID of the dashboard owner. Defaults to the current user if not set.",
			},
			"pages": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"display_period": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
						"sort_order": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"widgets": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Widget type (e.g., system.clock, item.graph, problems)",
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"x": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"y": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"width": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"height": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"view_mode": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"fields": {
										Type: schema.TypeList,
										Required: true,
										Description: "Widget-specific configuration fields",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type: schema.TypeInt,
													Required: true,
													Description: "Field type identifier (integer)",
												},
												"name": {
													Type: schema.TypeString,
													Required: true,
													Description: "Field name (e.g., 'resourceid', 'graphid')",
												},
												"value": {
													Type: schema.TypeString,
													Required: true,
													Description: "Field value (always stored as string in TF state)",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceZabbixDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*ZabbixGraphAPI)
	dashboard := buildDashboardObject(d)

	dashboardsToCreate := Dashboards{*dashboard}
	err := api.DashboardsCreate(dashboardsToCreate)
	if err != nil {
		return fmt.Errorf("Error creating Zabbix dashboard %s: %v", dashboard.Name, err)
	}

	dashboard.DashboardID = dashboardsToCreate[0].DashboardID
	d.SetId(dashboard.DashboardID)

	log.Printf("[DEBUG] Created dashboard with ID %s", d.Id())
	return resourceZabbixDashboardRead(d, meta)
}

func resourceZabbixDashboardRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*ZabbixGraphAPI)
	dashboardID := d.Id()
	log.Printf("[DEBUG] Reading dashboard with ID %s", dashboardID)

	dashboards, err := api.DashboardsGet(zabbixapi.Params{
		"dashboardids": dashboardID,
		"selectPages":  "extend",
		"selectWidgets": "extend",
		"selectFields":  "extend",
	})
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") || strings.Contains(err.Error(), "No dashboard found") {
			log.Printf("[WARN] Zabbix Dashboard (%s) not found, removing from state", dashboardID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Zabbix dashboard %s: %v", dashboardID, err)
	}

	if len(dashboards) != 1 {
		log.Printf("[WARN] Expected one dashboard %s but got %d. Removing from state", dashboardID, len(dashboards))
		d.SetId("")
		return nil
	}
	dashboard := dashboards[0]

	d.Set("name", dashboard.Name)
	d.Set("userid", dashboard.UserID)

	pages := make([]interface{}, len(dashboard.Pages))
	for i, page := range dashboard.Pages {
		widgets := make([]interface{}, len(page.Widgets))
		for j, widget := range page.Widgets {
			fields := make([]interface{}, len(widget.Fields))
			for k, field := range widget.Fields {
				fields[k] = map[string]interface{}{
					"type":  field.Type,
					"name":  field.Name,
					"value": fmt.Sprintf("%v", field.Value),
				}
			}
			widgets[j] = map[string]interface{}{
				"type":      widget.Type,
				"name":      widget.Name,
				"x":         widget.X,
				"y":         widget.Y,
				"width":     widget.Width,
				"height":    widget.Height,
				"view_mode": widget.ViewMode,
				"fields":    fields,
			}
		}
		pages[i] = map[string]interface{}{
			"name":           page.Name,
			"display_period": page.DisplayPeriod,
			"sort_order":     page.SortOrder,
			"widgets":        widgets,
		}
	}
	if err := d.Set("pages", pages); err != nil {
		return fmt.Errorf("error setting pages for dashboard %s: %v", dashboardID, err)
	}

	return nil
}

func resourceZabbixDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*ZabbixGraphAPI)
	dashboardID := d.Id()

	dashboard := buildDashboardObject(d)
	dashboard.DashboardID = dashboardID

	dashboardsToUpdate := Dashboards{*dashboard}
	err := api.DashboardsUpdate(dashboardsToUpdate)
	if err != nil {
		return fmt.Errorf("Error updating Zabbix dashboard %s: %v", dashboardID, err)
	}
	log.Printf("[DEBUG] Updated dashboard with ID %s", dashboardID)
	return resourceZabbixDashboardRead(d, meta)
}

func resourceZabbixDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*ZabbixGraphAPI)
	dashboardID := d.Id()

	err := api.DashboardsDelete([]string{dashboardID})
	if err != nil {
		if strings.Contains(err.Error(), "No dashboard found") || strings.Contains(err.Error(), "does not exist") {
			log.Printf("[WARN] Zabbix Dashboard (%s) already deleted, removing from state", dashboardID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting Zabbix dashboard %s: %v", dashboardID, err)
	}
	log.Printf("[DEBUG] Deleted dashboard with ID %s", dashboardID)
	d.SetId("")
	return nil
}

func resourceZabbixDashboardExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*ZabbixGraphAPI)
	dashboardID := d.Id()

	_, err := api.DashboardsGet(zabbixapi.Params{
		"dashboardids": dashboardID,
		"output":       "dashboardid",
	})
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") || strings.Contains(err.Error(), "No dashboard found") {
			log.Printf("[DEBUG] Dashboard %s doesn't exist", dashboardID)
			return false, nil
		}
		return false, fmt.Errorf("Error checking existence of dashboard %s: %v", dashboardID, err)
	}
	return true, nil
}

func buildDashboardObject(d *schema.ResourceData) *Dashboard {
	dashboard := &Dashboard{
		Name:   d.Get("name").(string),
		UserID: d.Get("userid").(string),
	}

	if v, ok := d.GetOk("pages"); ok {
		pagesRaw := v.([]interface{})
		dashboard.Pages = make([]DashboardPage, len(pagesRaw))
		for i, pageRaw := range pagesRaw {
			pageMap := pageRaw.(map[string]interface{})
			widgetsRaw := pageMap["widgets"].([]interface{})
			widgets := make([]DashboardWidget, len(widgetsRaw))
			for j, widgetRaw := range widgetsRaw {
				widgetMap := widgetRaw.(map[string]interface{})
				fieldsRaw := widgetMap["fields"].([]interface{})
				fields := make([]DashboardWidgetField, len(fieldsRaw))
				for k, fieldRaw := range fieldsRaw {
					fieldMap := fieldRaw.(map[string]interface{})
					fields[k] = DashboardWidgetField{
						Type:  fieldMap["type"].(int),
						Name:  fieldMap["name"].(string),
						Value: fieldMap["value"].(string),
					}
				}
				widgets[j] = DashboardWidget{
					Type:     widgetMap["type"].(string),
					Name:     widgetMap["name"].(string),
					X:        widgetMap["x"].(int),
					Y:        widgetMap["y"].(int),
					Width:    widgetMap["width"].(int),
					Height:   widgetMap["height"].(int),
					ViewMode: widgetMap["view_mode"].(int),
					Fields:   fields,
				}
			}
			dashboard.Pages[i] = DashboardPage{
				Name:          pageMap["name"].(string),
				DisplayPeriod: pageMap["display_period"].(int),
				SortOrder:     pageMap["sort_order"].(int),
				Widgets:       widgets,
			}
		}
	}

	return dashboard
} 