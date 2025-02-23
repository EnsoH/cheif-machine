package account

// import (
// 	"context"
// 	"dw/models"
// 	"sync"

// 	"golang.org/x/sync/errgroup"
// )

// type Account struct {
// 	API_secret string
// 	API_key    string
// 	Coins      string
// }

// func AccountFactory(addresses []string, withdrawConfig *models.WithdrawConfig) ([]Account, error) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	g, ctx := errgroup.WithContext(ctx)
// 	var (
// 		mu   sync.Mutex
// 		accs []*Account
// 	)

// 	accs = make([]*Account, len(addresses))

// 	for _, addr := range addresses {
// 		addr := addr

// 		g.Go(func() error {
// 			if ctx.Err() != nil {
// 				return ctx.Err()
// 			}

// 			account
// 		})
// 	}

// 	return nil, nil
// }
