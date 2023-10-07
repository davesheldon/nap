#!/bin/bash
set -e

GOOS=windows GOARCH=amd64 go build -o ./dist/windows/amd64/nap.exe
perl -e '
  use strict;
  use warnings;
  use autodie;
  use IO::Compress::Zip qw(:all);
  zip [
    "./dist/windows/amd64/nap.exe"
  ] => "dist/windows/amd64/nap.zip",
       FilterName => sub { s[^src/][] },
       Zip64 => 0,
  or die "Zip failed: $ZipError\n";
'
rm -rf ./dist/windows/amd64/nap.exe

GOOS=darwin GOARCH=amd64 go build -o ./dist/macos/amd64/nap
tar -czvf ./dist/macos/amd64/nap.tar.gz ./dist/macos/amd64/nap
rm -rf ./dist/macos/amd64/nap
GOOS=darwin GOARCH=arm64 go build -o ./dist/macos/arm64/nap
tar -czvf ./dist/macos/arm64/nap.tar.gz ./dist/macos/arm64/nap
rm -rf ./dist/macos/arm64/nap
GOOS=linux GOARCH=amd64 go build -o ./dist/linux/amd64/nap
tar -czvf ./dist/linux/amd64/nap.tar.gz ./dist/linux/amd64/nap
rm -rf ./dist/linux/amd64/nap