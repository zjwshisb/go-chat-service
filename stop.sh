#!/bin/bash
pid=$(head -1 ./ws.pid)
kill -9 $pid