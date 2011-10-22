#!/bin/bash

grep 'perf' /var/log/locserver/server.log | cut -f2-5 > inPerf.log
grep 'perf' /var/log/locserver/server.log | cut -f2-5 > outPerf.log
sed -i 's/sNotVisible/notVis/g' outPerf.log
