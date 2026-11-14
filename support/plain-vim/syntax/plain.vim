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
syntax keyword plainStatement If Then Else ElseIf Select Case For To Step In Return End
syntax keyword plainDeclaration Var Const Function Struct As

" Types
syntax keyword plainType Numeric String Boolean List Struct

" Standard Library Modules
syntax keyword plainModule Console System

" Link to standard Vim highlight groups
hi def link plainComment    Comment
hi def link plainString     String
hi def link plainNumber     Number
hi def link plainConstant   Boolean
hi def link plainStatement  Statement
hi def link plainDeclaration Keyword
hi def link plainType       Type
hi def link plainModule     Structure

let b:current_syntax = "plain"