# IBM CPLEX as a solver

We use [IBM ILOG CPLEX](https://www.ibm.com/products/ilog-cplex-optimization-studio) as a solver to the resource assignment Mixed Integer Linear Programming (MILP) problem.

## Notes

- CPLEX (including opl) is assumed to be installed.
- Since we use the Go language, and given that there is no (as far as we know) reliable Go library to interface to CPLEX, we use the `oplrun` CLI command directly. This involves:
  - creating an opl model file for the particular optimization problem,
  - creating a data file populated with the input parameters of the problem, and
  - processing the output of the opl command to extract the solution.
- Two environment variables are assumed to be set, in order to provide information about various paths (should include `/` at the end):
  - `CPLEX_MODEL_PATH` path to the opl models, and
  - `CPLEX_DATA_PATH` path to the input and output data files.
