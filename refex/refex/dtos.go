package refex

type match struct {
	startPos int
	endPos   int
	pattern  string
	parts    []*part
}

type part struct {
	code  string
	isArg bool
	index int
}
