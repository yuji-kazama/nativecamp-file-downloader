#!/bin/bash

PAGE_URLS="
https://nativecamp.net/textbook/page-detail/1/38330
https://nativecamp.net/textbook/page-detail/2/38530
"


rm -rf ./out
go run main.go $PAGE_URLS
# ./ncfiledownloader $PAGE_URLS

open ./out
