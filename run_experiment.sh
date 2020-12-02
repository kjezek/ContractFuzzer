#!/bin/sh 
DIR=${PWD}
NEXT_TASK=$(curl $SERVER_HOST:9999/task)

if [ NEXT_TASK -eq "DONE" ]; then

    echo "ALL is done - no more experiments to run"
    sleep 360

else

    CONTRACT_DIR="/contracts/$NEXT_TASK"
    mkdir -p "/reporter/$NEXT_TASK/"

    export CONTRACT_DIR
    echo "Testing contracts from " $CONTRACT_DIR
    # KJ: RESEARCH - we do not need GETH at all
    #nohup ./geth_run.sh>>$CONTRACT_DIR/fuzzer/reporter/geth_run.log 2>&1 &
    #sleep 60
    cd $DIR
    nohup ./tester_run.sh>>"/reporter/$NEXT_TASK/tester_run.log" 2>&1 &
    sleep 10
    cd $DIR
    ./fuzzer_run.sh>>"/reporter/$NEXT_TASK/fuzzer_run.log" 2>&1 
    echo "Test finished!"
    echo "v_v..."
    echo "Please go to /reporter/$NEXT_TASK/ to see the results."

    # Send results back to the server

fi
