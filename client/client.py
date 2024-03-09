from context import Context
import uuid
import time
import requests


class BidResponse:
    def __init__(self, req_id="", price_to_bid=0.0, optimized_price=0.0, status="error"):
        self.req_id = req_id
        self.price_to_bid = price_to_bid
        self.optimized_price = optimized_price
        self.status = status


class Client:
    def __init__(self, context: Context, host: str, port: int):
        self.context = context
        self.headers = {'Content-Type': 'application/json'}
        self.url = 'http://' + host + ':' + str(port)
        self.last_ts = None
        self.context.show()

    def send_bid_request(self):
        price_to_bid, ts = self.context.get_price_in_time()

        if self.last_ts is not None:
            time.sleep((ts - self.last_ts) / 1000)
        self.last_ts = ts

        req_id = uuid.uuid4().hex
        json_body = {
            "id": str(req_id),
            "price": price_to_bid,
            "floor_price": self.context.floor_price,
            "data_center": self.context.dc,
            "app_publisher_id": self.context.pub_id,
            "bundle_id": self.context.bundle_id,
            "tag_id": self.context.tag_id,
            "device_geo_country": self.context.cc,
            "ext_ad_format": self.context.ad_format,
        }

        response = requests.post(url=self.url + '/optimize', json=json_body, headers=self.headers)
        if response.status_code != 200:
            return BidResponse()

        resp_json = response.json()
        optimized_price = resp_json['optimized_price']
        status = resp_json['status']

        bid_response = BidResponse(req_id=req_id,
                                   price_to_bid=price_to_bid,
                                   optimized_price=optimized_price,
                                   status=status)
        return bid_response

    def send_impression(self, req_id: str, price: float, imp: bool):
        json_body = {"id": req_id,
                     "price": price,
                     "impression": imp,
                     }
        response = requests.post(url=self.url + '/feedback', json=json_body, headers=self.headers)
        if response.status_code != 200:
            print("Feedback status code: ", response.status_code)

        resp_json = response.json()
        ack = resp_json['ack']
        return response.status_code == 200 and ack
