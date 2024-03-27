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
        self.params = {'context': self.context.context_hash}
        self.fig, (self.ax_true, self.ax_points) = plt.subplots(nrows=2, ncols=1)
        # self.fig, (self.ax_true, self.ax_lr, self.ax_net, self.ax_ex) = plt.subplots(nrows=4, ncols=1)

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
        # response = requests.get(url=self.url + '/quantities', params=self.params, headers=self.headers)
        # if response.status_code != 200:
        #     return
        # resp_json = response.json()
        # prices = resp_json['price']
        # exploration_qty = resp_json['exploration_qty']
        # exploitation_qty = resp_json['exploitation_qty']
        # update_qty = resp_json['update_qty']
        # probabilities = resp_json['probability']
        #
        # response = requests.get(url=self.url + '/winning-curve', params=self.params, headers=self.headers)
        # if response.status_code != 200:
        #     return
        # wc_prices, learned_probabilities = self.parse_winning_curve_response(response.json())
        #
        # response = requests.get(url=self.url + '/quality', params=self.params, headers=self.headers)
        # if response.status_code != 200:
        #     return
        # q_prices, qualities = self.parse_quality(response.json())
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

        # self.ax_net.cla()
        # if learned_probabilities is not None:
        #
        #     self.ax_net.plot(wc_prices, learned_probabilities, color='red', label='# learned')
        #     self.ax_net.plot(self.context.prices(), self.context.win_rate, color='blue', label='# true')
        #
        #     net_revenue = self.net_revenue(self.context.prices())
        #     e, v, i = self.max_p(net_revenue, self.context.win_rate)
        #     self.ax_net.plot(prices, net_revenue, label='# net revenue')
        #     self.ax_net.scatter(prices, e, alpha=0.7, color='yellow')
        #     self.ax_net.scatter([prices[i]], [v], alpha=1.0, color='red', marker='x', s=100, label=' # best bid price')
        #
        #     net_revenue = self.net_revenue(wc_prices)
        #     _, v, i = self.max_p(net_revenue, learned_probabilities)
        #     self.ax_net.scatter([prices[i]], [v], alpha=1.0, color='black', marker='x', s=100, label=f'# our bid price')
        # self.ax_net.legend()

        # self.ax_ex.cla()
        # self.ax_ex.plot(prices, exploration_qty, label='# explorations')
        # self.ax_ex.plot(prices, update_qty, label='# updates')
        # self.ax_ex.legend()
        #
        # self.ax_exp.cla()
        # self.ax_exp.plot(prices, exploitation_qty, label='# exploitations')
        # self.ax_exp.legend()

        # self.ax_ex.cla()
        # self.ax_ex.plot(q_prices, qualities, color='black', label='# quality')
        # self.ax_ex.legend()


def animate_call(i, animate):
    animate.quantities()


def run_animate(host: str, port: int, context: context.Context, simulator: simulator.Simulator, log: logging.Logger):
    animate = Animate(host, port, context, simulator, log)
    # Plot every 5 sec
    global_any = FuncAnimation(animate.fig, animate_call, interval=5000, fargs=(animate,))
    plt.tight_layout()
    plt.show(block=True)
