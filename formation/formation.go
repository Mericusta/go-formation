package formation

import (
	"fmt"
	"go-formation/utility"
)

type Formation struct {
	File           string
	Field          string
	HasDecoration  bool
	DecorationNode *PerpendicularNode
	FormationNode  Node
}

func (f *Formation) GetRelateFileFieldMap() map[string]string {
	if f.HasDecoration {
		return f.DecorationNode.GetRelateFileFieldMap()
	}
	return f.FormationNode.GetRelateFileFieldMap()
}

func (f *Formation) RelationCheck(gameDataJsonObjectMap map[string]*GameDataJsonObject) (bool, []error) {
	gameDataJsonObject, hasGameDataJsonObject := gameDataJsonObjectMap[f.File]
	if gameDataJsonObject == nil || !hasGameDataJsonObject {
		return false, []error{fmt.Errorf("file %v game data json object is nil", f.File)}
	}
	checkDataIndex, hasField := gameDataJsonObject.Format[f.Field]
	if checkDataIndex < 0 || !hasField {
		return false, []error{fmt.Errorf("file %v field %v index %v is invalid", f.File, f.Field, checkDataIndex)}
	}

	if f.HasDecoration {
		return f.relationWithDecorationCheck(gameDataJsonObjectMap)
	}
	// fmt.Printf("DEBUG: relationCheck = %v\n", gameDataJsonObjectMap)
	return f.relationCheck(gameDataJsonObjectMap)
}

func (f *Formation) relationWithDecorationCheck(gameDataJsonObjectMap map[string]*GameDataJsonObject) (bool, []error) {
	gameDataJsonObject := gameDataJsonObjectMap[f.File]
	checkDataIndex := gameDataJsonObject.Format[f.Field]
	relationCheckErrorSlice := make([]error, 0)
	relateFileFieldContentSliceMap, traitErrorSlice := traitRelateFileFieldContentSliceMapWithDecorationNode(
		gameDataJsonObjectMap,
		gameDataJsonObject,
		checkDataIndex,
		f.DecorationNode.RefKeyFormationNode.GetKey(),
		f.DecorationNode.RefKeyFormationNode.GetValue(),
		f.DecorationNode.RefValueSubFormationMap,
		f.File, f.Field,
	)
	relationCheckErrorSlice = append(relationCheckErrorSlice, traitErrorSlice...)
	// fmt.Printf("DEBUG: ref file %v, field %v\n", f.DecorationNode.RefKeyFormationNode.GetKey(), f.DecorationNode.RefKeyFormationNode.GetValue())

	ok, checkErrorSlice := relationCheckHandle(relateFileFieldContentSliceMap, gameDataJsonObjectMap)
	relationCheckErrorSlice = append(relationCheckErrorSlice, checkErrorSlice...)
	return ok, relationCheckErrorSlice
}

func traitRelateFileFieldContentSliceMapWithDecorationNode(
	gameDataJsonObjectMap map[string]*GameDataJsonObject,
	gameDataJsonObject *GameDataJsonObject,
	checkDataIndex int,
	refFile, refField string,
	refValueSubFormationMap map[string]*ColonNode,
	traitFile, traitField string,
) (map[string]map[string][]string, []error) {
	relateFileFieldContentSliceMap := make(map[string]map[string][]string)
	traitRelateFileFieldContentSliceMapErrorSlice := make([]error, 0)
	refGameDataJsonObject, hasRefFile := gameDataJsonObjectMap[refFile]
	if refGameDataJsonObject == nil || !hasRefFile {
		return nil, []error{fmt.Errorf("file %v field %v reference file %v field %v game data json object is nil", traitFile, traitField, refFile, refField)}
	}
	refIndex, hasRefField := refGameDataJsonObject.Format[refField]
	if !hasRefField {
		return nil, []error{fmt.Errorf("file %v field %v reference file %v field %v index does not exist in Format %v", traitFile, traitField, refFile, refField, refGameDataJsonObject.Format)}
	}

	// fmt.Printf("DEBUG: checkDataIndex is = %v, refFile = %v, refField = %v, refIndex = %v, traitFile = %v, traitField = %v\n", checkDataIndex, refFile, refField, refIndex, traitFile, traitField)

	for _, rowDataSlice := range gameDataJsonObject.Data {
		refValue := fmt.Sprintf("%v", rowDataSlice[refIndex])
		checkData := fmt.Sprintf("%v", rowDataSlice[checkDataIndex])
		// fmt.Printf("DEBUG: row %v data is %v\n", row, rowDataSlice)
		if len(checkData) == 0 {
			continue
		}
		refNode, hasRefNode := refValueSubFormationMap[refValue]
		if !hasRefNode {
			fmt.Printf("Error: %v.%v reference value %v from %v.%v does not exist\n", traitFile, traitField, refValue, refFile, refField)
			traitRelateFileFieldContentSliceMapErrorSlice = append(traitRelateFileFieldContentSliceMapErrorSlice, fmt.Errorf("%v.%v reference value %v from %v.%v does not exist", traitFile, traitField, refValue, refFile, refField))
			continue
		}
		rowContentResult, parseContentErrorSlice := refNode.ParseContent(checkData)
		// fmt.Printf("DEBUG: %v.%v row %v check data index is %v, data is '%v', rowContentResult is '%v'\n", traitFile, traitField, row, checkDataIndex, checkData, rowContentResult)
		relateFileFieldContentSliceMap = MergeFileFieldContentSliceMap(relateFileFieldContentSliceMap, rowContentResult)
		traitRelateFileFieldContentSliceMapErrorSlice = append(traitRelateFileFieldContentSliceMapErrorSlice, parseContentErrorSlice...)
	}
	return relateFileFieldContentSliceMap, traitRelateFileFieldContentSliceMapErrorSlice
}

func (f *Formation) relationCheck(gameDataJsonObjectMap map[string]*GameDataJsonObject) (bool, []error) {
	gameDataJsonObject := gameDataJsonObjectMap[f.File]
	checkDataIndex := gameDataJsonObject.Format[f.Field]
	relationCheckErrorSlice := make([]error, 0)

	relateFileFieldContentSliceMap, traitErrorSlice := traitRelateFileFieldContentSliceMap(gameDataJsonObject, checkDataIndex, f.FormationNode, f.File, f.Field)
	relationCheckErrorSlice = append(relationCheckErrorSlice, traitErrorSlice...)

	// fmt.Printf("DEBUG: relateFileFieldContentSliceMap = %v\n", relateFileFieldContentSliceMap)

	ok, checkErrorSlice := relationCheckHandle(relateFileFieldContentSliceMap, gameDataJsonObjectMap)
	relationCheckErrorSlice = append(relationCheckErrorSlice, checkErrorSlice...)

	return ok, relationCheckErrorSlice
}

func traitRelateFileFieldContentSliceMap(
	gameDataJsonObject *GameDataJsonObject,
	checkDataIndex int,
	formationNode Node,
	traitFile, traitField string,
) (map[string]map[string][]string, []error) {
	relateFileFieldContentSliceMap := make(map[string]map[string][]string)
	traitRelateFileFieldContentSliceMapErrorSlice := make([]error, 0)
	for _, rowDataSlice := range gameDataJsonObject.Data {
		// fmt.Printf("DEBUG: row %v data is %v\n", row, rowDataSlice)
		checkData := fmt.Sprintf("%v", rowDataSlice[checkDataIndex])
		// NOTE: 零值引用忽略，由空 json 解出来的 int 值为0，需要忽略
		if len(checkData) == 0 || checkData == "0" || checkData == "-1" {
			// fmt.Printf("DEBUG: row %v continue\n", row)
			continue
		}
		// fmt.Printf("DEBUG: formationNode.GetFormation() = %v\n", formationNode.GetFormation())
		rowContentResult, parseContentErrorSlice := formationNode.ParseContent(checkData)
		// fmt.Printf("DEBUG: %v.%v row %v check data index is %v, data is '%v', rowContentResult is '%v'\n", traitFile, traitField, row, checkDataIndex, checkData, rowContentResult)
		relateFileFieldContentSliceMap = MergeFileFieldContentSliceMap(relateFileFieldContentSliceMap, rowContentResult)
		traitRelateFileFieldContentSliceMapErrorSlice = append(traitRelateFileFieldContentSliceMapErrorSlice, parseContentErrorSlice...)
	}
	return relateFileFieldContentSliceMap, traitRelateFileFieldContentSliceMapErrorSlice
}

func relationCheckHandle(relateFileFieldContentSliceMap map[string]map[string][]string, gameDataJsonObjectMap map[string]*GameDataJsonObject) (bool, []error) {
	// fmt.Printf("DEBUG: relateFileFieldContentSliceMap = %v\n", relateFileFieldContentSliceMap)
	relationCheckErrorSlice := make([]error, 0)
	for relateFilename, relateFieldContentSliceMap := range relateFileFieldContentSliceMap {
		relateGameDataJsonObject, hasRelateFile := gameDataJsonObjectMap[relateFilename]
		if relateGameDataJsonObject == nil || !hasRelateFile {
			// return false, fmt.Errorf("relate file %v game data json object is nil", relateFilename)
			relationCheckErrorSlice = append(relationCheckErrorSlice, fmt.Errorf("relate file %v game data json object is nil", relateFilename))
			continue
		}
		for relateField, contentSlice := range relateFieldContentSliceMap {
			relateFieldIndex, hasRelateField := relateGameDataJsonObject.Format[relateField]
			if relateFieldIndex < 0 || !hasRelateField {
				// fmt.Printf("DEBUG: relate %v.%v index %v is invalid, format = %v\n", relateFilename, relateField, relateFieldIndex, relateGameDataJsonObject.Format)
				relationCheckErrorSlice = append(relationCheckErrorSlice, fmt.Errorf("relate %v.%v index %v is invalid, format = %v", relateFilename, relateField, relateFieldIndex, relateGameDataJsonObject.Format))
				continue
			}
			for _, content := range contentSlice {
				exists := false
				for _, relateDataSlice := range relateGameDataJsonObject.Data {
					// fmt.Printf("DEBUG: check content '%v' from file %v field %v index %v from relateDataSlice '%v', relateDataSlice[%v] = '%v'\n", content, relateFilename, relateField, relateFieldIndex, relateDataSlice, relateFieldIndex, relateDataSlice[relateFieldIndex])
					if utility.CompareGameDataJsonObjectData(relateDataSlice[relateFieldIndex], content) {
						// fmt.Printf("DEBUG: content '%v' exists\n", content)
						exists = true
						break
					}
				}
				if !exists {
					fmt.Printf("Error: %v.%v can not find content %v\n", relateFilename, relateField, content)
					relationCheckErrorSlice = append(relationCheckErrorSlice, fmt.Errorf("%v.%v can not find content %v", relateFilename, relateField, content))
					continue
				}
				// fmt.Printf("DEBUG: content %v can be found in relate file %v field %v index %v\n", content, relateFilename, relateField, relateFieldIndex)
			}
		}
	}

	return len(relationCheckErrorSlice) == 0, relationCheckErrorSlice
}

func testDecoration() {
	formation := &Formation{
		File:  "RewardsPoolCfg",
		Field: "fixed_jackpot",
	}

	originFormationValue := `format(
		RewardsPoolCfg.reward_type(1):RewardsGroupCfg.group_id,PH;RewardsGroupCfg.group_id,PH|
		RewardsPoolCfg.reward_type(2):RewardsGroupCfg.group_id,PH;RewardsGroupCfg.group_id,PH|
		RewardsPoolCfg.reward_type(3):RewardsGroupCfg.group_id,PH;RewardsGroupCfg.group_id,PH|
		RewardsPoolCfg.reward_type(4):RewardsGroupCfg.group_id,PH,PH;RewardsGroupCfg.group_id,PH,PH|
		RewardsPoolCfg.reward_type(5):RewardsGroupCfg.group_id,PH,PH;RewardsGroupCfg.group_id,PH,PH
	)`

	formationValue := TraitFormation(originFormationValue)

	formationValueWithoutSpace, err := TrimSpaceInString(formationValue)
	if err != nil {
		fmt.Printf("Error: trim space in string but occurs error: %v\n", err)
		return
	}

	if HasDecoration(formationValueWithoutSpace) {
		formation.HasDecoration = true
		formation.DecorationNode = ParseDecoration(formationValueWithoutSpace)
		if formation.DecorationNode == nil {
			fmt.Printf("Error: formation '%v' has decoration but decoration node is nil\n", formationValueWithoutSpace)
			return
		}

		// for k, n := range formation.DecorationNode.KeySubFormationMap {
		// 	fmt.Printf("DEBUG: k = %v, n = %v : %v, formation = %v\n", k, n.GetKeyRelateFormation(), n.GetKeyRelateValue(), n.Value)
		// }
	} else {
		formation.HasDecoration = false
		ParseFormation(formationValueWithoutSpace)
	}
}

// testFormation(`A.b,PH;A.b,PH`, `1001,10;1002,20;1003,30`)
// testFormation(`A.b,PH;A.b,PH`, `1001,10`)
// testFormation(`A.b,PH,B.c,PH`, `1001,10,1002,20`)
// testFormation(`A.b,PH`, `1001,10`)
// testFormation(`A.b,A.b`, `1001,1002`)
// testFormation(`A.b,A.b`, `1001`)
// testFormation(`A.b`, `1001`)

// func testFormation(configFormation, configContent string) {
// 	fmt.Printf("DEBUG: config formation is '%v'\n", configFormation)
// 	fmt.Printf("DEBUG: config content is '%v'\n", configContent)

// 	var configFormationNode Node

// 	if nodeMatcherMap[SEMICOLON].CanMatch(configFormation) {
// 		configFormationNode = nodeMatcherMap[SEMICOLON].NewNode()
// 	} else if nodeMatcherMap[COMMA].CanMatch(configFormation) {
// 		configFormationNode = nodeMatcherMap[COMMA].NewNode()
// 	} else if nodeMatcherMap[FULLSTOP].CanMatch(configFormation) {
// 		configFormationNode = nodeMatcherMap[FULLSTOP].NewNode()
// 	}

// 	configFormationNode.ParseFormation(configFormation)
// 	configFormationNode.ParseContent(configContent)
// }
