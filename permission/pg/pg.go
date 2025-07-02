package pg

import (
	"context"
	"sync"
	"errors"

	"github.com/dronm/gobizap/permission"

	"github.com/jackc/pgx/v5/pgxpool"
)

var manager = &Manager{}

type Manager struct {
	dbPool *pgxpool.Pool
	mx sync.RWMutex
	rules permission.PermRules
}

func (mng *Manager) Reload() error{
	if mng.dbPool == nil {
		return errors.New("InitManager dbPool must initialized with *pgxpool.Pool")
	}
	mng.mx.Lock()
	defer mng.mx.Unlock()

	mng.rules = make(permission.PermRules)
	if err := mng.dbPool.QueryRow(context.Background(), `SELECT rules FROM permissions LIMIT 1`).Scan(&mng.rules); err != nil {
		return err
	}
	return nil
}

//controller=no _Controller postfix!!!
func (mng *Manager) IsAllowed(role, controller, method string) bool{
	return mng.rules.IsAllowed(role, controller, method)
}

// InitManager initializes permission manager
// First parameter: ConnectionString string in pg format: postgresql://{USER_NAME}@{HOST}:{PORT}/{DATABASE}
func (mng *Manager) InitManager(mngParams []interface{}) (err error) {
	if len(mngParams) < 1 {
		return errors.New("InitManager missing parameter: *pgxpool.Pool, string")
	}	
	ok := false
	mng.dbPool, ok = mngParams[0].(*pgxpool.Pool)
	if !ok {
		return errors.New("InitManager dbPool parameter must be of type *pgxpool.Pool")
	}
	mng.Reload()
	return nil
}

func init() {
	permission.Register("pg", manager)
}
