package persistence

import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// TODO: in-memory persistence ...
type IDeviceRepository interface {
	CreateDevice(device *domain.Device) error
	CountDevices() int
	GetDevice(id string) (*domain.Device, error)
	IncrementSignatureCounter(deviceID string) error
	GetDeviceMutex(deviceID string) *sync.Mutex
}

type DeviceRepository struct {
	mutex sync.RWMutex
	devices map[string]*domain.Device
	devicesMutexes map[string]*sync.Mutex //For locking per device
}

func NewDeviceRepository() IDeviceRepository {
	return &DeviceRepository{
		mutex:   sync.RWMutex{},
		devices: make(map[string]*domain.Device),
		devicesMutexes: make(map[string]*sync.Mutex),
	}
}

func (m *DeviceRepository) CreateDevice(device *domain.Device) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.devices[device.ID] = device
	return nil
}

func (m *DeviceRepository) CountDevices() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.devices)
}

func (m *DeviceRepository) GetDevice(id string) (*domain.Device, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	device, exists := m.devices[id]
	if !exists {
		return nil, fmt.Errorf("device with id %s not found", id)
	}

	return device, nil
}

func (m *DeviceRepository) IncrementSignatureCounter(deviceID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	device, exists := m.devices[deviceID]
	if !exists {
		return fmt.Errorf("device with id %s not found", deviceID)
	}

	device.SignatureCounter++
	return nil
}

func (m *DeviceRepository) GetDeviceMutex(deviceID string) *sync.Mutex {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if m.devicesMutexes[deviceID] == nil {
		m.devicesMutexes[deviceID] = &sync.Mutex{}
	}
	return m.devicesMutexes[deviceID]
}