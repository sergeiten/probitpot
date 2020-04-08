package probitpot

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"probitpot/probit"
	"strings"
	"time"

	"github.com/fatih/color"
)

const tokenMargin = 10

type Opts struct {
	ClientID        string  `long:"client_id" description:"Client ID"`
	ClientSecretKey string  `long:"client_secret_key" description:"Client Secret Key"`
	MarketID        string  `long:"market_id" description:"Market ID"`
	MinPrice        float64 `long:"min_price" description:"Minimal price that can be generated"`
	MaxPrice        float64 `long:"max_price" description:"Maximal price that can be generated"`
	MinQuantity     float64 `long:"min_quantity" description:"Minimal quantity of tokens that can be generated"`
	MaxQuantity     float64 `long:"max_quantity" description:"Maximal quantity of tokens that can be generated"`
	Transactions    int     `long:"transactions" description:"Number of transactions that will be generated"`
	SellDelay       int     `long:"sell_delay" description:"Delay after sell action (in milliseconds)"`
	BuyDelay        int     `long:"buy_delay" description:"Delay after buy action (in milliseconds)"`
}

type Bot struct {
	opts             Opts
	client           *probit.Probit
	ticker           *time.Ticker
	lastTokenRefresh time.Time
	tickerDone       chan struct{}
	runDone          chan struct{}
	AllDone          chan struct{}
}

func NewBot(opts Opts) (*Bot, error) {
	rand.Seed(time.Now().UnixNano())

	client, err := probit.NewProbit(opts.ClientID, opts.ClientSecretKey)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		opts:       opts,
		client:     client,
		ticker:     time.NewTicker(1 * time.Second),
		tickerDone: make(chan struct{}),
		runDone:    make(chan struct{}),
		AllDone:    make(chan struct{}),
	}

	go b.tick()

	return b, nil
}

func (b *Bot) Run() error {
	err := b.client.Token()
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	b.lastTokenRefresh = time.Now()

	go func() {
		i := 1
		for {
			select {
			case <-b.runDone:
				return
			default:
				if i > b.opts.Transactions {
					close(b.AllDone)
					return
				}
				limitPrice := round(randF(b.opts.MinPrice, b.opts.MaxPrice), 1)
				quantity := fmt.Sprintf("%.3f", randF(b.opts.MinQuantity, b.opts.MaxQuantity))

				newSellOrder, err := b.client.Sell(b.opts.MarketID, probit.TypeLimit, fmt.Sprintf("%.1f", limitPrice), quantity)
				if err != nil {
					log.Fatalf("failed to sell: %v", err)
				}
				printOrderEvent(newSellOrder)
				b.sleep(b.opts.SellDelay)

				newBuyOrder, err := b.client.Buy(b.opts.MarketID, probit.TypeLimit, fmt.Sprintf("%.1f", limitPrice), quantity, newSellOrder.Data.ClientOrderID)
				if err != nil {
					log.Fatalf("failed to buy: %v", err)
				}
				printOrderEvent(newBuyOrder)

				// don't sleep for last order
				if i != b.opts.Transactions {
					b.sleep(b.opts.BuyDelay)
				}
				i++
			}
		}
	}()

	return nil
}

func (b *Bot) tick() {
	for {
		select {
		case <-b.tickerDone:
			return
		case t := <-b.ticker.C:
			diff := t.Unix() - b.lastTokenRefresh.Unix()

			if diff >= int64(b.client.ExpiredIn)-tokenMargin {
				err := b.client.Token()
				if err != nil {
					log.Fatalf("failed to refresh token: %v", err)
				}

				b.lastTokenRefresh = time.Now()
			}
		}
	}
}

func (b *Bot) Stop() {
	b.ticker.Stop()
	close(b.tickerDone)
	close(b.runDone)
}

func (b *Bot) sleep(delay int) {
	randDelay := randI(1, delay)
	printDelayEvent(randDelay)
	time.Sleep(time.Duration(randDelay) * time.Millisecond)
}

func printOrderEvent(order *probit.NewOrderResponse) {
	datetime := time.Now().Format("2006-01-02 15:04:05")

	side := ""
	switch order.Data.Side {
	case probit.SideBuy:
		side = color.GreenString(strings.ToUpper(order.Data.Side))
	case probit.SideSell:
		side = color.RedString(strings.ToUpper(order.Data.Side))
	}

	limitPrice := color.BlueString(order.Data.LimitPrice)
	quantity := color.MagentaString(order.Data.Quantity)
	log.Printf("Date: %s Action: %s Price: %s Quantity: %s\n", datetime, side, limitPrice, quantity)
}

func printDelayEvent(delay int) {
	log.Printf("Sleep: %.2f seconds\n", float64(delay)/1000)
}

func randF(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randI(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func round(n float64, p int) float64 {
	s := math.Pow10(p)
	return math.Round(n*s) / s
}
