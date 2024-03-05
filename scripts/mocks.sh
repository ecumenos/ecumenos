#!/bin/sh

echo "starting mocks generation"
go mod tidy
mockery --all
echo "mocks were generared"

for ignored in "ResponseBuildOption.go"
  do
    for i in $(find * -type f); do 
      if [ -f "$i" ]; then
        if echo "$i" | grep -iq "$ignored"; then
          rm -rfv $i
        fi
      fi
    done
  done
echo "deleted ignored mocks"
go mod tidy
