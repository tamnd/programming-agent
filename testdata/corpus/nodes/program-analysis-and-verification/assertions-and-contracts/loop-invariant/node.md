---
tag: 7NVRNT
kind: concept
slug: loop-invariant
name: Loop invariant
aliases: [invariant]
area: program-analysis-and-verification/assertions-and-contracts
level: intro
summary: A statement about the variables of a loop that is true before each step. It is the main tool to prove that a loop is correct.
sources:
  - book: isbn:9780262046305
    title: Cormen, Leiserson, Rivest and Stein, Introduction to Algorithms, fourth edition
    locator: Section 2.1
  - paper: doi:10.1145/363235.363259
    title: C. A. R. Hoare, An axiomatic basis for computer programming (1969)
---

## Idea

A loop changes its variables at each step. A loop invariant is a fact about these variables that does not change. It is true before the first step. Each step keeps it true. When the loop stops, the invariant and the reason for the stop together tell you that the result is correct.

## How it works

To prove a loop with an invariant $P$, you show three things:

1. **Start.** $P$ is true before the first step.
2. **Step.** If $P$ is true and the loop condition is true, then $P$ is true after one step.
3. **End.** When the loop condition is false, $P$ gives the result that you want.

This is a proof by induction on the number of steps. You also need a proof that the loop stops. For this, find a *variant*: a non-negative integer that becomes smaller at each step.

In a search of a sorted list, as in [search a sorted list](/search-a-sorted-list/), the invariant is: if the key is in the list, it is in $A[lo \ldots hi-1]$. The variant is $hi - lo$.

## In languages

Programmers write an invariant as an assertion, so that a test fails when the invariant is false. Some languages and tools let you state the invariant in the code: Dafny has the clause `invariant`, Ada SPARK has the pragma `Loop_Invariant`, and Eiffel has the clause `invariant` on a loop.

## References
