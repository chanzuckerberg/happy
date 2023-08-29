
struct HvmConfig {

    GithubPAT *string

}


func getHvmConfig() (*HvmConfig, error){
	home, err := os.UserHomeDir()

	if err != nil {
		return nil, errors.Wrap(err, "getting current user home directory")
	}

	configPath := path.Join(home, ".czi", "etc", "hvmconfig.json")

    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        fmt.Println("Config file does not exist")
        return
    }

    file, err := os.Open(configPath)
    if err != nil {
        return nil, errors.Wrap(err, "opening config file")
    }
    defer file.Close()

    // Parse json from file into HvmConfig struct

    output := &HvmConfig{}
    err = json.NewDecoder(file).Decode(&output)

    if err != nil {
        return nil, errors.Wrap(err, "parsing config file")
    }

    // Return HvmConfig struct

    return output, nil

}
