#access_key_id       = "MyAccessKey"
#secret_key_id       = "MySecretKey"
#region              = "eu-west-2"

image_id        = "ami-cdfcddb7" # Ubuntu-20.04-2021.02.10-5 on eu-west-2
vm_type         = "tinav4.c1r1p2"
allowed_cidr    = ["0.0.0.0/0"]
net_ip_range    = "10.0.0.0/16"
subnet_ip_range = "10.0.0.0/24"
