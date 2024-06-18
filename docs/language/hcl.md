---
title: HCL expressions
type: docs
weight: 90
---

# HCL expressions

Fabric configuration language supports a set of native [HCL](https://github.com/hashicorp/hcl) (HashiCorp Configuration Language)
expressions.

See a snippet below for the examples:

```hcl
document "example" {

  vars {

    arithmetic = "1 + 2 = ${1 + 2}"
    # "arithmetic": "1 + 2 = 3",

    logic = "true and false is ${true && false}"
    # "logic": "true and false is false",

    conditionals = "2 is ${2 % 2 == 0 ? "even" : "odd"}"
    # "conditionals": "2 is even"

    # technically, this is a tuple
    loop_over_list = [ for el in [1, 2, 3]: el * 2 ]
    # "loop_over_list": [
    #     2,
    #     4,
    #     6
    # ],

    loop_over_tuple = [ for el in [1, "two", 3]: "value is ${el}" ]
    # "loop_over_tuple": [
    #     "value is 1",
    #     "value is two",
    #     "value is 3"
    # ],

    # technically, this is an object
    loop_over_map = [ for k, v in {"a": 1, "b": 2, "c": 3}: "key ${k}: value ${v}" ]
    # "loop_over_map": [
    #     "key a: value 1",
    #     "key b: value 2",
    #     "key c: value 3"
    # ],

    loop_over_object = [ for k, v in {"a": 1, "b": "two", "c": 3}: "key ${k}: value ${v}" ]
    # "loop_over_object": [
    #     "key a: value 1",
    #     "key b: value two",
    #     "key c: value 3"
    # ],

    loop_creating_object = { for v in [1, 2, 3]: "${v}" => v * 2 }
    # "loop_creating_object": {
    #     "1": 2,
    #     "2": 4,
    #     "3": 6
    # },

    loop_with_filter = { for v in [1, 2, 3, 4]: "${v}" => v * 2 if v % 2 == 0 }
    # "loop_with_filter": {
    #     "2": 4,
    #     "4": 8
    # },

    loop_with_grouping = { for v in [1, 2, 3, 4]: (v%2 == 0 ? "evens" : "odds") => v... }
    # "loop_with_grouping": {
    #     "evens": [2, 4],
    #     "odds": [1, 3]
    # },

    splat_expression = ([
      {
        id: 1,
        name: "foo",
      },
      {
        id: 2,
        name: "bar",
      },
      {
        id: 3,
        name: "baz",
      },
    ])[*].name
    # "splat_expression": [
    #     "foo",
    #     "bar",
    #     "baz"
    # ],

    template_directives = <<EOT
      %{~ for v in [1, 2, 3, 4] }
        %{~ if v % 2 == 0 ~}
          ${v} is even
        %{ else ~}
          ${v*2} is doubled ${v}
        %{ endif ~}
      %{ endfor ~}
    EOT
    # "template_directives": "2 is doubled 1\n2 is even\n6 is doubled 3\n4 is even\n"
  }

  content text {
    value = <<-EOT
      arithmetic: {{ .vars.arithmetic }}

      logic: {{ .vars.logic }}

      conditionals: {{ .vars.conditionals }}

      loop_over_list: {{ .vars.loop_over_list | toJson }}

      loop_over_tuple: {{ .vars.loop_over_tuple | toJson }}

      loop_over_map: {{ .vars.loop_over_map | toJson }}

      loop_over_object: {{ .vars.loop_over_object | toJson }}

      loop_creating_object: {{ .vars.loop_creating_object | toJson }}

      loop_with_filter: {{ .vars.loop_with_filter | toJson }}

      loop_with_grouping: {{ .vars.loop_with_grouping | toJson }}

      splat_expression: {{ .vars.splat_expression | toJson }}

      template_directives: {{ .vars.template_directives }}
    EOT
  }
}
```

renders into:

```text
arithmetic: 1 + 2 = 3

logic: true and false is false

conditionals: 2 is even

loop_over_list: [2,4,6]

loop_over_tuple: ["value is 1","value is two","value is 3"]

loop_over_map: ["key a: value 1","key b: value 2","key c: value 3"]

loop_over_object: ["key a: value 1","key b: value two","key c: value 3"]

loop_creating_object: {"1":2,"2":4,"3":6}

loop_with_filter: {"2":4,"4":8}

loop_with_grouping: {"evens":[2,4],"odds":[1,3]}

splat_expression: ["foo","bar","baz"]

template_directives:
          2 is doubled 1

          2 is even

          6 is doubled 3

          4 is even
```
