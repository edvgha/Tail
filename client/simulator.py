import logging

import context
import client
import auction


class Simulator:
    def __init__(self, cln: client.Client, ctx: context.Context, log: logging.Logger):
        self.client = cln
        self.context = ctx
        self.signal = False
        self.log = log
        self.auction = auction.Auction(self.context.min_price, self.context.max_price, log)

    def stop(self):
        self.signal = True

    def run(self):
        while not self.signal:
            bid_response = self.client.send_bid_request()
            if bid_response.status == "error":
                self.log.debug(" ---> send_bid_request(): status == error")
                continue

            if bid_response.optimized_price > bid_response.price_to_bid:
                raise Exception("simulator.run(): optimized price > bid price")
            if bid_response.status != "explored":
                continue
            impression = self.auction.step(bid_response.optimized_price)

            if impression:
                res = self.client.send_impression(bid_response.req_id, bid_response.optimized_price, impression)
                if not res:
                    self.log.debug(" ---> send_impression(): ack == false")

