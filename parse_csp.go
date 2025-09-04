package csp

import (
	"fmt"
	orderedmap "github.com/shacha086/ordered-map/v2"
	"slices"
	"strings"
)

type CSP struct {
	commands *orderedmap.OrderedMap[string, []string]
}

func NewCSPFromHeader(cspHeader string) (*CSP, error) {
	parsedCommands := orderedmap.New[string, []string]()
	csp := &CSP{commands: parsedCommands}
	e := csp.AddHeader(cspHeader)
	return csp, e
}

func (c *CSP) Encoded() string {
	header := ""
	for k, v := range c.commands.FromOldest() {
		header += k + " " + strings.Join(v, " ") + "; "
	}
	return strings.TrimSpace(header)
}

func (c *CSP) AddHeader(cspHeader string) error {
	var warnings []string

	for _, command := range strings.Split(cspHeader, ";") {
		if command == "" {
			continue
		}
		err := c.AddCommand(command)
		if err != nil {
			warnings = append(warnings, err.Error())
		}
	}
	if len(warnings) > 0 {
		return fmt.Errorf("%s", strings.Join(warnings, "; "))
	}
	return nil
}

func (c *CSP) AddCommand(cspCommand string) error {
	cspCommand, _ = strings.CutSuffix(strings.TrimSpace(cspCommand), ";")
	tokens := strings.Fields(cspCommand)
	if len(tokens) < 2 {
		return fmt.Errorf("invalid command: %s", cspCommand)
	}
	directive := tokens[0]
	sources := tokens[1:]
	c.Add(directive, sources...)
	return nil
}

func (c *CSP) Add(directive string, sources ...string) {
	if directive == "" {
		return
	}
	if len(sources) < 1 {
		return
	}

	orig, _ := c.commands.Get(directive)
	c.commands.Set(directive, append(orig, sources...))
}

func (c *CSP) AddAllDirective(sources ...string) {
	for k := range c.commands.KeysFromOldest() {
		c.Add(k, sources...)
	}
}

func (c *CSP) Remove(directive string, sources ...string) {
	if directive == "" {
		return
	}
	if len(sources) < 1 {
		return
	}

	orig, _ := c.commands.Get(directive)
	c.commands.Set(directive, slices.DeleteFunc(orig, func(s string) bool {
		return slices.Contains(sources, s)
	}))
}

func (c *CSP) RemoveAllDirective(sources ...string) {
	for k := range c.commands.KeysFromOldest() {
		c.Remove(k, sources...)
	}
}

func (c *CSP) Set(directive string, sources ...string) {
	if directive == "" {
		return
	}
	if len(sources) < 1 {
		return
	}

	c.commands.Set(directive, sources)
}

func (c *CSP) SetAllDirective(sources ...string) {
	for k := range c.commands.KeysFromOldest() {
		c.Set(k, sources...)
	}
}
