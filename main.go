package main

import (
	"fmt"
	"github.com/dark-enstein/crontable/pkg/meaning"
	"github.com/dark-enstein/crontable/pkg/reader"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Println("please pass in the location of the crontab file to be read. \n usage: crontable <file>")
		os.Exit(1)
		return
	}
	fileLoc := args[0]
	//fileLoc := "/Users/ayobamibamigboye/GolandProjects/crontable/pkg/reader/.crontab"
	cronFile, err := reader.OpenCrontableFile(fileLoc)
	if err != nil {
		os.Exit(1)
		return
	}

	isValid, err := cronFile.Validate()
	if !isValid {
		log.Printf("crontab is not invalid. \nsample cronfile: %v\n", reader.SampleCronFile)
		os.Exit(1)
		return
	}

	cExpr, err := cronFile.MarshalIntoCronExpression()
	if err != nil {
		log.Printf("crontab is not invalid. reading failed with: %s", err.Error())
		os.Exit(1)
		return
	}

	cExprDecode := cronFile.Decode()
	if cExprDecode == nil {
		log.Printf("crontab is not invalid. reading failed with: %s", err.Error())
		os.Exit(1)
		return
	}

	fmt.Printf("cron expression read: %#v\n", cExpr)
	fmt.Printf("cron expression decoded: %#v\n", cExprDecode)

	meaning.Write
}
