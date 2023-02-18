package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"cloud-disk/internal/auth"
	"cloud-disk/internal/constants"
)

var AppCfg *CloudDiskConfig

func InitConfig() error {
	confFilePath := filepath.Join(constants.RootPath, constants.ConfigFilePath)
	cloudDiskCfg, err := parseConfig(confFilePath)
	if err != nil {
		return err
	}
	AppCfg = cloudDiskCfg
	auth.Auth.SecretKey = []byte(AppCfg.AuthCfg.SecretKey)
	return nil
}

func parseConfig(confFilePath string) (*CloudDiskConfig, error) {
	file, err := os.Open(confFilePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	decoder := yaml.NewDecoder(file)
	cfg := new(CloudDiskConfig)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
