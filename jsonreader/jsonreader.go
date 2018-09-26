package jsonreader

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func FromFile(file string) (string, error) {
	fb, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(fb), nil
}

func FromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FromStdin() (string, error) {
	b := strings.Builder{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		_, err := b.Write(scanner.Bytes())

		if err != nil {
			return "", err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func FromPipe(_ string) (string, error) {
	json, err := FromStdin()
	if err != nil {
		return "", err
	}
	return json, nil
}

type JsonFunc func(s string) (string, error)
