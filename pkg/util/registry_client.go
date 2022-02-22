package util

type RegistryClient interface {
	Login(userName string, password string, endpoint string) error
}

type DefaultRegistryClient struct{}

func (e DefaultRegistryClient) Login(userName string, password string, endpoint string) error {
	return RegistryLogin(userName, password, endpoint)
}

func NewDefaultRegistryClient() RegistryClient {
	return DefaultRegistryClient{}
}

type DummyRegistryClient struct{}

func (e DummyRegistryClient) Login(_ string, _ string, _ string) error {
	return nil
}

func NewDummyRegistryClient() RegistryClient {
	return DummyRegistryClient{}
}
