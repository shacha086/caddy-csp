package csp

import (
	"fmt"
	orderedmap "github.com/shacha086/ordered-map/v2"
	"strings"
)

type CSP struct {
	commands *orderedmap.OrderedMap[string, []string]
}

func NewCSPFromHeader(cspHeader string) (*CSP, error) {
	parsedCommands := orderedmap.New[string, []string]()
	csp := &CSP{commands: parsedCommands}

	var warnings []string
	for _, command := range strings.Split(cspHeader, ";") {
		if command == "" {
			continue
		}
		tokens := strings.Fields(command)
		if len(tokens) < 2 {
			warnings = append(warnings, "ignored directive: "+command)
			continue
		}
		directive := tokens[0]
		sources := tokens[1:]
		parsedCommands.Set(directive, sources)
	}

	if len(warnings) > 0 {
		return csp, fmt.Errorf("%s", strings.Join(warnings, "; "))
	}
	return csp, nil
}

func (c *CSP) Encoded() string {
	header := ""
	for k, v := range c.commands.FromOldest() {
		header += k + " " + strings.Join(v, " ") + "; "
	}
	return strings.TrimSpace(header)
}
