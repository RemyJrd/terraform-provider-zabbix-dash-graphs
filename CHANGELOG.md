## Unreleased

- **New Feature:** Added support for Zabbix graphs through the new `zabbix_graph` resource, allowing management of graphs from template items
- **New Feature:** Added support for Zabbix dashboards through the new `zabbix_dashboard` resource, enabling creation and management of dashboards with various widget types
- **New Example:** Added an example for using the new graph and dashboard resources in `examples/graph_dashboard/`

## 0.3.1 (October 20, 2020)

NOTES:

- Tested with Zabbix 3.2, 3.4, 4.0 and 4.4
- First release of the Claranet fork to the terraform registry

FEATURES:

- **Documentation:** Added some documentation ([#10](https://github.com/claranet/terraform-provider-zabbix/pull/10))
- **License:** License file added by original author

## 0.3.0 (March 27, 2020)

NOTES:

- Tested with Zabbix 3.2, 3.4, 4.0 and 4.4
- First release of the Claranet fork

FEATURES:

- **New Data Source:** `zabbix_server` ([#XX]())
- **New Resource:** `zabbix_item` ([#XX]())
- **New Resource:** `zabbix_trigger` ([#XX]())
- **New Resource:** `zabbix_template` ([#XX]())
- **New Resource:** `zabbix_template_link` ([#XX]())
- **New Resource:** `zabbix_lld_rule` ([#XX]())
- **New Resource:** `zabbix_item_prototype` ([#XX]())
- **New Resource:** `zabbix_trigger_prototype` ([#XX]())

IMPROVEMENTS:

- Retry create and delete operations on transient failures

BUG FIXES:

- Resolve issue when api.Version() is called concurrently to other methods

## 0.2.7 (May 23, 2019)

NOTES:

- Support terraform 0.12 (upstream release)

## 0.2.6 (May 23, 2019)

NOTES:

- Support terraform 0.10.x (upstream release)

## 0.2.5 (July 05, 2017)

NOTES:

- Support terraform 0.9.11 (upstream release)

## 0.2.4 (February 21, 2017)

NOTES:

- First release (upstream release)
