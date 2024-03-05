package appsettings

type Country struct {
	CountryCode string   `yaml:"country_code"`
	Enabled     bool     `yaml:"enabled"`
	Regions     []string `yaml:"regions"`
}

type Language struct {
	LangaugeCode string `yaml:"language_code"`
	Enabled      bool   `yaml:"enabled"`
}

type configurations struct {
	Countries []*Country  `yaml:"countries"`
	Languages []*Language `yaml:"languages"`
}

func New() *configurations {
	return &configurations{}
}
