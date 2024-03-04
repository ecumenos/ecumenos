#!/bin/sh

function deps()
{
  go mod tidy
}

function clean()
{
  find . -name "*.go" -path "**/mocks/*" | while read file; do rm $$file; done;
}

function deleteIgnored()
{
  ignored=( "ResponseBuildOption" )
  for i in "${!ignored[@]}"
  do
    rm -rfv \$i
  done
}

clean
deps
mockery --all
deleteIgnored
deps
