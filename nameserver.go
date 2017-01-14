package goinwx

type NameserverService interface {
}

type NameserverServiceOp struct {
	client *Client
}

var _ NameserverService = &NameserverServiceOp{}
