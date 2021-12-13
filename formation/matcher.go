package formation

type Matcher struct {
	t        MarkerType
	CanMatch func(string) bool
	NewNode  func() Node
}
