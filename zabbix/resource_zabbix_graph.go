package zabbix

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	zabbixapi "github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// GraphTypeMap maps string values to GraphType constants
var GraphTypeMap = map[string]int{
	"normal":    0,
	"stacked":   1,
	"pie":       2,
	"exploded":  3,
}

// GraphTypeStringMap is a mapping of numeric graph types to strings
var GraphTypeStringMap = map[int]string{
	0: "normal",
	1: "stacked",
	2: "pie",
	3: "exploded",
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
				MinItems: 1,
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
	api := meta.(*ZabbixGraphAPI)

	graph := buildGraphObject(d)

	graphsToCreate := Graphs{*graph}
	err := api.GraphsCreate(graphsToCreate)
	if err != nil {
		return fmt.Errorf("Error creating Zabbix graph %s: %v", graph.Name, err)
	}

	graph.GraphID = graphsToCreate[0].GraphID
	d.SetId(graph.GraphID)

	log.Printf("[DEBUG] Created graph with ID %s", d.Id())

	return resourceZabbixGraphRead(d, meta)
}

func resourceZabbixGraphRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*ZabbixGraphAPI)
	graphID := d.Id()
	log.Printf("[DEBUG] Will read graph with ID %s", graphID)

	graphs, err := api.GraphsGet(zabbixapi.Params{
		"graphids":        graphID,
		"selectGraphItems": "extend",
	})

	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") || strings.Contains(err.Error(), "No graph found") {
			log.Printf("[WARN] Zabbix Graph (%s) not found, removing from state", graphID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Zabbix graph %s: %v", graphID, err)
	}

	if len(graphs) != 1 {
		log.Printf("[WARN] Expected one graph with ID %s but got %d graphs. Removing from state.", graphID, len(graphs))
		d.SetId("")
		return nil
	}

	graph := graphs[0]

	log.Printf("[DEBUG] Read graph: %+v", graph)

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

	readGraphItems := make([]map[string]interface{}, len(graph.GraphItems))
	for i, item := range graph.GraphItems {
		drawTypeInt, _ := strconv.Atoi(item.DrawType)
		yAxisSideInt, _ := strconv.Atoi(item.YAxisSide)
		calcFncInt, _ := strconv.Atoi(item.CalcFnc)
		typeInt, _ := strconv.Atoi(item.Type)

		readGraphItems[i] = map[string]interface{}{
			"item_id":   item.ItemID,
			"color":     item.Color,
			"calc_fnc":  calcFncInt,
			"type":      typeInt,
			"draw_type": DrawTypeStringMap[drawTypeInt],
			"yaxisside": YAxisSideStringMap[yAxisSideInt],
		}
	}
	if err := d.Set("graph_items", readGraphItems); err != nil {
		return fmt.Errorf("Error setting graph_items for graph %s: %v", graphID, err)
	}

	return nil
}

func resourceZabbixGraphUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*ZabbixGraphAPI)
	graphID := d.Id()

	graph := buildGraphObject(d)
	graph.GraphID = graphID

	graphsToUpdate := Graphs{*graph}
	err := api.GraphsUpdate(graphsToUpdate)
	if err != nil {
		return fmt.Errorf("Error updating Zabbix graph %s: %v", graphID, err)
	}

	log.Printf("[DEBUG] Updated graph with ID %s", graphID)

	return resourceZabbixGraphRead(d, meta)
}

func resourceZabbixGraphDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*ZabbixGraphAPI)
	graphID := d.Id()

	graphIDs := []string{graphID}
	err := api.GraphsDelete(graphIDs)
	if err != nil {
		if strings.Contains(err.Error(), "No graph found") || strings.Contains(err.Error(), "does not exist") {
			log.Printf("[WARN] Zabbix Graph (%s) already deleted, removing from state", graphID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting Zabbix graph %s: %v", graphID, err)
	}

	log.Printf("[DEBUG] Deleted graph with ID %s", graphID)

	d.SetId("")
	return nil
}

func resourceZabbixGraphExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	api := meta.(*ZabbixGraphAPI)
	graphID := d.Id()

	_, err := api.GraphsGet(zabbixapi.Params{
		"graphids": graphID,
		"output":   "graphid",
	})

	if err != nil {
		if strings.Contains(err.Error(), "Expected exactly one result") || strings.Contains(err.Error(), "No graph found") {
			log.Printf("[DEBUG] Graph with ID %s doesn't exist", graphID)
			return false, nil
		}
		return false, fmt.Errorf("Error checking existence of Zabbix graph %s: %v", graphID, err)
	}
	log.Printf("[DEBUG] Graph with ID %s exists", graphID)
	return true, nil
}

func buildGraphObject(d *schema.ResourceData) *Graph {
	graph := Graph{
		Name:        d.Get("name").(string),
		Width:       d.Get("width").(int),
		Height:      d.Get("height").(int),
		YAxisMin:    d.Get("yaxismin").(string),
		YAxisMax:    d.Get("yaxismax").(string),
		PercentLeft: d.Get("percent_left").(string),
		PercentRight: d.Get("percent_right").(string),
		Type:        GraphTypeMap[d.Get("type").(string)],
		Show: GraphShow{
			Legend:           boolToInt(d.Get("show_legend").(bool)),
			WorkPeriod:       boolToInt(d.Get("show_work_period").(bool)),
			Triggers:         boolToInt(d.Get("show_triggers").(bool)),
			ThreeDimensional: boolToInt(d.Get("show_3d").(bool)),
		},
	}

	graphItemsRaw := d.Get("graph_items").([]interface{})
	items := make(GraphItems, len(graphItemsRaw))

	for i, v := range graphItemsRaw {
		itemMap := v.(map[string]interface{})

		drawTypeStr := itemMap["draw_type"].(string)
		yAxisSideStr := itemMap["yaxisside"].(string)

		items[i] = GraphItem{
			ItemID:     itemMap["item_id"].(string),
			Color:      itemMap["color"].(string),
			CalcFnc:    fmt.Sprint(itemMap["calc_fnc"].(int)),
			Type:       fmt.Sprint(itemMap["type"].(int)),
			DrawType:   fmt.Sprint(DrawType[drawTypeStr]),
			YAxisSide:  fmt.Sprint(YAxisSide[yAxisSideStr]),
		}
	}

	graph.GraphItems = items

	return &graph
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}