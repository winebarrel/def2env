package def2env

import (
	"bufio"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"
)

type AllowList struct {
	names []string
	all   bool
}

func NewAllowList(options *AllowListOptions) (*AllowList, error) {
	nameSet := map[string]struct{}{}

	for _, fileOrName := range options.Only {
		if u, err := url.Parse(fileOrName); err == nil && u.Scheme == "file" {
			filePath := path.Join(u.Host, u.Path)
			f, err := os.Open(filePath)

			if err != nil {
				if u.Query().Get("required") == "false" {
					continue
				} else {
					return nil, err
				}
			}

			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				name := strings.TrimSpace(scanner.Text())

				if name == "" || strings.HasPrefix(name, "#") {
					continue
				}

				nameSet[name] = struct{}{}
			}

			f.Close()
		} else {
			nameSet[fileOrName] = struct{}{}
		}
	}

	names := []string{}

	for n := range nameSet {
		names = append(names, n)
	}

	allowlist := &AllowList{
		names: names,
		all:   options.All,
	}

	return allowlist, nil
}

func (allowlist AllowList) IsAllowed(name string) bool {
	return allowlist.all || slices.Contains(allowlist.names, name)
}
