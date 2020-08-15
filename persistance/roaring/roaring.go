package roaring

type entry struct {
	key       uint16
	container container
}

type Roaring struct {
	containers []entry
}
