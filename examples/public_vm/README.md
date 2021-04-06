This more complex example shows how to create a public VM with:
- client-side generated keypair
- security group and rules
- public IP
- resized bootdisk volume
- additional volume attached

It will also generate a `connect.sh` script to ease your connection.
Once `terraform apply` done, VM will need some time to boot before beeing able to connect.
