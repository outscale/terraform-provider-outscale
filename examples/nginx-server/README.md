# Nginx server

This example shows how to provision a public Linux virtual machine and bootstrap it with `user_data` to run an Nginx web server.

It demonstrates how to:

- configure the OUTSCALE provider
- register an SSH keypair
- create a security group
- open SSH and HTTP access
- provision a virtual machine
- attach a public IP
- initialize the instance with cloud-init
- generate `user_data` from a Terraform template with `templatefile()`

After deployment, the instance serves a simple web page over HTTP.

## Requirements

- Terraform
- An OUTSCALE account
- A valid OUTSCALE image ID
- An SSH public key available locally

## Files

- `provider.tf`: Terraform and provider configuration
- `variables.tf`: input variables
- `locals.tf`: computed values and rendered cloud-init content
- `keypair.tf`: SSH keypair registration
- `network.tf`: security group and rules
- `vm.tf`: virtual machine resource
- `public_ip.tf`: public IP allocation and association
- `outputs.tf`: useful outputs
- `templates/cloud-init.yaml.tftpl`: cloud-init template used as `user_data`

## Usage

Initialize the working directory:

```sh
terraform init
```

Create a variable file from the example:

```sh
cp terraform.tfvars.example terraform.tfvars
```

Edit terraform.tfvars and set your own values, especially:
* access_key
* secret_key
* image_id
* public_key_path

Apply the configuration:

```sh
terraform apply
```

Retrieve the web server URL:

```sh
terraform output web_url
```

Retrieve the public IP:

```sh
terraform output public_ip
```

Open the URL in your browser. You should see the demo page served by Nginx.

## Connect with SSH

You can retrieve the suggested SSH command with:
```sh
terraform output ssh_command
```
> The private key is generated automatically and saved locally in the example directory.

## Cleanup

To remove all resources created by this example:
```sh
terraform destroy
```

## Notes

> This example generates an SSH keypair with the `tls_private_key` resource and saves the private key locally.
> 
> This makes the example easier to run, but the generated private key is stored in Terraform state. 
>
> This is acceptable for a demo or local test environment, but it is not recommended for production use.