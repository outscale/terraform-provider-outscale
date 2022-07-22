0.5.4 (July 26, 2022)
========================

BUG FIXES:
----------
* Fix ```public_ip_id``` attribute when imported outscale_nat_service resource ([GH-95](https://github.com/outscale-dev/terraform-provider-outscale/issues/95))
* Fix issue with importing server_certificate ([GH-97](https://github.com/outscale-dev/terraform-provider-outscale/issues/97))
* Fix StartVM action when updating VM tags ([GH-86](https://github.com/outscale-dev/terraform-provider-outscale/issues/86))
* Fix ```secondary_private_ip_count``` parameter in outscale_nic_private_ip resource ([GH-100](https://github.com/outscale-dev/terraform-provider-outscale/issues/100))

IMPROVEMENT:
-----------

* Update retrying when api call throttled ([GH-106](https://github.com/outscale-dev/terraform-provider-outscale/issues/106))
* Improve integration and acceptance tests
* Add ```frieza-clean``` to clean account after running tests

0.5.3 (Mars 25, 2022)
========================

FEATURES:
---------

* Add "expiration_date" parameter to "outscale_access_key" resource and datasources (TPD-1987)

BUG FIXES:
----------

* Fix acceptance tests
* Fix public_ip datasource filters #64
* Fix tags on outscale_resources #68
* iops value is set to 0 for standard volumes (TPD-2053)

IMPROVEMENT:
-----------

* Update sdk
* Add credential checking
* Make the generation of the doc more automatic

0.5.2 (February 02, 2022)
========================

FEATURES:
---------

* Enhance User-Agent #56

BUG FIXES:
----------

* Fix Resource outscale_route_table_link import issue (TPD-2011)
* Update old wiki links to new docs address #49

0.5.1 (December 21, 2021)
========================

BUG FIXES:
----------

* Fix OMI used in examples
* Update dependency due to security alert ([CVE](https://github.com/advisories/GHSA-25xm-hr59-7c27))

0.5.0 (November 12, 2021)
========================

FEATURES:
---------

* Support terraform client version 1.0.X
* New import: outscale_security_group_rule (TPD-1892, TPD-1157)
* New Data Source: outscale_server_certificate (TPD-1979)
* New Data Source: outscale_server_certificates (TPD-1979)
* New Resource: outscale_server_certificate (TPD-1979)
* New Data Source: outscale_server_certificate (TPD-1923)
* New Data Source: outscale_snapshot_export_task (TPD-1825)
* New Data Source: outscale_snapshot_export_tasks (TPD-1825)
* New Resource: outscale_snapshot_export_task (TPD-1825)
* New Data Source: outscale_image_export_task (TPD-1923)
* New Data Source: outscale_image_export_tasks (TPD-1923)
* New Resource: outscale_image_export_task (TPD-1923)
* Support tags for the outscale_subnet resource (TPD-1976)
* Support new filters for the outscale_nic & outscale_nics data sources (TPD-1989)
* Support new filters for the outscale_security_group & outscale_security_groups data sources (TPD-1990)
* Support new filters for the outscale_image & outscale_images data sources (TPD-1991)
* Support new filters for the outscale_nat_service & outscale_nat_services data sources (TPD-1992)
* Support new filters for the outscale_client_gateway & outscale_client_gateways data sources (TPD-1993)
* Support new filters for the outscale_subnet & outscale_subnets data sources (TPD-1994)
* Support new filters for the outscale_vpn_connection & outscale_vpn_connections data sources (TPD-1995)
* Support new filters for the outscale_vm_state & outscale_vm_states data sources (TPD-1998)
* Add state argument to outscale_vm (TPD-2007)

BUG FIXES:
----------

* Fix the update of a route when only modifying the target (TPD-1963)
* Fix the descriptions filter of outscale_image & outscale_images data sources (TPD-1991)
* From the outscale_dhcp_option resource and data source, and the outscale_dhcp_options data source (TPD-1997)
* Fix the update of the outscale_nat_service resource with new subnet_id/public_ip_id values (TPD-2013)
* Remove the request_id from all data sources and resources (TPD-2015)
* Rename attributes of outscale_quota and outscale_quotas data sources to comply with the API (TPD-2024)

0.4.1 (July 19, 2021)
========================

NOTES:
------
Add arm64 binaries for linux and macOS

0.4.0 (july 9, 2021)
========================

FEATURES:
---------

* New Data Source: outscale_product_type
* New Data Source: outscale_product_types
* New Data Source: outscale_quota
* New Data Source: outscale_quotas

BUG FIXES:
----------

* Create Dockerfile to build documentation (TPD-1978)
* Create "outscale_quota" and "outscale_quotas" datasources in oAPI client (TPD-1980)
* Create "outscale_product_type" and "outscale_product_types" datasource in oAPI client (TPD-1981)
* Update osc-sdk-go (TPD-1982)
* Check route state after creation (TPD-1983)
* Check Windows admin password (TPD-1984)
* "load_balancer" resource: SecurityGroups can now be updated (TPD-2000)
* LBU getting deleted while adding a second SG (TPD-2004)
* "outscale_route" unnecessary Update action when updating "await_active_state" parameter (TPD-2005)
* "outscale_vm" unnecessary StartVms action when updating "get_admin_password" parameter (TPD-2006)


NOTES:
------

**WARNING:** When creating access keys, the secret key is stored in the Terraform state. For security reasons, it is strongly recommended to create access keys using the API instead of the Terraform resource. For more information on how to create access keys using the API, see our [official API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

It is recommended to set tags inside the resources rather than using outscale_tag.

0.3.1 (April 16, 2021)
========================

BUG FIXES:
----------
* add goreleaser files

0.3.0 (April 06, 2021)
========================

FEATURES:
---------

* New Data Source: outscale_flexible_gpu
* New Data Source: outscale_flexible_gpu_catalog
* New Data Source: outscale_flexible_gpus
* New Data Source: outscale_net_access_point
* New Data Source: outscale_net_access_points
* New Data Source: outscale_net_access_point_services
* New Data Source: outscale_regions
* New Data Source: outscale_subregion
* New Data Source: outscale_subregions
* New Data Source: outscale_vm_types

* New Resource: outscale_net_access_point
* New Resource: outscale_flexible_gpu
* New Resource: outscale_flexible_gpu_link

BUG FIXES:
----------

* Issue when importing the resource "outscale_load_balancer_listener_rule" (TPD-1941)
* terraform crash when "outscale_net_peering" resource cannot be found (TPD-1943)
* Create "outscale_net_access_point" resource and datasource(s) in oAPI client (TPD-1945)
* Create "outscale_flexible_gpu" resource and datasource(s) in oAPI client (TPD-1946)
* Create "outscale_flexible_gpu_link" resource in oAPI client (TPD-1947)
* Create "outscale_flexible_gpu_catalog" datasource in oAPI client (TPD-1948)
* Create "outscale_vm_types" datasource in oAPI client (TPD-1949)
* Add "outscale_regions" datasource in oAPI client (TPD-1950)
* Add "outscale_subregion" and "outscale_subregions" datasource in oAPI client (TPD-1951)
* Create "outscale_net_access_point_services" datasource in oAPI client (TPD-1952)
* "outscale_load_balancer_listener_rule" cannot be updated (TPD-1953)
* Terraform hangs when changing instance type on VMs with shutdown behaviour set to "restart" (TPD-1954)
* Terraform crash when "outscale_route_table" resource cannot be found (TPD-1960)
* "outscale_route_table" datasource(s) is not sending all filters (TPD-1961)
* terraform crash when "outscale_nat_service" resource cannot be found (TPD-1962)
* "outscale_internet_service" datasource(s) is not sending all filters (TPD-1964)
* Filters should not be mandatory in "outscale_vm_types" datasource in oAPI client (TPD-1968)
* "dhcp_options_set_id" attribute is missing in "outscale_dhcp_option" datasource(s) (TPD-1969)
* Examples rework (TPD-1970)
* Add CONTRIBUTING.md (TPD-1971)
* Integrate QA tests (TPD-1973)


KNOWN INCOMPATIBILITIES:
------------------------

* outscale_load_balancer datasource: When applying the same configuration file twice in a row (with non change), terraform asks fo the user confirmation to read the datasource again (TPD-1942).


NOTES:
------

**WARNING:** When creating access keys, the secret key is stored in the Terraform state. For security reasons, it is strongly recommended to create access keys using the API instead of the Terraform resource. For more information on how to create access keys using the API, see our [official API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

It is recommended to set tags inside the resources rather than using outscale_tag.


0.2.0 (November 30, 2020)
========================

FEATURES:
---------

* New Data Source: outscale_load_balancer
* New Data Source: outscale_load_balancer_listener_rule
* New Data Source: outscale_load_balancer_listener_rules
* New Data Source: outscale_load_balancer_tags
* New Data Source: outscale_load_balancer_vm_health
* New Data Source: outscale_load_balancers

* Changed Data Source: outscale_vms_states is replaced by outscale_vm_state

* New Resource: outscale_load_balancer
* New Resource: outscale_load_balancer_attributes
* New Resource: outscale_load_balancer_listener_rule
* New Resource: outscale_load_balancer_policy
* New Resource: outscale_load_balancer_vms


BUG FIXES:
----------

* oAPI outscale_load_balancer resource is using fcu type call (TPD-1739)
* Add "outscale_load_balancer" datasource and datasources (TPD-1906)
* Wrong attributes for "outscale_load_balancer_vms" resource (TPD-1907)
* "image_id" should be mandatory in "outscale_vm" resource (TPD-1911)
* terraform crash when client gateway and vpn resources cannot be found (TPD-1914)
* Changed Data Source: "outscale_vms_state" is replaced by "outscale_vm_states" (TPD-1915)
* Issue with listeners on "outscale_load_balancer" (TPD-1916)
* Wrong attributes for "outscale_load_balancer" resource (TPD-1917)
* "outscale_load_balancer" resource: terraform destroy fails for private LBU (TPD-1918)
* "outscale_load_balancer" resource: issue when creating a load balancer with multiple securtiy groups (TPD-1919)
* "outscale_load_balancer_listeners" resource: "server_certificate_id" is not sent in the request (TPD-1920)
* "outscale_load_balancer_listeners" resource: missing attributes in the state (TPD-1921)
* Issue when creating a load balancer policy (TPD-1922)
* Add "outscale_load_balancer_listener_rule" resource and datasource(s) (TPD-1925)
* "outscale_load_balancer_vm_health" datasource(s) (TPD-1926)
* Regressions on "refactor-osc-client" branch (TPD-1927)
* Missing Health check attributes in "outscale_load_balancer" resource (TPD-1928)
* Terraform crashes when creating "outscale_load_balancer_ssl_certificate" resource (TPD-1930)
* "outscale_load_balancer_tags" datasources are not supported (TPD-1931)
* Migrate all LBU attributes modifications to "outscale_load_balancer_attributes" (TPD-1932)
* Issue on terraform refresh/destroy when using "outscale_nic" and "outscale_nic_private_ip" (TPD-1933)
* Regression on ""outscale_vpn_connection"" on "refactor-osc-client" branch (TPD-1934)
* x509 client certificate authentication (TPD-1936)
* Issue with custom endpoints (TPD-1938)
* Regression on "outscale_access_keys" datasource on develop-oapi (TPD-1939)
* Issue when deactivating LBU access logs (TPD-1940)


KNOWN INCOMPATIBILITIES:
------------------------

* outscale_load_balancer datasource: When applying the same configuration file twice in a row (with non change), terraform asks fo the user confirmation to read the datasource again (TPD-1942).
* outscale_load_balancer_listener_rule: The resource cannot be imported correctly because of listener block (TPD-1941).


NOTES:
------

**WARNING:** When creating access keys, the secret key is stored in the Terraform state. For security reasons, it is strongly recommended to create access keys using the API instead of the Terraform resource. For more information on how to create access keys using the API, see our [official API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

It is recommended to set tags inside the resources rather than using outscale_tag.


0.1.2-beta (unpublished)
========================

BUG FIXES:
----------

* x509 client authentication support (TPD-1936)


0.1.1 (October 02, 2020)
========================

BUG FIXES:
----------

* Support proxy for Outscale Terraform provider (TPD-1924)


0.1.0 (August 07, 2020)
========================

NOTES: Rename of the previous release

0.1.1 (October 02, 2020)
========================

BUG FIXES:
----------

* Support proxy for Outscale Terraform provider (TPD-1924)


0.1.0 (August 07, 2020)
========================

NOTES:
Rename of the previous release


0.1.0RC9 (June 22, 2020)
========================

FEATURES:
---------

* New Data Source: outscale_access_key
* New Data Source: outscale_access_keys
* New Data Source: outscale_client_gateway
* New Data Source: outscale_client_gateways
* New Data Source: outscale_dhcp_option
* New Data Source: outscale_dhcp_options
* New Data Source: outscale_virtual_gateway
* New Data Source: outscale_virtual_gateways
* New Data Source: outscale_vpn_connection
* New Data Source: outscale_vpn_connections

* New Resource: outscale_access_key
* New Resource: outscale_client_gateway
* New Resource: outscale_dhcp_option
* New Resource: outscale_virtual_gateway
* New Resource: outscale_virtual_gateway_link
* New Resource: outscale_virtual_gateway_route_propagation
* New Resource: outscale_vpn_connection
* New Resource: outscale_vpn_connection_route


BUG FIXES:
----------

* outscale_vpn_connection resource is using fcu type call (TPD-1738)
* Wrong attributes for outscale_net_attributes resource (TPD-1803)
* Wrong attributes for outscale_vpn_gateway_link resource (TPD-1827)
* Missing attributes in outscale_vpn_gateways datasources (TPD-1828)
* outscale_vpn_gateway tags are not updated when updating the configuration file (TPD-1829)
* Migrate outscale_vpn_gateway_route_propagation resource (TPD-1830)
* Migrate outscale_dhcp_option resource and datasource(s) to oAPI client (TPD-1832)
* Missing attribute when importing outscale_virtual_gateway (TPD-1836)
* Terraform crash when importing 'outscale_internet_service_link (TPD-1838)
* outscale_net_attributes cannot be imported (TPD-1840)
* add accepter_net_id  and source_net_id attributes when importing outscale_net_peering (TPD-1844)
* Terraform crash when importing 'outscale_virtual_gateway_link (TPD-1847)
* Migrate outscale_client_gateway resource (datasource) to oAPI client (TPD-1850)
* outscale_vm resource: bloc_device_mappings arguments are not ForceNew type argument (TPD-1856)
* outscale_volumes datasources are not sending the filters (TPD-1872)
* outscale_security_group_rule: cannot create a rule with security_group_name (TPD-1874)
* Regression: outscale_route_table cannot be imported (TPD-1875)
* outscale_vm resource: performance cannot be set (TPD-1876)
* outscale_net datasource(s) is not sending the filters (TPD-1883)
* Regression: Missing attribute on outscale_route_table_link (TPD-1888)
* outscale_subnet datasource(s) is not sending available_ips_counts filter (TPD-1889)
* outscale_public_ip datasource(s) is not sending  LinkPublicIpIds filter (TPD-1890)
* Regression: outscale_snapshot cannot be imported (TPD-1891)
* outscale_nic_private_ip cannot be imported (TPD-1893)
* outscale_nic_link cannot be imported (TPD-1894)
* Migrate outscale_access_key resource and datasource(s) to oAPI client (TPD-1895)
* Regression: outscale_security_group_rule is not sending from_port and to_port arguments when they are set to 0 (TPD-1896)
* Issue with outscale_client_gateway datasource on develop-oapi branch (TPD-1899)
* Missing attributes when importing outscale_vpn_connection (TPD-1900)
* outscale_vpn_connection_route cannot be imported (TPD-1902)
* outscale_vm: issue when creating a VM with blockDeviceMapping (TPD-1904)
* outscale_subnet resource and datasource are missing map_public_ip_on_launch attribute (TPD-1905)
* Regression on outscale_net_peering_acceptation on develop-oapi branch (TPD-1908)


NOTES:
------

**WARNING:** When creating access keys, the secret key is stored in the Terraform state. For security reasons, it is strongly recommended to create access keys using the API instead of the Terraform resource. For more information on how to create access keys using the API, see our [official API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

It is recommended to set tags inside the resources rather than using outscale_tag.


0.1.0RC8.2 (April 20, 2020)
===========================

BUG FIXES:
----------

* outscale_image datasource(s) is very slow (TPD-1881)
* Issue with tags order on different resources (TPD-1882)
* outscale_image resource: tags cannot be set (TPD-1884)
* Missing tags on outscale_public_ip datasource(s) (TPD-1885)
* Issue when updating outscale_public_ip resource tags (TPD-1886)
* Missing tags on multiple datasource(s) (TPD-1887)


NOTES:
------

It is recommended to set tags inside the resources rather than using outscale_tag.



0.1.0RC8.1 (March 2, 2020)
==========================

BUG FIXES:
----------

* outscale_vm cannot be imported (TPD-1833)
* Missing attributes when importing outscale_route_table_link (TPD-1843)
* outscale_volume_link cannot be imported (TPD-1848)
* Regression on Outscale_security_group on clean-code (TPD-1867)
* Regression on Outscale_nic on clean-code branch (TPD-1868)
* Regression on Outscale_nat_service on clean-code branch (TPD-1869)
* Regression on Outscale_route_table_link on clean-code (TPD-1871)
* outscale_security_group_rule disappears after a terraform refresh (TPD-1873)
* outscale_public_ip datasource and datasources are not sending the filters (TPD-1877)


NOTES:
------

Release sent to Hasicorp for full review



0.1.0RC7 (February 21, 2020)
============================

BUG FIXES:
----------

* Remove legacy API (TPD-1752)
* outscale_vm tags are not updated when updating the configuration (TPD-1793)
* Missing attribute when importing outscale_volume (TPD-1834)
* Missing attribute when importing outscale_nat_service (TPD-1839)
* Support custom endpoints (TPD-1845)
* outscale_vm resource is always sending DeleteOnVmDeletion in nics block attributes (TPD-1846)
* Performance attribute is missing in outscale_vm resource and datasource (TPD-1853)
* Terraform crashs when route table is linked to multiple subnets (TPD-1857)
* Terraform crash when updating destination_ip_range attribute in outscale_route resource (TPD-1863)
* outscale_route resource cannot be updated (TPD-1864)



0.1.0RC6.1 (December 23, 2019)
============================== 

BUG FIXES:
----------

* outscale_Snapshot resource: cannot import snapshot (TPD-1732)
* outscale_Snapshot resource: cannot copy snapshot (TPD-1733)
* outscale_image resource: cannot copy an image (TPD-1734)
* outscale_image resource: cannot register an image (TPD-1735)
* outscale_vm resource is missing Nics attributes (TPD-1756)
* net_id and state always empty for outscale_internet_service (TPD-1770)
* outscale_nics datasources is using fcu type call (TPD-1773)
* outscale_nat_services datasource is not sending the filter to oAPI (TPD-1778)
* NetAccessPoint route attributes are empty in route_table datasource (TPD-1779)
* Regression on Nic datasource and datasources on oapi-client-external-library branch (TPD-1789)
* Regression on Images datasource and datasources on oapi-client-external-library branch (TPD-1791)
* Dependency issue on Destroy outscale_vm on oapi-client-external-library branch (TPD-1792)
* Tags are not updated when updating the configuration file for different resources (TPD-1794)
* Regression on outscale_vm when creating a vm in a net (TPD-1806)
* Regression: vm_initiated_shutdown_behavior cannot be set on outscale_vm resource in a net (TPD-1807)
* Regression on creating outscale_volume resource (IO1 volumes ) (TPD-1808)
* Cannot remove tags when updating the configuration file for different resources (TPD-1809)
* Regression: outscale_volumes_link resource fails (TPD-1810)
* Impossible to create a Nat service route in outscale_route resource (TPD-1811)
* Regression on destroy outscale_route (TPD-1812)
* nat_service_id attribute is missing in outscale_route_table datasource and datasources (TPD-1813)
* Missing Inbound_rules in outscale_security_group datasource (TPD-1814)
* outscale_security_group datasource fails if the security group has a rule to another security rule (TPD-1815)
* Regression on outscale_snapshot datasource et datasources (TPD-1816)
* outscale_image: cannot register an image anymore (TPD-1817)
* outscale_image: cannot copy an image anymore (TPD-1818)
* outscale_snapshot: cannot import a snapshot (TPD-1819)
* Missing attributes in outscale_snapshot resource (TPD-1820)
* outscale_volume datasource et datasources fail if the volume is linked to a VM (TPD-1821)
* Missing Password Attribute in outscale_vm (TPD-1822)



0.1.0RC5 (October 16, 2019)
=========================== 

BUG FIXES:
----------

* outscale_image Incorrect attribute on datasource (TPD-1776)
* Wrong attribute for Tags in outscale_subnets datasources (TPD-1763)
* outscale_nat_services datasources is missing request_id attribute (TPD-1762)
* Regression on outscale_nat_service (destroy) (TPD-1761)
* Regression: Placement attributes cannot be set in oAPI outscale_vm resource (TPD-1759)
* Regression on security_group_rule resource (TPD-1758)
* Regression on outscale_volume_link resource (TPD-1757)
* outscale_vm resource is not sending block_device_mappings attributes (TPD-1755)
* outscale_subnet datasource: wrong attribute for tags (TPD-1749)
* outscale_public_ip datasource is not sending the filter to oAPI (TPD-1748)
* image_id and subnet_id are not ForceNew type arguments (TPD-1746)
* outscale_snapshot resource: tags cannot be set (TPD-1745)
* outscale_internet_service resource: tags cannot be set (TPD-1744)
* outscale_subnet resource: tags cannot be set (TPD-1743)
* outscale_vm resource: tags cannot be set (TPD-1742)
* outscale_vm resource: private_ips and nic_id arguments cannot be set (TPD-1736)
* outscale_vm_attributes: VM are always stopped before updating VM attributes (TPD-1730)
* scale_Net_Peering datasource is using fcu type call (TPD-1729)
* outscale_nic resource: private_ips and security_group_id arguments cannot be set (TPD-1781)
* outscale_route_table fails (TPD-1726)
* outscale_vms datasources fails (TPD-1760)
* outscale_public_ips datasource is using fcu calls (TPD-1695)
* outscale_nic_private_ip fails at unlink (TPD-1727)
* outscale_nic resource and datasource are having attribute errors (TPD-1717)
* outscale_nets is reporting invalid attributes (TPD-1711)
* outscale_net datasource is not sending the filter to oAPI (TPD-1721)
* outscale_image_launch_permission fails for Permissions removals (TPD-1782)
* outscale_Snapshot_Update fails (TPD-1728)
* outscale_vm doesn't modify deletion protection attribute (TPD-1725)
* Wrong attributes for block_device_mappings in outscale_vm resource and datasource (TPD-1771)
* Wrong attribute for Tags in outscale_image resource and datasource (TPD-1768)
* Unnecessary vms.x.block_device_mappings.# attributes in outscale_vms datasources (TPD-1777)
* Unnecessary virtual_gateway_id attribute in route_table datasource (TPD-1780)
* Remove all secret information from code (TPD-1750)
* Missing Tags attributes in outscale_vm resource (regression) (TPD-1772)
* Add password atttribute to outscale_vm for windows vm (TPD-1747)
* request_id attribute is missing in outscale_security_groups datasources (TPD-1766)
* outscale_vm resource is always sending DeleteOnVmDeletion in block_device_mappings attribute (TPD-1774)
* outscale_route_tables datasources fails (TPD-1764)
* outscale_images datasources fails (TPD-1769)
* outscale_net_peerings datasources is using fcu type call (TPD-1767)



0.1.0RC4 (June 28, 2019)
======================== 

BUG FIXES:
----------

* outscale_net_peering_acceptation fails (TPD-1722)
* outscale_net_attributes resource and datasource are missing request_id attribute (TPD-1724)
* outscale_image_launch_permission fails (TPD-1719)
* outscale_vm resource is having issues with arguments (TPD-1731)
* outscale_nat_service resource is having argument issue (TPD-1702)
* outscale_nic_link is having issues with arguments (TPD-1718)
* outscale_keypairs is crashing (TPD-1710)
* outscale_volumes_link resource is generating fcu type call(TPD-1688)



0.1.0RC3 (March 29, 2019)
========================= 

BUG FIXES:
----------

* outscale_public_ip datasource is not displayed in terraform show (TPD-1684)
* outscale_public_ip_link is missing attributes (TPD-1686)
* outscale_volume resource and datasource provider error in api branch (TPD-1685)
* outscale_public_ip resource is missing request_id attribute (TPD-1683)
* make test in error on develop (TPD-1680)
* TestAccOutscaleOAPIPublicIP_basic fails on dv1 server (TPD-1681)


EXPERIMENTS:
------------

* TERPIN-0327-FCURES outscale_vpn_gateway_link (TPD-779)
* TERPIN-0480-FCUEXTRES outscale_snapshot_export_tasks (TPD-789)
* TERPIN-0352-LBURES outscale_load_balancer_cookiepolicy (TPD-861)
* TERPIN-0501-FCUDS outscale_image (TPD-640)
* TERPIN-0395-DLRES outscale_directlink (TPD-993)
* TERPIN-0417-EIMRES outscale_group_user (TPD-1477)
* TERPIN-0375-LBURES outscale_load_balancer_vms (TPD-876)
* TERPIN-0448-EIMRES outscale_policy_group (TPD-1008)
* TERPIN-0158-FCURES outscale_image_launch_permission (TPD-812)
* TERPIN-0317-FCURES outscale_vpn_connection_route (TPD-742)
* TERPIN-0347-LBURES outscale_load_balancer_health_check (TPD-1483)
* TERPIN-0379-LBURES outscale_load_balancer_attributes (TPD-870)
* TERPIN-0306-FCURES outscale_lin_api_access (TPD-901)
* TERPIN-0385-LBURES outscale_load_balancer_ssl_certificate (TPD-1489)
* TERPIN-0211-FCURES outscale_nic (TPD-822)
* TERPIN-0485-FCUEXTRES outscale_snapshot_import (TPD-896)
* TERPIN-0433-EIMRES outscale_policy_user_link (TPD-968)
* TERPIN-0153-FCURES outscale_image_copy (TPD-754)
* TERPIN-0453-EIMRES outscale_policy_user (TPD-963)
* TERPIN-0443-EIMRES outscale_policy_version (TPD-1024)
* TERPIN-0557-DLDS outscale_directlink_vpn_gateways (TPD-1612)
* TERPIN-0475-FCUEXTRES outscale_image_tasks (TPD-764)
* TERPIN-0405-EIMRES outscale_user_api_keys (TPD-1018)
* TERPIN-0400-DLRES outscale_directlink_interface (TPD-1013)
* TERPIN-0269-FCURES outscale_snapshot_copy (TPD-891)
* TERPIN-0312-FCURES outscale_vpn_connection (TPD-737)
* TERPIN-0438-EIMRES outscale_policy (TPD-1464)
* TERPIN-0232-FCURES outscale_reserved_vms_offer_purchase (TPD-769)
* TERPIN-0422-EIMRES outscale_group (TPD-1470)
* TERPIN-0322-FCURES outscale_vpn_gateway (TPD-774)
* TERPIN-0428-EIMRES outscale_policy_group_link (TPD-1003)
* TERPIN-0361-LBURES outscale_load_balancer (TPD-817)
* TERPIN-0279-FCURES outscale_snapshot_attributes (TPD-885)
* TERPIN-0148-FCURES outscale_image_register (TPD-759)
* TERPIN-0332-FCURES outscale_vpn_gateway_route_propagation (TPD-784)



0.1.0RC2 (March 29, 2019)
========================= 

FEATURES:
---------

* New Data Source: outscale_vm
* New Data Source: outscale_vms
* New Data Source: outscale_vm_attributes
* New Data Source: outscale_vm_state
* New Data Source: outscale_public_ip
* New Data Source: outscale_public_ips
* New Data Source: outscale_volume
* New Data Source: outscale_volumes
* New Data Source: outscale_internet_service
* New Data Source: outscale_internet_services
* New Data Source: outscale_nat_service
* New Data Source: outscale_nat_services
* New Data Source: outscale_subnet
* New Data Source: outscale_subnets
* New Data Source: outscale_route_table
* New Data Source: outscale_route_tables
* New Data Source: outscale_route
* New Data Source: outscale_routes
* New Data Source: outscale_net
* New Data Source: outscale_nets
* New Data Source: outscale_net_attributes
* New Data Source: outscale_net_peering
* New Data Source: outscale_net_peerings
* New Data Source: outscale_image
* New Data Source: outscale_images
* New Data Source: outscale_image_launch_permission
* New Data Source: outscale_keypair
* New Data Source: outscale_keypairs
* New Data Source: outscale_security_group
* New Data Source: outscale_security_groups
* New Data Source: outscale_security_group_rule
* New Data Source: outscale_tag
* New Data Source: outscale_tags
* New Data Source: outscale_nic
* New Data Source: outscale_nics
* New Data Source: outscale_nic_private_ip
* New Data Source: outscale_nic_private_ips
* New Data Source: outscale_snapshot
* New Data Source: outscale_snapshots
* New Data Source: outscale_snapshot_attributes

* New Resource: outscale_vm
* New Resource: outscale_vm_attributes
* New Resource: outscale_vm_state
* New Resource: outscale_public_ip
* New Resource: outscale_public_ip_link
* New Resource: outscale_volume
* New Resource: outscale_volumes_link
* New Resource: outscale_internet_service
* New Resource: outscale_internet_service_link
* New Resource: outscale_nat_service
* New Resource: outscale_subnet
* New Resource: outscale_route_table
* New Resource: outscale_route_table_link
* New Resource: outscale_route
* New Resource: outscale_net
* New Resource: outscale_net_attributes
* New Resource: outscale_net_peering
* New Resource: outscale_net_peering_acceptation
* New Resource: outscale_image
* New Resource: outscale_image_launch_permission
* New Resource: outscale_keypair
* New Resource: outscale_security_group
* New Resource: outscale_security_group_rule
* New Resource: outscale_tag
* New Resource: outscale_nic
* New Resource: outscale_nic_link
* New Resource: outscale_nic_private_ip
* New Resource: outscale_snapshot
* New Resource: outscale_snapshot_attributes

NOTES:
------

**WARNING:** When creating keypairs, the private key and fingerprint are stored in the Terraform state. For security reasons, it is strongly recommended to create keypairs using the API instead of the Terraform resource. For more information on how to create keypairs using the API, see our [official API documentation](https://docs.outscale.com/api#3ds-outscale-api-keypair).

It is recommended to set tags inside the resources rather than using outscale_tag.


0.1.0RC1 (Februry 23, 2018)
=========================== 

FEATURES:
---------

* New Data Source: outscale_vm
* New Data Source: outscale_vms


*  New Resource: outscale_vm


NOTES:
------

One resource/datasource delivery for initial Hashicorp review. 

