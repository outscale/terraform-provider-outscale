#!/bin/bash

for f in $(find provider_outscale_oapi/test/ -name '*.py' -type f); do
    cat=$(echo $f | cut -d '/' -f4 | sed 's/outscale_//g')
    test_name=$(echo $f | cut -d '/' -f5 | sed 's/.*outscale_//g;s/.py//g')
    tf_file_list=$(grep "ln -s" $f | sed 's/.*ln -s //g' | sed 's/ test.tf.*//g')
    echo "mkdir -p qa_provider_oapi/$cat/TF00_"$test_name
    mkdir -p "qa_provider_oapi/$cat/TF00_"$test_name
    origin="qa_provider_oapi/$cat/TF00_"$test_name"/origin.txt" 
    echo "" > $origin
    echo "python file: $f" >> $origin
    i=1
    if [ -z "$tf_file_list" ]; then
        echo "WARNING: No tf file" >> $origin
    fi
    for tf_file in $tf_file_list; do
        step_name=$(echo $tf_file | sed 's/.*config_outscale_//g')
        if [ ! -f "$tf_file" ]; then
            echo "WARNING: tf file not found '$tf_file'" >> $origin
            continue
        fi
        nb_lin=$(grep "outscale_lin" $tf_file | grep -c -v "#")
        if [ "$nb_lin" != "0" ]; then
            echo "WARNING: LIN found in tf file '$tf_file'" >> $origin
        fi
        cp $tf_file "qa_provider_oapi/$cat/TF00_"$test_name"/step"$i"."$step_name".tf"
        echo "tf file: $tf_file" >> $origin

        IFS=$'\n'
        check_list=$(grep "check.py" $f | grep -v "#runCmde" | sed 's/.*\.json //g;s/", .*//g')
        if [ -z "$check_list" ]; then
            echo "WARNING: No check file !!!" >> $origin
        fi
        for check in $check_list; do
            check_file=$(echo $check | cut -d ' ' -f1)
            res_name=$(echo $check | cut -d ' ' -f2)
            if [ ! -f "$check_file" ]; then
                echo "WARNING: check file not found '$check_file'" >> $origin
                continue
            fi
            cp $check_file "qa_provider_oapi/$cat/TF00_"$test_name"/step"$i"."$res_name".check"
            echo "check file: $check_file" >> $origin
        done
        unset IFS

        i=$((i+1))
    done

done


#for f in $(grep -r "var." ./qa_provider_oapi | cut -d ':' -f1 | sort -u); do
for f in $(find ./qa_provider_oapi -type f -name "*.tf"); do
    echo $f
    sed -i '' 's/"${//g;s/}"//g;s/var.region}a"/format("%s%s", var.region, "a")/g;s/var.region}b"/format("%s%s", var.region, "b")/g' $f
done

