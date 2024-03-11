import logging
import random
import numpy as np

N = 51


class Auction:
    def __init__(self, min_price: float, max_price: float, log: logging.Logger):
        self.feasible_prices = np.linspace(min_price, max_price, num=N)
        self.log = log
        self.count = 0
        self.__gen_winning_curve()

    def step(self, price: float) -> bool:
        self.count += 1
        if self.count % 50 == 0:
            self.__gen_winning_curve()
        return random.uniform(0, 1) <= self.__win_probability(price)

    def __gen_winning_curve(self):
        self.curve = 1 - np.random.default_rng().exponential(scale=1, size=N-1)
        self.curve.sort()
        self.curve = (self.curve - self.curve[0]) / (self.curve[N - 2] - self.curve[0])

    def __win_probability(self, price: float) -> float:
        for i in range(N - 1):
            if self.feasible_prices[i] <= price <= self.feasible_prices[i + 1]:
                return float(self.curve[i])

        self.log.debug(f"unfeasible price: {price} [{self.feasible_prices[0], self.feasible_prices[N-1]}]")

