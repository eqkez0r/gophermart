package authinspector

import (
	"errors"
	"github.com/eqkez0r/gophermart/internal/storage"
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

type AuthInspector struct {
	logger *zap.SugaredLogger

	userstorage storage.Storage
	m           sync.Mutex
	authmap     map[string]time.Time
}

func New(
	logger *zap.SugaredLogger,
	userstorage storage.Storage,
) *AuthInspector {

	return &AuthInspector{
		logger:      logger,
		userstorage: userstorage,
		m:           sync.Mutex{},
		authmap:     make(map[string]time.Time),
	}
}

func (AI *AuthInspector) Auth(user obj.User, authtime time.Time) error {
	const op = "Auth Inspector error: "
	u, err := AI.userstorage.GetUser(user.Login)
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
	AI.authmap[user.Login] = authtime
	return nil
}

func (AI *AuthInspector) CheckAuth() (bool, error) {

	return false, nil
}

func (AI *AuthInspector) CheckInDatabase() (bool, error) {

	return false, nil
}

func (AI *AuthInspector) Observe() {

}
