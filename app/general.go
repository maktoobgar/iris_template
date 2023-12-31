package app

import (
	"fmt"
	"log"
	"runtime"

	g "service/global"
	"service/pkg/colors"
)

func runCronJobs() {
	// For Reference use:
	// https://crontab.guru/every-minute
	//
	// Example:
	// g.Cron.AddJob("* * * * *", some_function)
}

func info() {
	if IsChild() {
		fmt.Println(colors.Cyan, fmt.Sprintf("==%sClone Number %s Started%s==%s", colors.Green, GetChildNumber(), colors.Cyan, colors.Reset))
		return
	}
	fmt.Println(colors.Cyan, fmt.Sprintf("\n==%sSystem Info%s==%s\n", colors.Yellow, colors.Cyan, colors.Reset))
	fmt.Printf("Name:\t\t\t%s%s%s\n", colors.Blue, g.Name, colors.Reset)
	fmt.Printf("Version:\t\t%s%s%s\n", colors.Blue, g.Version, colors.Reset)
	if g.CFG.ClonesCount != 0 {
		corsCapacity := runtime.GOMAXPROCS(0)
		cloneColor := colors.Green
		if g.CFG.ClonesCount > corsCapacity {
			cloneColor = colors.Red
		}
		fmt.Printf("Clones:\t\t\t%s%d%s (%d cors)\n", cloneColor, g.CFG.ClonesCount, colors.Reset, corsCapacity)
	}
	mainOrTest := "test"
	mainOrTestColor := colors.Red + mainOrTest + colors.Reset
	if !g.CFG.Debug {
		mainOrTest = "main"
		mainOrTestColor = colors.Green + mainOrTest + colors.Reset
	}
	for name, database := range g.CFG.Gateway.Databases {
		if name == mainOrTest {
			if database.Type == "sqlite3" {
				fmt.Printf("Main Database:\t\t%v, %v (%v)\n", database.Type, database.DbName, mainOrTestColor)
			} else {
				fmt.Printf("Main Database:\t\t%v, %v, %v:%v (%v)\n", database.Type, database.DbName, database.Host, database.Port, mainOrTestColor)
			}
			if g.DB == nil {
				log.Fatal("default database connection is not assigned as main database")
			}
			break
		}
	}
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
