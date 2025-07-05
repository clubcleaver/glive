package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type configFormat struct {
	Command   string   `json:"command"`
	Directory string   `json:"watch"`
	Skip      []string `json:"skip"`
}

func GetConfig() (config configFormat, errr error) {
	// Get WD
	workingDir, err := os.Getwd()
	if err != nil {
		errr = fmt.Errorf("config error, Ref: Could not read working Directory. Error: %v", err.Error())
		return
	}
	// Get Config File
	configPath := fmt.Sprintf("%s/glive.json", workingDir)
	configFile, err := os.Open(configPath)
	if err != nil {
		errr = fmt.Errorf("config error, Ref: Could not find:\"glive.json\" in directory: %v. Error: %v", workingDir, err.Error())
		return
	}
	defer configFile.Close()

	// Read File Content
	fileStat, err := configFile.Stat()
	if err != nil {
		errr = fmt.Errorf("config error, Ref: No Access to file info, Check permissions. Error: %v", err.Error())
		return
	}

	var content = make([]byte, fileStat.Size())
	_, err = configFile.Read(content)
	if err != nil {
		errr = fmt.Errorf("config error, Ref: Read contents of file:\"glive.json\". Error: %v", err.Error())
		return
	}

	// Parse Config File in JSON
	err = json.Unmarshal(content, &config)
	if err != nil {
		errr = fmt.Errorf("config error, Ref: Parsing JSON from file:\"glive.json\". Error: %v", err.Error())
		return
	}

	// validate Input
	if err = config.Validate(); err != nil {
		errr = err
		return
	}

	// Check Config Directory to Watch is a directory and is valid
	dir, err := os.Open(config.Directory)
	if err != nil {
		errr = fmt.Errorf("config error, Ref: Could Not find for watch directory: \"%v\", Error: %v", config.Directory, err.Error())
		return
	}
	info, err := dir.Stat()
	if err != nil {
		errr = fmt.Errorf("config error, Ref: Could Not find Info for watch directory: \"%v\", Error: %v", config.Directory, err.Error())
		return
	}
	if !info.IsDir() {
		errr = fmt.Errorf("config error, Ref: Value provided for \"watch\" should be a directory")
		return
	}

	// Handle config validation in calling function
	return
}

func (conf *configFormat) Validate() error {
	// Sanitize Spaces
	conf.Command = strings.TrimSpace(conf.Command)
	conf.Directory = strings.TrimSpace(conf.Directory)
	if len(conf.Skip) > 0 {
		for ind, val := range conf.Skip {
			conf.Skip[ind] = strings.TrimSpace(val)

		}
	}

	// Check Value Lengths
	if conf.Command == "" {
		return fmt.Errorf("config error, Ref: \"command\" must be atlease 1 string value")
	}

	if conf.Directory == "" || len(strings.Split(conf.Directory, " ")) > 1 {
		return fmt.Errorf("config error, Ref: \"watch\" must be single string value")
	}

	return nil
}
