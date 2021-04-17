package main

import (
	"GraduateDesign/data"
	"GraduateDesign/engine"
	"fmt"
	"log"
)

const port = 9999

func main() {
	router := engine.SecKillEngine()
	defer data.Close()

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Panic("router运行时出错： " + err.Error())
	}

}
