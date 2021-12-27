#!/bin/bash

export TF_LOG_PATH=./log
export TF_LOG=DEBUG
export TF_ACC=1

# List all tests
temp_file=$(mktemp ./terraform_test.XXXXX)
if [ $? -ne 0 ]; then
    echo "Error while creating the temporary file"
    exit 1
fi
go test -list ^.*$  github.com/terraform-providers/terraform-provider-outscale/outscale > "${temp_file}"
if [ $? -ne 0 ]; then
    echo "Error while generating the list of all the tests"
    exit 1
fi

while read -r test_name; do
    echo -n "[INFO] Executing the test ${test_name}: "

    if [ "$CHECK_RESOURCE_LEAK" = "1" ]; then
        # First create the snapshot
        frieza snapshot rm ${test_name} 2>/dev/null
        frieza snapshot new ${test_name} opensource
        if [ $? -ne 0 ]; then
            echo "Error while generating the snapshot before the test ${test_name}"
            exit 1
        fi
    fi
    # Launch the test
    go test -timeout 1h -run ${test_name} github.com/terraform-providers/terraform-provider-outscale/outscale -count=1 -v 1>./${test_name} 2>&1  
    if [ $? -ne 0 ]; then
        echo "KO"
        cat ./${test_name}
    fi

    if [ "$CHECK_RESOURCE_LEAK" = "1" ]; then
        # Compute the snapshot after
        frieza clean --plan ${test_name} > ./${test_name}_plan
        if [ $? -ne 0 ]; then
            echo "Error while generating the snapshot before the test ${test_name}"
            exit 1
        fi
        
        grep "No new object to delete in profile" ./${test_name}_plan 1>/dev/null
        RC=$?
        if [ $RC -eq 1 ]; then
            echo "KO, some resource are not deleted"
            cat ./${test_name}_plan
        elif [ $RC -eq 0 ]; then
            echo "OK"
        else
            echo "KO, internal error"
        fi
    else
        echo "OK"
    fi

done < "${temp_file}"