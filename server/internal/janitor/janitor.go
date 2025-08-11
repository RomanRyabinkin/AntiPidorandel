package janitor

import (
	"context"
	"log"
	"time"

	"github.com/RomanRyabinkin/antipidorandel/internal/config"
	"github.com/RomanRyabinkin/antipidorandel/internal/store"
)

func Start(ctx context.Context, st store.Store, cfg config.Config) {
	t := time.NewTicker(cfg.JanitorInterval)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-t.C:
				dBefore := time.Now().Add(-cfg.RetainDelivered)
				now := time.Now()
				d1, d2, err := st.DeleteExpired(context.Background(), dBefore, now)
				if err != nil {
					log.Printf("janitor error: %v", err) // TODO: Убрать принты в будущем. Сейчас пока черновая версия
				} else if d1+d2 > 0 {
					log.Printf("janitor: removed delivered=%d expired=%d", d1, d2) // TODO: Убрать принты в будущем. Сейчас пока черновая версия
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}