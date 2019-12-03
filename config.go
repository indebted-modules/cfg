package cfg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/rs/zerolog/log"
)

// SysMailFromName for internal notifications
const SysMailFromName = "Engineering"

// SysMailFromAddress for internal notifications
const SysMailFromAddress = "engineering@indebted.co"

// Sess session within AWS
var Sess = session.Must(session.NewSession())

// SecretValue gets secret value by ID
func SecretValue(id string) string {
	sm := secretsmanager.New(session.New(), aws.NewConfig())
	res, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("Failed getting secrets")
	}
	var s string
	if res.SecretString != nil {
		s = *res.SecretString
	} else {
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(res.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decoded, res.SecretBinary)
		if err != nil {
			log.
				Fatal().
				Err(err).
				Msg("Failed decoding secrets")
		}
		s = string(decoded[:len])
	}
	return s
}

// SecretJSON gets secret JSON by ID
func SecretJSON(id string) map[string]string {
	s := SecretValue(id)
	secrets := map[string]string{}
	err := json.Unmarshal([]byte(s), &secrets)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("Failed parsing secrets")
	}
	return secrets
}

// MustEnv gets environment variable or panics if not present
func MustEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.
			Fatal().
			Str("var", key).
			Msg("Missing environment variable")
	}
	return val
}

// DatabaseURL from env var DB_URL if present, or SecretsManager/db-secrets
func DatabaseURL() string {
	fromEnv, ok := os.LookupEnv("DB_URL")
	if ok {
		return fromEnv
	}
	secrets := SecretJSON("db-secrets")
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		secrets["username"],
		secrets["password"],
		secrets[MustEnv("DB_ENDPOINT")],
		MustEnv("SVC_NAME"),
	)
}
