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
  "While"
  "Return"
  "End"
  "Break"
] @keyword.control

[
  "Var"
  "Const"
  "Function"
  "Struct"
  "As"
] @keyword.storage

[
  "And"
  "Or"
] @keyword.operator

; --- Types ---
[
  "Numeric"
  "String"
  "Boolean"
  "List"
  "Struct"
  "Function"
] @type

; --- Standard Library Modules ---
[
  "Console"
  "System"
  "List"
  "String"
  "Struct"
  "Math"
  "Random"
  "File"
  "Time"
  "Decimal"
  "Environment"
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