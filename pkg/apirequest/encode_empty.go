package apirequest

func EncodeEmpty() Request {
	return Request{
		Method: "GET",
		MIME:   "text/plain",
		Body:   nil,
	}
}
