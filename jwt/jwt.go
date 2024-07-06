package jwt

import (
	"encoding/json"
	"strings"
	"time"

	utils "github.com/theveterandev/htmx-go-template/utilities"
)

type Claim struct {
	Username  string
	ExpiresAt int64
	IssuedAt  int64
	Issuer    string
	Roles     []string
}

func GenerateToken(username string) (string, error) {
	claim := &Claim{
		Username:  username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "htmx-go-template",
		Roles:     []string{"admin"},
	}

	claimJson, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}
	encodedClaim := utils.EncodeBase64(claimJson)
	encodedSignature := utils.CreateHMAC(encodedClaim)
	token := strings.Join([]string{encodedClaim, encodedSignature}, ".")
	return token, nil
}

func ValidateToken(token string) *Claim {
	splitToken := strings.Split(token, ".")
	if len(splitToken) != 2 {
		return nil
	}

	encodedClaim, encodedSignature := splitToken[0], splitToken[1]

	expectedSignature := utils.CreateHMAC(encodedClaim)

	if expectedSignature != encodedSignature {
		return nil
	}

	decodedClaim, err := utils.DecodeBase64(splitToken[0])
	if err != nil {
		return nil
	}

	var claim Claim
	json.Unmarshal(decodedClaim, &claim)

	if time.Now().Unix() > claim.ExpiresAt {
		return nil
	}

	return &claim
}
