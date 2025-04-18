package zabbix

import (
	"fmt"
	"log"
	"strings"

	zabbixapi "github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// GraphType defines the graph types
var GraphType = map[string]zabbixapi.GraphType{
	"normal":    zabbixapi.Normal,
	"stacked":   zabbixapi.Stacked,
	"pie":       zabbixapi.Pie,
	"exploded":  zabbixapi.Exploded,
}

// GraphTypeStringMap is a mapping of numeric graph types to strings
var GraphTypeStringMap = map[zabbixapi.GraphType]string{
	zabbixapi.Normal:    "normal",
	zabbixapi.Stacked:   "stacked",
	zabbixapi.Pie:       "pie",
	zabbixapi.Exploded:  "exploded",
}

// YAxisSide defines the possible positions of the Y axis
var YAxisSide = map[string]int{
	"left":  0,
	"right": 1,
}

// YAxisSideStringMap is a mapping of numeric Y axis positions to strings
var YAxisSideStringMap = map[int]string{
	0: "left",
	1: "right",
}

// DrawType defines the graph item draw types
var DrawType = map[string]int{
	"line":          0,
	"filled_region": 1,
	"bold_line":     2,
	"dot":           3,
	"dashed_line":   4,
	"gradient_line": 5,
}

// DrawTypeStringMap is a mapping of numeric draw types to strings
var DrawTypeStringMap = map[int]string{
	0: "line",
	1: "filled_region",
	2: "bold_line",
	3: "dot",
	4: "dashed_line",
	5: "gradient_line",
}

// ResourceZabbixGraph creates the Zabbix graph resource
func resourceZabbixGraph() *schema.Resource {
	return &schema.Resource{
		Create: resourceZabbixGraphCreate,
		Read:   resourceZabbixGraphRead,
		Update: resourceZabbixGraphUpdate,
		Delete: resourceZabbixGraphDelete,
		Exists: resourceZabbixGraphExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"host_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"width": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  900,
			},
			"height": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  200,
			},
			"yaxismin": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"yaxismax": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "100",
			},
			"show_work_period": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_triggers": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_legend": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"show_3d": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"percent_left": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"percent_right": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "normal",
				ValidateFunc: validation.StringInSlice([]string{"normal", "stacked", "pie", "exploded"}, false),
			},
			"graph_items": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"item_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"color": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"calc_fnc": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2, // avg
						},
						"type": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0, // simple
						},
						"draw_type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "line",
							ValidateFunc: validation.StringInSlice([]string{"line", "filled_region", "bold_line", "dot", "dashed_line", "gradient_line"}, false),
						},
						"yaxisside": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "left",
							ValidateFunc: validation.StringInSlice([]string{"left", "right"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceZabbixGraphCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbixapi.API)

	graph := buildGraphObject(d)

	graphs := zabbixapi.Graphs{*graph}
	err := api.GraphsCreate(graphs)
	if err != nil {
		return err
	}

	d.SetId(graphs[0].GraphID)

	log.Printf("[DEBUG] Created graph with ID %s", d.Id())

	return resourceZabbixGraphRead(d, meta)
}

func resourceZabbixGraphRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbixapi.API)

	log.Printf("[DEBUG] Will read graph with ID %s", d.Id())

	graphs, err := api.GraphsGet(zabbixapi.Params{
		"graphids":        d.Id(),
		"selectGraphItems": "extend",
	})

	if err != nil {
		return err
	}

	if len(graphs) != 1 {
		return fmt.Errorf("Expected one graph with ID %s and got %d graphs", d.Id(), len(graphs))
	}

	graph := graphs[0]

	log.Printf("[DEBUG] Graph name is %s", graph.Name)

	d.Set("name", graph.Name)
	d.Set("width", graph.Width)
	d.Set("height", graph.Height)
	d.Set("yaxismin", graph.YAxisMin)
	d.Set("yaxismax", graph.YAxisMax)
	d.Set("show_work_period", graph.Show.WorkPeriod == 1)
	d.Set("show_triggers", graph.Show.Triggers == 1)
	d.Set("show_legend", graph.Show.Legend == 1)
	d.Set("show_3d", graph.Show.ThreeDimensional == 1)
	d.Set("percent_left", graph.PercentLeft)
	d.Set("percent_right", graph.PercentRight)
	d.Set("type", GraphTypeStringMap[graph.Type])

	// Get the host ID associated with this graph
	hostID, err := getGraphHostID(api, d.Id())
	if err != nil {
		return err
	}
	d.Set("host_id", hostID)

	// Build the graph items
	graphItems := make([]map[string]interface{}, len(graph.GraphItems))
	for i, item := range graph.GraphItems {
		graphItems[i] = map[string]interface{}{
			"item_id":   item.ItemID,
			"color":     item.Color,
			"calc_fnc":  item.CalcFunction,
			"type":      item.Type,
			"draw_type": DrawTypeStringMap[item.DrawType],
			"yaxisside": YAxisSideStringMap[item.YAxisSide],
		}
	}
	d.Set("graph_items", graphItems)

	return nil
}

func resourceZabbixGraphUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbixapi.API)

	graph := buildGraphObject(d)
	graph.GraphID = d.Id()

	graphs := zabbixapi.Graphs{*graph}
	err := api.GraphsUpdate(graphs)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updated graph with ID %s", d.Id())

	return resourceZabbixGraphRead(d, meta)
}

func resourceZabbixGraphDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*zabbixapi.API)

	graphIDs := []string{d.Id()}
	err := api.GraphsDelete(graphIDs)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleted graph with ID %s", d.Id())

	d.SetId("")
	return nil
}

func resourceZabbixGraphExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*zabbixapi.API)

	graphs, err := api.GraphsGet(zabbixapi.Params{
		"graphids": d.Id(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") {
			log.Printf("[DEBUG] Graph with ID %s doesn't exist", d.Id())
			return false, nil
		}
		return false, err
	}
	return len(graphs) > 0, nil
}

func buildGraphObject(d *schema.ResourceData) *zabbixapi.Graph {
	graph := zabbixapi.Graph{
		Name:        d.Get("name").(string),
		Width:       d.Get("width").(int),
		Height:      d.Get("height").(int),
		YAxisMin:    d.Get("yaxismin").(string),
		YAxisMax:    d.Get("yaxismax").(string),
		PercentLeft: d.Get("percent_left").(string),
		PercentRight: d.Get("percent_right").(string),
		Type:        GraphType[d.Get("type").(string)],
		Show: zabbixapi.GraphShow{
			Legend:           boolToInt(d.Get("show_legend").(bool)),
			WorkPeriod:       boolToInt(d.Get("show_work_period").(bool)),
			Triggers:         boolToInt(d.Get("show_triggers").(bool)),
			ThreeDimensional: boolToInt(d.Get("show_3d").(bool)),
		},
	}

	graphItems := d.Get("graph_items").([]interface{})
	items := make(zabbixapi.GraphItems, len(graphItems))
	
	for i, v := range graphItems {
		item := v.(map[string]interface{})
		items[i] = zabbixapi.GraphItem{
			ItemID:       item["item_id"].(string),
			Color:        item["color"].(string),
			CalcFunction: item["calc_fnc"].(int),
			Type:         item["type"].(int),
			DrawType:     DrawType[item["draw_type"].(string)],
			YAxisSide:    YAxisSide[item["yaxisside"].(string)],
		}
	}
	
	graph.GraphItems = items

	return &graph
}

func getGraphHostID(api *zabbixapi.API, graphID string) (string, error) {
	// Get graph items 
	graphs, err := api.GraphsGet(zabbixapi.Params{
		"graphids":         graphID,
		"selectGraphItems": []string{"itemid"},
	})
	
	if err != nil {
		return "", err
	}
	
	if len(graphs) != 1 || len(graphs[0].GraphItems) < 1 {
		return "", fmt.Errorf("Failed to get graph host ID: no graph items found")
	}
	
	// Get first item to determine host ID
	items, err := api.ItemsGet(zabbixapi.Params{
		"itemids":     graphs[0].GraphItems[0].ItemID,
		"selectHosts": []string{"hostid"},
	})
	
	if err != nil {
		return "", err
	}
	
	if len(items) != 1 || len(items[0].Hosts) < 1 {
		return "", fmt.Errorf("Failed to get graph host ID: no host found for item")
	}
	
	return items[0].Hosts[0].HostID, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}