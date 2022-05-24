package api_common

type MicroserviceConfiguration struct {
	Infrastructure Infrastructure `yaml:"infrastructure"`
	Application    Application    `yaml:"application"`
}

type MicroService struct {
	SslEnabled             bool   `yaml:"sslEnabled"`
	SslPrivateKeyFilepath  string `yaml:"sslPrivateKeyFilepath"`
	SslCertificateFilepath string `yaml:"sslCertificateFilepath"`
	Port                   int    `yaml:"port"`
}

type Database struct {
	Name                  string `yaml:"name"`
	SslEnabled            bool   `yaml:"sslEnabled"`
	Address               string `yaml:"address"`
	Port                  int    `yaml:"port"`
	Username              string `yaml:"username"`
	PasswordFilepath      string `yaml:"passwordFilepath"`
	EncryptionKeyFilepath string `yaml:"encryptionKeyFilepath"`
}

type Common struct {
	InternalCACertFilepath string `yaml:"internalCACertFilepath"`
}

type Rabbit struct {
	Url          string     `yaml:"url"`
	Producer     RabbitInfo `yaml:"producer"`
	Consumer     RabbitInfo `yaml:"consumer"`
	Monitor      RabbitInfo `yaml:"monitor"`
	Notification RabbitInfo `yaml:"notification"`
}
type Application struct {
	Name            string     `yaml:"name"`
	BaseUrl         string     `yaml:"baseUrl"`
	Jwt             Jwt        `yaml:"jwt"`
	CorsPolicy      CorsPolicy `yaml:"corsPolicy"`
	MaxFailedLogins int        `yaml:"maxFailedLogins"`
	Password        Password   `yaml:"password"`
	Template        Template   `yaml:"template"`
	Config          Config     `yaml:"config"`
}

type Config struct {
	Import string `yaml:"import"`
	Update string `yaml:"update"`
}

type Template struct {
	AfterCreation       string `yaml:"afterCreation"`
	AfterForgotPassword string `yaml:"afterForgotPassword"`
}

type Infrastructure struct {
	MicroService MicroService `yaml:"microservice"`
	Database     Database     `yaml:"database"`
	Common       Common       `yaml:"common"`
	Rabbit       Rabbit       `yaml:"rabbit"`
}

type RabbitInfo struct {
	Queue    string `yaml:"queue"`
	Key      string `yaml:"key"`
	Exchange string `yaml:"exchange"`
}

type Jwt struct {
	Api Api `yaml:"api"`
}

type CorsPolicy struct {
	Enabled        bool     `yaml:"enabled"`
	AllowedDomains []string `yaml:"allowedDomains"`
}

type Api struct {
	Kid                string `yaml:"kid"`
	Audience           string `yaml:"audience"`
	Issuer             string `yaml:"issuer"`
	PublicKeyFilepath  string `yaml:"publicKeyFilepath"`
	PrivateKeyFilepath string `yaml:"privateKeyFilepath"`
	AccessToken        Token  `yaml:"accessToken"`
	RefreshToken       Token  `yaml:"refreshToken"`
}

type Password struct {
	PwdDuration int `yaml:"pwdDuration"`
	PwdWarning  int `yaml:"pwdWarning"`
}

type Token struct {
	Claims        []string `yaml:"claims"`
	ExpiryMinutes int      `yaml:"expiryMinutes"`
}
