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

			req := httptest.NewRequest("POST", "/api/v0/device", strings.NewReader(`{"algorithm": "RSA", "label": "test-device"}`))
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
MIGJAoGBAM/tvE/dja6Y8T8TbYSZHpve3ytzv1yiDwhVlF7avZRdiFRU7srNkaRR
8r746jm/VYYA5rLftyhteEHzZHZgXKHjS+ehavTAtFe4BEcUsk7PudebgD+cFC4E
F9Sa+aRvyTn0Rg3NFtf9s+MiixfdkDfybuqQ8lN+SqK7uOMqpnFJAgMBAAE=
-----END RSA_PUBLIC_KEY-----`,
				PrivateKey:       `-----BEGIN RSA_PRIVATE_KEY-----
MIICXQIBAAKBgQDP7bxP3Y2umPE/E22EmR6b3t8rc79cog8IVZRe2r2UXYhUVO7K
zZGkUfK++Oo5v1WGAOay37cobXhB82R2YFyh40vnoWr0wLRXuARHFLJOz7nXm4A/
nBQuBBfUmvmkb8k59EYNzRbX/bPjIosX3ZA38m7qkPJTfkqiu7jjKqZxSQIDAQAB
AoGAAWJve0Iyo2Toi92DVCyf6hcr5lOhrAfZJGRfVdoMZ2v3F+MfWmQK82BOxsqb
NUbEpWbvDqEWQr6TjZJoKuvGG/bGMCX3yt3KNEh53IqaykKjQLfhZnu7zNRJwBUD
x2heTYxDaDiD9ZyZs1lqg4cpIuConKOCJTso02p30MtKzfUCQQDwWYcUaBn4Eua2
FvGsmTrhOhYdrFHqU9C56xfGMnlHLJs4behCa2n6/2O2Ouk5a6TGTM3rMWiCC+Vi
eNRGDUlrAkEA3XfFIq8QAbgBu9oN+CviZy5R87e8p26JXwtDOAvrIO38fUfpHkLu
999C+d0ZN37wnQWwRkIWCYYbfVnC+QzZGwJBAIdrGeWQhdk05RKBOOdzai5OKPnN
BlZNpROrdriv5Y8Jfec8XZlWpd7KmCaraI52rN8hlP/H1cc35qUlyQwzHkMCQQC2
Z0PVSiw7ziqXZoPk53gEFXFn8ueNWwwHXMZTLfXNXFV9dbG5u9UIEDkghAqV25Yf
LaU+aIWv+GVBu6FK8FsLAkAYlnxMV1l6AEKbHbGWgK38gL6+vQ6jk4Heo7DJS3fQ
pAQzbR1h7p39hZYB5AGCnslhokTBQiuGzUts9VVGqZ38
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