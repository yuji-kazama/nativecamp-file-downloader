#!/bin/bash

PAGE_URLS="https://nativecamp.net/textbook/page-detail/2/20513 https://nativecamp.net/textbook/page-detail/2/20511 https://nativecamp.net/textbook/page-detail/2/20525"

rm -rf ./out
go run main.go $PAGE_URLS
open ./out
