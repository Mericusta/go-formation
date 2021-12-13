package main

import (
	"fmt"
	"go-formation/formation"
)

func main() {
	testDecoration()
}

func testDecoration() {
	formationExample := &formation.Formation{
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

	formationValue := formation.TraitFormation(originFormationValue)

	formationValueWithoutSpace, err := formation.TrimSpaceInString(formationValue)
	if err != nil {
		fmt.Printf("Error: trim space in string but occurs error: %v\n", err)
		return
	}

	if formation.HasDecoration(formationValueWithoutSpace) {
		formationExample.HasDecoration = true
		formationExample.DecorationNode = formation.ParseDecoration(formationValueWithoutSpace)
		if formationExample.DecorationNode == nil {
			fmt.Printf("Error: formation '%v' has decoration but decoration node is nil\n", formationValueWithoutSpace)
			return
		}

		// for k, n := range formationExample.DecorationNode.KeySubFormationMap {
		// 	fmt.Printf("DEBUG: k = %v, n = %v : %v, formation = %v\n", k, n.GetKeyRelateFormation(), n.GetKeyRelateValue(), n.Value)
		// }
	} else {
		formationExample.HasDecoration = false
		formation.ParseFormation(formationValueWithoutSpace)
	}

	testFormation(`A.b,PH;A.b,PH`, `1001,10;1002,20;1003,30`)
	testFormation(`A.b,PH;A.b,PH`, `1001,10`)
	testFormation(`A.b,PH,B.c,PH`, `1001,10,1002,20`)
	testFormation(`A.b,PH`, `1001,10`)
	testFormation(`A.b,A.b`, `1001,1002`)
	testFormation(`A.b,A.b`, `1001`)
	testFormation(`A.b`, `1001`)
}

func testFormation(configFormation, configContent string) {
	fmt.Printf("DEBUG: config formation is '%v'\n", configFormation)
	fmt.Printf("DEBUG: config content is '%v'\n", configContent)

	var configFormationNode formation.Node

	if formation.GetNodeMatcher(formation.SEMICOLON).CanMatch(configFormation) {
		configFormationNode = formation.GetNodeMatcher(formation.SEMICOLON).NewNode()
	} else if formation.GetNodeMatcher(formation.COMMA).CanMatch(configFormation) {
		configFormationNode = formation.GetNodeMatcher(formation.COMMA).NewNode()
	} else if formation.GetNodeMatcher(formation.FULLSTOP).CanMatch(configFormation) {
		configFormationNode = formation.GetNodeMatcher(formation.FULLSTOP).NewNode()
	}

	configFormationNode.ParseFormation(configFormation)
	configFormationNode.ParseContent(configContent)
}
