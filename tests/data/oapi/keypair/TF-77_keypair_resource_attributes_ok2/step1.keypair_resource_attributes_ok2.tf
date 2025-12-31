resource "outscale_keypair" "outscale_keypair" {
    keypair_name = "test-keypair-${random_string.suffix[0].result}"
    public_key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCow1EsH10/Xs1/8ZWiOttAaIJnSC2IlKxWTZftyotGU7CwNUtsw9g6MjTLYGTbOeSilbhAxRyuvrm2h4AwLo9Me4Oc6UhAzQwhANHYGMP0YA7t0R6qJ+d4/sPojnwcI3cT6ysIw0y7GwDeWHVq88AazdXXqK/JhNAE+pbdUV4eESj2z/jgDXaZepdCphd0cBA0vgz4+m2ONq41TWcAZl3/2ZzqGn8931a8QlmZCiUw3EPnB9O3GYJimk6pqV651PokWcgkYhuCDUjuGAbSGtPfHKTGcZc6DCVXBXapbdejZ3CqBGfvH8Zt+0XklaesmvwYc/UTgGK1+cEbTZ5n+H07 meriem.zouari.ext@mariem-zouari.local"
}
