---
layout: "zabbix"
page_title: "Zabbix: zabbix_graph"
sidebar_current: "docs-zabbix-resource-graph"
description: |-
  Creates a Zabbix graph.
---

# zabbix_graph

The graph resource allows one to manage Zabbix graphs.

## Example Usage

```hcl
resource "zabbix_graph" "system_performance" {
  name     = "System Performance"
  host_id  = zabbix_template.example_template.id
  width    = 900
  height   = 300
  type     = "normal"
  yaxismin = "0"
  
  # First graph item - CPU load
  graph_items {
    item_id   = zabbix_item.cpu_load_avg1.id
    color     = "FF5555"
    draw_type = "line"
    yaxisside = "left"
  }
  
  # Second graph item - Memory utilization
  graph_items {
    item_id   = zabbix_item.memory_util.id
    color     = "5555FF"
    draw_type = "filled_region"
    yaxisside = "right"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the graph.
* `host_id` - (Required) The ID of the host or template that the graph belongs to.
* `width` - (Optional) The width of the graph in pixels. Default is 900.
* `height` - (Optional) The height of the graph in pixels. Default is 200.
* `yaxismin` - (Optional) The minimum value of the Y axis. Default is "0".
* `yaxismax` - (Optional) The maximum value of the Y axis. Default is "100".
* `show_work_period` - (Optional) Whether to show the working period on the graph. Default is true.
* `show_triggers` - (Optional) Whether to show triggers on the graph. Default is true.
* `show_legend` - (Optional) Whether to show the legend on the graph. Default is true.
* `show_3d` - (Optional) Whether to show the graph as 3D. Default is false.
* `percent_left` - (Optional) The left percentage. Default is "0".
* `percent_right` - (Optional) The right percentage. Default is "0".
* `type` - (Optional) The type of the graph. Possible values are "normal", "stacked", "pie", "exploded". Default is "normal".
* `graph_items` - (Required) One or more graph item blocks. The graph item block is described below.

Graph item blocks support the following:

* `item_id` - (Required) The ID of the item to include in the graph.
* `color` - (Required) The color of the item in the graph, in hexadecimal RGB format.
* `calc_fnc` - (Optional) The data processing method for the item. Default is 2 (average).
* `type` - (Optional) The type of the graph item. Default is 0 (simple).
* `draw_type` - (Optional) The draw type of the item. Possible values are "line", "filled_region", "bold_line", "dot", "dashed_line", "gradient_line". Default is "line".
* `yaxisside` - (Optional) The side of the Y axis that the item is displayed on. Possible values are "left", "right". Default is "left".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the graph.

## Import

Graphs can be imported using the graph ID, e.g.

```
$ terraform import zabbix_graph.system_performance 12345
```