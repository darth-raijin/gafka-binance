package binance

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/darth-raijin/gafka-binance/internal/database"
	"go.uber.org/zap"
)

type BinanceInterface interface {
	RotateBasepath(basepath string)
	Ping() error

	GetExchangeInfo() (GetExchangeInfoResponse, error)
	GetAggregateTradeStreams(symbol string, tradesChan chan<- GetAgregateTradeStreamsResponse) error
}

type BinanceClient struct {
	// Dependencies
	Db               *database.Database
	Basepath         string
	log              *zap.Logger
	AndenServiceeHer string
}

func NewBinanceClient(db *database.Database, basepath string, logger *zap.Logger) *BinanceClient {
	return &BinanceClient{
		Db:       db,
		Basepath: basepath,
		log:      logger,
	}
}

func (b *BinanceClient) RotateBasepath(basepath string) {
	b.Basepath = basepath
}

func (b *BinanceClient) Ping() error {
	res, err := http.Get(b.Basepath + "/api/v3/ping")
	if err != nil {
		b.log.Warn("Failed to ping Binance API", zap.Error(err))
		return err
	}

	var parsedResponse PingResponse

	err = json.NewDecoder(res.Body).Decode(&parsedResponse)
	if err != nil {
		b.log.Warn("Failed to parse Binance API ping response", zap.Error(err))
		return err
	}

	b.log.Info("Binance API pinged", zap.Int64("server_time", parsedResponse.Servertime))
	return nil
}

type PingResponse struct {
	Servertime int64 `json:"serverTime"`
}

func (b *BinanceClient) GetExchangeInfo() (GetExchangeInfoResponse, error) {
	res, err := http.Get(b.Basepath + "/api/v3/exchangeInfo")
	if err != nil {
		return GetExchangeInfoResponse{}, err
	}

	if res.StatusCode != 200 {
		b.log.Warn("Binance API returned non-200 status code", zap.Int("status_code", res.StatusCode))
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return GetExchangeInfoResponse{}, err
	}

	var response GetExchangeInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return GetExchangeInfoResponse{}, err
	}

	return response, nil
}

type GetExchangeInfoResponse struct {
	Timezone        string                        `json:"timezone"`
	ServerTime      int64                         `json:"serverTime"`
	RateLimits      []any                         `json:"rateLimits"`
	ExchangeFilters []any                         `json:"exchangeFilters"`
	Symbols         GetExchangeInfoResponseSymbol `json:"symbols"`
}

type GetExchangeInfoResponseSymbol struct {
	Symbol                          string   `json:"symbol"`
	Status                          string   `json:"status"`
	BaseAsset                       string   `json:"baseAsset"`
	BaseAssetPrecision              int      `json:"baseAssetPrecision"`
	QuoteAsset                      string   `json:"quoteAsset"`
	QuotePrecision                  int      `json:"quotePrecision"`
	QuoteAssetPrecision             int      `json:"quoteAssetPrecision"`
	OrderTypes                      []string `json:"orderTypes"`
	IcebergAllowed                  bool     `json:"icebergAllowed"`
	OcoAllowed                      bool     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed      bool     `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool     `json:"allowTrailingStop"`
	CancelReplaceAllowed            bool     `json:"cancelReplaceAllowed"`
	IsSpotTradingAllowed            bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool     `json:"isMarginTradingAllowed"`
	Filters                         []any    `json:"filters"`
	Permissions                     []string `json:"permissions"`
	DefaultSelfTradePreventionMode  string   `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string `json:"allowedSelfTradePreventionModes"`
}

func (b *BinanceClient) GetAggregateTradeStreams(symbol string, tradesChan chan<- GetAgregateTradeStreamsResponse) error {
	u := url.URL{Scheme: "wss", Host: "fstream.binance.com", Path: fmt.Sprintf("/ws/%s@aggTrade", symbol)}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				close(tradesChan) // Signal the receiver that the stream is closed
				return
			}

			var trade GetAgregateTradeStreamsResponse
			if err := json.Unmarshal(message, &trade); err != nil {
				// Log error or handle it according to your error policy
				continue
			}

			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				continue
			}
			quantity, err := strconv.ParseFloat(trade.Quantity, 64)
			if err != nil {
				continue
			}

			b.log.Info("Received trade", zap.String("symbol", strings.ToUpper(symbol)), zap.Float64("price", price), zap.Float64("quantity", quantity))
			tradesChan <- trade // Send the trade to the channel
		}
	}()

	return nil
}

type GetAgregateTradeStreamsResponse struct {
	EventType     string `json:"e"`
	EventTim      int64  `json:"E"`
	Symbol        string `json:"s"`
	TradeID       int    `json:"t"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	BuyerOrderID  int    `json:"b"`
	SellerOrderID int    `json:"a"`
	TradeTime     int64  `json:"T"`
	MarketBuyer   bool   `json:"m"`
}
