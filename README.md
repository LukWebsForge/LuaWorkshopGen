# Workshop Lua

Generates a lua file for your gmod server which forces every client to download all required mods.

I've created this to automate my TTT server. 
If you just looking for an fast online version, try the 
[Workshop Collection Generator](https://csite.io/tools/gmod-universal-workshop). 
After you've got the generated file, place it into `lua\autorun\server\workshop.lua` and voila.

### Arguments

1. **api-key** - Your Steam api key. Get it [here](https://steamcommunity.com/dev/apikey).
2. **collection** - The id of your collections. It's the number at the end of the url.
3. **output** - Where the file should be stored. E.g. `lua\autorun\server\workshop.lua`

### Build

Download and install [Go](https://golang.org/), then run `go build .`