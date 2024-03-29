package updater

import (
	"encoding/json"
	"os"
	"path"
)

const npmrc = "//registry.npmjs.org/:_authToken=${NPM_TOKEN}\n"

var FUVERSION = "dev"

type Updater struct{}

func (u *Updater) Init(_ map[string]string) error {
	return nil
}

func (u *Updater) Name() string {
	return "npm"
}

func (u *Updater) Version() string {
	return FUVERSION
}

func (u *Updater) ForFiles() string {
	return "package\\.json"
}

func updateJSONFile(fName, newVersion string) error {
	file, err := os.OpenFile(fName, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	var data map[string]json.RawMessage
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	data["version"] = json.RawMessage("\"" + newVersion + "\"")
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (u *Updater) Apply(file, newVersion string) error {
	if err := updateJSONFile(file, newVersion); err != nil {
		return err
	}

	packageLockPath := path.Join(path.Dir(file), "package-lock.json")
	if _, err := os.Stat(packageLockPath); err == nil {
		if err := updateJSONFile(packageLockPath, newVersion); err != nil {
			return err
		}
	}

	if os.Getenv("NPM_CONFIG_USERCONFIG") != "" {
		return nil
	}

	var err error
	npmrcPath := path.Join(path.Dir(file), ".npmrc")
	if _, err = os.Stat(npmrcPath); os.IsNotExist(err) {
		return os.WriteFile(npmrcPath, []byte(npmrc), 0o644)
	}

	return err
}
