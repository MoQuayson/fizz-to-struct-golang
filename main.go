package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	readFile, err := os.Open("fizz.txt")

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	var res string

	for _, line := range fileLines {
		//check and get table name
		if strings.Contains(line, "create_table") {
			expData := Explode("create_table(", line)
			tableName := expData[0]
			if len(strings.Trim(expData[0], " ")) == 0 {
				tableName = expData[1]
			}

			tableName = Explode(`"`, tableName)[1]

			tableName = TrimColumnName(tableName)

			res = fmt.Sprintf("type %s struct{\n", tableName)
		}
		//get column name
		if strings.Contains(line, "t.Column") {

			expData := Explode("t.Column(", line)
			columnAttr := expData[0]
			if len(strings.Trim(expData[0], " ")) == 0 {
				columnAttr = expData[1]
			}

			expData = Explode(",", columnAttr)
			column := strings.Trim(expData[0], `" "`)
			attr := strings.Trim(expData[1], `" ")`)

			res = res + GenerateStructProps(column, attr)
		}

	}

	//add timestaps
	res = res + GenerateStructProps("created_at", "created_at")
	res = res + GenerateStructProps("updated_at", "updated_at")

	//closing bracket
	res = res + "}"

	//Output
	WriteOutputToFile(res)
}

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

func GenerateStructProps(column, attr string) string {
	title := TrimColumnName(column)
	attr = GetStructPropDataType(attr)
	dataAttr := fmt.Sprintf(`json:"%s" db:"%s"`, column, column)
	dataAttr = fmt.Sprintf("`%s`", dataAttr)
	res := fmt.Sprintf("%s %s %s\n", title, attr, dataAttr)

	return res
}

func TrimColumnName(column string) string {
	expData := Explode("_", column)
	var res string

	if column == "id" {
		return strings.ToUpper(column)
	}
	for _, val := range expData {
		res = res + strings.Title(val)
	}

	return res
}

func GetStructPropDataType(attr string) string {
	if attr == "datetime" || attr == "timestamp" || attr == "time" {
		return "time.Time"
	} else if attr == "uuid" {
		return "uuid.UUID"
	} else if attr == "integer" {
		return "int32"
	}

	return attr
}

func WriteOutputToFile(content string) {
	f, err := os.Create("./output.txt")
	if err != nil {
		fmt.Println(err)
	}

	n, err := f.WriteString(content + "\n")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()
}
