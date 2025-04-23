---
page_title: "zabbix_graph Resource - terraform-provider-zabbix-dash-graphs"
subcategory: ""
description: |-
  Manages Zabbix custom graphs.
---

# zabbix_graph Resource

Provides a Zabbix custom graph resource. This can be used to create and manage custom graphs, which are not linked to a template.

## Example Usage

```terraform
resource "zabbix_host_group" "example" {
  name = "Example Group"
}

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

resource "zabbix_item" "memory_usage" {
  name    = "Memory Usage"
  key     = "vm.memory.size[pavailable]"
  host_id = zabbix_host.example.id
  type    = 0 // Zabbix Agent
  value_type = 3 // Numeric unsigned
  delay   = "1m"
  units   = "%"
}

resource "zabbix_graph" "example_graph" {
  name   = "CPU and Memory Graph"
  width  = 1000
  height = 300
  yaxis_type = 0 // Calculated

  // Left Y-axis for CPU Load
  graph_items {
    item_id     = zabbix_item.cpu_load.id
    color       = "1A7C11"
    calc_fnc    = 2 // Average
    draw_type   = 5 // Filled region
    sort_order  = 1
    y_axis_side = 0 // Left
  }

  // Right Y-axis for Memory Usage (%)
  graph_items {
    item_id     = zabbix_item.memory_usage.id
    color       = "F63100"
    calc_fnc    = 4 // Last
    draw_type   = 0 // Line
    sort_order  = 0
    y_axis_side = 1 // Right
  }
}
```

## Argument Reference

The following arguments are supported:

*   `name` - (Required) The technical name of the graph.
*   `width` - (Optional) Width of the graph in pixels. Defaults to `900`.
*   `height` - (Optional) Height of the graph in pixels. Defaults to `200`.
*   `yaxis_type` - (Optional) Y-axis maximum calculation method. `0` - calculated, `1` - fixed, `2` - item. Defaults to `0`.
*   `yaxis_min_type` - (Optional) Y-axis minimum calculation method. `0` - calculated, `1` - fixed, `2` - item. Defaults to `0`.
*   `yaxis_max_type` - (Optional) Alias for `yaxis_type`. Supported for compatibility. Defaults to `0`.
*   `yaxis_min_itemid` - (Optional) Item ID used for the Y-axis minimum value if `yaxis_min_type` is `2`. Defaults to `"0"`.
*   `yaxis_max_itemid` - (Optional) Item ID used for the Y-axis maximum value if `yaxis_type` is `2`. Defaults to `"0"`.
*   `yaxix_min_value` - (Optional) Fixed minimum value for the Y-axis if `yaxis_min_type` is `1`. Defaults to `0.0`.
*   `yaxix_max_value` - (Optional) Fixed maximum value for the Y-axis if `yaxis_type` is `1`. Defaults to `100.0`.
*   `graph_items` - (Required) List of graph items (lines, regions, etc.) to display on the graph. At least one item is required.
    *   `item_id` - (Required) The ID of the Zabbix Item to display.
    *   `color` - (Optional) Hexadecimal color code for the graph item. Defaults to `"000000"` (black).
    *   `calc_fnc` - (Optional) Value calculation function. `1` - minimum, `2` - average, `4` - maximum, `7` - all. Defaults to `2` (average).
    *   `draw_type` - (Optional) Drawing style. `0` - line, `1` - filled region, `2` - bold line, `3` - dot, `4` - dashed line, `5` - gradient line. Defaults to `0`.
    *   `sort_order` - (Optional) Drawing order, lower numbers drawn first. Defaults to `0`.
    *   `type` - (Optional) Type of graph item. `0` - simple, `1` - graph sum. Defaults to `0`.
    *   `y_axis_side` - (Optional) Y-axis side. `0` - left, `1` - right. Defaults to `0`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

*   `graph_id` - The ID of the graph in Zabbix.

## Import

Graphs can be imported using their graph ID, e.g.

```bash
terraform import zabbix_graph.example_graph 456
``` 