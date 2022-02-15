package workspace_repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/jeremywohl/flatten"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const tfrcFileName = ".terraform.d/credentials.tfrc.json"

func GetTfeToken(tfeUrl string, executor util.Executor) (string, error) {
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

	token, err := readTerraformTokenFile(hostName)
	if err == nil {
		return token, nil
	}

	composeArgs := []string{"terraform", "login", hostName}

	tf, err := exec.LookPath("terraform")
	if err != nil {
		return "", errors.Wrap(err, "please set env var TFE_TOKEN")
	}

	cmd := &exec.Cmd{
		Path:   tf,
		Args:   composeArgs,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	err = executor.Run(cmd)
	if err != nil {
		return "", errors.Wrap(err, "please set env var TFE_TOKEN")
	}
	token, err = readTerraformTokenFile(hostName)
	if err != nil {
		return "", errors.Wrap(err, "please set env var TFE_TOKEN")
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
		query := fmt.Sprintf("credentials[%s][token]", terraformHostName)
		token, ok := tfeConfig[query]
		if ok {
			return token.(string), nil
		}
	}

	log.Println("Cannot read a token from the tfrc file.")
	return "", errors.New("unable to read the TFE token")
}
