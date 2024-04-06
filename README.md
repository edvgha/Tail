<a name="readme-top"></a>

# Tail
Online/adaptive learning bid optimizer
To be optimal when participating to online biddings, we should be able to solve the following problem:
What is the best price to bid? in other words what is the optimal price to bid?

This project is trying to answer to that question, for details pelase see doc.

## Build and run
### Server
```
make tidy
```
```
make audit
```
```
make build
```
```
/tmp/bin/a.out
```

### Client
```
pipenv shell
```
```
python main.py
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### TODO
- [ ] Use exploitation feedback for free update.
- [ ] Use exploitation feedback for online evaluation.
- [ ] For evaluation and quality use Chernoff bound.
- [ ] Add weighted average.
- [ ] Add demo to show gain.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
