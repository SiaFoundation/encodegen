package imported

//go:generate encodegen -t HashAlias,Imported

type Hash [32]byte

type HashAlias Hash

type Imported struct {
	A int
	B string
	C bool
}
