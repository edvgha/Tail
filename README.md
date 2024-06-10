<a name="readme-top"></a>

# Online Bid Optimizer for DSPs and AdExchange Paltforms
Online/adaptive learning bid optimizer

It's a key to be able to adopt your bid optimizer very fast since upstream auction winning price 
distribution are changing very fast and classic offline trained models can't catch that change in time.

![plot](./doc/dist_change.png)

if we knew upstream winning curve distribution we could estimate best bid price by calculating maximum expected 
bid price under the winning curve distribute. 
So task is reduced to estimate upstream winning curve as accurate as possible in the short period of time and later adopt it in runtime.

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
