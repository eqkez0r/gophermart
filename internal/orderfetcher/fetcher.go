package orderfetcher

import (
	"context"
	"encoding/json"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// pool order
// 429
// time.sleep
// try again this order

type OrdersProvider interface {
	GetUnfinishedOrders(ctx context.Context) ([]*obj.Order, error)
}

type OrderFetcher struct {
	storage    OrdersProvider
	logger     *zap.SugaredLogger
	accrualuri string
	client     *resty.Client
}

func New(
	logger *zap.SugaredLogger,
	accrualuri string,
	s OrdersProvider,
) *OrderFetcher {

	return &OrderFetcher{
		storage:    s,
		client:     resty.New(),
		logger:     logger,
		accrualuri: accrualuri,
	}
}

func (or *OrderFetcher) Run(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-ctx.Done():
			{
				or.logger.Infof("order fetcher stopped")
				wg.Done()
				return
			}
		default:
			{
				orders, err := or.storage.GetUnfinishedOrders(ctx)
				if err != nil {
					or.logger.Warnw("failed to get unfinished orders", "error", err)
					continue
				}
				for _, o := range orders {
					orderNum := int(*o.Number)
				retry:
					or.logger.Infof("Send request to order number %d", orderNum)
					res, err := or.client.R().
						Get("http://" + or.accrualuri + "/api/orders/" + strconv.Itoa(orderNum))
					if err != nil {
						or.logger.Warnw("failed to get orders", "error", err)
						continue
					}
					or.logger.Infof("Successfully request to order number %d."+
						" Recieved status code %d", orderNum, res.StatusCode())

					switch res.StatusCode() {
					case http.StatusTooManyRequests:
						{
							rt := res.Header().Get("Retry-After")
							rtInt, err := strconv.Atoi(rt)
							if err != nil {
								or.logger.Warnw("failed to parse retry after", "error", err)
							}
							time.Sleep(time.Millisecond * time.Duration(rtInt))
							goto retry
						}
					case http.StatusOK:
						{
							or.logger.Infof("Successfully request to order number %d. Response: %v", orderNum, res.String())

						}
					}
				}
			}
		}
	}
}

func (or *OrderFetcher) Post(ctx context.Context, order *obj.Withdraw) error {
	b, err := json.Marshal(order)
	if err != nil {
		return err
	}

	resp, err := or.client.R().
		SetBody(b).
		SetHeader("Content-Type", "application/json").
		Post("http://" + or.accrualuri + "/api/orders")
	if err != nil {
		return err
	}
	or.logger.Infof("Order response status %s, resposne: %s", resp.Status(), resp.String())
	return nil
}
