package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type SignTransactionRequest struct {
	DeviceID string `json:"device_id" validate:"required"`
	Data string `json:"data" validate:"required"`
}

type SignatureResponse struct {
	Signature string `json:"signature"`
	SignedData string `json:"signed_data"`
}

func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var req SignTransactionRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"Invalid JSON format",
		})
		return
	}

	// Validate the request
	if validationErrors := validateRequest(req); validationErrors != nil {
		WriteErrorResponse(response, http.StatusBadRequest, validationErrors)
		return
	}

	device, err := s.DeviceRepository.GetDevice(req.DeviceID)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			err.Error(),
		})
		return
	}

	//For locking per device to avoid race conditions when incrementing the signature counter
	deviceMutex := s.DeviceRepository.GetDeviceMutex(req.DeviceID)
	deviceMutex.Lock()
	defer deviceMutex.Unlock()

	signatureResponse, signatureCounter, err := s.signData(device, req.Data)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			err.Error(),
		})
		return
	}

	// Save the signature to the repository
	signatureRecord := &domain.Signature{
		ID:               uuid.New().String(),
		DeviceID:         device.ID,
		SignatureCounter: signatureCounter,
		SignatureValue:   signatureResponse.Signature,
	}
	err = s.SignatureRepository.CreateSignature(signatureRecord)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			err.Error(),
		})
		return
	}

	err = s.DeviceRepository.IncrementSignatureCounter(device.ID)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			err.Error(),
		})
		return
	}

	WriteAPIResponse(response, http.StatusOK, signatureResponse)
}

func (s *Server) ShowAllSignaturesByDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	deviceID := request.URL.Query().Get("device_id")
	if deviceID == "" {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"device_id is required",
		})
		return
	}
	
	signatures, err := s.SignatureRepository.GetAllSignaturesByDeviceID(deviceID)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			err.Error(),
		})
		return
	}

	WriteAPIResponse(response, http.StatusOK, wrapSignatureListResponse(signatures))
}

func (s *Server) signData(device *domain.Device, data string) (SignatureResponse, int, error) {
	var signer crypto.Signer
	var err error

	// Create the appropriate signer based on algorithm
	switch device.Algorithm {
	case "RSA":
		marshaler := crypto.NewRSAMarshaler()
		keyPair, err := marshaler.Unmarshal([]byte(device.PrivateKey))
		if err != nil {
			return SignatureResponse{}, 0, err
		}
		signer = crypto.NewRSASigner(keyPair)
	case "ECC":
		marshaler := crypto.NewECCMarshaler()
		keyPair, err := marshaler.Decode([]byte(device.PrivateKey))
		if err != nil {
			return SignatureResponse{}, 0, err
		}
		signer = crypto.NewECCSigner(keyPair)
	default:
		return SignatureResponse{}, 0, fmt.Errorf("unsupported algorithm: %s", device.Algorithm)
	}

	// Build the raw string format
	var rawStringFormat string
	var signatureCounter int

	if device.SignatureCounter == 0 {
		signatureCounter = 0
		rawStringFormat = fmt.Sprintf("0_%s_%s", data, base64.StdEncoding.EncodeToString([]byte(device.ID)))
	} else {
		latestSignature, err := s.SignatureRepository.GetLatestSignature(device.ID)
		if err != nil {
			return SignatureResponse{}, 0, err
		}
		signatureCounter = latestSignature.SignatureCounter + 1
		rawStringFormat = fmt.Sprintf("%d_%s_%s", latestSignature.SignatureCounter+1, data, latestSignature.SignatureValue)
	}

	// Sign the data
	signature, err := signer.Sign([]byte(rawStringFormat))
	if err != nil {
		return SignatureResponse{}, 0, err
	}

	// Create the response
	signatureResponse := SignatureResponse{
		Signature:  base64.StdEncoding.EncodeToString(signature),
		SignedData: rawStringFormat,
	}

	return signatureResponse, signatureCounter, nil
}

type GetSignatureResponse struct {
	ID string `json:"id"`
	DeviceID string `json:"device_id"`
	SignatureCounter int `json:"signature_counter"`
	SignatureValue string `json:"signature_value"`
}

func wrapSignatureListResponse(signatures []*domain.Signature) []GetSignatureResponse {
	signatureResponses := make([]GetSignatureResponse, 0, len(signatures))
	for _, signature := range signatures {
		signatureResponses = append(signatureResponses, GetSignatureResponse{
			ID: signature.ID,
			DeviceID: signature.DeviceID,
			SignatureCounter: signature.SignatureCounter,
			SignatureValue: signature.SignatureValue,
		})
	}
	return signatureResponses
}