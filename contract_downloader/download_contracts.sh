#!/bin/bash

MAPPING_DIR=${1}
OUTPUT_DIR=${2}
BATH_SIZE=${3}

ADDR_COUNT=$(ls ${MAPPING_DIR} | wc -l)

rm -rf ${OUTPUT_DIR}
rm last_line.txt

for (( i = 0; i < ${ADDR_COUNT}; i += ${BATH_SIZE} )); do
    BATH_DIR="${OUTPUT_DIR}/$i"
    mkdir -p "${BATH_DIR}/fuzzer/config"

    # copy seeds 
    cp -r ../examples/delegatecall_dangerous/fuzzer/config/*.json "${BATH_DIR}/fuzzer/config"
    node downloader-main.js ${MAPPING_DIR} ${BATH_DIR} $BATH_SIZE

    COUNT_CHECK_ABIS=$(ls "${BATH_DIR}/verified_contract_abis" | wc -l)
    COUNT_CHECK_CONTRACTS=$(ls "${BATH_DIR}/fuzzer/config/contracts.list" | wc -l)
    echo "Finished Bath ${i}/${ADDR_COUNT} - ABIs: ${COUNT_CHECK_ABIS}, Contracts.list: ${COUNT_CHECK_CONTRACTS}"
done
