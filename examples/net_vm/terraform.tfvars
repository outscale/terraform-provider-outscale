#access_key_id       = "MyAccessKey"
#secret_key_id       = "MySecretKey"
#region              = "eu-west-2"
#image_id            = "OUTSCALE_IMAGEID" #using environment variable

vm_type         = "tinav6.c1r1p2"
allowed_cidr    = ["0.0.0.0/0"]
net_ip_range    = "10.0.0.0/16"
subnet_ip_range = "10.0.0.0/24"
