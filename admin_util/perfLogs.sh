#!/bin/bash

grep 'inPerf' /var/log/locserver/server.log | cut -f2-5 > inPerf.log
grep 'outPerf' /var/log/locserver/server.log | cut -f2-5 > outPerf.log
