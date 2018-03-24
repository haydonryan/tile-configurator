#!/bin/bash -e

if [ "$#" -ne 3 ]; then
  echo "Need three parameters:"
  echo "1: Ops Man URL (without https)"
  echo "2: username"
  echo "3: password"
  exit 1
fi

TARGET=$1
USERNAME=$2
PASSWORD=$3

# Get the list of GUIDs for the products staged and save them into a variable
TILELIST=$(om-linux -t $TARGET -u $USERNAME -p $PASSWORD -k curl -p /api/v0/staged/products --silent | grep guid | awk -F '"' '{print $4}')


# For each guid, call OM and get the manifest
while read -r line
do
 echo "getting properties for $line and saving it to $line.json"
 om-linux -t $TARGET -u $USERNAME -p $PASSWORD -k curl -p /api/v0/staged/products/$line/properties --silent > $line.json
done <<<"$TILELIST"


