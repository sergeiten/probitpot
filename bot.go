package probitpot

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"probitpot/probit"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Opts struct {
	ClientID        string  `long:"client_id" description:"Client ID"`
	ClientSecretKey string  `long:"client_secret_key" description:"Client Secret Key"`
	MarketID        string  `long:"market_id" description:"Market ID"`
	MinPrice        float64 `long:"min_price" description:"Minimal price that can be generated"`
	MaxPrice        float64 `long:"max_price" description:"Maximal price that can be generated"`
	MinQuantity     int     `long:"min_quantity" description:"Minimal quantity of tokens that can be generated"`
	MaxQuantity     int     `long:"max_quantity" description:"Maximal quantity of tokens that can be generated"`
	Transactions    int     `long:"transactions" description:"Number of transactions that will be generated"`
	Delay           int     `long:"delay" description:"Delay between transactions (in seconds)"`
}

type Bot struct {
	opts   Opts
	client *probit.Probit
}

func NewBot(opts Opts) (*Bot, error) {
	rand.Seed(time.Now().UnixNano())

	client, err := probit.NewProbit(opts.ClientID, opts.ClientSecretKey)
	if err != nil {
		return nil, err
	}

	return &Bot{
		opts:   opts,
		client: client,
	}, nil
}

func (b *Bot) Run() error {
	for i := 1; i <= b.opts.Transactions; i++ {
		err := b.client.Token()
		if err != nil {
			return fmt.Errorf("failed to get token: %v", err)
		}

		limitPrice := round(randF(b.opts.MinPrice, b.opts.MaxPrice), 1)
		quantity := strconv.Itoa(randI(b.opts.MinQuantity, b.opts.MaxQuantity))

		newSellOrder, err := b.client.Sell(b.opts.MarketID, probit.TypeLimit, fmt.Sprintf("%.1f", limitPrice), quantity)
		if err != nil {
			log.Fatalf("failed to sell: %v", err)
		}
		printOrderEvent(newSellOrder)
		b.sleep()

		newBuyOrder, err := b.client.Buy(b.opts.MarketID, probit.TypeLimit, fmt.Sprintf("%.1f", limitPrice), quantity, newSellOrder.Data.ClientOrderID)
		if err != nil {
			log.Fatalf("failed to buy: %v", err)
		}
		printOrderEvent(newBuyOrder)

		// don't sleep for last order
		if i != b.opts.Transactions {
			b.sleep()
		}
	}

	return nil
}

func (b *Bot) sleep() {
	delay := randI(1, b.opts.Delay)
	printDelayEvent(delay)
	time.Sleep(time.Duration(delay) * time.Second)
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
	log.Printf("Sleep: %d seconds\n", delay)
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
