// +build go1.15

package puppet_strings_go

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

type RubyNodeHandler func(p *RubyParser, n *sitter.Node, ctx interface{}) bool

type RubyParser struct {
	handlers map[int]RubyNodeHandler
	defaultHandler RubyNodeHandler
	depth int
}

func NewRubyParser() RubyParser {
	return RubyParser{
		handlers: make(map[int]RubyNodeHandler),
		defaultHandler: nil,
	}
}

func (p *RubyParser) RegisterHandler(token int, handler RubyNodeHandler) {
	p.handlers[token] = handler
}

func (p *RubyParser) UnregisterHandler(token int) {
	p.handlers[token] = nil
}

func (p *RubyParser) RegisterDefaultHandler(handler RubyNodeHandler) {
	p.defaultHandler = handler
}

func (p *RubyParser) UnregisterDefaultHandler() {
	p.defaultHandler = nil
}

func (p *RubyParser) ParseAndVisit(content []byte, ctx interface{}) {
	parser := sitter.NewParser()
	parser.SetLanguage(ruby.GetLanguage())

	tree := parser.Parse(nil, content)
	p.visit(tree.RootNode(), ctx)
}

func (p *RubyParser) Visit(node *sitter.Node, ctx interface{}) bool {
	p.depth = 0
	return p.visit(node, ctx)
}

func (p *RubyParser) visit(node *sitter.Node, ctx interface{}) bool {
	// DEBUG
	it := ""
	for i := 0; i < p.depth; i++ {
		it = it + "--"
	}
	fmt.Printf("%s(%d) %s\n", it, node.Symbol(), node.Type())

	handled := false
	if handler, ok := p.handlers[int(node.Symbol())]; ok && handler != nil {
		p.depth = p.depth + 1
		handled = handler(p, node, ctx)
		p.depth = p.depth - 1
	}
	if !handled { handled = p.VisitDefaultHandler(node, ctx) }

	//TODO throw on unhandled?
	return handled
}

func (p *RubyParser) VisitDefaultHandler(node *sitter.Node, ctx interface{}) bool {
	if p.defaultHandler != nil {
		p.depth = p.depth + 1
		r := p.defaultHandler(p, node, ctx)
		p.depth = p.depth - 1
		return r
	}
	return false
}

func (p *RubyParser) VisitChildrenOfHandler(_p *RubyParser, node *sitter.Node, ctx interface{}) bool {
	for i := 0; i < int(node.ChildCount()); i++ {
		if !p.visit(node.Child(i), ctx) { return false }
	}
	return true
}

func (p *RubyParser) FirstChildWithSymbol(node *sitter.Node, symbol int) *sitter.Node {
	for i := 0; i < int(node.ChildCount()); i++ {
		if (int(node.Child(i).Symbol()) == symbol) { return node.Child(i) }
	}
	return nil
}


// Constants from https://github.com/smacker/go-tree-sitter/blob/master/ruby/parser.c
const (
	token_sym_identifier = 1
	token_anon_sym__END = 2
	token_sym_uninterpreted = 3
	token_anon_sym_BEGIN = 4
	token_anon_sym_LBRACE = 5
	token_anon_sym_RBRACE = 6
	token_anon_sym_END = 7 // Yes there is a difference between this and 2
	token_anon_sym_def = 8
	token_anon_sym_LPAREN = 9
	token_anon_sym_RPAREN = 10
	token_anon_sym_DOT = 11
	token_anon_sym_COLON_COLON = 12
	token_anon_sym_COMMA = 13
	token_anon_sym_PIPE = 14
	token_anon_sym_SEMI = 15
	token_anon_sym_STAR = 16
	token_anon_sym_STAR_STAR = 17
	token_anon_sym_AMP = 18
	token_anon_sym_COLON = 19
	token_anon_sym_EQ = 20
	token_anon_sym_class = 21
	token_anon_sym_LT = 22
	token_anon_sym_module = 23
	token_anon_sym_end = 24
	token_anon_sym_return = 25
	token_anon_sym_yield = 26
	token_anon_sym_break = 27
	token_anon_sym_next = 28
	token_anon_sym_redo = 29
	token_anon_sym_retry = 30
	token_anon_sym_if = 31
	token_anon_sym_unless = 32
	token_anon_sym_while = 33
	token_anon_sym_until = 34
	token_anon_sym_rescue = 35
	token_anon_sym_for = 36
	token_anon_sym_in = 37
	token_anon_sym_do = 38
	token_anon_sym_case = 39
	token_anon_sym_when = 40
	token_anon_sym_elsif = 41
	token_anon_sym_else = 42
	token_anon_sym_then = 43
	token_anon_sym_begin = 44
	token_anon_sym_ensure = 45
	token_anon_sym_EQ_GT = 46
	token_anon_sym_LBRACK = 47
	token_anon_sym_RBRACK = 48
	token_anon_sym_COLON_COLON2 = 49
	token_anon_sym_AMP_DOT = 50
	token_anon_sym_LPAREN2 = 51
	token_anon_sym_PLUS_EQ = 52
	token_anon_sym_DASH_EQ = 53
	token_anon_sym_STAR_EQ = 54
	token_anon_sym_STAR_STAR_EQ = 55
	token_anon_sym_SLASH_EQ = 56
	token_anon_sym_PIPE_PIPE_EQ = 57
	token_anon_sym_PIPE_EQ = 58
	token_anon_sym_AMP_AMP_EQ = 59
	token_anon_sym_AMP_EQ = 60
	token_anon_sym_PERCENT_EQ = 61
	token_anon_sym_GT_GT_EQ = 62
	token_anon_sym_LT_LT_EQ = 63
	token_anon_sym_CARET_EQ = 64
	token_anon_sym_QMARK = 65
	token_anon_sym_COLON2 = 66
	token_anon_sym_DOT_DOT = 67
	token_anon_sym_DOT_DOT_DOT = 68
	token_anon_sym_and = 69
	token_anon_sym_or = 70
	token_anon_sym_PIPE_PIPE = 71
	token_anon_sym_AMP_AMP = 72
	token_anon_sym_LT_LT = 73
	token_anon_sym_GT_GT = 74
	token_anon_sym_LT_EQ = 75
	token_anon_sym_GT = 76
	token_anon_sym_GT_EQ = 77
	token_anon_sym_CARET = 78
	token_anon_sym_PLUS = 79
	token_anon_sym_SLASH = 80
	token_anon_sym_PERCENT = 81
	token_anon_sym_EQ_EQ = 82
	token_anon_sym_BANG_EQ = 83
	token_anon_sym_EQ_EQ_EQ = 84
	token_anon_sym_LT_EQ_GT = 85
	token_anon_sym_EQ_TILDE = 86
	token_anon_sym_BANG_TILDE = 87
	token_anon_sym_defined_QMARK = 88
	token_anon_sym_not = 89
	token_anon_sym_BANG = 90
	token_anon_sym_TILDE = 91
	token_sym_constant = 92
	token_sym_instance_variable = 93
	token_sym_class_variable = 94
	token_sym_global_variable = 95
	token_anon_sym_DASH = 96
	token_anon_sym_PLUS_AT = 97
	token_anon_sym_DASH_AT = 98
	token_anon_sym_LBRACK_RBRACK = 99
	token_anon_sym_LBRACK_RBRACK_EQ = 100
	token_anon_sym_BQUOTE = 101
	token_anon_sym_undef = 102
	token_anon_sym_alias = 103
	token_sym_comment = 104
	token_sym_integer = 105
	token_sym_float = 106
	token_sym_complex = 107
	token_anon_sym_r = 108
	token_sym_super = 109
	token_anon_sym_true = 110
	token_anon_sym_TRUE = 111
	token_anon_sym_false = 112
	token_anon_sym_FALSE = 113
	token_sym_self = 114
	token_anon_sym_nil = 115
	token_anon_sym_NIL = 116
	token_sym_character = 117
	token_anon_sym_POUND_LBRACE = 118
	token_aux_sym_string_array_token1 = 119
	token_sym_escape_sequence = 120
	token_anon_sym_LBRACK2 = 121
	token_anon_sym_DASH_GT = 122
	token_sym_line_break = 123
	token_sym_simple_symbol = 124
	token_sym_string_start = 125
	token_sym_symbol_start = 126
	token_sym_subshell_start = 127
	token_sym_regex_start = 128
	token_sym_string_array_start = 129
	token_sym_symbol_array_start = 130
	token_sym_heredoc_body_start = 131
	token_sym_string_content = 132
	token_sym_heredoc_content = 133
	token_sym_string_end = 134
	token_sym_heredoc_end = 135
	token_sym_heredoc_beginning = 136
	token_sym_block_ampersand = 137
	token_sym_splat_star = 138
	token_sym_unary_minus = 139
	token_sym_binary_minus = 140
	token_sym_binary_star = 141
	token_sym_singleton_class_left_angle_left_langle = 142
	token_sym_identifier_hash_key = 143
	token_sym_program = 144
	token_sym_statements = 145
	token_sym_begin_block = 146
	token_sym_end_block = 147
	token_sym_statement = 148
	token_sym_method = 149
	token_sym_singleton_method = 150
	token_sym_method_rest = 151
	token_sym_parameters = 152
	token_sym_bare_parameters = 153
	token_sym_block_parameters = 154
	token_sym_formal_parameter = 155
	token_sym_simple_formal_parameter = 156
	token_sym_splat_parameter = 157
	token_sym_hash_splat_parameter = 158
	token_sym_block_parameter = 159
	token_sym_keyword_parameter = 160
	token_sym_optional_parameter = 161
	token_sym_class = 162
	token_sym_superclass = 163
	token_sym_singleton_class = 164
	token_sym_module = 165
	token_sym_return_command = 166
	token_sym_yield_command = 167
	token_sym_break_command = 168
	token_sym_next_command = 169
	token_sym_return = 170
	token_sym_yield = 171
	token_sym_break = 172
	token_sym_next = 173
	token_sym_redo = 174
	token_sym_retry = 175
	token_sym_if_modifier = 176
	token_sym_unless_modifier = 177
	token_sym_while_modifier = 178
	token_sym_until_modifier = 179
	token_sym_rescue_modifier = 180
	token_sym_while = 181
	token_sym_until = 182
	token_sym_for = 183
	token_sym_in = 184
	token_sym_do = 185
	token_sym_case = 186
	token_sym_when = 187
	token_sym_pattern = 188
	token_sym_if = 189
	token_sym_unless = 190
	token_sym_elsif = 191
	token_sym_else = 192
	token_sym_then = 193
	token_sym_begin = 194
	token_sym_ensure = 195
	token_sym_rescue = 196
	token_sym_exceptions = 197
	token_sym_exception_variable = 198
	token_sym_body_statement = 199
	token_sym_expression = 200
	token_sym_arg = 201
	token_sym_primary = 202
	token_sym_parenthesized_statements = 203
	token_sym_element_reference = 204
	token_sym_scope_resolution = 205
	token_sym_call = 206
	token_sym_command_call = 207
	token_sym_method_call = 208
	token_sym_command_argument_list = 209
	token_sym_argument_list = 210
	token_sym_argument_list_with_trailing_comma = 211
	token_sym_argument = 212
	token_sym_splat_argument = 213
	token_sym_hash_splat_argument = 214
	token_sym_block_argument = 215
	token_sym_do_block = 216
	token_sym_block = 217
	token_sym_assignment = 218
	token_sym_command_assignment = 219
	token_sym_operator_assignment = 220
	token_sym_command_operator_assignment = 221
	token_sym_conditional = 222
	token_sym_range = 223
	token_sym_binary = 224
	token_sym_command_binary = 225
	token_sym_unary = 226
	token_sym_parenthesized_unary = 227
	token_sym_unary_literal = 228
	token_sym_right_assignment_list = 229
	token_sym_left_assignment_list = 230
	token_sym_mlhs = 231
	token_sym_destructured_left_assignment = 232
	token_sym_rest_assignment = 233
	token_sym_lhs = 234
	token_sym_variable = 235
	token_sym_operator = 236
	token_sym_method_name = 237
	token_sym_setter = 238
	token_sym_undef = 239
	token_sym_alias = 240
	token_sym_rational = 241
	token_sym_true = 242
	token_sym_false = 243
	token_sym_nil = 244
	token_sym_chained_string = 245
	token_sym_interpolation = 246
	token_sym_string = 247
	token_sym_subshell = 248
	token_sym_string_array = 249
	token_sym_symbol_array = 250
	token_sym_symbol = 251
	token_sym_regex = 252
	token_sym_heredoc_body = 253
	token_aux_sym_literal_contents = 254
	token_sym_array = 255
	token_sym_hash = 256
	token_sym_pair = 257
	token_sym_lambda = 258
	token_sym_empty_statement = 259
	token_sym_terminator = 260
	token_aux_sym_statements_repeat1 = 261
	token_aux_sym_parameters_repeat1 = 262
	token_aux_sym_block_parameters_repeat1 = 263
	token_aux_sym_case_repeat1 = 264
	token_aux_sym_case_repeat2 = 265
	token_aux_sym_when_repeat1 = 266
	token_aux_sym_exceptions_repeat1 = 267
	token_aux_sym_body_statement_repeat1 = 268
	token_aux_sym_command_argument_list_repeat1 = 269
	token_aux_sym_mlhs_repeat1 = 270
	token_aux_sym_undef_repeat1 = 271
	token_aux_sym_chained_string_repeat1 = 272
	token_aux_sym_string_array_repeat1 = 273
	token_aux_sym_symbol_array_repeat1 = 274
	token_aux_sym_heredoc_body_repeat1 = 275
	token_aux_sym_hash_repeat1 = 276
	token_anon_alias_sym_DQUOTE = 277
	token_alias_sym_bare_string = 278
	token_alias_sym_bare_symbol = 279
	token_alias_sym_destructured_parameter = 280
	token_alias_sym_lambda_parameters = 281
	token_alias_sym_method_parameters = 282
)

