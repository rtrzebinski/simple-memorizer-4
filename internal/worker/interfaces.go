package worker

type Writer interface {
	StoreResult(*Result) error
}
