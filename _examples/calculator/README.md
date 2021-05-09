# calculator

A simple calculator to demonstrate how to parse expressions with pre/in/postfix operators of varying associativity.

## How To Run

```
go run .
```

## Supported operations

* `q`: Quits the program.
* `+-*/%`: All basic math operators (`1+2*3`).
* `^` or `ˆ`: Power (`2ˆ3`).
* `!`: Factorial (`3!`).
* `√`: Square root (`√2`) or any other root (`3√9`). Not associative meaning `2√8√16` is an error, but `(2√8)√16` is okay.
* `()`: Parentheses (`(1+2)*3`).
* `π`: Produces `π` (`π*10ˆ2`).
* `ans`: Produces the last answer (`ans+1`).
