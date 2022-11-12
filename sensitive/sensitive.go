package sensitive

import (
	sensitive "github.com/importcjj/sensitive"
)

var filter *sensitive.Filter

func Init(dbfile string) {
	filter = sensitive.New()
	filter.UpdateNoisePattern(`[\|\s&%$@* ]+`)
	err := filter.LoadWordDict(dbfile)
	if err != nil {
		panic("未配置敏感词->" + dbfile)
	}
}

func Validate(text string) (bool, string) {
	return filter.Validate(text)
}

func Replace(text string, repl rune) string {
	return filter.Replace(text, repl)
}
