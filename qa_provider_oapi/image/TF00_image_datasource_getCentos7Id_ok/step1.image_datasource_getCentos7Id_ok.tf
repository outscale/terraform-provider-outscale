data "outscale_images" "outscale_image" {
    filter {
        name = "name"
        values = ["Centos 7 (*"]
    }
}

output "images"   {
    value = zipmap(formatlist("*%s",data.outscale_images.outscale_image.images_set.*.name), data.outscale_images.outscale_image.images_set)
    #value = zipmap(formatlist("*%s",data.outscale_images.outscale_image.images_set), formatlist("*%s",data.outscale_images.outscale_image.images_set))
    #value = zipmap(formatlist("*%s",data.outscale_images.outscale_image.images_set), formatlist("*%s",data.outscale_images.outscale_image.images_set))
    #value = formatlist("*%s",data.outscale_images.outscale_image.images_set.*.image_id)
    #value = length(data.outscale_images.outscale_image.images_set)
}

#output "network_info" {
#  value = zipmap(openstack_networking_network_v2.network.*.name, openstack_networking_subnet_v2.subnet.*.cidr)
#}

#my_expression = zipmap(random_shuffle.x.result,formatlist("*%s",random_shuffle.x.result))
