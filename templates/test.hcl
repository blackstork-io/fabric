data plugin_b "data_plugin_b" {
    parameter_z = ["a", "b", "c", "d"]
}

content text "external_block" {
    text = "External block body"
}

document "test-document" {

    data plugin_a "data_plugin_a" {
        parameter_x = 1
        parameter_y = 2
    }

    data ref "data_plugin_b" {
        // This should be automatically resolved to the referenced block by HCL parser
        ref = data.plugin_b.data_plugin_b
    }

    content text _ {
        query = ".data.plugin_a.data_plugin_a"
        text = "The value is {{ .data.plugin_a.data_plugin_a.result }}"
    }

    content generic _ {

        content ref _ {
            // This should be automatically resolved to the referenced block by HCL parser
            ref = content.text.external_block
        }

        content table _ {
            // JQ query
            query = ".data.plugin_b.data_plugin_b.result | length"
            text = "The length of the list is {{ .query_result }}"
            columns = ["ColumnA", "ColumnB", "ColumnC"]
        }
    }
}