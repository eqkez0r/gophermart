package authinspector

import (
	"context"
	"errors"
	e "github.com/eqkez0r/gophermart/pkg/error"
	obj "github.com/eqkez0r/gophermart/pkg/objects"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	TTL = time.Minute
)

var (
	errUserIsNotExist = errors.New("user is not exist")
	errInvalidPass    = errors.New("incorrect pass")
)

type GettingUserProvider interface {
	GetUser(context.Context, string) (*obj.User, error)
}

type AuthInspector struct {
	logger *zap.SugaredLogger

	storage GettingUserProvider
	m       sync.Mutex
	authmap map[string]time.Time
}

func New(
	logger *zap.SugaredLogger,
	userstorage GettingUserProvider,
) *AuthInspector {

	return &AuthInspector{
		logger:  logger,
		storage: userstorage,
		m:       sync.Mutex{},
		authmap: make(map[string]time.Time),
	}
}

func (ai *AuthInspector) Auth(ctx context.Context, user *obj.User, authtime time.Time) error {
	const op = "Auth Inspector error: "
	ai.logger.Infof("Auth user: %v", user)
	u, err := ai.storage.GetUser(ctx, user.Login)
	if err != nil {
		return e.Wrap(op, err)
	}
	if u == nil {
		return e.Wrap(op, errUserIsNotExist)
	}
	hashingPass := user.Password //TODO: CHECK HASHING PASS
	if u.Password != hashingPass {
		return e.Wrap(op, errInvalidPass)
	}
	ai.m.Lock()
	ai.authmap[user.Login] = authtime
	ai.m.Unlock()
	return nil
}

func (ai *AuthInspector) CheckAuth(user *obj.User) bool {
	ai.m.Lock()
	_, ok := ai.authmap[user.Login]
	defer ai.m.Unlock()
	return ok
}

func (ai *AuthInspector) CheckInDatabase(login string) (bool, error) {
	u, err := ai.storage.GetUser(context.Background(), login)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, errUserIsNotExist
	}
	return true, nil
}

func (ai *AuthInspector) Observe(ctx context.Context) {
	ticker := time.NewTicker(TTL / 3)
	for range ticker.C {
		select {
		case <-ctx.Done():
			{
				ai.logger.Info("auth inspector stopped")
				return
			}
		default:
			{
				ai.m.Lock()
				for k, v := range ai.authmap {
					if time.Since(v) > TTL {
						delete(ai.authmap, k)
					}
				}
				ai.m.Unlock()
			}
		}
	}
}
