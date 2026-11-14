; --- Comments ---
(comment) @comment

; --- Strings & Literals ---
(string_literal) @string
(numeric_literal) @number
(boolean_literal) @boolean
"None" @constant.builtin

; --- Keywords & Control Flow ---
[
  "If"
  "Then"
  "Else"
  "ElseIf"
  "Select"
  "Case"
  "For"
  "To"
  "Step"
  "In"
  "Return"
  "End"
] @keyword.control

[
  "Var"
  "Const"
  "Function"
  "Struct"
  "As"
] @keyword.storage

; --- Types ---
[
  "Numeric"
  "String"
  "Boolean"
  "List"
  "Struct"
] @type

; --- Standard Library Modules ---
[
  "Console"
  "System"
] @support.class

; --- Identifiers & Functions ---
(identifier) @variable
(function_statement name: (identifier) @function)
(call_expression function: (identifier) @function.call)

; --- Operators ---
[
  "+"
  "-"
  "*"
  "/"
  "=="
  "!="
  "<"
  "<="
  ">"
  ">="
  "="
] @operator