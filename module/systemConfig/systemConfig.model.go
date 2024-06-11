package systemConfig


type SystemConfigModel struct{
	DeviceId string `json:"device_id" toml:"device_id"`
	CollectionFrequency string `json:"collection_frequency" toml:"collection_frequency"`
	ImgNameRules string `json:"img_name_rules" toml:"img_name_rules"`
	ResolutionRatio string `json:"resolution_ratio" toml:"resolution_ratio"`

}

func(SystemConfigModel) DataFileName() string{
	return "system_config_model.toml"
}


type WifiConfig struct{
	Ssid string `json:"ssid" toml:"ssid"`
	Password string `json:"password" toml:"password"`
}

type Local struct{
	Enable bool `json:"enable" toml:"enable"`
	OutputPath string `json:"output_path" toml:"output_path"`
}

type Nfs struct{
	Enable bool `json:"enable" toml:"enable"`
	OutputPath string `json:"output_path" toml:"output_path"`
	NfsHost string `json:"nfs_host" toml:"nfs_host"`
}

type Cifs struct{
	Enable bool `json:"enable" toml:"enable"`
	OutputPath string `json:"output_path" toml:"output_path"`
	CifsHost string `json:"cifs_host" toml:"cifs_host"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
}

type OutputConfig struct{
	Local Local `json:"local" toml:"local"`
	Nfs Nfs `json:"nfs" toml:"nfs"`
	Cifs Cifs `json:"cifs" toml:"cifs"`
}

func(OutputConfig)DataFileName() string{
	return "output_config.toml"
}


type WifiInfo struct{
	
} 