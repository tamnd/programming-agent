---
tag: 8NNKCK
kind: feature
owner: python
slug: closures-and-nonlocal
name: Closures and nonlocal
area: languages/python/functions
level: intro
summary: A nested function in Python can read the variables of the function around it. The statement nonlocal lets it assign to them.
realizes: [closures]
since: "3.0"
sources:
  - url: https://docs.python.org/3.14/reference/simple_stmts.html#the-nonlocal-statement
    title: The Python Language Reference, 7.13 The nonlocal statement, Python 3.14
    locator: 7.13
---

## Synopsis

```text
nonlocal identifier ["," identifier]*
```

## Semantics

A nested function can read a variable of an enclosing function. An assignment makes a new local variable, unless the statement `nonlocal` names the variable. This makes a [closure](/closures/).

## Examples

```python run
def outer():
    x = 1
    def inner():
        nonlocal x
        x += 1
    inner()
    return x

print(outer())
```

```output
2
```

## References
