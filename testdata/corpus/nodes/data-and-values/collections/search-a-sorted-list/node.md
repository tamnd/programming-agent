---
tag: 2B7QDX
kind: task
slug: search-a-sorted-list
name: Search a sorted list
aliases: [binary search, find in a sorted array]
area: data-and-values/collections
level: intro
summary: Find the index of a value in a sorted list, or report that the value is not in the list.
category: searching
uses_algorithms: [binary-search]
complexity:
  vars: {n: the number of items in the list}
  time: {worst: Θ(log n), average: Θ(log n), best: Θ(1)}
  space: {extra: Θ(1)}
prereqs: [loop-invariant, asymptotic-notation]
see_also: [loop-invariant]
sources:
  - book: isbn:9780201896855
    title: Donald E. Knuth, The Art of Computer Programming, Volume 3, second edition
    locator: Section 6.2.1, Algorithm B
  - url: https://docs.python.org/3.14/library/bisect.html
    title: The Python Standard Library, bisect, Python 3.14
    locator: bisect_left
  - url: https://research.google/blog/extra-extra-read-all-about-it-nearly-all-binary-searches-and-mergesorts-are-broken/
    title: Joshua Bloch, Nearly all binary searches and mergesorts are broken (2006)
    note: the overflow of (lo + hi) / 2
---

## Task

The input is a list $A$ of $n$ items in non-decreasing order, and a value $x$. The output is an index $i$ with $A[i] = x$, or $-1$ when no item is equal to $x$. When the list has more than one copy of $x$, any index of a copy is a correct answer.

The solution only compares items. So it works for numbers, strings and any other type with a total order.

## Solutions

Each solution keeps a range of indices $[lo, hi)$. The value can only be at an index in this range. Each step compares $x$ with the middle item and keeps the half that can hold $x$.

```pseudo id=main
SEARCH(A, x)
L1  lo ← 0, hi ← n                    ▷ the range is A[lo .. hi − 1]
L2  while lo < hi do
L3      mid ← lo + ⌊(hi − lo) / 2⌋     ▷ no overflow, unlike ⌊(lo + hi) / 2⌋
L4      if A[mid] < x then
L5          lo ← mid + 1                ▷ x is not in A[lo .. mid]
L6      else if A[mid] > x then
L7          hi ← mid                    ▷ x is not in A[mid .. hi − 1]
L8      else return mid
L9  return −1
```

The proof uses a [loop invariant](/loop-invariant/): if $x$ is in $A$, then $x$ is in $A[lo \ldots hi-1]$. Each step removes at least half of the range, so the loop runs at most $\lfloor \lg n \rfloor + 1$ times. The time is $\Theta(\log n)$ in the worst case, in the sense of [asymptotic notation](/asymptotic-notation/).

## Edge cases

- The list is empty. The loop does not run, and the result is $-1$.
- The value is smaller than each item, or larger than each item.
- The list has many copies of the value.

## Pitfalls

- **Overflow.** The expression $\lfloor (lo + hi)/2 \rfloor$ can overflow a fixed-size integer when the list is very large. This bug was in the Java library for nine years. Use $lo + \lfloor (hi - lo)/2 \rfloor$, as line L3 does.
- **An infinite loop.** With a closed range $[lo, hi]$, the update "$lo \gets mid$" can keep the range the same size. Each update must remove $mid$, or the range must be half-open as here.
- **Unsorted input.** On a list that is not sorted, the result has no meaning. The solution does not detect this.

## References
