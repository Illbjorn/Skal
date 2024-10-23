package main

import (
	"os"
)

func printHelpTextAndExit() {
	println(helpText)
	os.Exit(0)
}

// func timestamp() string {
// 	return time.Now().Format("01-02-2006 3:04:05pm")
// }

// func prefix(lvl string) *pprint.Msg {
// 	m := pprint.New()
// 	switch lvl {
// 	case "ERR", "FTL":
// 		m.Red(lvl)
// 	case "INF":
// 		m.Green(lvl)
// 	case "WARN":
// 		m.Yellow(lvl)
// 	case "DBG":
// 		m.Gray(lvl)
// 	}
// 	return m.Add(" ").White(timestamp())
// }

// func errGeneric(msg string, ftl bool) {
// 	if ftl {
// 		prefix("FTL").White(msg).Println()
// 		os.Exit(1)
// 	}
// 	prefix("ERR").White(msg).Println()
// }
