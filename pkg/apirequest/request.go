package apirequest

type Request struct {
	Method string
	MIME   string
	Body   []byte
}
