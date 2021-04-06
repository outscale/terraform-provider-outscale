#access_key_id       = "MyAccessKey"
#secret_key_id       = "MySecretKey"
#region              = "eu-west-2"

volume_type     = "io1"
volume_iops     = 10000
volume_size_gib = 200
image_id        = "ami-cdfcddb7" # Ubuntu-20.04-2021.02.10-5 on eu-west-2
vm_type         = "tinav4.c1r1p2"
allowed_cidr    = ["0.0.0.0/0"]
