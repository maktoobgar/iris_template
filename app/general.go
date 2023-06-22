package app

import (
	"fmt"

	g "service/global"
	"service/pkg/colors"
)

func info() {
	fmt.Println(colors.Cyan, fmt.Sprintf("\n==%sSystem Info%s==%s\n", colors.Yellow, colors.Cyan, colors.Reset))
	fmt.Printf("Name:\t\t\t%s%s%s\n", colors.Blue, g.Name, colors.Reset)
	fmt.Printf("Version:\t\t%s%s%s\n", colors.Blue, g.Version, colors.Reset)
	// TODO: Check active/inactive microservices and print their status
	if g.CFG.Debug {
		fmt.Printf("Debug:\t\t\t%s%v%s\n", colors.Red, g.CFG.Debug, colors.Reset)
	} else {
		fmt.Printf("Debug:\t\t\t%s%v%s\n", colors.Green, g.CFG.Debug, colors.Reset)
	}
	fmt.Printf("Address:\t\thttp://%s:%s\n", g.CFG.Gateway.IP, g.CFG.Gateway.Port)
	fmt.Printf("Allowed Origins:\t%v\n", g.CFG.AllowOrigins)
	if g.CFG.AllowHeaders != "" {
		fmt.Printf("Extra Allowed Headers:\t%v\n", g.CFG.AllowHeaders)
	}
	fmt.Print(colors.Cyan, "===============\n\n", colors.Reset)
}
