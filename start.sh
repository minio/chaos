#!/bin/sh

# supervisord has child protection. On starting worker too from supervisor 
# its not possible stop Minio process again from it.
nohup sh -c /usr/local/bin/worker &


