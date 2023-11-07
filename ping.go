package routerosv7_restfull_api

import (
	"errors"
	"fmt"
	"github.com/go-ping/ping"
)

type Pinger interface {
	Run() error
	Statistics() *ping.Statistics
}

type PingManager struct {
	pinger Pinger
}

func NewPing(host string) *PingManager {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Printf("Error creating pinger: %v\n", err)
		return nil
	}

	pinger.Count = pingCount
	pinger.Timeout = pingTimeout
	pinger.Interval = pingInterval

	return &PingManager{pinger: pinger}
}

func (pm *PingManager) CheckAvailableDevice() error {
	pinger := pm.pinger

	err := pinger.Run()
	if err != nil {
		return err
	}

	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return errors.New("device is not available")
	}

	return nil
}
