# Declarative Network Topology Configuration Language with Lazy Instantiation

- Feature Name: `lazy-evaluated-configuration-language`
- Status: `Pending`

## Summary

[summary]: #summary

A one-paragraph explanation of the architectural decision, feature, or subsystem design.

## Motivation

[motivation]: #motivation

Any architectural change or new subsystem should focus on solving a concrete problem faced by users or developers of the system. This section should explain the problem in detail, including the necessary background and context.

It should also contain specific use cases or workloads where this design helps, and explain how it addresses them.

This section is one of the most important parts of any design document.

## Guide-level explanation

[guide-level-explanation]: #guide-level-explanation

Explain the proposal as if it were already implemented and you were teaching it to another user or developer of the system. That generally means:

- Introducing new named concepts.
- Explaining the design largely in terms of examples.
- Explaining how users or developers should _think_ about the feature, and how it should impact the way they use or develop the system. It should explain the impact as concretely as possible.

## Reference-level explanation

[reference-level-explanation]: #reference-level-explanation

This is the technical portion of the design document. Explain the design in sufficient detail that:

- Its interaction with other features and subsystems is clear.
- It is reasonably clear how the feature or subsystem would be implemented.
- Corner cases are dissected by example.

The section should return to the examples given in the previous section, and explain more fully how the detailed proposal makes those examples work.

### Lexical Structure and Syntax

#### Source Encoding and Indentation

A _source file_ is a UTF-8 encoded text representing a top-level _block body_. Within a line, tokens are separated by spaces (`0x20`).

_Blocks_ are enclosed in curly braces (`{` and `}`) and indented exclusively with tabs (`\t`, `0x09`). Spaces are not permitted for indentation.

Top-level _statements_ have no indentation. Each nested block increases the indentation level by one tab.

Non-empty lines inside a block must match that block's indentation level. The closing brace (`}`) must align with the enclosing block.

```hcl
// good: tab indentation matches block nesting depth
node "spine-1" {
 ⇥  profile = "router"
 ⇥  interface "lo0" {
 ⇥   ⇥  type = loopback
 ⇥  }
}

// bad: leading spaces instead of tabs
node "spine-1" {
....profile = "router"
}
```

#### Comments and Line Termination

```ebnf
Utf8Bytes  = valid utf-8 byte sequence without '\n' ;

Comment   = "//" Utf8Bytes ;
EndOfLine = [ Comment ] "\n" { [ Comment ] "\n" } ;
```

A _comment_ begins with `//` and extends to the end of the line. Multi-line block comments are not supported.

Blank lines and comment-only lines are ignored everywhere.

_Statements_ and _block headers_ end with an _end-of-line_. Semicolons are not used.

#### Identifiers and Keywords

```ebnf
DecDigit   = "0" .. "9" ;
LowerAlpha = "a" .. "z" ;

Ident = LowerAlpha { LowerAlpha | DecDigit | "_" } ;
```

_Identifiers_ consist of lowercase ASCII letters, decimal digits, and underscores. They must start with a letter (`lower_snake_case`).

The following _identifiers_ are _reserved keywords_. These keywords can only be used in their correct contexts. They cannot be used as _variable_, _function_, or _attribute_ names:

- `bind`
- `else`
- `false`
- `foreign`
- `include`
- `let`
- `true`
- `with`

The following identifiers are _weak keywords_. They have special meaning only in certain contexts:

- `path`

#### Literals and Interpolation

```ebnf
IntLit   = DecDigits ;
FloatLit = DecDigits "." DecDigits ;
BoolLit  = "true" | "false" ;

RateUnit     = "bps" | "Kbps" | "Mbps" | "Gbps" ;
DurationUnit = "ns" | "us" | "ms" | "s" | "m" ;
SizeUnit     = "B" | "KB" | "MB" | "GB" ;

RateLit     = DecDigits RateUnit ;
DurationLit = DecDigits DurationUnit ;
SizeLit     = DecDigits SizeUnit ;

Interp = "$" Ident
       | "$(" Expr ")" ;

StrLit = '"' { Utf8Str | Interp } '"' ;

EndpointLit = StrLit "." StrLit ;

IPv4Octet  = DecDigits | Interp ;
IPv4Addr   = IPv4Octet "." IPv4Octet "." IPv4Octet "." IPv4Octet ;
IPv4Prefix = IPv4Addr "/" ( IntLit | Interp ) ;

IPv6Hextet = HexDigits | Interp ;
IPv6Head   = IPv6Hextet { ":" IPv6Hextet } ;
IPv6Tail   = IPv6Hextet { ":" IPv6Hextet } ;
IPv6Addr   = IPv6Head [ "::" [ IPv6Tail ] ] | "::" [ IPv6Tail ] ;
IPv6Prefix = IPv6Addr "/" ( IntLit | Interp ) ;

Literal = IntLit
        | FloatLit
        | BoolLit
        | RateLit
        | DurationLit
        | SizeLit
        | StrLit
        | EndpointLit
        | IPv4Addr
        | IPv4Prefix
        | IPv6Addr
        | IPv6Prefix ;
```

_Unit suffixes_ for rate, duration, and size literals are case-sensitive and must immediately follow the digits without interleaving spaces.

Hexadecimal digits in _IPv6 literals_ must be lowercase (`0`–`9`, `a`–`f`).

An _IPv6 address_ must contain either 8 hextets, or a single `::` separator with 0 to 7 hextets in total.

An _interpolation_ (`$ident` or `$(expr)`) embeds a _variable_ or _expression_ into a _string_, _IPv4_, or _IPv6 literal_. The `$` prefix is only allowed inside these literals, where it distinguishes _variables_ and _expressions_ from plain text and hexadecimal digits (`a`–`f`).

```hcl
// good: lowercase hex and $-prefixed interpolation inside literals
ip   = fd00:0:ff:$s:$(r + 1)::1/64
name = "spine-$s"

// bad: uppercase hex digits in IPv6 literal
// ip = FD00:0:FF::1/64

// bad: '$' prefix used outside a literal
// let x = $s + 1
```

#### Expressions

```ebnf
PrimaryExpr = Ident
            | Literal
            | "(" Expr ")" ;
PostfixExpr = PrimaryExpr { "[" Expr "]" | "." Ident } ;
FnAppExpr   = Ident { PostfixExpr }
            | PostfixExpr ;
UnaryExpr   = [ "-" | "!" ] FnAppExpr ;
MulExpr     = UnaryExpr { ( "*" | "/" | "%" ) UnaryExpr } ;
AddExpr     = MulExpr { ( "+" | "-" ) MulExpr } ;
CmpExpr     = AddExpr [ ( "==" | "!=" | "<" | "<=" | ">" | ">=" ) AddExpr ] ;
AndExpr     = CmpExpr { "&&" CmpExpr } ;
OrExpr      = AndExpr { "||" AndExpr } ;
Expr        = OrExpr ;
```

Operator precedence follows the order of the grammar rules above, from highest (_primary_ and _postfix_ _expressions_) to lowest (_logical OR_).

Binary arithmetic (`*`, `/`, `%`, `+`, `-`) and logical (`&&`, `||`) operators are left-associative.

Comparison operators (`==`, `!=`, `<`, `<=`, `>`, `>=`) share a single precedence level and are non-associative.

_Function application_ is written in prefix juxtaposition form without parentheses around the argument list.

_Function application_ has higher precedence than unary operators (`-` and `!`). Thus, `- f x` means `-(f x)`, and negative function arguments must be enclosed in parentheses.

```hcl
// good: parenthesized negative argument and explicit comparison conjunction
let wrapped  = mod (-1) 8
let in_range = a < b && b < c

// bad: unparenthesized '-1' is parsed as subtraction '(mod) - (1)'
// let wrapped = mod -1 8

// bad: comparison operators cannot be chained
// let in_range = a < b < c
```

#### Statements and Blocks

```ebnf
BlockGuard  = "|" Expr ;
ClauseGuard = "|" ( "else" | Expr ) ;

LetDecl  = "let" Ident { Ident } [ ClauseGuard ] "=" Expr ;
AttrDecl = Ident [ ClauseGuard ] "=" Expr ;

Generator   = Ident "<-" Expr ".." Expr ;
ForeignDecl = Ident "foreign" ;

InlineWith = "[" Generator { "," Generator } "]" ;

WithDecl = "with" "{" EndOfLine { WithItem EndOfLine } "}" ;
WithItem = Generator
         | ForeignDecl
         | BlockGuard ;

BindDecl = "bind" "{" EndOfLine { BindItem EndOfLine } "}" ;
BindItem = Ident [ "=" Expr ] ;

IncludeDecl = "include" [ InlineWith ] [ BlockGuard ] "{" EndOfLine { IncludeItem EndOfLine } "}" ;
IncludeItem = "path" "=" '"' Utf8Str '"'
            | WithDecl
            | BindDecl ;

BlockDecl = Ident [ StrLit ] [ InlineWith ] [ BlockGuard ] "{" EndOfLine BlockBody "}" ;
BlockBody = { Statement EndOfLine } ;
Statement = WithDecl
          | IncludeDecl
          | BindDecl
          | LetDecl
          | AttrDecl
          | BlockDecl ;
```

A _block header_ must reside on a single line. The opening brace `{` must appear on the same line as the block header and be immediately followed by an _end-of-line_. The closing brace `}` must appear on its own line. Single-line blocks are not permitted.

Each _block_ and _include block_ may contain at most one _with block_ and at most one _bind block_.

_Inline generators_ (`[...]`) and a header _block guard_ (`| ...`) are shorthand for a _with block_. A block that uses header shorthand cannot also contain an inner _with block_.

Inside a _bind block_, a bare identifier `v` without an `= expr` is shorthand for `v = v`.

An _include block_ must contain exactly one `path` attribute.

When a _let binding_ is defined across multiple guarded clauses within a block, all clauses for that name must have the exact same parameter list.

```hcl
// Inline header shorthand:
link [h <- 1..3] | r < s {
	from = "leaf-$r"."uplink-$s-$h"
	to   = "spine-$s"."eth-$r-$h"
}

// Equivalent block form:
link {
	with {
		h <- 1..3
		| r < s
	}
	from = "leaf-$r"."uplink-$s-$h"
	to   = "spine-$s"."eth-$r-$h"
}
```

### Formal Grammar

```ebnf
(* --- Lexical terminals --- *)

DecDigit  = "0" .. "9" ;
DecDigits = DecDigit { DecDigit } ;

HexDigit  = "0" .. "9" | "a" .. "f" ;
HexDigits = HexDigit { HexDigit } ;

RateUnit     = "bps" | "Kbps" | "Mbps" | "Gbps" ;
DurationUnit = "ns" | "us" | "ms" | "s" | "m" ;
SizeUnit     = "B" | "KB" | "MB" | "GB" ;

LowerAlpha = "a" .. "z" ;
Ident = LowerAlpha { LowerAlpha | DecDigit | "_" } ;

Utf8Bytes  = valid utf-8 byte sequence without '\n' ;
Utf8Str    = valid utf-8 byte sequence without '"', '$', and '\n' ;

Comment   = "//" Utf8Bytes ;
EndOfLine = [ Comment ] "\n" { [ Comment ] "\n" } ;

(* --- Literals and typed interpolation --- *)

IntLit   = DecDigits ;
FloatLit = DecDigits "." DecDigits ;
BoolLit  = "true" | "false" ;

RateLit     = DecDigits RateUnit ;
DurationLit = DecDigits DurationUnit ;
SizeLit     = DecDigits SizeUnit ;

Interp = "$" Ident
       | "$(" Expr ")" ;

StrLit = '"' { Utf8Str | Interp } '"' ;

EndpointLit = StrLit "." StrLit ;

IPv4Octet  = DecDigits | Interp ;
IPv4Addr   = IPv4Octet "." IPv4Octet "." IPv4Octet "." IPv4Octet ;
IPv4Prefix = IPv4Addr "/" ( IntLit | Interp ) ;

IPv6Hextet = HexDigits | Interp ;
IPv6Head   = IPv6Hextet { ":" IPv6Hextet } ;
IPv6Tail   = IPv6Hextet { ":" IPv6Hextet } ;
IPv6Addr   = IPv6Head [ "::" [ IPv6Tail ] ] | "::" [ IPv6Tail ] ;
IPv6Prefix = IPv6Addr "/" ( IntLit | Interp ) ;

Literal = IntLit
        | FloatLit
        | BoolLit
        | RateLit
        | DurationLit
        | SizeLit
        | StrLit
        | EndpointLit
        | IPv4Addr
        | IPv4Prefix
        | IPv6Addr
        | IPv6Prefix ;

(* --- Expressions --- *)

PrimaryExpr = Ident
            | Literal
            | "(" Expr ")" ;
PostfixExpr = PrimaryExpr { "[" Expr "]" | "." Ident } ;
FnAppExpr   = Ident { PostfixExpr }
            | PostfixExpr ;
UnaryExpr   = [ "-" | "!" ] FnAppExpr ;
MulExpr     = UnaryExpr { ( "*" | "/" | "%" ) UnaryExpr } ;
AddExpr     = MulExpr { ( "+" | "-" ) MulExpr } ;
CmpExpr     = AddExpr [ ( "==" | "!=" | "<" | "<=" | ">" | ">=" ) AddExpr ] ;
AndExpr     = CmpExpr { "&&" CmpExpr } ;
OrExpr      = AndExpr { "||" AndExpr } ;
Expr        = OrExpr ;

(* --- Blocks --- *)

BlockGuard  = "|" Expr ;
ClauseGuard = "|" ( "else" | Expr ) ;

LetDecl  = "let" Ident { Ident } [ ClauseGuard ] "=" Expr ;
AttrDecl = Ident [ ClauseGuard ] "=" Expr ;

Generator   = Ident "<-" Expr ".." Expr ;
ForeignDecl = Ident "foreign" ;

InlineWith = "[" Generator { "," Generator } "]" ;

WithDecl = "with" "{" EndOfLine { WithItem EndOfLine } "}" ;
WithItem = Generator
         | ForeignDecl
         | BlockGuard ;

BindDecl = "bind" "{" EndOfLine { BindItem EndOfLine } "}" ;
BindItem = Ident [ "=" Expr ] ;

IncludeDecl = "include" [ InlineWith ] [ BlockGuard ] "{" EndOfLine { IncludeItem EndOfLine } "}" ;
IncludeItem = "path" "=" '"' Utf8Str '"'
            | WithDecl
            | BindDecl ;

BlockDecl = Ident [ StrLit ] [ InlineWith ] [ BlockGuard ] "{" EndOfLine BlockBody "}" ;
BlockBody = { Statement EndOfLine } ;
Statement = WithDecl
          | IncludeDecl
          | BindDecl
          | LetDecl
          | AttrDecl
          | BlockDecl ;

(* --- Top-level file is a block body --- *)
File = [ EndOfLine ] BlockBody ;
```

### Host Schema Contract

The language syntax is domain-agnostic. Valid _block types_, _attributes_, and _enumeration types_ are defined by an external _host schema_ supplied by the embedding application.

#### Block Kinds

Every _block_ must match a _block type_ from the host schema. For each block type, the schema defines the permitted _attributes_ and nested _child block types_.

The host schema classifies each block type into one of two kinds:

- A _list block_ does not accept a string label and collects instantiated elements into a sequence.
- A _dictionary block_ requires a string label and stores instantiated elements in a map indexed by that label.

Only _dictionary blocks_ can be indexed by key in _postfix expressions_.

```hcl
// good: dictionary block (with label) and list block (no label)
node "spine-1" {
	profile = "router"
}
link {
	from = "spine-1"."eth0"
	to   = "leaf-1"."uplink0"
}

// bad: missing string label on a dictionary block
// node {
// 	profile = "router"
// }

// bad: unexpected string label on a list block
// link "fabric-1" {
// 	from = "spine-1"."eth0"
// 	to   = "leaf-1"."uplink0"
// }
```

#### Attribute Modifiers

Every _attribute_ inside a block must match an attribute definition for that block type in the host schema. The schema specifies the attribute's expected _type_ and its _completeness modifier_:

- A _required attribute_ must produce a value for every instance of the block. Completely omitting a required attribute from a block is a static error.
- A _default attribute_ falls back to a schema-defined default value when omitted or when none of its _clause guards_ match.
- A _nullable attribute_ may remain unset when omitted or when none of its _clause guards_ match. However, reading an unset nullable attribute through a cross-reference from another block is an evaluation error.

```hcl
node "spine-1" {
	profile = "router" // required attribute

	interface "eth0" {
		// 'type' and 'mtu' are default attributes (fall back to 'ethernet' and 1500)
		// 'ip' is a nullable attribute and may be omitted if unused by other blocks
	}
}
```

#### Schema Enumerations

The host schema may define named _enumeration types_, each consisting of a finite set of valid variant _identifiers_.

_Enumeration variants_ are written as bare _identifiers_ without quotes.

## Drawbacks

[drawbacks]: #drawbacks

Why should we _not_ do this?

## Rationale and alternatives

[rationale-and-alternatives]: #rationale-and-alternatives

- Why is this design the best in the space of possible designs?
- What other designs have been considered and what is the rationale for not choosing them?
- What is the impact of not doing this?

## Prior art

[prior-art]: #prior-art

Discuss prior art, both the good and the bad, in relation to this proposal.

## Unresolved questions

[unresolved-questions]: #unresolved-questions

- What parts of the design do you expect to resolve through the review process before this gets merged?
- What parts of the design do you expect to resolve through the implementation of this feature?
- What related issues do you consider out of scope for this document that could be addressed in the future independently of the solution that comes out of this design?

## Future possibilities

[future-possibilities]: #future-possibilities

Think about what the natural extension and evolution of your proposal would be and how it would affect the project as a whole in a holistic way.

This is also a good place to "dump ideas", if they are out of scope for the document you are writing but otherwise related.

If you have tried and cannot think of any future possibilities, you may simply state that you cannot think of anything.
