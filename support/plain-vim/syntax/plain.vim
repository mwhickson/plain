" Vim syntax file
" Language: PLAIN Programming Language

if exists("b:current_syntax")
  finish
endif

" PLAIN keywords are case-insensitive
syntax case ignore

" Comments
syntax match plainComment "'.*$" display

" Strings
syntax region plainString start=+"+ skip=+\\"+ end=+"+

" Numbers
syntax match plainNumber "\b\d\+\(\.\d\+\)\?\b"

" Constants
syntax keyword plainConstant True False None

" Keywords & Control Flow
syntax keyword plainStatement If Then Else ElseIf Select Case For To Step While Return End Break
syntax keyword plainLogicalOperator And Or
syntax keyword plainDeclaration Var Const Function Struct As

" Types
syntax keyword plainType Numeric String Boolean List Struct Function

" Standard Library Modules
syntax match plainModule "\<\(Console\|System\|List\|String\|Struct\|Math\|Random\|File\|Time\|Decimal\|Environment\)\>\ze\s*\."

" Operators and delimiters
syntax match plainOperator "==\|!=\|<=\|>=\|[=<>+*/-]"
syntax match plainLineContinuation "\\[ \t]*$"
syntax match plainPunctuation "[,.\[\]()]"

" Link to standard Vim highlight groups
hi def link plainComment    Comment
hi def link plainString     String
hi def link plainNumber     Number
hi def link plainConstant   Boolean
hi def link plainStatement  Statement
hi def link plainLogicalOperator Operator
hi def link plainDeclaration Keyword
hi def link plainType       Type
hi def link plainModule     Structure
hi def link plainOperator   Operator
hi def link plainLineContinuation Special
hi def link plainPunctuation Delimiter

let b:current_syntax = "plain"