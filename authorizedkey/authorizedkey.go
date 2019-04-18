package authorizedkey

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	authorizedKeyFile = "/authorized_keys"
)

// File is the base structure for your auth key struct
type File struct {
	UID string
	Key string
}

func writtenKey(id string, key string) string {
	return "# <<< start " + id + " >>>\n" + key + "\n# <<< end " + id + " >>>\n"
}

func appendToFile(key string) error {
	f, err := os.OpenFile(authorizedKeyFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(key); err != nil {
		return err
	}
	return nil
}

func existInFile(id string) bool {
	b, err := ioutil.ReadFile(authorizedKeyFile)
	if err != nil {
		return false
	}
	s := string(b)
	// //check whether s contains substring text
	if strings.Contains(s, "<<< start "+id) {
		return true
	}
	return false
}

func removeFromFile(id string) error {
	b, err := ioutil.ReadFile(authorizedKeyFile)
	if err != nil {
		return err
	}
	s := string(b)
	r, _ := regexp.Compile("# <<< start " + id + " >>>\n.*\n# <<< end " + id + " >>>\n")
	ns := r.ReplaceAllString(s, "")
	nb := []byte(ns)
	err = ioutil.WriteFile(authorizedKeyFile, nb, 0644)
	if err != nil {
		return err
	}
	return nil
}

// WriteKey write the key as long as it doesn't exist
func (f File) WriteKey() error {
	if existInFile(f.UID) {
		return errors.New("Key " + f.UID + " already exists.")
	}

	str := writtenKey(f.UID, f.Key)
	return appendToFile(str)
}

// UpdateKey updates the key for the id in the authorzed keys file
func (f File) UpdateKey() error {
	err := f.DeleteKey()
	if err != nil {
		return err
	}

	err = f.WriteKey()
	if err != nil {
		return err
	}
	return nil
}

// DeleteKey updates the key for the id in the authorzed keys file
func (f File) DeleteKey() error {
	if existInFile(f.UID) {
		err := removeFromFile(f.UID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Key " + f.UID + " did not exist in authorized_keys")
}

func (f File) Sync(remove bool) error {
	if remove && existInFile(f.UID) {
		return f.DeleteKey()
	}

	if existInFile(f.UID) {
		return f.UpdateKey()
	}

	return f.WriteKey()
}
