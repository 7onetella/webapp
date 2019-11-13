#!/bin/sh

rm -f urls.txt
touch urls.txt

for i in {1..1000} 
do
  echo "/redis/set/k${i}/v${i}" >> urls.txt 
done

for i in {1..1000} 
do
  echo "/redis/get/k${i}" >> urls.txt 
done

