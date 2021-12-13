package formation

import (
	"fmt"
	"regexp"
	"strings"
)

type Node interface {
	CanMatch(string) bool
	ParseFormation(string) bool
	ParseContent(string) (map[string]map[string][]string, []error)
	GetKey() string
	GetValue() string
	GetFormation() string
	GetRelateFileFieldMap() map[string]string
}

type BaseNode struct {
	Key   string
	Value string
}

func (n *BaseNode) GetKey() string {
	return n.Key
}

func (n *BaseNode) GetValue() string {
	return n.Value
}

type SemicolonNode struct {
	BaseNode
	Formation string
	SubNode   Node
}

func (n *SemicolonNode) CanMatch(c string) bool {
	return nodeMatcherMap[SEMICOLON].CanMatch(c)
}

func (n *SemicolonNode) ParseFormation(c string) bool {
	n.Formation = c
	if len(n.Formation) == 0 {
		fmt.Printf("Error: semicolon node formation is empty\n")
		return false
	}

	for _, subFormation := range strings.Split(n.Formation, ";") {
		if !nodeMatcherMap[COMMA].CanMatch(subFormation) {
			fmt.Printf("Error: semicolon node can not match sub content '%v'\n", subFormation)
		} else {
			subNode := nodeMatcherMap[COMMA].NewNode()
			subNode.ParseFormation(subFormation)
			n.SubNode = subNode
			break
		}
	}
	return true
}

func (n *SemicolonNode) ParseContent(c string) (map[string]map[string][]string, []error) {
	parseContentErrorSlice := make([]error, 0)
	fileFieldContentSliceMap := make(map[string]map[string][]string)
	for _, subContent := range strings.Split(c, ";") {
		subNodeResultMap, errorSlice := n.SubNode.ParseContent(subContent)
		parseContentErrorSlice = append(parseContentErrorSlice, errorSlice...)
		fileFieldContentSliceMap = MergeFileFieldContentSliceMap(fileFieldContentSliceMap, subNodeResultMap)
	}
	return fileFieldContentSliceMap, parseContentErrorSlice
}

func (n *SemicolonNode) GetFormation() string {
	return n.Formation
}

func (n *SemicolonNode) GetRelateFileFieldMap() map[string]string {
	return n.SubNode.GetRelateFileFieldMap()
}

type CommaNode struct {
	BaseNode
	Formation    string
	SubNodeSlice []Node
}

func (n *CommaNode) CanMatch(c string) bool {
	return nodeMatcherMap[COMMA].CanMatch(c)
}

func (n *CommaNode) ParseFormation(c string) bool {
	n.Formation = c
	if len(n.Formation) == 0 {
		fmt.Printf("Error: comma node formation is empty\n")
		return false
	}

	same := true
	var lastSubFormation string
	for _, subFormation := range strings.Split(n.Formation, ",") {
		if !nodeMatcherMap[FULLSTOP].CanMatch(subFormation) && subFormation != "PH" {
			fmt.Printf("Error: comma node can not match sub content '%v'\n", subFormation)
		} else {
			if len(lastSubFormation) != 0 {
				if same {
					same = lastSubFormation == subFormation
				}
			}
			subNode := nodeMatcherMap[FULLSTOP].NewNode()
			subNode.ParseFormation(subFormation)
			n.SubNodeSlice = append(n.SubNodeSlice, subNode)
			lastSubFormation = subFormation
		}
	}

	if same {
		n.SubNodeSlice = n.SubNodeSlice[:1]
	}

	return true
}

func (n *CommaNode) ParseContent(c string) (map[string]map[string][]string, []error) {
	parseContentErrorSlice := make([]error, 0)
	fileFieldContentSliceMap := make(map[string]map[string][]string)
	subContentSlice := strings.Split(c, ",")
	if len(n.SubNodeSlice) == 1 {
		for _, subContent := range subContentSlice {
			// fmt.Printf("DEBUG: sub content '%v' formation is '%v'\n", subContent, n.SubNodeSlice[0].GetFormation())
			if n.SubNodeSlice[0].GetFormation() == "PH" {
				continue
			}
			if _, hasFile := fileFieldContentSliceMap[n.SubNodeSlice[0].GetKey()]; !hasFile {
				fileFieldContentSliceMap[n.SubNodeSlice[0].GetKey()] = make(map[string][]string)
			}
			if _, hasField := fileFieldContentSliceMap[n.SubNodeSlice[0].GetKey()][n.SubNodeSlice[0].GetValue()]; !hasField {
				fileFieldContentSliceMap[n.SubNodeSlice[0].GetKey()][n.SubNodeSlice[0].GetValue()] = make([]string, 0)
			}
			fileFieldContentSliceMap[n.SubNodeSlice[0].GetKey()][n.SubNodeSlice[0].GetValue()] = append(fileFieldContentSliceMap[n.SubNodeSlice[0].GetKey()][n.SubNodeSlice[0].GetValue()], subContent)
		}
	} else if len(subContentSlice) == len(n.SubNodeSlice) {
		for index, subContent := range subContentSlice {
			// fmt.Printf("DEBUG: sub content '%v' formation is '%v'\n", subContent, n.SubNodeSlice[index].GetFormation())
			if n.SubNodeSlice[index].GetFormation() == "PH" {
				continue
			}
			if _, hasFile := fileFieldContentSliceMap[n.SubNodeSlice[index].GetKey()]; !hasFile {
				fileFieldContentSliceMap[n.SubNodeSlice[index].GetKey()] = make(map[string][]string)
			}
			if _, hasField := fileFieldContentSliceMap[n.SubNodeSlice[index].GetKey()][n.SubNodeSlice[index].GetValue()]; !hasField {
				fileFieldContentSliceMap[n.SubNodeSlice[index].GetKey()][n.SubNodeSlice[index].GetValue()] = make([]string, 0)
			}
			fileFieldContentSliceMap[n.SubNodeSlice[index].GetKey()][n.SubNodeSlice[index].GetValue()] = append(fileFieldContentSliceMap[n.SubNodeSlice[index].GetKey()][n.SubNodeSlice[index].GetValue()], subContent)
		}
	} else {
		// fmt.Printf("Error: sub node %v length %v is not equal 1 or sub content slice %v length %v\n", func() string {
		// 	s := make([]string, 0)
		// 	for _, n := range n.SubNodeSlice {
		// 		s = append(s, fmt.Sprintf("'%v'", n.GetFormation()))
		// 	}
		// 	return strings.Join(s, " ")
		// }(), len(n.SubNodeSlice), subContentSlice, len(subContentSlice))
		parseContentErrorSlice = append(parseContentErrorSlice, fmt.Errorf("sub node length %v is not equal 1 or sub content slice %v length %v", len(n.SubNodeSlice), subContentSlice, len(subContentSlice)))
	}
	return fileFieldContentSliceMap, parseContentErrorSlice
}

func (n *CommaNode) GetFormation() string {
	return n.Formation
}

func (n *CommaNode) GetRelateFileFieldMap() map[string]string {
	relateFileFiledMap := make(map[string]string)
	for _, subNode := range n.SubNodeSlice {
		for filename, field := range subNode.GetRelateFileFieldMap() {
			relateFileFiledMap[filename] = field
		}
	}
	return relateFileFiledMap
}

type FullstopNode struct {
	BaseNode
	Formation     string
	IsPlaceHolder bool
}

func (n *FullstopNode) CanMatch(c string) bool {
	return nodeMatcherMap[FULLSTOP].CanMatch(c)
}

func (n *FullstopNode) ParseFormation(c string) bool {
	n.Formation = c
	if c == "PH" {
		n.IsPlaceHolder = true
	} else {
		fullstopIndex := strings.IndexRune(c, markerRuneMap[FULLSTOP])
		n.Key = c[:fullstopIndex]
		n.Value = c[fullstopIndex+1:]
	}
	return true
}

func (n *FullstopNode) ParseContent(c string) (map[string]map[string][]string, []error) {
	// fmt.Printf("DEBUG: content '%v' formation is '%v.%v'\n", c, n.Key, n.Value)
	fileFieldContentSliceMap := make(map[string]map[string][]string)
	fileFieldContentSliceMap[n.Key] = make(map[string][]string)
	fileFieldContentSliceMap[n.Key][n.Value] = append(fileFieldContentSliceMap[n.Key][n.Value], c)
	return fileFieldContentSliceMap, nil
}

func (n *FullstopNode) GetFormation() string {
	return n.Formation
}

func (n *FullstopNode) GetRelateFileFieldMap() map[string]string {
	if n.IsPlaceHolder {
		return nil
	}
	return map[string]string{n.Key: n.Value}
}

type ColonNode struct {
	BaseNode
	Formation string
	KeyNode   *BracketsNode
	ValueNode Node
}

func (n *ColonNode) CanMatch(c string) bool {
	return nodeMatcherMap[COLON].CanMatch(c)
}

func (n *ColonNode) ParseFormation(c string) bool {
	n.Formation = c

	index := strings.IndexRune(c, markerRuneMap[COLON])
	if index == -1 {
		return false
	}

	if !regexp.MustCompile(`(?ms)^(?P<KEY>[^\(\)]+)\((?P<VALUE>[^\(\)]+)\)$`).MatchString(c[:index]) {
		fmt.Printf("Error: colon key '%v' does not match brackets marker\n", c[:index])
		return false
	}

	n.KeyNode = &BracketsNode{}
	n.KeyNode.ParseFormation(c[:index])
	n.ValueNode = ParseFormation(c[index+1:])

	return true
}

func (n *ColonNode) ParseContent(c string) (map[string]map[string][]string, []error) {
	// fmt.Printf("DEBUG: key %v : content '%v'\n", n.KeyNode.Formation, c)
	return n.ValueNode.ParseContent(c)
}

func (n *ColonNode) GetFormation() string {
	return n.Formation
}

func (n *ColonNode) GetKeyRelateValue() string {
	return n.KeyNode.GetRelateValue()
}

// func (n *ColonNode) GetKeyRelateFormation() string {
// 	return n.KeyNode.GetRelateFormation()
// }

func (n *ColonNode) GetKeyRelateFormationNode() Node {
	return n.KeyNode.Key
}

func (n *ColonNode) GetRelateFileFieldMap() map[string]string {
	relateFileFieldMap := make(map[string]string)
	for filename, field := range n.KeyNode.Key.GetRelateFileFieldMap() {
		relateFileFieldMap[filename] = field
	}
	for filename, field := range n.ValueNode.GetRelateFileFieldMap() {
		relateFileFieldMap[filename] = field
	}
	return relateFileFieldMap
}

type PerpendicularNode struct {
	Formation               string
	RefKeyFormationNode     Node
	RefValueSubFormationMap map[string]*ColonNode
}

func (n *PerpendicularNode) ParseFormation(c string) bool {
	n.Formation = c
	n.RefValueSubFormationMap = make(map[string]*ColonNode)
	subContentSlice := strings.Split(c, "|")
	var relateFormation string
	for _, subContent := range subContentSlice {
		if !nodeMatcherMap[COLON].CanMatch(subContent) {
			fmt.Printf("Error: perpendicular sub-content '%v' does not match colon matcher\n", subContent)
			return false
		}
		subNode := &ColonNode{}
		subNode.ParseFormation(subContent)
		if _, hasKey := n.RefValueSubFormationMap[subNode.GetKeyRelateValue()]; hasKey {
			fmt.Printf("Error: in RefValueSubFormationMap, key '%v' already exists in decoration '%v'\n", subNode.GetKeyRelateValue(), n.Formation)
			return false
		}
		// 暂时以引用的 value 作为 key，意味着一个字段只能引用一个字段的值作为分类的键
		// fmt.Printf("DEBUG: subNode.GetKeyRelateValue() = %v, subNode = %+v\n", subNode.GetKeyRelateValue(), subNode)
		// fmt.Printf("DEBUG: subContent = '%v'\n", subContent)
		n.RefValueSubFormationMap[subNode.GetKeyRelateValue()] = subNode
		if len(relateFormation) == 0 {
			n.RefKeyFormationNode = subNode.GetKeyRelateFormationNode()
		} else if relateFormation != n.RefKeyFormationNode.GetFormation() {
			fmt.Printf("Error: in RefValueSubFormationMap, key '%v' relate formation '%v' is not same as relate formation '%v'\n", subNode.GetKeyRelateValue(), n.RefKeyFormationNode.GetFormation(), relateFormation)
			return false
		}
	}
	return true
}

func (n *PerpendicularNode) GetFormation() string {
	return n.Formation
}

func (n *PerpendicularNode) GetFormationNodeByKey(key string) Node {
	return nil
}

func (n *PerpendicularNode) GetRelateFileFieldMap() map[string]string {
	relateFileFieldMap := make(map[string]string)
	for _, subNode := range n.RefValueSubFormationMap {
		for filename, field := range subNode.GetRelateFileFieldMap() {
			relateFileFieldMap[filename] = field
		}
	}
	return relateFileFieldMap
}

type BracketsNode struct {
	Formation string
	Key       Node
	Value     string
}

func (n *BracketsNode) ParseFormation(c string) {
	n.Formation = c
	bracketsRegexp := regexp.MustCompile(`(?ms)^(?P<KEY>[^\(\)]+)\((?P<VALUE>[^\(\)]+)\)$`)

	subMatchSlice := bracketsRegexp.FindStringSubmatch(c)
	for subMatchIndex, subMatchName := range bracketsRegexp.SubexpNames() {
		if subMatchName == "KEY" {
			if !nodeMatcherMap[FULLSTOP].CanMatch(subMatchSlice[subMatchIndex]) {
				fmt.Printf("Error: bracket key '%v' does not match fullstop marker\n", subMatchSlice[subMatchIndex])
				return
			}
			n.Key = nodeMatcherMap[FULLSTOP].NewNode()
			n.Key.ParseFormation(subMatchSlice[subMatchIndex])
		} else if subMatchName == "VALUE" {
			n.Value = subMatchSlice[subMatchIndex]
		}
	}
}

func (n *BracketsNode) GetRelateFormation() string {
	return n.Key.GetFormation()
}

func (n *BracketsNode) GetRelateValue() string {
	return n.Value
}

func (n *BracketsNode) GetRelateFileFieldMap() map[string]string {
	relateFileFieldMap := make(map[string]string)
	for filename, field := range n.Key.GetRelateFileFieldMap() {
		relateFileFieldMap[filename] = field
	}
	return relateFileFieldMap
}

func MergeFileFieldContentSliceMap(o, n map[string]map[string][]string) map[string]map[string][]string {
	for filename, fieldContentSliceMap := range n {
		if _, hasFile := o[filename]; !hasFile {
			o[filename] = make(map[string][]string)
		}
		for field, contentSlice := range fieldContentSliceMap {
			if _, hasField := o[filename][field]; !hasField {
				o[filename][field] = make([]string, 0)
			}
			o[filename][field] = append(o[filename][field], contentSlice...)
		}
	}
	return o
}
