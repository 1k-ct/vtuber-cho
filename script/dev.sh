#!/bin/sh

set -eu
# read NAME

# echo  -e ok \n $NAME

CMDNAME=`basename $0`

while getopts t:l: OPT
do
  case $OPT in
    "t" ) FLG_T="TRUE" ; VALUE_T="$OPTARG" ;;
    "l" ) FLG_L="TRUE" ; VALUE_L="$OPTARG" ;;
      * ) echo "Usage: $CMDNAME [-t target-file] [-l load-file]" 1>&2
          exit 1 ;;
  esac
done

source ../.env

echo ${KEY} $VALUE_L ec2-user@000.000:/home/ec2-user/$VALUE_T 

read -p "ok? (y/N): " yn
if [[ $yn = [yY] ]]; then
    scp -i ${KEY} -C -r $VALUE_L ${EC2}:/home/ec2-user/$VALUE_T
    echo Done
else
    echo cancel
    exit 1
fi