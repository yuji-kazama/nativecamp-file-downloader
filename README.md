# nativecamp-file-downloader

## Requirement
* Go 1.17 or higher

## Usage
 go run main.go [-p us/uk/ca] <NativeCamp_DailyNews_Page_URLs>"

## Troubleshooting
* If the program is not working, please try folloing steps:
  1. Open a NativeCamp_DailyNews_Page_URLs (e.g. https://nativecamp.net/textbook/page-detail/1/38974) in your chrome browser
  2. Open the developer tools in your chrome browser
  3. Select the audio button element in the developer tools
  4. Copy the audio button element's HTML code (Copy full XPath)
  5. Past the copied XPath in the xxXPath const value in the main.go file
