resource "outscale_access_key" "my_access_key" {
  count = 3
  expiration_date = "2024-02-20T10:04:05.000Z"
}

resource "outscale_access_key" "my_access_key_2" {
  expiration_date = "2024-02-18"
}

data "outscale_access_keys" "my_access_keys"{
 filter {
   name ="access_key_ids"
   values = [outscale_access_key.my_access_key[0].access_key_id,outscale_access_key.my_access_key_2.access_key_id]
   }
}

