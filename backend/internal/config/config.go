package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Rates   Rates   `json:"rates"`
	Timing  Timing  `json:"timing"`
	Scoring Scoring `json:"scoring"`
	Units   Units   `json:"units"`
	Weapons Weapons `json:"weapons"`
}

type Rates struct {
	SpawnRateMs            int `json:"spawnRateMs"`
	IndividualWeaponRateMs int `json:"individualWeaponRateMs"`
	GroupWeaponRateMs      int `json:"groupWeaponRateMs"`
	TtlCheckRateMs         int `json:"ttlCheckRateMs"`
}

type Timing struct {
	BattlefieldTtlSec float64 `json:"battlefieldTtlSec"`
	TransitTtlSec     float64 `json:"transitTtlSec"`
}

type Scoring struct {
	Points    PointValues `json:"points"`
	Penalties Penalties   `json:"penalties"`
}

type PointValues struct {
	Soldier    int `json:"soldier"`
	Tank       int `json:"tank"`
	Helicopter int `json:"helicopter"`
}

type Penalties struct {
	Escape int `json:"escape"`
	Breach int `json:"breach"`
}

type Units struct {
	Hitpoints Hitpoints `json:"hitpoints"`
}

type Hitpoints struct {
	Soldier    int `json:"soldier"`
	Tank       int `json:"tank"`
	Helicopter int `json:"helicopter"`
}

type Weapons struct {
	ProcessingPower ProcessingPower `json:"processingPower"`
}

type ProcessingPower struct {
	Individual int `json:"individual"`
	Group      int `json:"group"`
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
