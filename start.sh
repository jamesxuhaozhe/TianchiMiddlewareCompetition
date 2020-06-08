#!/bin/bash

if   [  $SERVER_PORT  ];
then
   /usr/local/src/tailSamp -p $SERVER_PORT &
else
   /usr/local/src/tailSamp -p 8000 &
fi
tail -f /usr/local/src/start.sh

