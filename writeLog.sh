#!/bin/bash

mkdir -p ./logs
rm -f logs/*.log
IFS='
'
for ll in `docker ps | grep antblockchain`;
do
        ee=${ll/%\ */}
        docker logs $ee >& ./logs/$ee.log
done
