# Tail
Online/adaptive learning bid optimizer
To be optimal when participating to online biddings, we should be able to solve the following problem:
What is the best price to bid? in other words what is the optimal price to bid?

This project is trying to answer to that question, for details pelase see doc.

## Build and run
### Server
1. make tidy
2. make audit
3. make build
5. /tmp/bin/a.out

### Client
1. pipenv shell
2. python main.py

### TODO
- [] Use exploitation feedback for free update.
- [] Use exploitation feedback for online evaluation.
- [] For evaluation and quality use Chernoff bound.
- [] Add weighted average.
- [] Add demo to show gain.
