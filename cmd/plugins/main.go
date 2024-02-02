package main

import "github.com/blackstork-io/fabric/plugins/content/text"

func main() {
	// call like ServePlugins(&text.Plugin{})
	ServePlugins(&text.Plugin{})
}
