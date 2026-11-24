const kw = (word) => alias(token(prec(1, new RegExp(word.split('').map(c => c).join(''), 'i'))), `${word.toLowerCase()}_keyword`);

module.exports = grammar({
  name: 'plain',

  extras: $ => [/[ \t\r]+/, /\\[ \t]*\r?\n/, $.comment],

  word: $ => $.identifier,

  conflicts: $ => [[$.return_statement]],

  rules: {
    source_file: $ => repeat(choice($._statement, $._newline)),

    _statement: $ => choice(
      $.variable_statement, $.assignment_statement, $.expression_statement,
      $.if_statement, $.select_statement, $.for_statement, $.while_statement,
      $.function_statement, $.struct_statement, $.return_statement, $.break_statement
    ),

    _newline: $ => /\n+/,
    comment: $ => token(seq("'", /[^\n]*/)),
    identifier: $ => /[A-Za-z_][A-Za-z0-9_]*/,
    type_name: $ => choice($.identifier, kw('Numeric')),

    variable_statement: $ => choice(
      seq(kw('Var'), field('name', $.identifier), kw('As'), field('type', $.type_name), optional(seq('=', field('value', $._expression)))),
      seq(kw('Var'), field('name', $.identifier), '=', field('value', $._expression)),
      seq(kw('Const'), field('name', $.identifier), optional(seq(kw('As'), field('type', $.type_name))), '=', field('value', $._expression))
    ),
    assignment_statement: $ => seq(field('target', $._assignable), '=', field('value', $._expression)),
    expression_statement: $ => $._expression,
    _assignable: $ => choice($.identifier, $.property_access_expression, $.index_expression),

    function_statement: $ => seq(kw('Function'), optional(seq('(', field('receiver', $.parameter), ')')), field('name', $.identifier), '(', optional($.parameter_list), ')', optional(seq(kw('As'), field('return_type', $.type_name))), $._newline, repeat(choice($._statement, $._newline)), kw('End'), kw('Function')),
    parameter_list: $ => seq($.parameter, repeat(seq(',', $.parameter))),
    parameter: $ => seq(field('name', $.identifier), kw('As'), field('type', $.type_name)),

    struct_statement: $ => seq(kw('Struct'), field('name', $.identifier), $._newline, repeat(choice($.struct_member, $._newline)), kw('End'), kw('Struct')),
    struct_member: $ => seq(field('name', $.identifier), kw('As'), field('type', $.type_name)),

    if_statement: $ => seq(kw('If'), $._condition_expression, kw('Then'), $._newline, repeat(choice($._statement, $._newline)), repeat(seq(kw('ElseIf'), $._condition_expression, kw('Then'), $._newline, repeat(choice($._statement, $._newline)))), optional(seq(kw('Else'), $._newline, repeat(choice($._statement, $._newline)))), kw('End'), kw('If')),
    for_statement: $ => seq(kw('For'), field('variable', $.identifier), kw('As'), kw('Numeric'), '=', $._expression, kw('To'), $._expression, optional(seq(kw('Step'), $._expression)), $._newline, repeat(choice($._statement, $._newline)), kw('End'), kw('For')),
    while_statement: $ => seq(kw('While'), $._expression, $._newline, repeat(choice($._statement, $._newline)), kw('End'), kw('While')),
    select_statement: $ => seq(kw('Select'), kw('Case'), $._expression, $._newline, repeat1($.case_clause), optional($.case_else_clause), kw('End'), kw('Select')),
    case_clause: $ => seq(kw('Case'), $.expression_list, $._newline, repeat(choice($._statement, $._newline))),
    case_else_clause: $ => seq(kw('Case'), kw('Else'), $._newline, repeat(choice($._statement, $._newline))),
    return_statement: $ => seq(kw('Return'), optional($._expression)),
    break_statement: $ => kw('Break'),

    _expression: $ => choice($.binary_expression, $.call_expression, $.property_access_expression, $.index_expression, $.identifier, $.numeric_literal, $.string_literal, $.boolean_literal, $.none_literal, $.list_literal, $.parenthesized_expression),
    binary_expression: $ => choice(
      ...['*','/'].map(op => prec.left(6, seq($._expression, op, $._expression))),
      ...['+','-'].map(op => prec.left(5, seq($._expression, op, $._expression))),
      ...['<','<=','>','>='].map(op => prec.left(4, seq($._expression, op, $._expression))),
      ...['==','!='].map(op => prec.left(3, seq($._expression, op, $._expression))),
      prec.left(2, seq($._expression, kw('And'), $._expression)),
      prec.left(1, seq($._expression, kw('Or'), $._expression))
    ),
    _condition_expression: $ => choice($._expression, prec.left(3, seq($._expression, '=', $._expression))),
    call_expression: $ => prec(8, seq(field('function', choice($.identifier, $.property_access_expression)), '(', optional($.expression_list), ')')),
    property_access_expression: $ => prec(8, seq(field('object', choice($.identifier, $.property_access_expression, $.call_expression, $.index_expression)), '.', field('property', $.identifier))),
    index_expression: $ => prec(8, seq(field('object', choice($.identifier, $.property_access_expression, $.call_expression, $.index_expression)), '[', $._expression, ']')),
    expression_list: $ => seq($._expression, repeat(seq(',', $._expression))),
    parenthesized_expression: $ => seq('(', $._expression, ')'),
    list_literal: $ => seq('[', optional($.list_elements), ']'),
    list_elements: $ => seq($._expression, repeat(seq(',', $._expression)), optional(',')),
    numeric_literal: $ => /[0-9]+(\.[0-9]+)?/,
    string_literal: $ => token(seq('"', /[^"\n]*/, '"')),
    boolean_literal: $ => choice(kw('True'), kw('False')),
    none_literal: $ => kw('None'),
  }
});
