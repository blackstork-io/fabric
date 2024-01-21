package parser

// func parseEmbededConfig(block *hclsyntax.Block, parent FabricBlock) (cfg *Config, diags diagnostics.Diag) {
// 	diags.Append(validateLabelsLength(block, 0, ""))

// 	cfg = &Config{
// 		block: block,
// 	}

// 	// TODO: get plugin config spec from the plugin and parse the config body according to it
// 	// call same func as parseConfig
// 	return
// }

// func parseRef(block *hclsyntax.Block, parent FabricBlock) (diags diagnostics.Diag) {
// 	diags.Append(validateLabelsLength(block, 0, ""))
// 	// diags.Append(validateNesting(block, parent, BlockKindDocument, BlockKindSection))

// 	if diags.HasErrors() {
// 		return
// 	}

// 	return
// }
