package zabbix

import (
	"fmt"
	"strings"
	"testing"

	zabbixapi "github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Basic acceptance test
func TestAccZabbixGraph_Basic(t *testing.T) {
	resourceName := "zabbix_graph.test_graph"
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
				Config: testAccZabbixGraphConfigBasic(groupName, hostName, templateName, itemKey, graphName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZabbixGraphExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", graphName),
					resource.TestCheckResourceAttr(resourceName, "width", "900"),
					resource.TestCheckResourceAttr(resourceName, "height", "200"),
					resource.TestCheckResourceAttr(resourceName, "type", "normal"),
					resource.TestCheckResourceAttr(resourceName, "graph_items.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "graph_items.0.color", "FF5555"),
				),
			},
			// Add update step if needed
		},
	})
}

// Test destroy function
func testAccCheckZabbixGraphDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*ZabbixGraphAPI)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_graph" {
			continue
		}

		_, err := api.GraphsGet(zabbixapi.Params{
			"graphids": rs.Primary.ID,
		})
		if err == nil {
			return fmt.Errorf("Graph still exists: %s", rs.Primary.ID)
		}
		// Consider Zabbix API error indicating not found as success for destroy
		if !strings.Contains(err.Error(), "No graph found") && !strings.Contains(err.Error(), "does not exist") {
		    return fmt.Errorf("Received unexpected error checking graph %s: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

// Test exists function
func testAccCheckZabbixGraphExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Graph ID is set")
		}

		api := testAccProvider.Meta().(*ZabbixGraphAPI)
		graphs, err := api.GraphsGet(zabbixapi.Params{
			"graphids": rs.Primary.ID,
			"output":   "graphid",
		})

		if err != nil {
             // Handle cases where 'not found' is the error
             if strings.Contains(err.Error(), "No graph found") || strings.Contains(err.Error(), "does not exist") {
                 return fmt.Errorf("Graph not found in Zabbix: %s", rs.Primary.ID)
             }
			 return fmt.Errorf("Error retrieving graph %s: %v", rs.Primary.ID, err)
		}
        if len(graphs) < 1 {
            return fmt.Errorf("No graph found with ID %s", rs.Primary.ID)
        }

		return nil
	}
}

// Basic graph config requiring dependent resources
func testAccZabbixGraphConfigBasic(groupName, hostName, templateName, itemKey, graphName string) string {
	return fmt.Sprintf(`
	provider "zabbix" {}

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
	  valuetype = 0 // Float
	  type = 0 // Zabbix agent
	  host_id = zabbix_template.test.id
	}

	resource "zabbix_graph" "test_graph" {
	  name = "%s"
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