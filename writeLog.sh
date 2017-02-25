#!/bin/bash

mkdir -p ./logs
IFS='
'
for ll in `docker ps | grep antblockchain`;
do
        ee=${ll/%\ */}
        docker logs $ee >& ./logs/$ee.log
done 