#!/bin/bash
set -e 

HOST=localhost:8080

function ctrl_c() {
    pkill -P $$ || true
    exit 0
}

./server -listen $HOST > /dev/null &

trap ctrl_c INT ERR

for i in {1..10}; do
    sleep 5
    echo -e "Students ${i}\n$(curl -o /dev/stdout -s $HOST/students)\n"
    echo -e "Exams ${i}\n$(curl -o /dev/stdout -s $HOST/exams)\n"
done

pkill -P $$ || true
