locals {
  # Prefix used to build resource names consistently.
  name_prefix = "${var.project_name}-${var.instance_name}"

  # Common tags applied to the VM.
  common_tags = {
    Name    = local.name_prefix
    Project = var.project_name
    Example = "nginx-server"
  }

  # Render the cloud-init template and inject Terraform values into it.
  user_data = base64encode(templatefile("${path.module}/templates/cloud-init.yaml.tftpl", {
    project_name  = var.project_name
    instance_name = var.instance_name
    region        = var.region
  }))
}