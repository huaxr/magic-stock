package gokms

// Public Body
type PublicForm struct {
	SecretId  string `json:"secretid"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
	PSM       string `json:"psm"`
	Token     string `json:"token"`
}

// Body interface
type ParamForm interface {
}

// Encrypt Form
type EncryptForm struct {
	PublicForm
	Keyid     string `json:"keyid"`
	Plaintext string `json:"plaintext"`
}

// Decrypt Form
type DecryptForm struct {
	PublicForm
	CiphertextBlob string `json:"ciphertextBlob"`
}

// Batch Decrypt Form
type BatchDecryptForm struct {
	PublicForm
	Ciphertexts []string `json:"ciphertextBlobs"`
}

// Create master key form
type NewKeyForm struct {
	PublicForm
	Alias       string `json:"alias"`
	Description string `json:"description"`
}

// List master key form
type ListKeyForm struct {
	PublicForm
}

// Get master key info form
type InfoKeyForm struct {
	PublicForm
	Keyid string `json:"keyid"`
}

// Generate master key form
type GenDataKeyForm struct {
	PublicForm
	Keyid string `json:"keyid"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	P     string `json:"p"`
}

// Master key info struct
type KeyInfo struct {
	Keyid       string
	Alias       string
	Description string
	Status      string
	Createdtime string
}

// Decrypt Data Key Form
type DecryptDataKeyForm struct {
	PublicForm
	CiphertextBlob string `json:"ciphertextBlob"`
}

// Decrypt Data Key By Id Form
type DecryptDataKeyByIdForm struct {
	PublicForm
	Keyid string `json:"keyid"`
}

//  Decrypt Data Key By Psm Form
type DecryptDataKeyByPSMForm struct {
	PublicForm
	P string `json:"p"`
}

// Share Data Key Form
type ShareDataKeyForm struct {
	PublicForm
	P           string `json:"p"`
	Description string `json:"description"`
}

// Share Master Key Form
type ShareMasterKeyForm struct {
	PublicForm
	Text     string `json:"text"`
	Username string `json:"username"`
}
