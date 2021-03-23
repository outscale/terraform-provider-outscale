This more complex example shows how to setup a simple http public load balancer with public VMs as backends.

You can configure the number of backend VM in `terraform.tfvars` as well as other parameters.

Once `terraform apply` done, terraform outputs the load balancer URL. You should be able to see on which backend VM you ended.

This example also generates a `connect_*.sh` script to ease your connection to different VM for further testing.
