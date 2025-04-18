provider "zabbix" {
  user       = "Admin"
  password   = "zabbix"
  server_url = "http://localhost/api_jsonrpc.php"
}

resource "zabbix_host_group" "example_group" {
  name = "Example Group"
}

resource "zabbix_template" "example_template" {
  host        = "example_template"
  name        = "Example Template"
  description = "Template for demonstrating graph and dashboard support"
  groups      = [zabbix_host_group.example_group.name]
  macro = {
    CPU_THRESHOLD = "80"
  }
}

# Create CPU load item
resource "zabbix_item" "cpu_load_avg1" {
  name        = "CPU Load AVG 1min"
  key         = "system.cpu.load[,avg1]"
  delay       = 60
  history     = "90d"
  trends      = "365d"
  host_id     = zabbix_template.example_template.id
}

# Create memory usage item
resource "zabbix_item" "memory_util" {
  name        = "Memory Utilization"
  key         = "vm.memory.utilization"
  delay       = 60
  history     = "90d"
  trends      = "365d"
  host_id     = zabbix_template.example_template.id
}

# Create a graph with CPU and Memory items
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
    color     = "FF5555"  # Red
    draw_type = "line"
    yaxisside = "left"
  }
  
  # Second graph item - Memory utilization
  graph_items {
    item_id   = zabbix_item.memory_util.id
    color     = "5555FF"  # Blue
    draw_type = "filled_region"
    yaxisside = "right"
  }
}

# Create a dashboard with the graph
resource "zabbix_dashboard" "system_dashboard" {
  name = "System Performance Dashboard"
  
  page {
    name = "System Metrics"
    display_period = 30
    
    # Add the graph as a widget
    widget {
      type   = "graph"
      name   = "CPU and Memory"
      x      = 0
      y      = 0
      width  = 12
      height = 8
      graph_id = zabbix_graph.system_performance.id
    }
    
    # Add a text widget with some description
    widget {
      type   = "plain_text"
      name   = "Description"
      x      = 0
      y      = 8
      width  = 12
      height = 4
      text   = "This dashboard displays system performance metrics including CPU load and memory utilization."
    }
  }
  
  # Second page with a URL to Zabbix documentation
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

# Link everything to the template
resource "zabbix_template_link" "example_template_link" {
  template_id = zabbix_template.example_template.id
  
  dynamic "item" {
    for_each = [
      zabbix_item.cpu_load_avg1.id,
      zabbix_item.memory_util.id,
    ]
    
    content {
      item_id = item.value
    }
  }
  
  # Link the graph to the template
  graph {
    graph_id = zabbix_graph.system_performance.id
  }
}

# Output the IDs
output "template_id" {
  value = zabbix_template.example_template.id
}

output "graph_id" {
  value = zabbix_graph.system_performance.id
}

output "dashboard_id" {
  value = zabbix_dashboard.system_dashboard.id
}