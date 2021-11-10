resource "outscale_volume" "outscale_volume_snap" {
    subregion_name   = "${var.region}a"
    size             = 5
}
resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume_snap.volume_id
}
resource "outscale_snapshot_export_task" "outscale_snapshot_export_task" {
    snapshot_id                     = outscale_snapshot.outscale_snapshot.snapshot_id
    osu_export {
         disk_image_format = "raw"
         osu_bucket        = var.osu_bucket_name 
         osu_prefix        = "prefix-194"
         osu_api_key {
           api_key_id      = var.access_key_id
           secret_key      = var.secret_key_id
          }
     }
  tags {
      key = "test"
      value = "test"
    }
}

data "outscale_snapshot_export_task" "data-1" {
   filter {
        name   = "task_ids"
        values = [outscale_snapshot_export_task.outscale_snapshot_export_task.task_id]
    }
}
