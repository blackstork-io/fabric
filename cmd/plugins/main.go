package main

import "github.com/blackstork-io/fabric/plugins/content/text"

func main() {
	ServePlugins(&text.Plugin{})
}
