#!/bin/bash
pid=$(head -1 ./wss.pid)
kill -9 $pid