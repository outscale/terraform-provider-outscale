#access_key_id       = "MyAccessKey"
#secret_key_id       = "MySecretKey"
#region              = "eu-west-2"

volume_type     = "io1"
volume_iops     = 10000
volume_size_gib = 60
image_id        = "ami-3e5dc2a7" # Windows 2019
vm_type         = "tinav4.c8r8p2"
allowed_cidr    = ["0.0.0.0/0"]
