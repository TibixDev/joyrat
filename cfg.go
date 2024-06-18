package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func SaveCfg(cfg *Settings) {
	// Get the app data directory
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	// Make the joyrat directory in the app data directory
	joyratCfgDir := path.Join(cfgDir, "joyrat")
	if _, err := os.Stat(joyratCfgDir); os.IsNotExist(err) {
		os.Mkdir(joyratCfgDir, 0755)
	}
	joyratCfgFile := path.Join(joyratCfgDir, "config.json")

	// Save the config file
	cfgData, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("Error marshalling config data:", err)
		return
	}
	err = os.WriteFile(joyratCfgFile, cfgData, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}

	fmt.Println("Config saved to", joyratCfgFile)
}

func LoadCfg(cfg *Settings) {
	// Get the app data directory
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	// Make the joyrat directory in the app data directory
	joyratCfgDir := path.Join(cfgDir, "joyrat")
	if _, err := os.Stat(joyratCfgDir); os.IsNotExist(err) {
		os.Mkdir(joyratCfgDir, 0755)
	}
	joyratCfgFile := path.Join(joyratCfgDir, "config.json")

	// Load the config file
	cfgData, err := os.ReadFile(joyratCfgFile)
	if err != nil {
		fmt.Println("Config file does not exist, we'll save the default config for now")
		// Save the default config
		SaveCfg(cfg)
		return
	}
	err = json.Unmarshal(cfgData, cfg)
	if err != nil {
		fmt.Println("Error unmarshalling config data:", err)
		return
	}

	fmt.Println("Config loaded from", joyratCfgFile)
}

func CopyConfig(src *Settings, dest *Settings) {
	dest.MOUSE_SPEED = src.MOUSE_SPEED
	dest.MOUSE_SPEED_LOW = src.MOUSE_SPEED_LOW
	dest.MOUSE_SPEED_HIGH = src.MOUSE_SPEED_HIGH
	dest.SCROLL_SPEED = src.SCROLL_SPEED
	dest.JOYSTICK_DEADZONE = src.JOYSTICK_DEADZONE
	dest.AXIS_LT = src.AXIS_LT
	dest.AXIS_RT = src.AXIS_RT
	dest.AXIS_LS_X = src.AXIS_LS_X
	dest.AXIS_LS_Y = src.AXIS_LS_Y
	dest.AXIS_RS_X = src.AXIS_RS_X
	dest.AXIS_RS_Y = src.AXIS_RS_Y
}
