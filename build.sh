#!/bin/bash
set -e

GOOS=windows GOARCH=amd64 go build -o ./dist/nap.exe
perl -e '
  use strict;
  use warnings;
  use autodie;
  use IO::Compress::Zip qw(:all);
  zip [
    "dist/nap.exe"
  ] => "dist/windows-amd64.zip",
       FilterName => sub { s[^dist/][] },
       Zip64 => 0,
  or die "Zip failed: $ZipError\n";
'
rm -rf ./dist/nap.exe

GOOS=darwin GOARCH=amd64 go build -o ./dist/nap
tar -czvf ./dist/macos-amd64.tar.gz ./dist/nap
rm -rf ./dist/nap
GOOS=darwin GOARCH=arm64 go build -o ./dist/nap
tar -czvf ./dist/macos-arm64.tar.gz ./dist/nap
rm -rf ./dist/nap
GOOS=linux GOARCH=amd64 go build -o ./dist/nap
tar -czvf ./dist/linux-amd64.tar.gz ./dist/nap
rm -rf ./dist/nap