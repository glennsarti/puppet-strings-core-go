require 'yard'


c = "An overview for the first overload.\n" +
"@raise SomeError this is some error\n" +
"@param param1 The first parameter.\n" +
"@param param2 The second parameter.\n" +
"@option param2 [String] :option an option\n" +
"@option param2 [String] :option2 another option\n" +
"@param param3 The third parameter.\n" +
"@param param4 The fourth parameter.\n" +
"@enum param4 :one Option one.\n" +
"@enum param4 :two Option two.\n" +
"@return Returns nothing.\n" +
"@return [Undef]\n" +
"@example Calling the function foo\n" +
"  $result = func4x(1, 'foooo')\n" +
"\n"

p = YARD::DocstringParser.new()
p.parse(c)

puts p.text
