#!/bin/bash

PAGE_URLS="
https://nativecamp.net/textbook/page-detail/1/38855
https://nativecamp.net/textbook/page-detail/1/38928
https://nativecamp.net/textbook/page-detail/1/38931
https://nativecamp.net/textbook/page-detail/1/38974
"


rm -rf ./out
go run main.go $PAGE_URLS
# ./ncfiledownloader $PAGE_URLS

open ./out
