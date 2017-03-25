package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func recordJob(j job, pathPrefix string) error {

	checksum := sha1.Sum([]byte(j.url.String()))
	parts := []string{pathPrefix}
	parts = append(parts, j.url.Host)
	parts = append(parts, fmt.Sprintf("%x", checksum))

	p := path.Join(parts...)

	if _, err := os.Stat(path.Dir(p)); os.IsNotExist(err) {
		err = os.MkdirAll(path.Dir(p), 0750)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(p, []byte(j.String()), 0640)
	if err != nil {
		return err
	}

	return nil
}

func (j job) String() string {
	buf := &bytes.Buffer{}

	buf.WriteString(j.url.String())
	buf.WriteString("\n\n")
	buf.WriteString(j.resp.status)
	buf.WriteString("\n")

	for name, values := range j.resp.headers {
		buf.WriteString(
			fmt.Sprintf("%s: %s\n", name, strings.Join(values, ", ")),
		)
	}

	buf.WriteString("\n\n")
	buf.Write(j.resp.body)
	buf.WriteString("\n")

	return buf.String()
}
