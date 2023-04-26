#access_key_id       = "MyAccessKey"
#secret_key_id       = "MySecretKey"
#region              = "eu-west-2"
#image_id            = "OUTSCALE_IMAGEID" #using environment variable

volume_type     = "io1"
volume_iops     = 10000
volume_size_gib = 200

vm_type         = "tinav4.c1r1p2"
allowed_cidr    = ["0.0.0.0/0"]
