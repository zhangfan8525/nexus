package myutils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type MyConfig map[string]string

var (
	configFile *os.File
	Configs    MyConfig
)

func parseConfig(file *os.File, configs MyConfig) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") {
			continue
		}

		lineSplit := strings.Split(line, "=")
		if len(lineSplit) == 2 {
			configs[lineSplit[0]] = lineSplit[1]
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (config MyConfig) getStrOrDefault(key, def string) string {
	if val, ok := config[key]; !ok {
		return def
	} else {
		return val
	}
}

func (config MyConfig) getIntOrDefault(key string, def int) int {
	if val, ok := config[key]; !ok {
		return def
	} else {
		v, err := strconv.Atoi(val)
		if err != nil {
			return def
		} else {
			return v
		}
	}
}

func init() {
	configFile, err := os.OpenFile("config", os.O_RDONLY, 0666)
	if err != nil {
		MyLogger.Println(err.Error())
	}

	defer func() {
		if err = configFile.Close(); err != nil {
			MyLogger.Println(err.Error())
		}
	}()

	Configs = make(map[string]string)
	parseConfig(configFile, Configs)
}
