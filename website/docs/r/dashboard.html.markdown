---
layout: "zabbix"
page_title: "Zabbix: zabbix_dashboard"
sidebar_current: "docs-zabbix-resource-dashboard"
description: |-
  Creates a Zabbix dashboard.
---

# zabbix_dashboard

The dashboard resource allows one to manage Zabbix dashboards.

## Example Usage

```hcl
resource "zabbix_dashboard" "system_dashboard" {
  name = "System Performance Dashboard"
  
  page {
    name = "System Metrics"
    display_period = 30
    
    # Add a graph widget
    widget {
      type   = "graph"
      name   = "CPU and Memory"
      x      = 0
      y      = 0
      width  = 12
      height = 8
      graph_id = zabbix_graph.system_performance.id
    }
    
    # Add a text widget
    widget {
      type   = "plain_text"
      name   = "Description"
      x      = 0
      y      = 8
      width  = 12
      height = 4
      text   = "This dashboard displays system performance metrics."
    }
  }
  
  # Second page with a URL widget
  page {
    name = "Help"
    display_period = 30
    
    widget {
      type   = "url"
      name   = "Zabbix Documentation"
      x      = 0
      y      = 0
      width  = 12
      height = 8
      url    = "https://www.zabbix.com/documentation/current/"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the dashboard.
* `owner_userid` - (Optional) The ID of the user that the dashboard belongs to. If not specified, defaults to the current API user.
* `page` - (Required) One or more page blocks. The page block is described below.

Page blocks support the following:

* `name` - (Required) The name of the page.
* `display_period` - (Optional) The display period of the page in seconds. Default is 30.
* `widget` - (Optional) One or more widget blocks. The widget block is described below.

Widget blocks support the following:

* `type` - (Required) The type of the widget. Possible values are "graph", "clock", "sysmap", "plain_text", "url", "trigger_info", "trigger_over", "problems", "problems_by_sv".
* `name` - (Required) The name of the widget.
* `x` - (Required) The X coordinate of the widget on the dashboard.
* `y` - (Required) The Y coordinate of the widget on the dashboard.
* `width` - (Required) The width of the widget.
* `height` - (Required) The height of the widget.
* `graph_id` - (Optional) The ID of the graph to display in the widget. Only applicable for graph widgets.
* `text` - (Optional) The text to display in the widget. Only applicable for plain_text widgets.
* `url` - (Optional) The URL to display in the widget. Only applicable for url widgets.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the dashboard.

## Import

Dashboards can be imported using the dashboard ID, e.g.

```
$ terraform import zabbix_dashboard.system_dashboard 12345
```