package zabbix

import (
	"fmt"
	"log"
	"strings"

	zabbixapi "github.com/claranet/go-zabbix-api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// WidgetTypeStringMap is a mapping of numeric widget types to strings
var WidgetTypeStringMap = map[int]string{
	0:  "graph",
	1:  "clock",
	2:  "sysmap",
	3:  "plain_text",
	10: "url",
	4:  "trigger_info",
	5:  "trigger_over",
	9:  "problems",
	8:  "problems_by_sv",
}

// ... rest of file unchanged ...