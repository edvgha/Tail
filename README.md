<a name="readme-top"></a>

# Tail
Online/adaptive learning bid optimizer

It's a key to be able to adopt your model very fast since auction winning price distribution 
are changing very fast and classic offline learning models can't catch that change.

During participation on upstream auction (in AdTech) we are learning maximum expected value 
under bid price distribution, which will provide best price to bid (do not overpay and do not underpay).

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
