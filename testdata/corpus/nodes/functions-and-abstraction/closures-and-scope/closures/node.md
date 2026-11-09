---
tag: 7QKX2M
kind: concept
slug: closures
name: Closure
aliases: [lexical closure, function closure]
area: functions-and-abstraction/closures-and-scope
level: intro
summary: A closure is a function together with the variables of the scope where it was made. The function can use these variables after that scope ends.
applies_to: core-languages
viz: optional
formal: none
sources:
  - paper: doi:10.1093/comjnl/6.4.308
    title: P. J. Landin, The mechanical evaluation of expressions (1964)
    locator: Section 5, closures
  - url: https://docs.python.org/3.14/reference/executionmodel.html
    title: The Python Language Reference, 4. Execution model, Python 3.14
    locator: 4.2.2 Resolution of names
---

## Idea

A function can use a variable that it does not define. The variable comes from the scope around the function. When the function is a value, it can live longer than that scope. A closure keeps the variables of that scope alive, so the function still works.

## How it works

The language stores the function together with references to the free variables that it uses. Each call of the outer function makes new variables, so each closure has its own copy of the state.

## In languages

The example files show the same counter in each language.

## References
