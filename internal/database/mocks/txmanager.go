package database_mocks

import "context"

type MockTxManager struct {
	BeginFn    func() (context.Context, error)
	CommitFn   func() error
	RollbackFn func() error
}

func (m MockTxManager) Begin(ctx context.Context) (context.Context, error) {
	return m.BeginFn()
}
func (m MockTxManager) Commit(ctx context.Context) error {
	return m.CommitFn()
}
func (m MockTxManager) Rollback(ctx context.Context) error {
	return m.RollbackFn()
}
