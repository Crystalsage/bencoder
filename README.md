# Bencoder
A [bencode](https://en.wikipedia.org/wiki/Bencode) parser written in Go.

Attempts to parse Bencode based on the following augmented BNF grammar.
```
 <BE>    ::= <DICT> | <LIST> | <INT> | <STR>
 <DICT>  ::= "d" 1 * (<STR> <BE>) "e"
 <LIST>  ::= "l" 1 * <BE>         "e"
 <INT>   ::= "i"     <SNUM>       "e"
 <STR>   ::= <NUM> ":" n * <CHAR>; where n equals the <NUM>

 <SNUM>  ::= "-" <NUM> / <NUM>
 <NUM>   ::= 1 * <DIGIT>
 <CHAR>  ::= %
 <DIGIT> ::= "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
```
## Progress
- [x] Strings
- [x] Integers
- [ ] Lists
- [ ] Dictionaries
