resource "outscale_access_key" "my_access_key_1"{
 }
resource "outscale_access_key" "my_access_key_2"{
 }
data "outscale_access_keys" "my_access_keys"{
filter {
 name ="access_key_ids"
 values = [outscale_access_key.my_access_key_1.access_key_id, outscale_access_key.my_access_key_2.access_key_id]
  }
}
