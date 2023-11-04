#!/bin/bash

PAGE_URLS="
 https://nativecamp.net/textbook/page-detail/2/20468
 https://nativecamp.net/textbook/page-detail/2/20433
 https://nativecamp.net/textbook/page-detail/2/20493
 https://nativecamp.net/textbook/page-detail/2/20553
"
rm -rf ./out
go run main.go $PAGE_URLS
# ./ncfiledownloader $PAGE_URLS

open ./out
