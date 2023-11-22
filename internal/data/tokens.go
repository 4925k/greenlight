package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"github.com/4925k/greenlight/internal/validator"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

type TokenModel struct {
	DB *sql.DB
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Encode the byte slice to a base-32-encoded string and assign it to the token
	//Plaintext field. This will be the token string that we send to the user in their
	//welcome email. They will look similar to this:
	//
	// Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
	//
	// Note that by default base-32 strings may be padded at the end with the =
	//character. We don't need this padding character for the purpose of our tokens, so
	//we use the WithPadding(base32.NoPadding) method in the line below to omit
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate an SHA-256 hash of the plaintext token string. This will be the value
	//that we store in the `hash` field of our database table. Note that the
	//sha256.Sum256() function returns an *array* of length 32, so to make it easier to
	//work with we convert it to a slice using the [:] operator before storing i
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, plaintext string) {
	v.Check(plaintext != "", "token", "must be provided")
	// It’s important to point out that the plaintext token strings we’re creating here like
	// Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
	//are not 16 characters long — but rather they have an underlying entropy of 16 bytes of randomness.
	//The length of the plaintext token string itself depends on how those 16 random bytes are encoded to create a string.
	//In our case we encode the random bytes to a base-32 string,which results in a string with 26 characters.
	//In contrast, if we encoded the random bytes using hexadecimal (base-16) the string would be 32 characters long instead.
	v.Check(len(plaintext) == 26, "token", "must be 26 bytes long")
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
				VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}