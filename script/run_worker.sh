#!/bin/sh
filename=worker-$2-$(date +%Y%m%d)_$(date +%H%M%S).log
nohup ./worker -host $1 -port $2 -serverAddr $3  >./$filename 2>&1 &

