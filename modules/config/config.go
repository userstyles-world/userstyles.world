// Package config provides configuration options.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	OpenAuthConfig struct {
		GitHubID       string
		GitHubSecret   string
		GitLabID       string
		GitLabSecret   string
		CodebergID     string
		CodebergSecret string
	}

	SecretsConfig struct {
		PasswordCost           int
		ScrambleStepSize       int
		ScrambleBytesPerInsert int

		SessionTokenKey  string
		RecoverTokenKey  string
		ProviderTokenKey string

		CryptoKey        string
		StatsKey         string
		OAuthClientKey   string
		OAuthProviderKey string
		EmailAddress     string
		EmailPassword    string
		EmailServer      string
	}

	StorageConfig struct {
		DataDir   string
		CacheDir  string
		ImageDir  string
		StyleDir  string
		ProxyDir  string
		PublicDir string
		LogFile   string
	}

	config struct {
		App      AppConfig
		Database DatabaseConfig
		OpenAuth OpenAuthConfig
		Secrets  SecretsConfig
		Storage  StorageConfig
	}
)

var (
	// App stores general configuration.
	App *AppConfig

	// Database stores database configuration.
	Database *DatabaseConfig

	// OpenAuth stores configuration needed for connecting to external services.
	OpenAuth *OpenAuthConfig

	// Secrets stores cryptographic keys and related configuration.
	Secrets *SecretsConfig

	// Storage stores paths to directories and files.
	Storage *StorageConfig
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
			BuildSignature: GitSignature,
		},
		Database: DatabaseConfig{
			Name:         "dev.db",
			Debug:        "info",
			Colorful:     true,
			MaxOpenConns: 10,
		},
		Secrets: SecretsConfig{
			PasswordCost:           10,
			ScrambleStepSize:       2,
			ScrambleBytesPerInsert: 3,
			SessionTokenKey:        "ABigSecretPassword",
			RecoverTokenKey:        "OhNoWeCantUseTheSameAsJWTBeCaUseSeCuRiTy1337",
			ProviderTokenKey:       "ImNotACatButILikeUnicorns",
			CryptoKey:              "ABigSecretPasswordWhichIsExact32",
			StatsKey:               "KeyUsedForHashingStatistics",
			OAuthClientKey:         "AnotherStringLetstrySomethΦΦΦ",
			OAuthProviderKey:       "(✿◠‿◠＾◡＾)っ✂❤",
		},
		Storage: StorageConfig{
			DataDir:   "data",
			CacheDir:  "data/cache",
			ImageDir:  "data/images",
			StyleDir:  "data/styles",
			ProxyDir:  "data/proxy",
			PublicDir: "data/public",
			LogFile:   "data/userstyles.log",
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
	OpenAuth = &c.OpenAuth
	Secrets = &c.Secrets
	Storage = &c.Storage

	return nil
}

var (
	GoVersion    string
	GitCommit    string
	GitSignature string

	PerformanceMonitor   = getEnvBool("PERFORMANCE_MONITOR", false)
	ProxyMonitor         = getEnv("PROXY_MONITOR", "unset")
	SearchReindex        = getEnvBool("SEARCH_REINDEX", false)

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
