#!/bin/sh
filename=server-$(date +%Y%m%d)_$(date +%H%M%S).log
nohup ./main --port $1 --host 0.0.0.0  >./$filename 2>&1 &
