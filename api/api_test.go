package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	mock_persistence "github.com/fiskaly/coding-challenges/signing-service-challenge/persistence/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAPISuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("Device Management", func() {
	var (
		ctrl *gomock.Controller
		mockDeviceRepository *mock_persistence.MockIDeviceRepository
		server *Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDeviceRepository = mock_persistence.NewMockIDeviceRepository(ctrl)
		server = &Server{
			DeviceRepository: mockDeviceRepository,
		}
	})

	Context("When creating a signature device", func() {
		It("should create a signature device", func() {
			mockDeviceRepository.EXPECT().CreateDevice(gomock.Any()).Return(nil)

			req := httptest.NewRequest("POST", "/api/v0/create-signature-device", strings.NewReader(`{"algorithm": "RSA", "label": "test-device"}`))
        	req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			server.CreateSignatureDevice(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})
	})
})

var _ = Describe("Transaction Signing", func() {
	var (
		ctrl *gomock.Controller
		mockDeviceRepository *mock_persistence.MockIDeviceRepository
		mockSignatureRepository *mock_persistence.MockISignatureRepository
		server *Server
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDeviceRepository = mock_persistence.NewMockIDeviceRepository(ctrl)
		mockSignatureRepository = mock_persistence.NewMockISignatureRepository(ctrl)
		server = &Server{
			DeviceRepository: mockDeviceRepository,
			SignatureRepository: mockSignatureRepository,
		}
	})

	Context("When signing a transaction", func() {
		It("should sign a transaction for RSA Algorithm", func() {
			// Create a mock device for the test
			mockDevice := &domain.Device{
				ID:               "test-device",
				Algorithm:        "RSA",
				PublicKey:        `-----BEGIN RSA_PUBLIC_KEY-----
MEgCQQDfWkuaBhDcpaMXUz4BJbthqp0HPyxQzyimYILXeItVoTO/hMMKS/GO+fKQ
j8V8jVirJFcyGA5JbUK3gX5LtPbFAgMBAAE=
-----END RSA_PUBLIC_KEY-----`,
				PrivateKey:       `-----BEGIN RSA_PRIVATE_KEY-----
MIIBOgIBAAJBAN9aS5oGENyloxdTPgElu2GqnQc/LFDPKKZggtd4i1WhM7+EwwpL
8Y758pCPxXyNWKskVzIYDkltQreBfku09sUCAwEAAQJAGTvmHTdwucED8UtwDqig
6DqipZIzU0joVo3CUo41rb2D1EpspD9LX48k6wzbKPSBz48TgPHp2h/dBgeYjNDz
UQIhAOLhvrgJT/KluKl33myhN+FGsqqOYTk39pM0MatuT9r1AiEA/ASZ1TeHOVst
nQCrztCk/NtyD6kr1OKoxUgL7L8+6pECIDRaqVroczVn/nPEwGPK1A089i+bSV4d
xt1zFt8bRnwdAiEAgpozyn5DUqMAyWtunfgceHmU667E60cnJU3H+EHH7jECIHMW
n7VMKQ2Z3tOM9tdxCIdYEBxeuWuFWvREyMTKbl8t
-----END RSA_PRIVATE_KEY-----`,
				SignatureCounter: 0,
				Label:           "test-device",
			}

			mutex := &sync.Mutex{}
			
			mockDeviceRepository.EXPECT().GetDeviceMutex(gomock.Any()).Return(mutex)
			mockDeviceRepository.EXPECT().GetDevice(gomock.Any()).Return(mockDevice, nil)
			mockDeviceRepository.EXPECT().IncrementSignatureCounter(gomock.Any()).Return(nil)
			mockSignatureRepository.EXPECT().CreateSignature(gomock.Any()).Return(nil)

			req := httptest.NewRequest("POST", "/api/v0/sign-transaction", strings.NewReader(`{"device_id": "test-device", "data": "test-data"}`))
        	req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			server.SignTransaction(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
		It("should sign a transaction for ECC Algorithm", func() {
			// Create a mock device for the test
			mockDevice := &domain.Device{
				ID:               "test-device",
				Algorithm:        "ECC",
				PublicKey:        `-----BEGIN PUBLIC_KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE6ZZszV54k4v8ssIlpj4dYFUo8SLBt11M
QR9EDRUeYaUYweqPG8j9pgBksx1wr9yIG7PjLkcg9dxPYst8zVQpk9ULCusPJ1d9
aOMrD5ANM8IoRAkYBUtJEirWGiCRk5/k
-----END PUBLIC_KEY-----`,
				PrivateKey:       `-----BEGIN PRIVATE_KEY-----
MIGkAgEBBDDjn7xR+VY2ST5b/WAZ5jO/tYik3vNANKdWSaYhggvJKolorpM0JcZu
Tqos5vIvuYqgBwYFK4EEACKhZANiAATplmzNXniTi/yywiWmPh1gVSjxIsG3XUxB
H0QNFR5hpRjB6o8byP2mAGSzHXCv3Igbs+MuRyD13E9iy3zNVCmT1QsK6w8nV31o
4ysPkA0zwihECRgFS0kSKtYaIJGTn+Q=
-----END PRIVATE_KEY-----`,
				SignatureCounter: 0,
				Label:            "test-device",
			}

			mutex := &sync.Mutex{}
			
			mockDeviceRepository.EXPECT().GetDeviceMutex(gomock.Any()).Return(mutex)
			mockDeviceRepository.EXPECT().GetDevice(gomock.Any()).Return(mockDevice, nil)
			mockDeviceRepository.EXPECT().IncrementSignatureCounter(gomock.Any()).Return(nil)
			mockSignatureRepository.EXPECT().CreateSignature(gomock.Any()).Return(nil)

			req := httptest.NewRequest("POST", "/api/v0/sign-transaction", strings.NewReader(`{"device_id": "test-device", "data": "test-data"}`))
        	req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			server.SignTransaction(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})
})