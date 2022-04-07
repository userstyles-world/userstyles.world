package utils

import (
	"crypto/cipher"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/chacha20poly1305"

	"userstyles.world/modules/config"
	"userstyles.world/modules/errors"
)

var (
	AEADCrypto       cipher.AEAD
	AEADOAuth        cipher.AEAD
	AEADOAuthp       cipher.AEAD
	VerifySigningKey = []byte(config.VerifyJWTSigningKey)
	OAuthPSigningKey = []byte(config.OAuthpJWTSigningKey)
	signingMethod    = "HS512"
)

func InitalizeCrypto() {
	var aead cipher.AEAD
	var err error

	aead, err = chacha20poly1305.NewX([]byte(config.CryptoKey))
	if err != nil {
		panic("Cannot create AEAD_CRYPTO chipher, due to " + err.Error())
	}
	AEADCrypto = aead

	aead, err = chacha20poly1305.NewX([]byte(config.OAuthKey))
	if err != nil {
		panic("Cannot create AEAD_OAUTH chipher, due to " + err.Error())
	}
	AEADOAuth = aead

	aead, err = chacha20poly1305.NewX([]byte(config.OAuthKey))
	if err != nil {
		panic("Cannot create AEAD_OAUTHP chipher, due to " + err.Error())
	}
	AEADOAuthp = aead
}

func sealText(text string, aead cipher.AEAD, nonceScrambling *config.ScrambleSettings) []byte {
	nonce := RandomBytes(aead.NonceSize())

	dest := aead.Seal(nil, nonce, UnsafeBytes(text), nil)
	return scrambleNonce(nonce, dest, nonceScrambling.StepSize, nonceScrambling.BytesPerInsert)
}

// scrambleNonce into string takes a nonce and a text
// And it will insert `bytesPerInsert`` bits of the nonce every `step` bytes.
// And it will paste the "rest" nonce if the text isn't big enough.
func scrambleNonce(nonce, text []byte, step, bytesPerInsert int) []byte {
	// Copy the text into encodedText.
	encodedText := make([]byte, 0, len(text)+len(nonce))
	encodedText = append(encodedText, text...)

	// Copy the nonce into restNonce.
	restNonce := make([]byte, len(nonce))
	copy(restNonce, nonce)

	var currentNonce byte

	// Loop trough the text every x `steps`.
mainLoop:
	for i := 0; i < len(encodedText); i += (step + bytesPerInsert) {
		// Insert `bytesPerInsert` bytes of nonce at the index.
		for j := 0; j < bytesPerInsert; j++ {
			if len(restNonce) == 0 {
				break mainLoop
			}

			// Create some space at the index we want to append the nonce to.
			encodedText = append(encodedText, 0)
			copy(encodedText[i+j+1:], encodedText[i+j:])

			// One-liner for front-shift
			currentNonce, restNonce = restNonce[0], restNonce[1:]

			// Add the current nonce to the text at the correct index.
			encodedText[i+j] = currentNonce
		}

		if len(restNonce) == 0 {
			break
		}
	}

	// Append the left overs from the nonce at the end of the text.
	encodedText = append(encodedText, restNonce...)
	return encodedText
}

// descrambleNonce will take a text, bytesPerInsert, step and the length of the nonce.
// And it will return the nonce and descrambled text.
func descrambleNonce(scrambledText []byte, nonceSize, step, bytesPerInsert int) ([]byte, []byte, error) {
	// Store the length of the scrambledText.
	textLen := len(scrambledText)

	if nonceSize >= textLen {
		return nil, nil, errors.TexTooShort(nonceSize, textLen)
	}

	// Because we don't want to modify the originial text, we copy it.
	text := make([]byte, textLen)
	copy(text, scrambledText)

	// nonce will be the nonce that was scrambled in `text`.
	// We already know the size of it. So we can allocate it right away.
	nonce := make([]byte, 0, nonceSize)

	// We need to store the amount of found bytes.
	// So we know when to stop looking for more bytes.
	foundBytes := 0

	var currentNonce byte

	// Because we can calculate which bytes where being appened to the text.
	// We can use this to to get the "appended nonce" if there is any.
	// So it won't interfere with the next step.
	var appendedNonce []byte

	// First we get the originial text length.
	originialTextLen := textLen - nonceSize

	// Then we calculate how many places that text had "places" to add bytes to.
	// We also add bytesPerInsert, because this calculation to don't take in account.
	// The very first few bytes in text are part of the scrambled nonce.
	placesToInsertByte := originialTextLen/step + bytesPerInsert

	// Now we know how many places the text had to append bytes to.
	// We check if the amount of places * amount of bytes per insert is larger.
	// That means it wouldn't have enough places to append all bytes at every x step.
	// And appended some bytes after the text.
	if (placesToInsertByte * bytesPerInsert) < nonceSize {
		// Now we calculate from where the specific nonce should be found.
		amountOfAppendedBytes := textLen - originialTextLen - (placesToInsertByte * bytesPerInsert)

		// Now we cut the appended nonce bytes from the text.
		// So it won't interfere with the next step.
		appendedNonce, text = text[textLen-amountOfAppendedBytes:], text[:textLen-amountOfAppendedBytes]

		// Also make sure to add the amount of appended bytes to the found bytes!
		foundBytes += amountOfAppendedBytes
	}

mainLoop:
	// Mikey don't you dare to optimize this loop!
	// The size of scrambledText changes after evert itteration.
	for i := 0; i < len(text); i += step {
		for j := 0; j < bytesPerInsert; j++ {
			// if we found the amount of bytes we need, we can break the loop.
			// i >= len(text) prevents that we will access a index that doesn't exist.
			if foundBytes == nonceSize || i >= len(text) {
				break mainLoop
			}

			// Also you would think we should use `i+j`, but like I SAID MIKEY.
			// The text will be changed all the time and the previous nonce will be already
			// gone and thus the `i` index will have the new nonce byte.

			currentNonce = text[i]

			// Cut the nonce from the text.
			text = append(text[:i], text[i+1:]...)

			// Append the correct nonce to the current nonce.
			nonce = append(nonce, currentNonce)
			foundBytes++
		}
	}
	nonce = append(nonce, appendedNonce...)

	return nonce, text, nil
}

func openText(encryptedMsg string, aead cipher.AEAD, nonceScrambling *config.ScrambleSettings) ([]byte, error) {
	if len(encryptedMsg) < aead.NonceSize() {
		return nil, errors.ErrMessageSmall
	}

	// Split nonce and ciphertext.
	nonce, ciphertext, err := descrambleNonce(UnsafeBytes(encryptedMsg), aead.NonceSize(),
		nonceScrambling.StepSize, nonceScrambling.BytesPerInsert)
	if err != nil {
		return nil, err
	}
	// Decrypt the message and check it wasn't tampered with.
	return aead.Open(nil, nonce, ciphertext, nil)
}

func VerifyJwtKeyFunction(t *jwt.Token) (any, error) {
	if t.Method.Alg() != signingMethod {
		return nil, errors.UnexpectedSigningMethod(t.Method.Alg())
	}
	return VerifySigningKey, nil
}

func OAuthPJwtKeyFunction(t *jwt.Token) (any, error) {
	if t.Method.Alg() != signingMethod {
		return nil, errors.UnexpectedSigningMethod(t.Method.Alg())
	}
	return OAuthPSigningKey, nil
}

func EncryptText(text string, aead cipher.AEAD, settings *config.ScrambleSettings) string {
	// We have to prepare the encrypted text for transport
	// Seal Text -> Base64(URL Version)
	sealedText := sealText(text, aead, settings)

	return EncodeToString(sealedText)
}

func DecryptText(preparedText string, aead cipher.AEAD, settings *config.ScrambleSettings) (string, error) {
	// Now we have to reverse the process.
	// Decode Base64(URL version) -> Unseal Text
	enryptedText, err := decodeBase64(preparedText)
	if err != nil {
		return "", err
	}

	decryptedText, err := openText(UnsafeString(enryptedText), aead, settings)
	if err != nil {
		return "", err
	}

	return UnsafeString(decryptedText), nil
}
