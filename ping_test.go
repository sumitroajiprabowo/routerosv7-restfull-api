package routerosv7_restfull_api

import (
	"errors"
	"github.com/go-ping/ping"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockPinger struct {
	runFunc        func() error
	statisticsFunc func() *ping.Statistics
}

func (mp *MockPinger) Run() error {
	return mp.runFunc()
}

func (mp *MockPinger) Statistics() *ping.Statistics {
	return mp.statisticsFunc()
}

func TestCheckAvailableDevice_DeviceAvailable(t *testing.T) {
	mockPinger := &MockPinger{
		runFunc: func() error {
			return nil
		},
		statisticsFunc: func() *ping.Statistics {
			return &ping.Statistics{PacketsRecv: 3}
		},
	}

	pingManager := PingManager{pinger: mockPinger}
	err := pingManager.CheckAvailableDevice()

	assert.NoError(t, err)
}

func TestCheckAvailableDevice_DeviceNotAvailable(t *testing.T) {
	mockPinger := &MockPinger{
		runFunc: func() error {
			return nil
		},
		statisticsFunc: func() *ping.Statistics {
			return &ping.Statistics{PacketsRecv: 0}
		},
	}

	pingManager := PingManager{pinger: mockPinger}
	err := pingManager.CheckAvailableDevice()

	assert.Error(t, err)
	assert.EqualError(t, err, "device is not available")
}

func TestCheckAvailableDevice_PingerError(t *testing.T) {
	mockPinger := &MockPinger{
		runFunc: func() error {
			return errors.New("ping error")
		},
		statisticsFunc: func() *ping.Statistics {
			return nil // Skenario tambahan: Statistik adalah nil ketika ada kesalahan
		},
	}

	pingManager := PingManager{pinger: mockPinger}
	err := pingManager.CheckAvailableDevice()

	assert.Error(t, err)
	assert.EqualError(t, err, "ping error")
}

func TestNewPing_Success(t *testing.T) {
	pingManager := NewPing("example.com")

	assert.NotNil(t, pingManager)
	assert.NotNil(t, pingManager.pinger)
}

func TestNewPing_PingerError(t *testing.T) {
	pingManager := NewPing("invalid host")

	assert.Nil(t, pingManager)
}
