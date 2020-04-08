probitpot: a bot for generating trading transactions
=========================================================

This application provides generating transactions between same account in [PROBIT](https://probit.com) currency exchange market.

Application accepts the following flags:
* Client ID in profile settings page (--client_id)
* Client Secret Key in profile settings page (--client_secret_key)
* Market ID is unique market id. Ex: "HCUT-KRW" (--market_id)
* MinPrice is minimal price bot can sell or buy tokens (--min_price)
* MaxPrice is maximal price bot can sell or buy tokens (--max_price)
* MinQuantity is minimal amount of token bot can sell or buy (--min_quantity)
* MaxQuantity is maximal amount of token bot can sell or buy (--max_quantity)
* Transactions is number of total sell/buy pair transactions (--transactions)
* SellDelay in milliseconds after sell action. Waits for random milliseconds between 0 and delay value (--sell_delay)
* BuyDelay in milliseconds after buy action. Waits for random milliseconds between 0 and delay value (--buy_delay)

### Run Example
```bash
probitpot --client_id=client_id --client_secret_key=client_secret_key --market_id=HCUT-KRW --min_price=4.3 --max_price=4.6 --min_quantity=100 --max_quantity=1000 --transactions=10 --sell_delay=100 --buy_delay=500
```

