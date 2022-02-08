package config

// type Secrets interface {
// 	GetTfeUrl() string
// 	GetTfeOrg() string
// 	GetClusterArn() string
// 	GetPrivateSubnets() []string
// 	GetSecurityGroups() []string
// 	GetServiceUrl(service string) (string, error)
// 	GetAllServicesUrl() map[string]*RegistryConfig
// }

// type SecretsBackend interface {
// 	GetSecrets(secretArn string) (Secrets, error)
// }

// type TfeSecrets struct {
// 	Url string `json:"url"`
// 	Org string `json:"org"`
// }

// type AwsSecretMgrSecrets struct {
// 	ClusterArn     string                     `json:"cluster_arn"`
// 	PrivateSubnets []string                   `json:"private_subnets"`
// 	SecurityGroups []string                   `json:"security_groups"`
// 	Services       map[string]*RegistryConfig `json:"ecrs"`
// 	Tfe            *TfeSecrets                `json:"tfe"`
// }

// type AwsSecretMgr struct {
// 	session         *session.Session
// 	awsConfig       *aws.Config
// 	secretmgrClient secretsmanageriface.SecretsManagerAPI
// 	secrets         *AwsSecretMgrSecrets
// }

// // TODO(el): don't use singletons
// var secretMgrSessInst SecretsBackend
// var creatSecretMgeOnce sync.Once

// func GetAwsSecretMgr(awsProfile string) SecretsBackend {
// 	creatSecretMgeOnce.Do(func() {
// 		awsConfig := &aws.Config{
// 			// TODO: don't hardcode region
// 			Region:     aws.String("us-west-2"),
// 			MaxRetries: aws.Int(2),
// 		}
// 		// TODO(el): share a session through the codebase
// 		// TODO(el): don't panic
// 		session := session.Must(session.NewSessionWithOptions(session.Options{
// 			Profile:           awsProfile,
// 			Config:            *awsConfig,
// 			SharedConfigState: session.SharedConfigEnable,
// 		}))
// 		secretmgrClient := secretsmanager.New(session)
// 		secretMgrSessInst = &AwsSecretMgr{
// 			session:         session,
// 			awsConfig:       awsConfig,
// 			secretmgrClient: secretmgrClient,
// 		}
// 	})
// 	return secretMgrSessInst
// }

// func GetAwsSecretMgrWithClient(client secretsmanageriface.SecretsManagerAPI) SecretsBackend {
// 	secretMgrSessInst = &AwsSecretMgr{
// 		secretmgrClient: client,
// 	}
// 	return secretMgrSessInst
// }

// func (s *AwsSecretMgr) GetSecrets(secretArn string) (Secrets, error) {
// 	if s.secrets != nil {
// 		return s.secrets, nil
// 	}
// 	configInput := secretArn
// 	input := &secretsmanager.GetSecretValueInput{
// 		SecretId: &configInput,
// 	}
// 	secretOutput, err := s.secretmgrClient.GetSecretValue(input)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "Error reading from AWS secrets manager")
// 	}

// 	s.secrets = &AwsSecretMgrSecrets{}
// 	secrets := *secretOutput.SecretString
// 	err = json.Unmarshal([]byte(secrets), s.secrets)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "could not parse JSON")
// 	}

// 	return s.secrets, nil
// }
