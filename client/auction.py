import random
import numpy as np


class Auction:
    def __init__(self, min_price: float, max_price: float):
        self.feasible_prices = np.linspace(min_price, max_price)
        self.__gen_winning_curve()
        self.count = 0
        self.n = 50

    def step(self, price: float) -> bool:
        self.count += 1
        if self.count % 50 == 0:
            self.__gen_winning_curve()
        return random.uniform(0, 1) <= self.__win_probability(price)

    def __gen_winning_curve(self):
        self.curve = 1 - np.random.default_rng().exponential(scale=1, size=self.n)
        self.curve.sort()
        self.curve = (self.curve - self.curve[0]) / (self.curve[self.n - 1] - self.curve[0])

    def __win_probability(self, price: float) -> float:
        for i in range(self.n - 1):
            if self.feasible_prices[i] <= price <= self.feasible_prices[i + 1]:
                return float(self.curve[i])

        raise Exception("__win_probability(): unfeasible price")

