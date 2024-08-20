# lpsolve

Solve resource assignment problem using Mixed Integer Linear Programming (MILP) solver

This is focused on the optimal assignment of heterogeneous accelarator types (GPUs) to multiple inference servers.

## Notes

- The [Golp](https://pkg.go.dev/github.com/draffensperger/golp) package is used as a Golang wrapper for the [lp_solve](https://lpsolve.sourceforge.net/5.5/) linear (and integer) programming library.
- Following installation instructions are for mac.
- Install lp_solve using brew.

```bash
brew install lp_solve
```

- Create a target directory `/opt/local/include/lpsolve` with a soft link to the installed include directory, e.g. `/homebrew/Cellar/lp_solve/5.5.2.11/include/`.
- Add soft links in directory `/opt/local/lib` to files in `/opt/homebrew/lib`.
