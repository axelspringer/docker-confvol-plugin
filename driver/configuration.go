package driver

// Configuration struct
type Configuration struct {
	Backend struct {
		Type string `json:"type"`
		Auth struct {
			User     string `json:"user"`
			Password string `json:"password"`
		} `json:"auth"`
		Prefix string `json:"prefix"`
	} `json:"backend"`
	Generator struct {
		Structure string `json:"structure"`
	} `json:"generator"`
}

// LoadConfigurationFromFile load the conf
/*
func LoadConfigurationFromFile(filePath string) (*Configuration, error) {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		logrus.Errorf("Configuration file error: %v\n", e)
		return nil, err
	}

	var c Configuration
	json.Unmarshal(file, &c)
}
*/
