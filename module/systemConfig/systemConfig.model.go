package systemConfig

type SystemConfigModel struct {
	DeviceId            string `json:"device_id" toml:"device_id"` //设备id
	CollectionFrequency string `json:"collection_frequency" toml:"collection_frequency"`//采集频率
	ImgNameRules        string `json:"img_name_rules" toml:"img_name_rules"` //命名规则
	ResolutionRatio     string `json:"resolution_ratio" toml:"resolution_ratio"` //分辨率
	MachineId           string `json:"machine_id" toml:"machine_id"` //机器id
}

func (SystemConfigModel) DataFileName() string {
	return "system_config_model.toml"
}

type WifiConfig struct {
	Ssid     string `json:"ssid" toml:"ssid"`
	Password string `json:"password" toml:"password"`
}

type Local struct {
	Enable     bool   `json:"enable" toml:"enable"`
	OutputPath string `json:"output_path" toml:"output_path"`
}

type Clouds struct {
	Enable bool   `json:"enable" toml:"enable"`
	Host   string `json:"host" toml:"host"`
}

type Cifs struct {
	Enable     bool   `json:"enable" toml:"enable"`
	OutputPath string `json:"output_path" toml:"output_path"`
	CifsHost   string `json:"cifs_host" toml:"cifs_host"`
	Username   string `json:"username" toml:"username"`
	Password   string `json:"password" toml:"password"`
}

type OutputConfig struct {
	Local  Local  `json:"local" toml:"local"`
	Clouds Clouds `json:"clouds" toml:"clouds"`
	// Cifs Cifs `json:"cifs" toml:"cifs"`
}

func (OutputConfig) DataFileName() string {
	return "output_config.toml"
}


type SchemaSaveImageIn struct {
	Imagebs64 string `json:"imagebs64"`
	DeviceID  string `json:"deviceID"`
	DeviceIP  string `json:"deviceIP"`
	Datetime  string `json:"dateTime"`
}

type SchemaSaveImageOut struct{
	Success bool `json:"success"`
	Info string `json:"info"`
}
