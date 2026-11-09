```python run
def counter():
    n = 0
    def step():
        nonlocal n
        n += 1
        return n
    return step

c = counter()
c()
print(c())
```

```output
2
```
