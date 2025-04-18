package zabbix

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccZabbixGraph_basic(t *testing.T) {
	var graph zabbix.Graph
	groupName := fmt.Sprintf("host_group_%s", acctest.RandString(5))
	hostName := fmt.Sprintf("host_%s", acctest.RandString(5))
	templateName := fmt.Sprintf("template_%s", acctest.RandString(5))
	itemKey := fmt.Sprintf("system.cpu.load[,avg%s]", acctest.RandString(3))
	graphName := fmt.Sprintf("graph_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixGraphConfig(groupName, hostName, templateName, itemKey, graphName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckZabbixGraphExists("zabbix_graph.test", &graph),
					resource.TestCheckResourceAttr("zabbix_graph.test", "name", graphName),
					resource.TestCheckResourceAttr("zabbix_graph.test", "width", "900"),
					resource.TestCheckResourceAttr("zabbix_graph.test", "height", "200"),
					resource.TestCheckResourceAttr("zabbix_graph.test", "type", "normal"),
					resource.TestCheckResourceAttr("zabbix_graph.test", "graph_items.#", "1"),
				),
			},
			{
				Config: testAccZabbixGraphUpdateConfig(groupName, hostName, templateName, itemKey, graphName+"_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckZabbixGraphExists("zabbix_graph.test", &graph),
					resource.TestCheckResourceAttr("zabbix_graph.test", "name", graphName+"_updated"),
					resource.TestCheckResourceAttr("zabbix_graph.test", "width", "800"),
					resource.TestCheckResourceAttr("zabbix_graph.test", "height", "300"),
					resource.TestCheckResourceAttr("zabbix_graph.test", "type", "stacked"),
					resource.TestCheckResourceAttr("zabbix_graph.test", "graph_items.#", "1"),
				),
			},
		},
	})
}

func testAccCheckZabbixGraphExists(n string, graph *zabbix.Graph) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No graph ID is set")
		}

		api := testAccProvider.Meta().(*zabbix.API)
		graphs, err := api.GraphsGet(zabbix.Params{
			"graphids": rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		if len(graphs) != 1 {
			return fmt.Errorf("Expected one graph, got %d", len(graphs))
		}

		*graph = graphs[0]

		return nil
	}
}

func testAccCheckZabbixGraphDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*zabbix.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_graph" {
			continue
		}

		graphs, err := api.GraphsGet(zabbix.Params{
			"graphids": rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		if len(graphs) > 0 {
			return fmt.Errorf("Graph still exists")
		}
	}

	return nil
}

func testAccZabbixGraphConfig(groupName, hostName, templateName, itemKey, graphName string) string {
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
`, groupName, hostName, templateName, itemKey, graphName)
}

func testAccZabbixGraphUpdateConfig(groupName, hostName, templateName, itemKey, graphName string) string {
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
  width = 800
  height = 300
  type = "stacked"
  graph_items {
    item_id = zabbix_item.test.id
    color = "FF5555"
    draw_type = "bold_line"
    yaxisside = "right"
  }
}
`, groupName, hostName, templateName, itemKey, graphName)
}