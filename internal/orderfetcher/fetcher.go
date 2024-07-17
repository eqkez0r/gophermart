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
	UpdateAccrual(context.Context, uint64, *obj.Accrual) error
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
	for {
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

				retry:
					or.logger.Infof("Send request to order number %d", o.Number)
					res, err := or.client.R().
						Get("http://" + or.accrualuri + "/api/orders/" + o.Number)
					if err != nil {
						or.logger.Warnw("failed to get orders", "error", err)
						continue
					}
					or.logger.Infof("Successfully request to order number %d."+
						" Recieved status code %d", o.Number, res.StatusCode())

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
							or.logger.Infof("Successfully request to order number %d. Response: %v", o.Number, res.String())
							d := res.Body()
							accrual := &obj.Accrual{}
							err = json.Unmarshal(d, accrual)
							if err != nil {
								or.logger.Warnw("failed to unmarshal withdraw accrual", "error", err)
							}
							//or.logger.Infof("Successfully request to order number %d. Response: %v", o.Number, res.String())
							if accrual.Status == obj.AccrualStatusInvalid || accrual.Status == obj.AccrualStatusProcessed {
								if err = or.storage.UpdateAccrual(ctx, o.UserID, accrual); err != nil {
									or.logger.Warnw("failed to update after accrual", "error", err)
								}
							}
						}
					default:
						{
							or.logger.Warnw("failed to parse orders", "error", res.String())
						}
					}
				}
			}
		}
	}
}

func (or *OrderFetcher) handle(
	accrual *obj.Accrual,

) error {

	return nil
}
