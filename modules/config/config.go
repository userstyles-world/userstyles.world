package config

import (
	"fmt"
	"os"
	"strconv"
)

type ScrambleSettings struct {
	StepSize       int
	BytesPerInsert int
}

var (
	Port                 = getEnv("PORT", ":3000")
	DB                   = getEnv("DB", "dev.db")
	DBDebug              = getEnv("DB_DEBUG", "silent")
	DBColor              = getEnv("DB_COLOR", "false")
	DBDrop               = getEnv("DB_DROP", "false")
	DBRandomData         = getEnv("DB_RANDOM_DATA", "false")
	Salt                 = getEnv("SALT", "10")
	JWTSigningKey        = getEnv("JWT_SIGNING_KEY", "ABigSecretPassword")
	VerifyJWTSigningKey  = getEnv("VERIFY_JWT_SIGNING_KEY", "OhNoWeCantUseTheSameAsJWTBeCaUseSeCuRiTy1337")
	OAuthpJWTSigningKey  = getEnv("OAUTHP_JWT_SIGNING_KEY", "ImNotACatButILikeUnicorns")
	CryptoKey            = getEnv("CRYPTO_KEY", "ABigSecretPasswordWhichIsExact32")
	StatsKey             = getEnv("STATS_KEY", "KeyUsedForHashingStatistics")
	OAuthKey             = getEnv("OAUTH_KEY", "AnotherStringLetstrySomethΦΦΦ")
	OAuthpKey            = getEnv("OAUTHP_KEY", "(✿◠‿◠＾◡＾)っ✂❤")
	EmailAddress         = getEnv("EMAIL_ADDRESS", "test@userstyles.world")
	EmailPassword        = getEnv("EMAIL_PWD", "hahah_not_your_password")
	GitCommit            = getEnv("GIT_COMMIT", "unset")
	GitHubClientID       = getEnv("GITHUB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitHubClientSecret   = getEnv("GITHUB_CLIENT_SECRET", "OurSecretHere?_www.youtube.com/watch?v=dQw4w9WgXcQ")
	GitlabClientID       = getEnv("GITLAB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitlabClientSecret   = getEnv("GITLAB_CLIENT_SECRET", "www.youtube.com/watch?v=dQw4w9WgXcQ")
	CodebergClientID     = getEnv("CODEBERG_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	CodebergClientSecret = getEnv("CODEBERG_CLIENT_SECRET", "IMgettinggboredd")

	// Production is used for various "feature flags".
	Production = DB != "dev.db"

	ScrambleConfig = &ScrambleSettings{
		StepSize:       getEnvInt("NONCE_SCRAMBLE_STEP", 2),
		BytesPerInsert: getEnvInt("NONCE_SCRAMBLE_BYTES_PER_INSERT", 3),
	}
)

func getEnvInt(name string, defaultValue int) int {
	envValue := getEnv(name, "__NOT_FOUND!")
	if envValue == "__NOT_FOUND!" {
		return defaultValue
	}

	envInt, err := strconv.Atoi(envValue)
	if err != nil {
		return defaultValue
	}
	return envInt
}

func getEnv(name, fallback string) string {
	if val, set := os.LookupEnv(name); set {
		return val
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf(`Env variable not found: %v`, name))
}

// OAuthURL returns the proper callback URL depending on the environment.
func OAuthURL() string {
	if Production {
		return "https://userstyles.world/api/callback/"
	}

	return "http://localhost" + Port + "/api/callback/"
}
