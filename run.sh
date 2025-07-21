#!/bin/bash

PAGE_URLS="
https://nativecamp.net/textbook/page-detail/1/40468
https://nativecamp.net/textbook/page-detail/1/40449
https://nativecamp.net/textbook/page-detail/1/40539
https://nativecamp.net/textbook/page-detail/1/40581
"

rm -rf ./out
go run cmd/ncfiledownloader/main.go $PAGE_URLS
open ./out
