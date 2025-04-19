package zabbix

import (
	"fmt"
	"strings"
	"testing"

	zabbixapi "github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccZabbixDashboard_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZabbixDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZabbixDashboardConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZabbixDashboardExists("zabbix_dashboard.test_dashboard"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test_dashboard", "name", "Test Dashboard - Basic"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test_dashboard", "pages.#", "1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test_dashboard", "pages.0.widgets.#", "1"),
					resource.TestCheckResourceAttr("zabbix_dashboard.test_dashboard", "pages.0.widgets.0.type", "system.clock"),
				),
			},
			// Add update step if needed
		},
	})
}

func testAccCheckZabbixDashboardDestroy(s *terraform.State) error {
	api := testAccProvider.Meta().(*ZabbixGraphAPI)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zabbix_dashboard" {
			continue
		}

		_, err := api.DashboardsGet(zabbixapi.Params{
			"dashboardids": rs.Primary.ID,
		})
		if err == nil {
			return fmt.Errorf("Dashboard still exists: %s", rs.Primary.ID)
		}
		// Consider Zabbix API error indicating not found as success for destroy
		if !strings.Contains(err.Error(), "No dashboard found") && !strings.Contains(err.Error(), "does not exist") {
		    return fmt.Errorf("Received unexpected error checking dashboard %s: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckZabbixDashboardExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Dashboard ID is set")
		}

		api := testAccProvider.Meta().(*ZabbixGraphAPI)
		 dashboards, err := api.DashboardsGet(zabbixapi.Params{
			"dashboardids": rs.Primary.ID,
			"output":       "dashboardid",
		 })

		if err != nil {
             // Handle cases where 'not found' is the error
             if strings.Contains(err.Error(), "No dashboard found") || strings.Contains(err.Error(), "does not exist") {
                 return fmt.Errorf("Dashboard not found in Zabbix: %s", rs.Primary.ID)
             }
			 return fmt.Errorf("Error retrieving dashboard %s: %v", rs.Primary.ID, err)
		}
        if len(dashboards) < 1 {
            return fmt.Errorf("No dashboard found with ID %s", rs.Primary.ID)
        }

		return nil
	}
}

func testAccZabbixDashboardConfigBasic() string {
	return `
	provider "zabbix" {}
	
	resource "zabbix_dashboard" "test_dashboard" {
	  name = "Test Dashboard - Basic"
	  pages {
	    name = "Page 1"
	    widgets {
	      type = "system.clock"
	      x = 0
	      y = 0
	      width = 4
	      height = 2
	      fields {
	        type = 2 // Location
	        name = "clocktype"
	        value = "0" // Server time
	      }
	    }
	  }
	}
	`
}