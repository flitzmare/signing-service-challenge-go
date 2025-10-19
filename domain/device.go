package domain

// TODO: signature device domain model ...
type Device struct {
	ID string
	Algorithm string
	PublicKey  string
	PrivateKey  string
	SignatureCounter int
	Label string
}