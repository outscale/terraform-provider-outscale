# Outscale terraform examples

This folder contains a number of examples to show how to use Outscale provider for terraform.

# How to test examples

First, make sure you have [installed terraform](https://www.terraform.io/downloads.html) > 0.13.x.

Each folder is self-contained example.
You will need to setup your credentials through environement variables:
```bash
export TF_VAR_access_key_id="myaccesskey"
export TF_VAR_secret_key_id="mysecretkey"
export TF_VAR_region="eu-west-2"
```

If you want to write your credentials in terraform variables, just edit `terraform.tfvars` file.

Once your credentials are configured, you can go to any example folder and test them:
```bash
cd volume
terraform init
# Check plan before applying
terraform plan
# Create volume
terraform apply
# Re-run plan to check that infrastructure is up-to-date
terraform plan
# Clean ressources
terraform destroy
```
