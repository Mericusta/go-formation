package formation

import (
	"fmt"
	"regexp"
)

var markerRuneMap map[MarkerType]rune
var nodeMatcherMap map[MarkerType]*Matcher

func init() {
	markerRuneMap = map[MarkerType]rune{
		COMMA:         ',',
		FULLSTOP:      '.',
		COLON:         ':',
		SEMICOLON:     ';',
		PERPENDICULAR: '|',
	}

	nodeMatcherMap = map[MarkerType]*Matcher{
		SEMICOLON: {
			t: SEMICOLON,
			CanMatch: func(s string) bool {
				return regexp.MustCompile(`(?ms)^[^;\|\s]+(;[^;\|\s]+)+$`).MatchString(s)
			},
			NewNode: func() Node {
				return &SemicolonNode{}
			},
		},
		COMMA: {
			t: COMMA,
			CanMatch: func(s string) bool {
				return regexp.MustCompile(`(?ms)^[-_\.\w]+(,[-_\.\w]+)+$`).MatchString(s)
			},
			NewNode: func() Node {
				return &CommaNode{}
			},
		},
		FULLSTOP: {
			t: FULLSTOP,
			CanMatch: func(s string) bool {
				return regexp.MustCompile(`(?ms)^(?P<KEY>[-_\w]+)\.(?P<VALUE>[-_\w]+)$`).MatchString(s)
			},
			NewNode: func() Node {
				return &FullstopNode{}
			},
		},
		COLON: {
			t: COLON,
			CanMatch: func(s string) bool {
				return regexp.MustCompile(`(?ms)^(?P<KEY>[^:\s]+):(?P<VALUE>[^:\s]+)$`).MatchString(s)
			},
			NewNode: func() Node {
				return &ColonNode{}
			},
		},
	}
}

// func HasFormation(c string) bool {
// 	regexp.MustCompile(`(?ms)^format\((?P<VALUE>.*)\)$`).MatchString()
// }

func GetNodeMatcher(t MarkerType) *Matcher {
	return nodeMatcherMap[t]
}

func HasDecoration(c string) bool {
	return regexp.MustCompile(`(?ms)^[^\|\s]+(\|[^\|\s]+)+$`).MatchString(c)
}

func ParseDecoration(c string) *PerpendicularNode {
	decorationNode := &PerpendicularNode{}
	decorationNode.ParseFormation(c)
	return decorationNode
}

func TrimSpaceInString(content string) (string, error) {
	SpaceExpression := `\s+`
	SpaceRegexp := regexp.MustCompile(SpaceExpression)
	if SpaceRegexp == nil {
		return content, fmt.Errorf("SpaceRegexp is nil")
	}
	return SpaceRegexp.ReplaceAllString(content, ""), nil
}

func TraitFormation(c string) string {
	formatRegexp := regexp.MustCompile(`(?ms)format\((?P<VALUE>.*)\)`)
	subMatchSlice := formatRegexp.FindStringSubmatch(c)
	// fmt.Printf("DEBUG: subMatchSlice = %v\n", subMatchSlice)
	if len(subMatchSlice) == 0 {
		return ""
	} else {
		for matchIndex, matchName := range formatRegexp.SubexpNames() {
			if matchIndex == 0 {
				continue
			} else if matchName == "VALUE" {
				return subMatchSlice[matchIndex]
			}
		}
	}
	return ""
}

func ParseFormation(c string) Node {
	var configFormationNode Node
	if nodeMatcherMap[SEMICOLON].CanMatch(c) {
		configFormationNode = nodeMatcherMap[SEMICOLON].NewNode()
	} else if nodeMatcherMap[COMMA].CanMatch(c) {
		configFormationNode = nodeMatcherMap[COMMA].NewNode()
	} else if nodeMatcherMap[FULLSTOP].CanMatch(c) {
		configFormationNode = nodeMatcherMap[FULLSTOP].NewNode()
	} else {
		fmt.Printf("Errorf: formation '%v' can not match any marker\n", c)
		return nil
	}
	configFormationNode.ParseFormation(c)
	return configFormationNode
}
