package workspace_repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const tfrcFileName = ".terraform.d/credentials.tfrc.json"

func GetTfeToken(tfeUrl string) (string, error) {
	token, ok := os.LookupEnv("TFE_TOKEN")
	if ok {
		return token, nil
	}

	hostName := tfeUrl
	if strings.Index(hostName, "http") == 0 {
		u, err := url.Parse(tfeUrl)
		if err != nil {
			log.Debugf("TFE URL %s is not valid: %s\n", tfeUrl, err.Error())
			return "", errors.Wrap(err, "please set env var TFE_TOKEN")
		}
		hostName = u.Host
	}

	return readTerraformTokenFile(hostName)
}

func readTerraformTokenFile(terraformHostName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "cannot locate home directory")
	}

	absolutePath := filepath.Join(homeDir, tfrcFileName)

	jsonFile, err := os.Open(absolutePath)
	if err != nil {
		return "", errors.Wrap(err, "cannot open terraform credentials file")
	}

	defer jsonFile.Close()

	var tfeConfig map[string]interface{}
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", errors.Wrap(err, "cannot read terraform credentials file")
	}

	err = json.Unmarshal(bytes, &tfeConfig)
	if err != nil {
		return "", errors.Wrap(err, "cannot read terraform credentials file")
	}

	tfeConfig, err = flatten.Flatten(tfeConfig, "", flatten.RailsStyle)
	if err == nil {
		query := fmt.Sprintf("credentials[%s][token]", terraformHostName)
		token, ok := tfeConfig[query]
		if ok {
			return token.(string), nil
		}
		return "", errors.New("credentials file contains no token")
	}

	log.Println("Cannot read a token from the tfrc file")
	return "", errors.New("unable to read the TFE token")
}
