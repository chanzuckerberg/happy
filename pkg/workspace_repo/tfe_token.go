package workspace_repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/jeremywohl/flatten"
	"github.com/pkg/errors"
)

const tfrcFileName = ".terraform.d/credentials.tfrc.json"

func GetTfeToken(tfeUrl string) (string, error) {
	token, ok := os.LookupEnv("TFE_TOKEN")
	if ok {
		return token, nil
	}

	u, err := url.Parse(tfeUrl)
	if err != nil {
		log.Debugf("TFE URL %s is not valid: %s\n", tfeUrl, err.Error())
		return "", errors.New("please set env var TFE_TOKEN")
	}

	token, err = readTerraformTokenFile(u.Host)
	if err == nil {
		return token, nil
	}

	composeArgs := []string{"terraform", "login", u.Host}

	tf, err := exec.LookPath("terraform")
	if err != nil {
		return "", errors.New("please set env var TFE_TOKEN")
	}

	cmd := &exec.Cmd{
		Path:   tf,
		Args:   composeArgs,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	err = cmd.Run()
	if err != nil {
		return "", errors.New("please set env var TFE_TOKEN")
	}
	token, err = readTerraformTokenFile(u.Host)
	if err != nil {
		return "", errors.New("please set env var TFE_TOKEN")
	}
	return token, nil
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
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", errors.Wrap(err, "cannot read terraform credentials file")
	}

	err = json.Unmarshal(bytes, &tfeConfig)
	if err != nil {
		return "", errors.Wrap(err, "cannot read terraform credentials file")
	}

	tfeConfig, err = flatten.Flatten(tfeConfig, "", flatten.RailsStyle)
	if err == nil {
		token, ok := tfeConfig[fmt.Sprintf("credentials[%s][token]", terraformHostName)]
		if ok {
			return token.(string), nil
		}
	}

	log.Println("Cannot read a token from the tfrc file.")
	return "", errors.New("unable to read the TFE token")
}
