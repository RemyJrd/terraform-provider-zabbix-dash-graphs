package zabbix

import (
	"fmt"
	"log"
	"strings"

	zabbixapi "github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// WidgetType defines the widget types for dashboards
var WidgetType = map[string]int{
	"graph":          0,
	"clock":          1,
	"sysmap":         2,
	"plain_text":     3,
	"url":            10,
	"trigger_info":   4,
	"trigger_over":   5, 
	"problems":       9,
	"problems_by_sv": 8,
}

// WidgetTypeStringMap is a mapping of numeric widget types to strings
var WidgetTypeStringMap = map[int]string{
	0:  "graph",
	1:  "clock",
	2:  "sysmap",
	3:  "plain_text",
	10: "url",
	4:  "trigger_info",
	5:  "trigger_over",
	9:  "problems",
	8:  "problems_by_sv",
}

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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"owner_userid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"page": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"display_period": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
						"widget": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{
												"graph", "clock", "sysmap", "plain_text", 
												"url", "trigger_info", "trigger_over", 
												"problems", "problems_by_sv",
											}, 
											false,
										),
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"x": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"y": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"width": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"height": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"graph_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ConflictsWith: []string{
											"page.0.widget.0.url", 
											"page.0.widget.0.text",
										},
									},
									"url": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ConflictsWith: []string{
											"page.0.widget.0.graph_id", 
											"page.0.widget.0.text",
										},
									},
									"text": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										ConflictsWith: []string{
											"page.0.widget.0.graph_id", 
											"page.0.widget.0.url",
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
	api := meta.(*zabbixapi.API)

	dashboard := buildDashboardObject(d, api)

	dashboards := Dashboards{*dashboard}
	err := api.DashboardsCreate(dashboards)
	if err != nil {
		return err
	}

	d.SetId(dashboards[0].DashboardID)

	log.Printf("[DEBUG] Created dashboard with ID %s", d.Id())

	return resourceZabbixDashboardRead(d, meta)
}

func resourceZabbixDashboardRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbixapi.API)

	log.Printf("[DEBUG] Will read dashboard with ID %s", d.Id())

	dashboards, err := api.DashboardsGet(zabbixapi.Params{
		"dashboardids": d.Id(),
		"selectPages": "extend",
		"selectUsers": "extend",
		"selectUserGroups": "extend",
	})

	if err != nil {
		return err
	}

	if len(dashboards) != 1 {
		return fmt.Errorf("Expected one dashboard with ID %s and got %d dashboards", d.Id(), len(dashboards))
	}

	dashboard := dashboards[0]

	log.Printf("[DEBUG] Dashboard name is %s", dashboard.Name)

	d.Set("name", dashboard.Name)
	d.Set("owner_userid", dashboard.UserID)

	// Handle dashboard pages
	pages := make([]map[string]interface{}, len(dashboard.Pages))
	for i, page := range dashboard.Pages {
		pageMap := map[string]interface{}{
			"name":           page.Name,
			"display_period": page.DisplayPeriod,
		}

		// Get widgets for this page
		widgets := page.Widgets
		widgetsList := make([]map[string]interface{}, len(widgets))
		for j, widget := range widgets {
			widgetMap := map[string]interface{}{
				"type":   WidgetTypeStringMap[widget.Type],
				"name":   widget.Name,
				"x":      widget.X,
				"y":      widget.Y,
				"width":  widget.Width,
				"height": widget.Height,
			}

			// Handle different widget types and their fields
			switch widget.Type {
			case WidgetType["graph"]:
				if resourceID, ok := widget.Fields.GetResourceID(); ok {
					widgetMap["graph_id"] = resourceID
				}
			case WidgetType["plain_text"]:
				if text, ok := widget.Fields.GetText(); ok {
					widgetMap["text"] = text
				}
			case WidgetType["url"]:
				if url, ok := widget.Fields.GetURL(); ok {
					widgetMap["url"] = url
				}
			}

			widgetsList[j] = widgetMap
		}

		pageMap["widget"] = widgetsList
		pages[i] = pageMap
	}

	d.Set("page", pages)

	return nil
}

func resourceZabbixDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbixapi.API)

	dashboard := buildDashboardObject(d, api)
	dashboard.DashboardID = d.Id()

	dashboards := Dashboards{*dashboard}
	err := api.DashboardsUpdate(dashboards)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updated dashboard with ID %s", d.Id())

	return resourceZabbixDashboardRead(d, meta)
}

func resourceZabbixDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbixapi.API)

	dashboardIDs := []string{d.Id()}
	err := api.DashboardsDelete(dashboardIDs)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleted dashboard with ID %s", d.Id())

	d.SetId("")
	return nil
}

func resourceZabbixDashboardExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbixapi.API)

	dashboards, err := api.DashboardsGet(zabbixapi.Params{
		"dashboardids": d.Id(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") {
			log.Printf("[DEBUG] Dashboard with ID %s doesn't exist", d.Id())
			return false, nil
		}
		return false, err
	}
	return len(dashboards) > 0, nil
}

func buildDashboardObject(d *schema.ResourceData, api *zabbixapi.API) *Dashboard {
	dashboard := Dashboard{
		Name: d.Get("name").(string),
	}

	// Set owner user ID if provided
	if userID, ok := d.GetOk("owner_userid"); ok {
		dashboard.UserID = userID.(string)
	} else {
		// Default to the current API user
		users, err := api.UsersGet(zabbixapi.Params{
			"output": []string{"userid"},
			"filter": map[string]interface{}{
				"alias": api.User(),
			},
		})
		if err == nil && len(users) > 0 {
			dashboard.UserID = users[0].UserID
		}
	}

	// Build pages
	pages := d.Get("page").([]interface{})
	dashboardPages := make(DashboardPages, len(pages))

	for i, p := range pages {
		page := p.(map[string]interface{})
		
		dashboardPage := DashboardPage{
			Name:          page["name"].(string),
			DisplayPeriod: page["display_period"].(int),
		}

		// Process widgets
		if widgets, ok := page["widget"].([]interface{}); ok && len(widgets) > 0 {
			dashboardWidgets := make(DashboardWidgets, len(widgets))
			
			for j, w := range widgets {
				widget := w.(map[string]interface{})
				
				dashboardWidget := DashboardWidget{
					Type:   WidgetType[widget["type"].(string)],
					Name:   widget["name"].(string),
					X:      widget["x"].(int),
					Y:      widget["y"].(int),
					Width:  widget["width"].(int),
					Height: widget["height"].(int),
					Fields: DashboardWidgetFields{},
				}

				// Set fields based on widget type
				switch dashboardWidget.Type {
				case WidgetType["graph"]:
					if graphID, ok := widget["graph_id"].(string); ok && graphID != "" {
						dashboardWidget.Fields = append(dashboardWidget.Fields, 
							DashboardWidgetField{
								Type:  0, // resourceid
								Name:  "graphid",
								Value: graphID,
							},
						)
					}
				case WidgetType["plain_text"]:
					if text, ok := widget["text"].(string); ok && text != "" {
						dashboardWidget.Fields = append(dashboardWidget.Fields, 
							DashboardWidgetField{
								Type:  1, // string
								Name:  "text",
								Value: text,
							},
						)
					}
				case WidgetType["url"]:
					if url, ok := widget["url"].(string); ok && url != "" {
						dashboardWidget.Fields = append(dashboardWidget.Fields, 
							DashboardWidgetField{
								Type:  1, // string
								Name:  "url",
								Value: url,
							},
						)
					}
				}

				dashboardWidgets[j] = dashboardWidget
			}
			
			dashboardPage.Widgets = dashboardWidgets
		}
		
		dashboardPages[i] = dashboardPage
	}

	dashboard.Pages = dashboardPages

	return &dashboard
}