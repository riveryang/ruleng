#!/bin/sh

socat -d TCP4-LISTEN:2375,fork UNIX-CONNECT:/var/run/docker.sock &
SOCAT_PID=$!

/ruleng --hostname=$HOSTNAME

if [ -f "/tmp/profile" ]; then
  source /tmp/profile
  cat /tmp/profile
fi

kill $SOCAT_PID

