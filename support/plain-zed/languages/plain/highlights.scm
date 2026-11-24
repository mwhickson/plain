; Comments and literals
(comment) @comment
(string_literal) @string
(numeric_literal) @number
(boolean_literal) @boolean
(none_literal) @constant.builtin

; Control flow and declarations
["if_keyword" "then_keyword" "else_keyword" "elseif_keyword"
 "select_keyword" "case_keyword" "for_keyword" "to_keyword" "step_keyword"
 "while_keyword" "return_keyword" "end_keyword" "break_keyword"] @keyword.control
["var_keyword" "const_keyword" "function_keyword" "struct_keyword"
 "as_keyword"] @keyword.storage
["and_keyword" "or_keyword"] @keyword.operator

; Names and types
(identifier) @variable
(function_statement name: (identifier) @function)
(call_expression function: [(identifier) (property_access_expression)] @function.call)
(parameter name: (identifier) @variable.parameter)
(struct_statement name: (identifier) @type)
(struct_member name: (identifier) @property)
(property_access_expression property: (identifier) @property)
(variable_statement type: (type_name) @type)
(parameter type: (type_name) @type)
(struct_member type: (type_name) @type)
(function_statement return_type: (type_name) @type)
"numeric_keyword" @type
(property_access_expression object: (identifier) @support.class
  (#match? @support.class "^(Console|System|List|String|Struct|Math|Random|File|Time|Decimal|Environment)$"))

; Operators
[("+") ("-") ("*") ("/") ("==") ("!=") ("<") ("<=") (">") (">=") ("=")] @operator
