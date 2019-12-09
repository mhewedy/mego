package conf

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const commentChar = "#"

var props map[string]string

type Source interface {
	read() (io.ReadCloser, error)
}

type source struct {
}

func (s source) read() (io.ReadCloser, error) {
	f, err := os.Open("mego.conf")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func init() {
	props = make(map[string]string)

	r, err := source{}.read()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) > 0 && !strings.HasPrefix(trimmedLine, commentChar) {
			kv := strings.Split(trimmedLine, "=")
			props[strings.TrimSpace(kv[0])] =
				strings.TrimSpace(strings.Split(kv[1], commentChar)[0]) // take value part before comment char
		}
	}
}

func Get(key string, defaultValue ...string) string {
	v, found := props[key]
	if !found || len(v) == 0 {
		fmt.Fprintln(os.Stderr, "key", key, "not found")
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return v
}

func GetBool(key string, defaultValue bool) bool {
	v, found := props[key]
	if !found {
		fmt.Fprintln(os.Stderr, "key", key, "not found")
		return defaultValue
	}
	b := v == "true" || v == "yes"
	return b
}

func GetInt(key string, defaultValue int) int {
	v, found := props[key]
	if !found {
		fmt.Fprintln(os.Stderr, "key", key, "not found")
		return defaultValue
	}

	d, err := strconv.Atoi(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return defaultValue
	}
	return d
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	v, found := props[key]
	if !found {
		fmt.Fprintln(os.Stderr, "key", key, "not found")
		return defaultValue
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return defaultValue
	}
	return d
}
