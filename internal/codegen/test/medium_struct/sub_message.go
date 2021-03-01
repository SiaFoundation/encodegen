package medium_struct

type SubMessage struct {
	Id             int
	Description    string
	Strings        []string
	PointerStrings []*string
}
