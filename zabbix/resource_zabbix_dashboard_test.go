package zabbix

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccZabbixDashboard_basic(t *testing.T) {
	var dashboard zabbix.Dashboard
	groupName := fmt.Sprintf("host_group_%s", acctest.RandString(5))
	hostName := fmt.Sprintf("host_%s", acctest.RandString(5))
	templateName := fmt.Sprintf("template_%s", acctest.RandString(5))
	itemKey := fmt.Sprintf("system.cpu.load[,avg%s]", acctest.RandString(3))
	graphName := fmt.Sprintf("graph_%s", acctest.RandString(5))
	dashboardName := fmt.Sprintf("dashboard_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixDashboardConfig(groupName, hostName, templateName, itemKey, graphName, dashboardName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckZabbixDashboardExists("zabbix_dashboard.test", &dashboard),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "name", dashboardName),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.#", "1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.0.name", "Page 1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.0.widget.#", "1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.0.widget.0.type", "graph"),
				),
			},
			{
				Config: testAccZabbixDashboardUpdateConfig(groupName, hostName, templateName, itemKey, graphName, dashboardName+"_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckZabbixDashboardExists("zabbix_dashboard.test", &dashboard),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "name", dashboardName+"_updated"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.#", "2"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.0.name", "Page 1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.1.name", "Page 2"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.0.widget.#", "1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.0.widget.0.type", "graph"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.1.widget.#", "1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test", "page.1.widget.0.type", "url"),
				),
			},
		},
	})
}

func testAccCheckZabbixDashboardExists(n string, dashboard *zabbix.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No dashboard ID is set")
		}

		api := testAccProvider.Meta().(*zabbix.API)
		dashboards, err := api.DashboardsGet(zabbix.Params{
			"dashboardids": rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		if len(dashboards) != 1 {
			return fmt.Errorf("Expected one dashboard, got %d", len(dashboards))
		}

		*dashboard = dashboards[0]

		return nil
	}
}

func testAccCheckZabbixDashboardDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*zabbix.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_dashboard" {
			continue
		}

		dashboards, err := api.DashboardsGet(zabbix.Params{
			"dashboardids": rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		if len(dashboards) > 0 {
			return fmt.Errorf("Dashboard still exists")
		}
	}

	return nil
}

func testAccZabbixDashboardConfig(groupName, hostName, templateName, itemKey, graphName, dashboardName string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test" {
  name = "%s"
}

resource "zabbix_template" "test" {
  host = "%s"
  groups = [zabbix_host_group.test.name]
  name = "%s"
}

resource "zabbix_item" "test" {
  name = "CPU Load"
  key = "%s"
  delay = "60"
  valuetype = "0"
  type = "0"
  host_id = zabbix_template.test.id
}

resource "zabbix_graph" "test" {
  name = "%s"
  host_id = zabbix_template.test.id
  width = 900
  height = 200
  type = "normal"
  graph_items {
    item_id = zabbix_item.test.id
    color = "FF5555"
    draw_type = "line"
  }
}

resource "zabbix_dashboard" "test" {
  name = "%s"
  
  page {
    name = "Page 1"
    display_period = 30
    
    widget {
      type = "graph"
      name = "CPU Graph"
      x = 0
      y = 0
      width = 12
      height = 5
      graph_id = zabbix_graph.test.id
    }
  }
}
`, groupName, hostName, templateName, itemKey, graphName, dashboardName)
}

func testAccZabbixDashboardUpdateConfig(groupName, hostName, templateName, itemKey, graphName, dashboardName string) string {
	return fmt.Sprintf(`
resource "zabbix_host_group" "test" {
  name = "%s"
}

resource "zabbix_template" "test" {
  host = "%s"
  groups = [zabbix_host_group.test.name]
  name = "%s"
}

resource "zabbix_item" "test" {
  name = "CPU Load"
  key = "%s"
  delay = "60"
  valuetype = "0"
  type = "0"
  host_id = zabbix_template.test.id
}

resource "zabbix_graph" "test" {
  name = "%s"
  host_id = zabbix_template.test.id
  width = 900
  height = 200
  type = "normal"
  graph_items {
    item_id = zabbix_item.test.id
    color = "FF5555"
    draw_type = "line"
  }
}

resource "zabbix_dashboard" "test" {
  name = "%s"
  
  page {
    name = "Page 1"
    display_period = 30
    
    widget {
      type = "graph"
      name = "CPU Graph"
      x = 0
      y = 0
      width = 12
      height = 5
      graph_id = zabbix_graph.test.id
    }
  }
  
  page {
    name = "Page 2"
    display_period = 60
    
    widget {
      type = "url"
      name = "Zabbix Documentation"
      x = 0
      y = 0
      width = 12
      height = 5
      url = "https://www.zabbix.com/documentation/current/manual"
    }
  }
}
`, groupName, hostName, templateName, itemKey, graphName, dashboardName)
}