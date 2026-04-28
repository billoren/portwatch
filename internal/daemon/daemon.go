package daemon

import (
	"context"
	"log"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/monitor"
	"portwatch/internal/notify"
	"portwatch/internal/rules"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
)

// Daemon orchestrates the port-watching loop.
type Daemon struct {
	cfg     *config.Config
	mon     *monitor.Monitor
	alerts  alert.Alerter
	ticker  *time.Ticker
}

// New constructs a Daemon from the given config path.
func New(cfgPath string) (*Daemon, error) {
	cfg, err := config.LoadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	rs, err := rules.LoadFile(cfg.RulesFile)
	if err != nil {
		return nil, err
	}

	sc := scanner.New(cfg.Ports)
	snap, _ := state.LoadSnapshot(cfg.StateFile)

	var notifiers []alert.Alerter
	notifiers = append(notifiers, alert.Logger(log.Writer()))
	for _, wh := range cfg.Webhooks {
		notifiers = append(notifiers, notify.NewWebhook(wh.URL, wh.Headers, wh.TimeoutSec))
	}

	mon := monitor.New(sc, rs, snap)

	return &Daemon{
		cfg:    cfg,
		mon:    mon,
		alerts: alert.Multi(notifiers...),
		ticker: time.NewTicker(time.Duration(cfg.IntervalSec) * time.Second),
	}, nil
}

// Run starts the daemon loop, blocking until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch: starting, interval=%ds", d.cfg.IntervalSec)
	for {
		select {
		case <-ctx.Done():
			d.ticker.Stop()
			log.Println("portwatch: shutting down")
			return ctx.Err()
		case <-d.ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("portwatch: tick error: %v", err)
			}
		}
	}
}

func (d *Daemon) tick() error {
	events, snap, err := d.mon.Scan()
	if err != nil {
		return err
	}
	for _, ev := range events {
		if err := d.alerts.Send(ev); err != nil {
			log.Printf("portwatch: alert error: %v", err)
		}
	}
	if len(events) > 0 {
		return state.SaveSnapshot(d.cfg.StateFile, snap)
	}
	return nil
}
