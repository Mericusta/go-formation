package utility

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func CompareGameDataJsonObjectData(data interface{}, content string) bool {
	return strings.Compare(fmt.Sprintf("%v", data), content) == 0
}

func TraitFileName(fullFilename, extendType string) string {
	return fullFilename[0:strings.Index(fullFilename, extendType)]
}

func ConvertFormationMapToJson(m map[string]map[string]string) map[string]string {
	formationJsonMap := make(map[string]string)
	for fileName, formationMap := range m {
		j, e := json.Marshal(formationMap)
		if e != nil {
			fmt.Printf("Error: json marshal formation map %v occurs error: %v\n", formationMap, e)
			return nil
		}
		formationJsonMap[fileName] = string(j)
	}

	return formationJsonMap
}

func ConvertFileContentToJson(r io.Reader) (string, map[string]string, error) {
	transReader := transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	fileReader := csv.NewReader(transReader)
	err, jsonString, formationMap := ProcessCsvAndFormation(fileReader)
	if err != nil {
		return "", nil, err
	}
	return jsonString, formationMap, nil
}

type KeyIndex struct {
	Name  string
	Type  string
	Index int
}

func ProcessCsvAndFormation(fileReader *csv.Reader) (error, string, map[string]string) {
	_, _ = fileReader.Read()                 //注释行
	formationArray, err := fileReader.Read() //策划注释行
	if err != nil {
		return err, "", nil
	}

	opsArray, err := fileReader.Read()
	if err != nil {
		return err, "", nil
	}

	keyArray, err := fileReader.Read()
	if err != nil {
		return err, "", nil
	}

	typeArray, err := fileReader.Read()
	if err != nil {
		return err, "", nil
	}

	keyMap := map[int]*KeyIndex{}
	format := map[string]int{}
	formationMap := make(map[string]string)

	arrayIndex := 0
	for i, v := range keyArray {
		ops := opsArray[i]
		if ops == "client" || ops == "none" {
			continue
		}

		key := &KeyIndex{}
		key.Name = v
		key.Type = typeArray[i]
		key.Index = arrayIndex

		keyMap[i] = key
		format[v] = arrayIndex
		formationMap[v] = formationArray[i]
		arrayIndex++
	}

	data := make([][]interface{}, 0, 128)

	for {

		line, err := fileReader.Read()
		if err == io.EOF {
			break
		}

		o := ProcessLine(line, keyMap)
		data = append(data, o)
	}

	jsonObject := &struct {
		Format map[string]int  `json:"Format"`
		Data   [][]interface{} `json:"Data"`
	}{}

	jsonObject.Format = format
	jsonObject.Data = data

	jsonString, err := json.Marshal(jsonObject)
	return err, string(jsonString), formationMap
}

func ProcessCsv(fileReader *csv.Reader) (error, string) {
	_, _ = fileReader.Read() //注释行
	_, _ = fileReader.Read() //策划注释行

	opsArray, err := fileReader.Read()
	if err != nil {
		return err, ""
	}

	keyArray, err := fileReader.Read()
	if err != nil {
		return nil, ""
	}

	typeArray, err := fileReader.Read()
	if err != nil {
		return nil, ""
	}

	keyMap := map[int]*KeyIndex{}
	format := map[string]int{}

	arrayIndex := 0
	for i, v := range keyArray {
		ops := opsArray[i]
		if ops == "client" || ops == "none" {
			continue
		}

		key := &KeyIndex{}
		key.Name = v
		key.Type = typeArray[i]
		key.Index = arrayIndex

		keyMap[i] = key
		format[v] = arrayIndex
		arrayIndex++
	}

	data := make([][]interface{}, 0, 128)

	for {

		line, err := fileReader.Read()
		if err == io.EOF {
			break
		}

		o := ProcessLine(line, keyMap)
		data = append(data, o)
	}

	jsonObject := &struct {
		Format map[string]int  `json:"Format"`
		Data   [][]interface{} `json:"Data"`
	}{}

	jsonObject.Format = format
	jsonObject.Data = data

	jsonString, err := json.Marshal(jsonObject)
	return err, string(jsonString)
}

func ProcessLine(dataArray []string, keyMap map[int]*KeyIndex) []interface{} {
	r := make([]interface{}, 0, len(dataArray))
	for i, v := range dataArray {
		key, ok := keyMap[i]
		if ok == false {
			continue
		}
		r = append(r, GetParseString(key.Type, v))
	}
	return r
}

func GetParseString(Type string, v string) interface{} {
	switch Type {
	case "int", "int64", "int32":
		s, _ := strconv.ParseInt(v, 0, 64)
		return s
	case "double":
		s, _ := strconv.ParseFloat(v, 64)
		return s
	case "string":
		return v
	default:
	}
	return ""
}
