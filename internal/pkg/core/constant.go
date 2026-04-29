package core

import (
	"regexp"
	"strings"
)

/** --- ID & NAME --- */

type IDNameInterface interface {
	OptionIDNames() map[int]string
	IDAndName(id int) map[string]interface{}
	Display(id int) string
}

type IDName struct{}

func (in IDName) Get(ind IDNameInterface) []interface{} {
	var results []interface{}
	for id, name := range ind.OptionIDNames() {
		results = append(results, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return results
}

func (in IDName) Display(id int, ind IDNameInterface) string {
	idNames := ind.OptionIDNames()
	if name, ok := idNames[id]; ok {
		return name
	}
	return ""
}

func (in IDName) IDAndName(id int, ind IDNameInterface) map[string]interface{} {
	return map[string]interface{}{
		"id":   id,
		"name": in.Display(id, ind),
	}
}

/** --- CODE & NAME --- */

type CodeNameInterface interface {
	OptionCodeNames() []string
}

type SnakeName struct{}

func (cn SnakeName) Get(cni CodeNameInterface) []interface{} {
	var results []interface{}
	for _, code := range cni.OptionCodeNames() {
		results = append(results, cn.CodeAndName(code))
	}

	return results
}

func (cn SnakeName) Display(code string) string {
	display := strings.Replace(code, "_", " ", -1)
	display = strings.Replace(display, "-", " ", -1)

	return strings.Title(display)
}

func (cn SnakeName) CodeAndName(code string) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"display": cn.Display(code),
	}
}

type CamelNameInterface interface {
	OptionCamelNames() []string
}

type CamelName struct{}

func (cn CamelName) Get(cni CamelNameInterface) []interface{} {
	var results []interface{}
	for _, code := range cni.OptionCamelNames() {
		results = append(results, cn.CodeAndName(code))
	}

	return results
}

func (cn CamelName) Display(code string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	display := re.ReplaceAllString(code, "${1} ${2}")

	return strings.Title(display)
}

func (cn CamelName) CodeAndName(code string) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"display": cn.Display(code),
	}
}
