#!/bin/sh 
DIR=${PWD}
NEXT_TASK=$(curl -s $SERVER_HOST:9999/task)

if [ "$NEXT_TASK" = "DONE" ]; then

    echo "ALL is done - no more experiments to run"
    sleep 360

else

    if [ "$NEXT_TASK" = "" ]; then
        echo "Master not up yet"
        exit 0
    fi

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