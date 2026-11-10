## Before you start

All work happens on GitHub. You need a GitHub account, Git and Go. To run the code of a language, you need its toolchain, or Docker.

Install the toolchain:

```sh
go install github.com/tamnd/programming-agent/cmd/prog@latest
```

Clone the content and check it:

```sh
git clone https://github.com/tamnd/programming
cd programming
prog check
prog site serve
```

`prog site serve` builds the site and serves it at http://localhost:8150.

## Write a new page

1. Find the page in the coverage plan, or open an issue that proposes it.
2. Run `prog new <kind> <slug> -area <area>`, for example `prog new concept closures -area functions-and-abstraction/closures-and-scope`. A feature, a release, a diagnostic, a proposal and an api entry also need `-owner`, such as `-owner python`. The command makes the directory, the front matter and the sections of the kind.
3. Write the sections in the order that the content standard gives.
4. Run `prog check`. Fix each finding.
5. Run `prog run <slug>`. Each code block of the page must give the output that the page shows.
6. Open a pull request. Use one pull request for each page.

## Write a code block

Each code block in a page runs in CI. Give the language after the fence, and give one mode:

- `run`: the block runs, and the next `output` block is its exact output.
- `fails expect=<signature>`: the block must fail in the way that the signature says, such as `expect=compile` or `expect=exception:ValueError`.
- `wrong`: the block runs, but it shows a mistake. The text after it says why.
- `norun reason=<reason>`: the block does not run. The reason is one of network, gui, platform, hardware, ub, vulnerable, fragment, pseudo, long or interactive.

`prog run <slug> -update` writes the output blocks from a real run. Do not write output by hand.

## Add a solution to a task

A task has a contract and shared cases. Each solution is one file in `impl/<language>/` with a native test file next to it. Each solution must pass the same cases.

```sh
prog test <slug>
```

## Promote a card

A card is one row of facts in a register under `cards/`, such as an error code, a proposal or a standard library function. The site shows each card on the page of its collection, for example `/rust/errors/#E0382`. A card has no full page.

To write the full page of a card, promote it:

```sh
prog promote rust/errors/E0382
```

The command makes the node directory, moves the facts of the card into the front matter, and removes the row from the register. The tag does not change, so old links still work. Then write the sections and open a pull request, as for a new page.

## Add a scene

A scene is `viz/scene.ts`. It imports only `@prog/viz`. It must give the same trace for the same input. Put at least 3 presets in `viz/presets.yaml`. Each step of the program must be reached by at least one preset.

## Writing rules

- Write short sentences. Use one idea in each sentence.
- Use the active voice.
- Use simple words. Define each technical term before you use it.
- Do not put a line break in the middle of a sentence.
- Do not use a dash as a pause in a sentence.
- Do not copy text from a book or from documentation. Write it in your own words and cite the source.
