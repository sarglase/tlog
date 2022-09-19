package hook

type Hook interface {
	Write([]byte)
}
