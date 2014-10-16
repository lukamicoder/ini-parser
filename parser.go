package iniparser

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	FileName string
	Sections []*Section
}

type Section struct {
	Name string
	Keys map[string]string
}

func (conf *Config) LoadFile(fileName string) error {
	if filepath.Dir(fileName) == "." {
		pos := strings.LastIndex(os.Args[0], string(filepath.Separator))
		path := os.Args[0][0 : pos+1]

		fileName = filepath.Join(path, fileName)
	}

	conf.FileName = fileName

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		err = conf.parseLine(line)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (conf *Config) parseLine(line string) error {
	if len(line) == 0 || line[0] == ';' || line[0] == '#' || line[0] == '\n' || line[0] == '\r' {
		return nil
	}

	if line[0] == '[' {
		pos := strings.Index(line, "]")
		if pos < 1 {
			return errors.New("Failed to parse section header: " + line)
		}

		section := new(Section)
		section.Name = line[1:pos]
		section.Keys = make(map[string]string)

		conf.Sections = append(conf.Sections, section)

		return nil
	}

	if conf.Sections == nil || len(conf.Sections) == 0 {
		return errors.New("Key is not under any section: " + line)
	}

	pos := strings.Index(line, "=")
	if pos < 1 {
		return errors.New("Failed to parse key line: " + line)
	}

	name := line[0:pos]
	val := line[pos+1 : len(line)]

	conf.Sections[len(conf.Sections)-1].Keys[name] = val

	return nil
}

func (conf *Config) GetSection(name string) (*Section, error) {
	for _, section := range conf.Sections {
		if section.Name == name {
			return section, nil
		}
	}

	return nil, errors.New("Section not found: " + name)
}

func (conf *Config) GetSections() []*Section {
	return conf.Sections
}

func (conf *Config) GetString(section string, key string) (string, error) {
	sec, err := conf.GetSection(section)
	if err != nil {
		return "", err
	}

	val := sec.Keys[key]
	if val == "" {
		return "", errors.New("Key " + key + " not found")
	}

	return val, nil
}

func (conf *Config) GetBool(section string, key string) (bool, error) {
	value, err := conf.GetString(section, key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(value)
}

func (conf *Config) GetInt(section string, key string) (int, error) {
	value, err := conf.GetString(section, key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func (conf *Config) GetInt64(section string, key string) (int64, error) {
	value, err := conf.GetString(section, key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(value, 10, 64)
}

func (conf *Config) GetFloat64(section string, key string) (float64, error) {
	value, err := conf.GetString(section, key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(value, 64)
}
