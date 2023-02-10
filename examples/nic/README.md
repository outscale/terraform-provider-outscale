This more complex example shows how to create a VM connected to multiple subnets (VPCs) with nics.

- net, subnet, internet_service, nic
- routes table and default route
- client-side generated keypair
- security group and rules
- public IP

It will also generate a `connect.sh` script to ease your connection.
Once `terraform apply` done, VM will need some time to boot before beeing able to connect.
