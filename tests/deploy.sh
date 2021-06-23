#!/bin/bash

pkg_path=$1
pkg_name=$(basename $1)
version_name=$(echo $pkg_name | sed 's/-py3.*//g;s/osc_qa_provider_oapi.*-//g' | tr '.' '_')
branch_name=$(echo $pkg_name | sed 's/-.*//g;s/osc_qa_provider_oapi//g;s/^_//g')

target_list="172.19.141.29"
target_user="jenkins"
target_key="$2"

venv_name="venv_"$version_name
if [ "$branch_name" != "" ]; then
    venv_name="venv_"$branch_name"_"$version_name
fi

for target in $target_list; do
    echo $target
    scp -i $target_key $pkg_path $target_user@$target:./$pkg_name
    ssh -i $target_key $target_user@$target << EOF
mkdir -p ~/install/osc_qa_provider_oapi
python3 -m venv ~/install/osc_qa_provider_oapi/$venv_name
source ~/install/osc_qa_provider_oapi/$venv_name/bin/activate
pip install ./$pkg_name
deactivate
if [ "$branch_name" != "" ]; then
    if [ -h "\$HOME/install/osc_qa_provider_oapi/venv_$branch_name" ]; then
        rm \$HOME/install/osc_qa_provider_oapi/venv_$branch_name
    fi
    ln -s \$HOME/install/osc_qa_provider_oapi/$venv_name \$HOME/install/osc_qa_provider_oapi/venv_$branch_name

    hist=5
    i=0
    for venv in \$(ls -rt \$HOME/install/osc_qa_provider_oapi | grep "venv_"$branch_name"_"); do
        if [ \$i -ge \$hist ]; then
            echo "TODO RM: \$venv"
        fi
        i=\$((i+1))
    done
fi
EOF
     ssh -i $target_key $target_user@$target rm ./$pkg_name
done
