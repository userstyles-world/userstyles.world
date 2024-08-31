// Package config provides configuration options.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type (
	AppConfig struct {
		Debug      bool
		Production bool
		Addr       string
		BaseURL    string

		Name         string
		Description  string
		Codename     string
		Copyright    string
		Started      time.Time `json:"-"`
		EmailRe      string
		PageMaxItems int

		Repository     string
		Discord        string
		Matrix         string
		OpenCollective string

		BuildCommit    string `json:"-"`
		BuildCommitSHA string `json:"-"`
		BuildSignature string `json:"-"`
	}

	DatabaseConfig struct {
		Name             string
		Debug            string
		Colorful         bool
		Migrate          bool
		Drop             bool
		RandomData       bool
		RandomDataAmount int
		MaxOpenConns     int
	}

	config struct {
		App      AppConfig
		Database DatabaseConfig
	}
)

var (
	// App stores general configuration.
	App *AppConfig

	// Database stores database configuration.
	Database *DatabaseConfig
)

func defaultConfig() *config {
	repo := "https://github.com/userstyles-world/userstyles.world"
	started := time.Now()

	return &config{
		App: AppConfig{
			Addr:         ":3000",
			BaseURL:      "http://localhost:3000",
			Name:         "UserStyles.world",
			Description:  "A free and open-source, community-driven website for browsing and sharing UserCSS userstyles.",
			Codename:     "Fennec Fox",
			Copyright:    started.Format("2006"),
			Started:      started,
			EmailRe:      `^[a-zA-Z0-9.!#$%&’*+/=?^_\x60{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+$`,
			PageMaxItems: 36,

			Repository:     repo,
			Discord:        "https://discord.gg/P5zra4nFS2",
			Matrix:         "https://matrix.to/#/#userstyles:matrix.org",
			OpenCollective: "https://opencollective.com/userstyles",

			BuildCommit:    repo + "/commit/" + GitCommit,
			BuildCommitSHA: fmt.Sprintf("%.8s", GitCommit),
			BuildSignature: GitVersion,
		},
		Database: DatabaseConfig{
			Name:         "dev.db",
			Debug:        "info",
			Colorful:     true,
			MaxOpenConns: 10,
		},
	}
}

// Load tries to load configuration from a given path.
func Load(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	c := defaultConfig()
	if err = json.Unmarshal(b, &c); err != nil {
		return err
	}

	if c.App.Debug {
		b, err := json.MarshalIndent(c, "", "\t")
		if err != nil {
			return err
		}

		log.Println("config:", string(b))
	}

	App = &c.App
	Database = &c.Database

	return nil
}

type ScrambleSettings struct {
	StepSize       int
	BytesPerInsert int
}

var (
	GitCommit  string
	GitVersion string
	Salt                 = getEnvInt("SALT", 10)
	JWTSigningKey        = getEnv("JWT_SIGNING_KEY", "ABigSecretPassword")
	VerifyJWTSigningKey  = getEnv("VERIFY_JWT_SIGNING_KEY", "OhNoWeCantUseTheSameAsJWTBeCaUseSeCuRiTy1337")
	OAuthpJWTSigningKey  = getEnv("OAUTHP_JWT_SIGNING_KEY", "ImNotACatButILikeUnicorns")
	CryptoKey            = getEnv("CRYPTO_KEY", "ABigSecretPasswordWhichIsExact32")
	StatsKey             = getEnv("STATS_KEY", "KeyUsedForHashingStatistics")
	OAuthKey             = getEnv("OAUTH_KEY", "AnotherStringLetstrySomethΦΦΦ")
	OAuthpKey            = getEnv("OAUTHP_KEY", "(✿◠‿◠＾◡＾)っ✂❤")
	EmailAddress         = getEnv("EMAIL_ADDRESS", "test@userstyles.world")
	EmailPassword        = getEnv("EMAIL_PWD", "hahah_not_your_password")
	GitHubClientID       = getEnv("GITHUB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitHubClientSecret   = getEnv("GITHUB_CLIENT_SECRET", "OurSecretHere?_www.youtube.com/watch?v=dQw4w9WgXcQ")
	GitlabClientID       = getEnv("GITLAB_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	GitlabClientSecret   = getEnv("GITLAB_CLIENT_SECRET", "www.youtube.com/watch?v=dQw4w9WgXcQ")
	CodebergClientID     = getEnv("CODEBERG_CLIENT_ID", "SOmeOneGiVeMeIdEaSwHaTtOpUtHeRe")
	CodebergClientSecret = getEnv("CODEBERG_CLIENT_SECRET", "IMgettinggboredd")
	PerformanceMonitor   = getEnvBool("PERFORMANCE_MONITOR", false)
	IMAPServer           = getEnv("IMAP_SERVER", "mail.userstyles.world:587")
	ProxyMonitor         = getEnv("PROXY_MONITOR", "unset")
	SearchReindex        = getEnvBool("SEARCH_REINDEX", false)

	ScrambleConfig = &ScrambleSettings{
		StepSize:       getEnvInt("NONCE_SCRAMBLE_STEP", 2),
		BytesPerInsert: getEnvInt("NONCE_SCRAMBLE_BYTES_PER_INSERT", 3),
	}

	DataDir   = path.Join(getEnv("DATA_DIR", "data"))
	CacheDir  = path.Join(DataDir, "cache")
	ImageDir  = path.Join(DataDir, "images")
	StyleDir  = path.Join(DataDir, "styles")
	ProxyDir  = path.Join(DataDir, "proxy")
	PublicDir = path.Join(DataDir, "public")

	LogFile = path.Join(DataDir, "userstyles.log")

	CachedCodeItems = getEnvInt("CACHED_CODE_ITEMS", 250)
	ProxyRealIP     = getEnv("PROXY_REAL_IP", "")
)

// OAuthURL returns an environment-specific callback URL used for OAuth services.
func OAuthURL() string {
	return App.BaseURL + "/api/callback/"
}

// raw tweaks allowed URLs to make them work seamlessly in both environments.
func raw(s string) string {
	if !App.Production {
		s += "|userstyles.world"
	}
	r := strings.NewReplacer("http://", "", "https://", "")
	return r.Replace(s)
}
