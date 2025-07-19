package def2env

import (
	"bufio"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/kayac/ecspresso/v2"
)

type AllowList struct {
	names []string
	all   bool
}

func NewAllowList(options *AllowListOptions) *AllowList {
	nameSet := map[string]struct{}{}

	for _, fileOrName := range options.Only {
		if u, err := url.Parse(fileOrName); err == nil && u.Scheme == "file" {
			f, err := os.Open(u.Host)

			if err != nil {
				ecspresso.LogWarn("file loading skipped: %s", err)
				continue
			}

			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				name := strings.TrimSpace(scanner.Text())

				if name == "" || strings.HasPrefix(name, "#") {
					continue
				}

				nameSet[name] = struct{}{}
			}
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

	return allowlist
}

func (allowlist AllowList) IsAllowed(name string) bool {
	return allowlist.all || slices.Contains(allowlist.names, name)
}
