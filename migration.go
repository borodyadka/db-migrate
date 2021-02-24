package migrate

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

type Entry struct {
	Name      string
	Timestamp int64
}

func (e *Entry) Filename(suffix string) string {
	return fmt.Sprintf("%d_%s.%s.sql", e.Timestamp, e.Name, suffix)
}

func (e *Entry) read(datadir, suffix string) (string, error) {
	log.WithField("file", e.Filename(suffix)).Debug("reading migration")
	data, err := ioutil.ReadFile(path.Join(datadir, e.Filename(suffix)))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (e *Entry) ReadUp(datadir string) (string, error) {
	return e.read(datadir, "up")
}

func (e *Entry) ReadDown(datadir string) (string, error) {
	return e.read(datadir, "down")
}

func (e *Entry) Create(datadir string) error {
	log.WithField("name", e.Name).Debug("create migration")
	fup, err := os.OpenFile(path.Join(datadir, e.Filename("up")), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fup.Close()
	fdown, err := os.OpenFile(path.Join(datadir, e.Filename("down")), os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fdown.Close()
	return nil
}
