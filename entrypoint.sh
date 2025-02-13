#!/bin/bash
yugabyted_args=""
if host "$YUGABYTED_JOIN" | grep " $(hostname -i)"
then
    echo "Ignoring join: $YUGABYTED_JOIN == $(host $YUGABYTED_JOIN)"
else
    until postgres/bin/pg_isready -h "$YUGABYTED_JOIN" ; do sleep 1 ; done | uniq
    yugabyted_args="--join=$YUGABYTED_JOIN $yugabyted_args"
fi
sleep $( host $(hostname -i) | cut -c1 ) # try to not get all at the same time
yugabyted start --background=false $yugabyted_args --tserver_flags=$TSERVER_FLAGS
