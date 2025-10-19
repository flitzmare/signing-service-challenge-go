package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

type CreateDeviceRequest struct {
    Algorithm string `json:"algorithm" validate:"required,oneof=RSA ECC"`
    Label     string `json:"label"`
}

type CreateDeviceResponse struct {
    ID               string `json:"id"`
    Algorithm        string `json:"algorithm"`
    PublicKey        string `json:"public_key"`
    SignatureCounter int    `json:"signature_counter"`
    Label           string `json:"label"`
}

// TODO: REST endpoints ...
func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var req CreateDeviceRequest
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

	device := domain.Device{
		ID: uuid.New().String(),
		Algorithm: req.Algorithm,
		SignatureCounter: 0,
		Label: req.Label,
	}

	switch req.Algorithm {
		case "RSA":
			rsa := crypto.RSAGenerator{}
			keyPair, err := rsa.Generate()
			if err != nil {
				WriteErrorResponse(response, http.StatusInternalServerError, []string{
					err.Error(),
				})
				return
			}

			rsaMarshaler := crypto.NewRSAMarshaler()
			public, private, err := rsaMarshaler.Marshal(*keyPair)
			if err != nil {
				WriteErrorResponse(response, http.StatusInternalServerError, []string{
					err.Error(),
				})
				return
			}

			device.PublicKey = string(public)
			device.PrivateKey = string(private)
		case "ECC":
			ecc := crypto.ECCGenerator{}
			keyPair, err := ecc.Generate()
			if err != nil {
				WriteErrorResponse(response, http.StatusInternalServerError, []string{
					err.Error(),
				})
				return
			}

			eccMarshaler := crypto.NewECCMarshaler()
			public, private, err := eccMarshaler.Encode(*keyPair)
			if err != nil {
				WriteErrorResponse(response, http.StatusInternalServerError, []string{
					err.Error(),
				})
				return
			}

			device.PublicKey = string(public)
			device.PrivateKey = string(private)
		default:
			WriteErrorResponse(response, http.StatusBadRequest, []string{
				"Invalid algorithm",
			})
			return
		}

	err := s.DeviceRepository.CreateDevice(&device)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			err.Error(),
		})
		return
	}

	WriteAPIResponse(response, http.StatusCreated, DeviceResponse(&device))
}

func DeviceResponse(device *domain.Device) CreateDeviceResponse {
	return CreateDeviceResponse{
		ID: device.ID,
		Algorithm: device.Algorithm,
		PublicKey: device.PublicKey,
		SignatureCounter: device.SignatureCounter,
		Label: device.Label,
	}
}