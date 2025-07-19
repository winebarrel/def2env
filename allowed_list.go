package def2env

import (
	"bufio"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/kayac/ecspresso/v2"
)

type AllowList []string

func NewAllowList(only []string) AllowList {
	nameSet := map[string]struct{}{}

	for _, fileOrName := range only {
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

	return names
}

func (allowlist AllowList) IsAllowed(name string) bool {
	return len(allowlist) == 0 || slices.Contains(allowlist, name)
}
