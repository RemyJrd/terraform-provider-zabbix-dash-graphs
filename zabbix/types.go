package zabbix

import (
	"encoding/json"
	"fmt"
	zabbixapi "github.com/claranet/go-zabbix-api"
)

// GraphType defines graph types
type GraphType int

const (
	Normal GraphType = iota
	Stacked
	Pie
	Exploded
)

// YAxisSideVal defines the Y axis position type
type YAxisSideVal int

const (
	Left YAxisSideVal = iota
	Right
)

// DrawTypeVal defines the graph item draw type
type DrawTypeVal int

const (
	Line DrawTypeVal = iota
	FilledRegion
	BoldLine
	Dot
	DashedLine
	GradientLine
)

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

// Graph represents a Zabbix graph
type Graph struct {
	GraphID      string           `json:"graphid"`
	Name         string           `json:"name"`
	Width        int              `json:"width,string"`
	Height       int              `json:"height,string"`
	YAxisMin     string           `json:"yaxismin"`
	YAxisMax     string           `json:"yaxismax"`
	PercentLeft  string           `json:"percent_left"`
	PercentRight string           `json:"percent_right"`
	Show         GraphShow        `json:"show"`
	Type         int              `json:"graphtype,string"`
	GraphItems   GraphItems       `json:"gitems"`
	Hosts        []zabbixapi.Host `json:"hosts"`
	ParentHosts  []zabbixapi.Host `json:"hosts"`
}

type GraphShow struct {
	Legend           int `json:"legend,string"`
	WorkPeriod       int `json:"work_period,string"`
	Triggers         int `json:"triggers,string"`
	ThreeDimensional int `json:"3d,string"`
}

type GraphItem struct {
	ItemID       string `json:"itemid"`
	Color        string `json:"color"`
	Type         string `json:"type,string"`
	YAxisSide    string `json:"yaxisside,string"`
	CalcFnc      string `json:"calc_fnc,string"`
	DrawType     string `json:"drawtype,string"`
}

type GraphItems []GraphItem
type Graphs []Graph

// ZabbixGraphAPI is a wrapper around the Zabbix API that provides additional functionality
type ZabbixGraphAPI struct {
	API *zabbixapi.API
}

// GraphsGet retrieves Zabbix graphs based on the provided parameters
func (z *ZabbixGraphAPI) GraphsGet(params zabbixapi.Params) (Graphs, error) {
	resp, err := z.API.CallWithError("graph.get", params)
	if err != nil {
		return nil, err
	}

	graphs := make(Graphs, 0)
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resultBytes, &graphs)
	return graphs, err
}

// GraphsCreate creates new graphs and populates their IDs
func (z *ZabbixGraphAPI) GraphsCreate(graphs Graphs) error {
	resp, err := z.API.CallWithError("graph.create", graphs)
	if err != nil {
		return err
	}

	// Extract the created graph IDs from the response
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected format for graph.create result")
	}
	graphIDsRaw, ok := result["graphids"].([]interface{})
	if !ok {
		return fmt.Errorf("could not find 'graphids' in graph.create result")
	}

	if len(graphIDsRaw) != len(graphs) {
		return fmt.Errorf("number of created graph IDs (%d) does not match input graphs (%d)", len(graphIDsRaw), len(graphs))
	}

	for i, idRaw := range graphIDsRaw {
		id, ok := idRaw.(string)
		if !ok {
			return fmt.Errorf("graph ID at index %d is not a string", i)
		}
		graphs[i].GraphID = id
	}
	return nil
}

// GraphsUpdate updates existing graphs
func (z *ZabbixGraphAPI) GraphsUpdate(graphs Graphs) error {
	_, err := z.API.CallWithError("graph.update", graphs)
	return err
}

// GraphsDelete deletes Zabbix graphs by their IDs
func (z *ZabbixGraphAPI) GraphsDelete(ids []string) error {
	_, err := z.API.CallWithError("graph.delete", ids)
	return err
}

// ZabbixItem represents a Zabbix item with additional fields not in the API's Item struct
type ZabbixItem struct {
	ItemID      string           `json:"itemid"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Key         string           `json:"key_"`
	HostID      string           `json:"hostid"`
	Hosts       []zabbixapi.Host `json:"hosts"`
}

// ItemsGet retrieves Zabbix items based on the provided parameters
func (z *ZabbixGraphAPI) ItemsGet(params zabbixapi.Params) ([]ZabbixItem, error) {
	resp, err := z.API.CallWithError("item.get", params)
	if err != nil {
		return nil, err
	}

	items := make([]ZabbixItem, 0)
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resultBytes, &items)
	return items, err
}

// ZabbixTemplate represents a Zabbix template with additional fields
type ZabbixTemplate struct {
	TemplateID string `json:"templateid"`
	Name       string `json:"name"`
	Items      []interface{} `json:"items"`
	Triggers   []interface{} `json:"triggers"`
	Graphs     []interface{} `json:"graphs"`
}

// TemplatesGet retrieves Zabbix templates based on the provided parameters
func (z *ZabbixGraphAPI) TemplatesGet(params zabbixapi.Params) ([]ZabbixTemplate, error) {
	resp, err := z.API.CallWithError("template.get", params)
	if err != nil {
		return nil, err
	}

	templates := make([]ZabbixTemplate, 0)
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resultBytes, &templates)
	return templates, err
}

// --- Dashboard API Methods --- 
// Note: These assume the underlying Zabbix API supports these calls.
// If not, the CallWithError will likely fail.

// Dashboard represents a Zabbix dashboard (structure might need adjustment based on actual API)
type Dashboard struct {
	DashboardID string `json:"dashboardid,omitempty"`
	Name        string `json:"name"`
	UserID      string `json:"userid,omitempty"` // Assuming dashboards are user-specific
	Pages       []DashboardPage `json:"pages"`
}

type DashboardPage struct {
	DashboardPageID string `json:"dashboard_pageid,omitempty"`
	Name            string `json:"name"`
	DisplayPeriod   int    `json:"display_period,omitempty"`
	SortOrder       int    `json:"sort_order,omitempty"`
	Widgets         []DashboardWidget `json:"widgets"`
}

type DashboardWidget struct {
	WidgetID string `json:"widgetid,omitempty"`
	Type     string `json:"type"`
	Name     string `json:"name,omitempty"` // Name might be optional depending on widget type
	X        int    `json:"x,string"`
	Y        int    `json:"y,string"`
	Width    int    `json:"width,string"`
	Height   int    `json:"height,string"`
	ViewMode int    `json:"view_mode,omitempty"`
	Fields   []DashboardWidgetField `json:"fields"`
}

// DashboardWidgetField represents dynamic fields within a widget
// Structure highly depends on the widget type
type DashboardWidgetField struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value interface{} `json:"value"`
}

type Dashboards []Dashboard

// DashboardsGet retrieves Zabbix dashboards
func (z *ZabbixGraphAPI) DashboardsGet(params zabbixapi.Params) (Dashboards, error) {
	resp, err := z.API.CallWithError("dashboard.get", params)
	if err != nil {
		return nil, err
	}
	dashboards := make(Dashboards, 0)
	resultBytes, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resultBytes, &dashboards)
	return dashboards, err
}

// DashboardsCreate creates new dashboards
func (z *ZabbixGraphAPI) DashboardsCreate(dashboards Dashboards) error {
	resp, err := z.API.CallWithError("dashboard.create", dashboards)
	if err != nil {
		return err
	}
	
	// Extract IDs similar to GraphsCreate
	result, ok := resp.Result.(map[string]interface{})
	if !ok { return fmt.Errorf("unexpected format for dashboard.create result") }
	dashboardIDsRaw, ok := result["dashboardids"].([]interface{})
	if !ok { return fmt.Errorf("could not find 'dashboardids' in dashboard.create result") }
	if len(dashboardIDsRaw) != len(dashboards) { return fmt.Errorf("ID count mismatch") }
	for i, idRaw := range dashboardIDsRaw {
		id, ok := idRaw.(string)
		if !ok { return fmt.Errorf("dashboard ID %d not string", i) }
		dashboards[i].DashboardID = id
	}
	return nil
}

// DashboardsUpdate updates existing dashboards
func (z *ZabbixGraphAPI) DashboardsUpdate(dashboards Dashboards) error {
	_, err := z.API.CallWithError("dashboard.update", dashboards)
	return err
}

// DashboardsDelete deletes Zabbix dashboards by their IDs
func (z *ZabbixGraphAPI) DashboardsDelete(ids []string) error {
	_, err := z.API.CallWithError("dashboard.delete", ids)
	return err
} 