package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/fatih/structtag"
	"github.com/siskinc/gogo-generate/common"
	"go/ast"
	"go/format"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	typeNames     = flag.String("type", "", "comma-separated list of type names; must be set")
	output        = flag.String("output", "", "output file name; default srcdir/<type>_generate_mgorm.go")
	client        = flag.String("client", "", "object of mgorm.MongoDBClient")
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("generate-mgorm")
	flag.Parse()

	if 0 == len(*typeNames) {
		flag.Usage()
		os.Exit(1)
	}

	if 0 == len(*client) {
		flag.Usage()
		os.Exit(1)
	}

	types := strings.Split(*typeNames, ",")
	//fmt.Println("types:", types)
	gofile := common.GetGoFile()
	//fmt.Println("gofile", gofile)
	gopkg := common.GetGoPackage()
	//fmt.Println("gopkg", gopkg)
	structFieldListMap, structDocumentListMap, err := common.ParseStruct(gofile, nil)
	if nil != err {
		log.Fatalln(" ParseStruct is err:", err)
	}
	for _, typeName := range types {
		var documentList []string
		var fieldList []*ast.Field
		var ok bool
		if documentList, ok = structDocumentListMap[typeName]; !ok {
			log.Fatalln(" not fount type: ", typeName)
		}
		//fmt.Println("documentList", documentList)
		if fieldList, ok = structFieldListMap[typeName]; !ok {
			log.Fatalln(" not fount type: ", typeName)
		}
		//fmt.Println("fieldList", fieldList[0])
		parseInfo := map[string]interface{}{
			"package": gopkg,
			"struct":  typeName,
			"client":  *client,
			"id":      "ID",
		}

		parseStr := parseBasic
		fieldMap := make(map[string]*ast.Field)
		for _, fieldObj := range fieldList {
			fieldName := fieldObj.Names[0].Name
			fieldMap[fieldName] = fieldObj
			if nil != fieldObj.Tag {
				tags, err := structtag.Parse(strings.Trim(fieldObj.Tag.Value, "`"))
				if nil != err {
					log.Fatalf(" Parse tag is err: %s, typeName: %s, filedName: %s", err, typeName, fieldName)
				}
				bsonTag, _ := tags.Get("bson")
				if nil != bsonTag {
					if "_id" == bsonTag.Name {
						parseInfo["id"] = fieldName
					}
				}
			}
		}

		for _, document := range documentList {
			if strings.HasPrefix(document, "@def") {
				commandList := strings.Split(document, " ")
				//fmt.Println("commandList", len(commandList), commandList)
				if 2 <= len(commandList) {
					fieldNameFromCommand := commandList[2]
					switch commandList[1] {
					case SoftDelete:
						softDeleteField, ok := fieldMap[fieldNameFromCommand]
						if !ok {
							log.Fatalf(" not found soft delete field %s", fieldNameFromCommand)
						}
						tagValue := ""
						if nil != softDeleteField.Tag {
							tagValue = strings.Trim(softDeleteField.Tag.Value, "`")
						}
						tags, err := structtag.Parse(tagValue)
						if nil != err {
							log.Fatalf(" Parse tag is err: %s, structName: %s, filedName: %s", err, typeName, fieldNameFromCommand)
						}
						bsonTag, _ := tags.Get("bson")
						bsonName := ""
						if nil != bsonTag {
							bsonName = bsonTag.Name
						} else {
							bsonName = ToSnakeCase(fieldNameFromCommand)
						}
						if SoftDelete == commandList[1] && 3 <= len(commandList) {
							if fieldName, ok := parseInfo[SoftDelete]; !ok {
								parseStr += parseSoftDelete
							} else {
								log.Fatalf(" soft delete have been double declared, %s, %s", fieldName, commandList[2])
							}
							parseInfo[SoftDelete] = commandList[2]
							parseInfo[SoftDeleteBsonName] = bsonName
						}
					case SoftDeleteAt:
						softDeleteAtField, ok := fieldMap[fieldNameFromCommand]
						if !ok {
							log.Fatalf(" not found soft delete field %s", fieldNameFromCommand)
						}
						tagValue := ""
						if nil != softDeleteAtField.Tag {
							tagValue = strings.Trim(softDeleteAtField.Tag.Value, "`")
						}
						tags, err := structtag.Parse(tagValue)
						if nil != err {
							log.Fatalf(" Parse tag is err: %s, structName: %s, filedName: %s", err, typeName, fieldNameFromCommand)
						}
						bsonName := ""
						bsonTag, _ := tags.Get("bson")
						if nil != bsonTag {
							bsonName = bsonTag.Name
						} else {
							bsonName = ToSnakeCase(fieldNameFromCommand)
						}
						if fieldName, ok := parseInfo[SoftDeleteAt]; ok {
							log.Fatalf(" soft delete at have been double declared, %s, %s", fieldName, fieldNameFromCommand)
						}
						parseInfo[SoftDeleteAt] = commandList[2]
						parseInfo[SoftDeleteAtBsonName] = bsonName
					}
				}
			}
		}

		tmpl, err := template.New("").Parse(parseStr)
		if nil != err {
			log.Fatalf(" New %s struct template is err: %s", typeName, err)
		}
		buff := bytes.NewBufferString("")
		err = tmpl.Execute(buff, parseInfo)
		if nil != err {
			log.Fatalf(" template Execute is err: %s, typeName is %s", err, typeName)
		}
		// 格式化
		src, err := format.Source(buff.Bytes())
		if nil != err {
			log.Fatalf(" format Source is err: %s, typeName is %s", err, typeName)
		}
		//fmt.Println("src:\n", string(src))
		outputName := *output
		if outputName == "" {
			baseName := fmt.Sprintf("%s_generate_mgorm.go", ToSnakeCase(typeName))
			outputName = filepath.Join(".", strings.ToLower(baseName))
		}
		err = ioutil.WriteFile(outputName, src, 0644)
		if err != nil {
			log.Fatalf(" write to file is err: %s", err)
		}
	}
}
