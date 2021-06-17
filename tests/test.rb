require 'yard'
require 'json'

c = "An overview for the first overload.\n" +
"@raise SomeError this is some error\n" +
"@param param1 The first parameter.\n" +
"@param The second parameter.\n" +
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

c = "@option foo [String] bar (nil) baz"

p = YARD::DocstringParser.new()
p.parse(c)

puts "-----"
pt = p.tags[0]
require 'pry'; binding.pry

puts "tag_name: #{pt.tag_name}"
puts "name: #{pt.name}"
puts "text: #{pt.text}"
puts "types: #{pt.types}"





# 20: p = YARD::DocstringParser.new()
# 21: p.parse(c)
# 22:
# 23: require 'pry'; binding.pry
# 24:
# => 25: puts p.tags[3].name
# 26: puts p.tags[3].text
# 27: puts p.tags[3].types
# 28: #puts p.tags[7].desc

# [1] pry(main)> xx = p.tags[3]
# => #<YARD::Tags::OptionTag:0x0000562e3b198740
# @name="[param2]",
# @pair=#<YARD::Tags::DefaultTag:0x0000562e3b1987b8 @defaults=nil, @name=":option", @tag_name="option", @text="an option", @types=nil>,
# @tag_name="option",
# @text=nil,
# @types=nil>
# [2] pry(main)> ls xx
# YARD::Tags::Tag#methods: explain_types  name  name=  object  object=  tag_name  tag_name=  text  text=  type  types  types=
# YARD::Tags::OptionTag#methods: pair  pair=
# instance variables: @name  @pair  @tag_name  @text  @types
# [3] pry(main)> xx.name
# => "[param2]"
# [4] pry(main)> xx.types
# => nil
# [5] pry(main)>

