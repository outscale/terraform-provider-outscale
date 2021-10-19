#!/bin/sh

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
docs_dir="${project_root}/docs"
output_dir="${project_root}/"

python3 -m venv "${docs_dir}/.venv"
. "${docs_dir}/.venv/bin/activate"

pip3 install -r "${docs_dir}/requirements.txt" 1>"${docs_dir}/generation_log.txt" 2>&1 \
&& python3 "${docs_dir}/generate_doc_terraform.py" \
            --provider_directory "${project_root}/outscale/" \
            --api "${docs_dir}/osc-api/outscale.yaml"  \
            --output_directory "$output_dir" \
            --template_directory "${docs_dir}/doc-terraform-template/"  1>"${docs_dir}/generation_log.txt" 2>&1
RES=$?
deactivate

if [ $RES -ne 0 ]; then
    echo "KO, see logs in ${docs_dir}/generation_log.txt"
else
    echo "OK"
fi

exit $RES