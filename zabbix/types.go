package zabbix

import (
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
	Type         GraphType        `json:"graphtype,string"`
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
	Type         int    `json:"type,string"`
	YAxisSide    int    `json:"yaxisside,string"`
	CalcFunction int    `json:"calc_fnc,string"`
	DrawType     int    `json:"drawtype,string"`
}

type GraphItems []GraphItem
type Graphs []Graph

// Dashboard represents a Zabbix dashboard
type Dashboard struct {
	DashboardID string         `json:"dashboardid"`
	Name        string         `json:"name"`
	UserID      string         `json:"userid"`
	Pages       DashboardPages `json:"pages"`
	Users       interface{}    `json:"users"`
	UserGroups  interface{}    `json:"usergroups"`
}

type DashboardPage struct {
	PageID        string           `json:"pageid"`
	Name          string           `json:"name"`
	DisplayPeriod int              `json:"display_period,string"`
	Widgets       DashboardWidgets `json:"widgets"`
}

type DashboardWidget struct {
	WidgetID string                `json:"widgetid"`
	Type     int                   `json:"type,string"`
	Name     string                `json:"name"`
	X        int                   `json:"x,string"`
	Y        int                   `json:"y,string"`
	Width    int                   `json:"width,string"`
	Height   int                   `json:"height,string"`
	Fields   DashboardWidgetFields `json:"fields"`
}

type DashboardWidgetField struct {
	Type  int    `json:"type,string"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DashboardWidgetFields []DashboardWidgetField

// GetResourceID gets the resource ID from widget fields (for graphs, etc.)
func (f DashboardWidgetFields) GetResourceID() (string, bool) {
	for _, field := range f {
		if field.Type == 0 && field.Name == "graphid" {
			return field.Value, true
		}
	}
	return "", false
}

// GetText gets text from widget fields
func (f DashboardWidgetFields) GetText() (string, bool) {
	for _, field := range f {
		if field.Type == 1 && field.Name == "text" {
			return field.Value, true
		}
	}
	return "", false
}

// GetURL gets URL from widget fields
func (f DashboardWidgetFields) GetURL() (string, bool) {
	for _, field := range f {
		if field.Type == 1 && field.Name == "url" {
			return field.Value, true
		}
	}
	return "", false
}

type DashboardWidgets []DashboardWidget
type DashboardPages []DashboardPage
type Dashboards []Dashboard

// GraphsGet gets information about graphs
func (api *zabbixapi.API) GraphsGet(params zabbixapi.Params) (Graphs, error) {
	var result Graphs
	response, err := api.CallWithError("graph.get", params)
	if err != nil {
		return nil, err
	}

	if err := api.ConvertResponse(response.Result, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GraphsCreate creates new graphs
func (api *zabbixapi.API) GraphsCreate(graphs Graphs) error {
	response, err := api.CallWithError("graph.create", graphs)
	if err != nil {
		return err
	}

	result := response.Result.(map[string]interface{})
	graphids := result["graphids"].([]interface{})

	for i, id := range graphids {
		graphs[i].GraphID = id.(string)
	}

	return nil
}

// GraphsUpdate updates graphs
func (api *zabbixapi.API) GraphsUpdate(graphs Graphs) error {
	response, err := api.CallWithError("graph.update", graphs)
	if err != nil {
		return err
	}

	return nil
}

// GraphsDelete deletes graphs
func (api *zabbixapi.API) GraphsDelete(ids []string) error {
	response, err := api.CallWithError("graph.delete", ids)
	if err != nil {
		return err
	}

	return nil
}

// DashboardsGet gets information about dashboards
func (api *zabbixapi.API) DashboardsGet(params zabbixapi.Params) (Dashboards, error) {
	var result Dashboards
	response, err := api.CallWithError("dashboard.get", params)
	if err != nil {
		return nil, err
	}

	if err := api.ConvertResponse(response.Result, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// DashboardsCreate creates new dashboards
func (api *zabbixapi.API) DashboardsCreate(dashboards Dashboards) error {
	response, err := api.CallWithError("dashboard.create", dashboards)
	if err != nil {
		return err
	}

	result := response.Result.(map[string]interface{})
	dashboardids := result["dashboardids"].([]interface{})

	for i, id := range dashboardids {
		dashboards[i].DashboardID = id.(string)
	}

	return nil
}

// DashboardsUpdate updates dashboards
func (api *zabbixapi.API) DashboardsUpdate(dashboards Dashboards) error {
	response, err := api.CallWithError("dashboard.update", dashboards)
	if err != nil {
		return err
	}

	return nil
}

// DashboardsDelete deletes dashboards
func (api *zabbixapi.API) DashboardsDelete(ids []string) error {
	response, err := api.CallWithError("dashboard.delete", ids)
	if err != nil {
		return err
	}

	return nil
}