---
page_title: "zabbix_dashboard Resource - terraform-provider-zabbix-dash-graphs"
subcategory: ""
description: |-
  Manages Zabbix dashboards.
---

# zabbix_dashboard Resource

Provides a Zabbix dashboard resource. This can be used to create and manage dashboards.

## Example Usage

```terraform
resource "zabbix_host" "example" {
  host = "example-host"
  name = "Example Host"
  interfaces {
    type = 1
    main = 1
    ip   = "192.168.1.100"
    port = 10050
    dns  = ""
  }
  groups = [zabbix_host_group.example.id]
}

resource "zabbix_item" "cpu_load" {
  name    = "CPU Load"
  key     = "system.cpu.load[percpu,avg1]"
  host_id = zabbix_host.example.id
  type    = 0 // Zabbix Agent
  value_type = 0 // Numeric float
  delay   = "30s"
}

resource "zabbix_graph" "cpu_graph" {
  name    = "CPU Load Graph"
  width   = 900
  height  = 200
  graph_items {
    item_id = zabbix_item.cpu_load.id
    color   = "1A7C11"
  }
}

resource "zabbix_dashboard" "example_dashboard" {
  name           = "My Example Dashboard"
  display_period = 7200 // 2 hours
  auto_start     = true
  private        = false // Make it public

  pages {
    name = "Overview"
    widgets {
      type      = "graph" // Display a graph
      width     = 12
      height    = 4
      x         = 0
      y         = 0
      graph_ids = [zabbix_graph.cpu_graph.id]
    }
  }

  pages {
     name = "Details"
     // Add more widgets for the second page if needed
  }
}
```

## Argument Reference

The following arguments are supported:

*   `name` - (Required) The name of the dashboard.
*   `display_period` - (Optional) Dashboard refresh interval in seconds. Defaults to `3600`.
*   `auto_start` - (Optional) Whether the dashboard slideshow should start automatically. Defaults to `true`.
*   `private` - (Optional) Whether the dashboard is private (accessible only by owner/admin) or public. Defaults to `true`.
*   `pages` - (Optional) A list of dashboard pages. If omitted, a single default page named "Page 1" will be created. At least one page block is required by the Zabbix API.
    *   `name` - (Required) The name of the dashboard page.
    *   `display_period` - (Optional) Page's refresh interval in seconds. Defaults to `0` (uses dashboard's interval).
    *   `widgets` - (Optional) A list of widgets displayed on the page.
        *   `type` - (Required) The type of the widget (e.g., "graph", "item", "text", etc.). Currently, only "graph" is shown in the example.
        *   `width` - (Required) Width of the widget in dashboard grid units.
        *   `height` - (Required) Height of the widget in dashboard grid units.
        *   `x` - (Optional) Horizontal position (column) of the widget in the grid. Defaults to `0`.
        *   `y` - (Optional) Vertical position (row) of the widget in the grid. Defaults to `0`.
        *   `view_mode` - (Optional) Widget-specific view mode (interpretation depends on widget type). Defaults to `0`.
        *   `graph_ids` - (Optional) A list of Zabbix Graph IDs to display if `type` is "graph".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

*   `dashboard_id` - The ID of the dashboard in Zabbix.
*   `pages.*.widgets.*.widget_id` - The ID of the created widget.

## Import

Dashboards can be imported using their dashboard ID, e.g.

```bash
terraform import zabbix_dashboard.example_dashboard 123
``` 