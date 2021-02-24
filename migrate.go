package migrate

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
)

var fr = regexp.MustCompile(`^(?P<timestamp>\d{14})_(?P<name>[\w\d_-]+)\.(?P<suffix>up|down)\.sql$`)

func parseName(regex *regexp.Regexp, str string) (Entry, string, error) {
	match := regex.FindStringSubmatch(str)
	names := regex.SubexpNames()

	var err error
	var result Entry
	var suffix string
	for i, value := range match {
		switch names[i] {
		case "timestamp":
			result.Timestamp, err = strconv.ParseInt(value, 10, 63)
			if err != nil {
				return Entry{}, "", err
			}
		case "name":
			result.Name = value
		case "suffix":
			suffix = value
		}
	}
	return result, suffix, nil
}

func ListMigrations(datadir string) ([]Entry, error) {
	items, err := ioutil.ReadDir(datadir)
	if err != nil {
		return nil, err
	}
	var result []Entry
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		log.WithField("file", item.Name()).Debug("found migration file")
		entry, suffix, err := parseName(fr, item.Name())
		if err != nil {
			log.WithError(err).Warn("cannot parse file name, skipped")
		}
		if suffix != "up" { // read only "up" migrations
			continue
		}
		result = append(result, entry)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return (result[i].Timestamp < result[j].Timestamp) && (result[i].Name < result[j].Name)
	})

	return result, nil
}
