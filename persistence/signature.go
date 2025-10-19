package persistence

import (
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type ISignatureRepository interface {
	CreateSignature(signature *domain.Signature) error
	GetLatestSignature(deviceID string) (*domain.Signature, error)
	GetAllSignatures() ([]*domain.Signature, error)
}

type SignatureRepository struct {
	mutex sync.RWMutex
	signatures map[string]*domain.Signature
}

func NewSignatureRepository() ISignatureRepository {
	return &SignatureRepository{
		mutex:      sync.RWMutex{},
		signatures: make(map[string]*domain.Signature),
	}
}

func (s *SignatureRepository) CreateSignature(signature *domain.Signature) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.signatures[signature.ID] = signature
	return nil
}

func (s *SignatureRepository) GetLatestSignature(deviceID string) (latestSignature *domain.Signature, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	maxCounter := -1

	// Iterate through all signatures to find the latest one for this device
	for _, signature := range s.signatures {
		if signature.DeviceID == deviceID && signature.SignatureCounter > maxCounter {
			latestSignature = signature
			maxCounter = signature.SignatureCounter
		}
	}

	if latestSignature == nil {
		return nil, fmt.Errorf("no signatures found for device %s", deviceID)
	}

	return latestSignature, nil
}

func (s *SignatureRepository) GetAllSignatures() ([]*domain.Signature, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	signatures := make([]*domain.Signature, 0, len(s.signatures))
	for _, signature := range s.signatures {
		signatures = append(signatures, signature)
	}

	return signatures, nil
}