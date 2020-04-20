package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AWS struct {
		Profile              string   `json:"profile"`
		Region               string   `json:"region"`
		CrossAccountRoleName string   `json:"cross-account-role-name"`
		Accounts             []string `json:"accounts"`
	} `json:"aws"`
}

type flags struct {
	Profile              string
	CrossAccountRoleName string
	Region               string
	Accounts             string
	ConfigPath           string
}

func missingArgError(arg ...string) error {
	flag.Usage()
	return fmt.Errorf("[error]: the following arguments are required: %v", arg)
}

func validateArgsAndFlags(cfg *Config, cliFlags flags) error {
	missingFlags := []string{}

	if len(flag.Args()) == 0 {
		missingFlags = append(missingFlags, "command")
	}

	if cfg.AWS.Accounts == nil || len(cfg.AWS.Accounts) == 0 {
		missingFlags = append(missingFlags, "accounts")
	}
	if cfg.AWS.CrossAccountRoleName == "" {
		missingFlags = append(missingFlags, "cross-account-role-name")
	}

	if len(missingFlags) != 0 {
		return missingArgError(missingFlags...)
	}

	return nil
}

func setAWSEnvVars(cfg *Config, flags *flags) {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	os.Setenv("AWS_PROFILE", cfg.AWS.Profile)
	os.Setenv("AWS_REGION", cfg.AWS.Region)
}

func fillConfigFromCLIFlags(cfg *Config, flags *flags) {
	if flags.Profile != "" {
		cfg.AWS.Profile = flags.Profile
	}

	if flags.Region != "" {
		cfg.AWS.Region = flags.Region
	}

	if flags.CrossAccountRoleName != "" {
		cfg.AWS.CrossAccountRoleName = flags.CrossAccountRoleName
	}

	if flags.Accounts != "" {
		// validate accounts flag input
		flags.Accounts = strings.Trim(flags.Accounts, ",")
		accountsList := strings.Split(flags.Accounts, ",")
		cfg.AWS.Accounts = accountsList
	}
}

func fillConfigFromEnvVars(cfg *Config) {
	envAWSProfile := os.Getenv("AWS_PROFILE")
	if envAWSProfile != "" {
		cfg.AWS.Profile = envAWSProfile
	}

	envAWSRegion := os.Getenv("AWS_REGION")
	if envAWSRegion != "" {
		cfg.AWS.Region = envAWSRegion
	}
}

func fillConfigFromFile(cfg *Config, configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err = decoder.Decode(cfg); err != nil {
		return err
	}

	return nil
}

func ConfigureApp() (*Config, error) {
	var cliFlags flags
	flag.StringVar(&cliFlags.Profile, "profile", "", "Use a specific profile from your credential file.")
	flag.StringVar(&cliFlags.CrossAccountRoleName, "cross-account-role-name", "", "Name of cross-account role used to authenticate to AWS accounts.")
	flag.StringVar(&cliFlags.Region, "region", "", "Set AWS region.")
	flag.StringVar(&cliFlags.Accounts, "accounts", "", "Comma-dilimted list of AWS accounts id to run command in.")
	flag.StringVar(&cliFlags.ConfigPath, "config", "", "Path to aws-bulk-cli configuration file.")
	flag.Parse()

	// config precedence is in this order:
	//   - cli flags
	//   - env vars
	//   - config file
	var cfg Config
	if cliFlags.ConfigPath != "" {
		if err := fillConfigFromFile(&cfg, cliFlags.ConfigPath); err != nil {
			return nil, err
		}
	}

	fillConfigFromEnvVars(&cfg)
	fillConfigFromCLIFlags(&cfg, &cliFlags)

	setAWSEnvVars(&cfg, &cliFlags)

	if err := validateArgsAndFlags(&cfg, cliFlags); err != nil {
		return nil, err
	}

	return &cfg, nil
}
