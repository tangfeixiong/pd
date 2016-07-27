#!/bin/bash

# pd-server -addr 127.0.0.1:1234 --etcd 127.0.0.1:2379 --cluster-id 1 --root /pd

pd-server $@
