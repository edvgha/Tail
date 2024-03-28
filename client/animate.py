import matplotlib as mpl
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation
import requests
import logging
import context
import simulator

mpl.use('Qt5Agg')
plt.style.use('seaborn')


class Animate:
    def __init__(self, host: str, port: int, context: context.Context, simulator: simulator.Simulator, log: logging.Logger):
        self.headers = {'Content-Type': 'application/json'}
        self.url = 'http://' + host + ':' + str(port)
        self.context = context
        self.simulator = simulator
        self.log = log
        self.params = {'ctx': self.context.context_hash}
        self.fig, (self.ax_true, self.ax_points, self.ax_space) = plt.subplots(nrows=3, ncols=1)

    @staticmethod
    def parse_winning_curve_response(response) -> (list[float], list[float]):
        learned_probabilities = response['learned_probability']
        prices = []
        if learned_probabilities is None:
            return None, None
        for p in response['probs']:
            prices.append(p['price'])
        return prices, learned_probabilities

    @staticmethod
    def parse_quality(response) -> (list[float], list[float]):
        return response['price'], response['quality']

    @staticmethod
    def net_revenue(prices) -> list[float]:
        net_revenue = []
        m = prices[len(prices) - 1]
        for p in prices:
            net_revenue.append(m - p)
        return net_revenue

    @staticmethod
    def max_p(net_revenue, probabilities) -> (list[float], float, int):
        e = [i * j for i, j in zip(probabilities, net_revenue)]
        m_id = -1
        m_val = 0.0
        for i, v in enumerate(e):
            if v > m_val:
                m_val = v
                m_id = i
        return e, m_val, m_id

    def quantities(self):
        response = requests.get(url=self.url + '/space', params=self.params, headers=self.headers)
        if response.status_code != 200:
            return
        resp_json = response.json()
        levels = resp_json['level']
        self.ax_true.cla()
        self.ax_true.plot(self.simulator.auction.prices(), self.simulator.auction.curve,
                          color='blue', label='# true win prob')
        self.ax_true.plot(self.simulator.auction.prices(), self.simulator.auction.net_revenue(),
                          color='blue', label='# true net revenue', linestyle='--')
        self.ax_true.plot(self.simulator.auction.prices(), self.simulator.auction.expectations(),
                          color='green', label='# true expectations', linestyle='-.')
        self.ax_true.scatter([self.simulator.auction.optimal_price()], [0.0],
                             alpha=1.0, color='red', marker='x', s=100, label='# optimal price')
        self.ax_true.legend()

        self.ax_points.cla()
        for i in range(len(self.simulator.d)):
            self.log.debug(f'true: {self.simulator.d[i][0]}, opt: {self.simulator.d[i][1]}, price: {self.simulator.d[i][2]}')
            self.ax_points.scatter([self.simulator.d[i][0]], [0.0], alpha=1.0, color='green', marker='x', s=100)
            self.ax_points.scatter([self.simulator.d[i][1]], [0.0], alpha=1.0, color='grey', marker='o', s=100)
        self.ax_points.legend()

        color = ['yellow', 'green', 'blue', 'grey', 'red']
        self.ax_space.cla()
        for i in range(len(levels)):
            self.ax_space.plot(levels[i]['price'], levels[i]['pr'], color=color[i])
        self.ax_space.legend()


def animate_call(i, animate):
    animate.quantities()


def run_animate(host: str, port: int, context: context.Context, simulator: simulator.Simulator, log: logging.Logger):
    animate = Animate(host, port, context, simulator, log)
    # Plot every 5 sec
    global_any = FuncAnimation(animate.fig, animate_call, interval=5000, fargs=(animate,))
    plt.tight_layout()
    plt.show(block=True)
