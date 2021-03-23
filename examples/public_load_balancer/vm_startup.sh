#!/bin/bash
node_number=$1
colors=("blue" "red" "green" "yellow" "black")
node_color=${colors[node_number]}
dest="/var/www/html/index.html"
apt-get update -y
apt-get install -y lighttpd
echo "<html><body style=\"background-color:${node_color};color:white;text-align:center;font-size:80px;\">Node ${node_number}</body></html>" > $dest
chmod a+r $dest
