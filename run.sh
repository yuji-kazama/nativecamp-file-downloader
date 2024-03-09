#!/bin/bash

PAGE_URLS="
https://nativecamp.net/textbook/page-detail/2/24193
https://nativecamp.net/textbook/page-detail/2/24163
https://nativecamp.net/textbook/page-detail/2/24177
https://nativecamp.net/textbook/page-detail/2/24172
https://nativecamp.net/textbook/page-detail/2/24189
"

rm -rf ./out
go run main.go $PAGE_URLS
# ./ncfiledownloader $PAGE_URLS

open ./out
