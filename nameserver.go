package goinwx

import (
	"errors"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

const (
	methodNameserverCheck        = "nameserver.check"
	methodNameserverCreate       = "nameserver.create"
	methodNameserverCreateRecord = "nameserver.createRecord"
	methodNameserverDelete       = "nameserver.delete"
	methodNameserverDeleteRecord = "nameserver.deleteRecord"
	methodNameserverInfo         = "nameserver.info"
	methodNameserverList         = "nameserver.list"
	methodNameserverUpdate       = "nameserver.update"
	methodNameserverUpdateRecord = "nameserver.updateRecord"
)

type NameserverService interface {
	Check(domain string, nameservers []string) (*NameserverCheckResponse, error)
	Info(domain string, roId int) (*NamserverInfoResponse, error)
	CreateRecord(*NameserverRecordRequest) (int, error)
	UpdateRecord(recordId int, request *NameserverRecordRequest) error
	DeleteRecord(recordId int) error
}

type NameserverServiceOp struct {
	client *Client
}

var _ NameserverService = &NameserverServiceOp{}

type NameserverCheckResponse struct {
	Details []string
	Status  string
}

type NameserverRecordRequest struct {
	RoId                   int    `structs:"roId,omitempty"`
	Domain                 string `structs:"domain,omitempty"`
	Type                   string `structs:"type"`
	Content                string `structs:"content"`
	Name                   string `structs:"name,omitempty"`
	Ttl                    int    `structs:"ttl,omitempty"`
	Prio                   int    `structs:"prio,omitempty"`
	UrlRedirectType        string `structs:"urlRedirectType,omitempty"`
	UrlRedirectTitle       string `structs:"urlRedirectTitle,omitempty"`
	UrlRedirectDescription string `structs:"urlRedirectDescription,omitempty"`
	UrlRedirectFavIcon     string `structs:"urlRedirectFavIcon,omitempty"`
	UrlRedirectKeywords    string `structs:"urlRedirectKeywords,omitempty"`
}

type NamserverInfoResponse struct {
}

func (s *NameserverServiceOp) Check(domain string, nameservers []string) (*NameserverCheckResponse, error) {
	req := s.client.NewRequest(methodNameserverCheck, map[string]interface{}{
		"domain": domain,
		"ns":     nameservers,
	})

	resp, err := s.client.Do(*req)
	if err != nil {
		return nil, err
	}

	var result NameserverCheckResponse
	err = mapstructure.Decode(*resp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *NameserverServiceOp) Info(domain string, roId int) (*NamserverInfoResponse, error) {
	requestMap := map[string]interface{}{
		"domain": domain,
	}
	if roId != 0 {
		requestMap["roId"] = roId
	}
	req := s.client.NewRequest(methodNameserverInfo, requestMap)

	resp, err := s.client.Do(*req)
	if err != nil {
		return nil, err
	}
	var result NamserverInfoResponse
	err = mapstructure.Decode(*resp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *NameserverServiceOp) CreateRecord(request *NameserverRecordRequest) (int, error) {
	req := s.client.NewRequest(methodNameserverCreateRecord, structs.Map(request))

	//fmt.Println("Args", req.Args)
	resp, err := s.client.Do(*req)
	if err != nil {
		return 0, err
	}

	var result map[string]int
	err = mapstructure.Decode(*resp, &result)
	if err != nil {
		return 0, err
	}

	return result["id"], nil
}

func (s *NameserverServiceOp) UpdateRecord(recordId int, request *NameserverRecordRequest) error {
	if request == nil {
		return errors.New("Request can't be nil")
	}
	requestMap := structs.Map(request)
	requestMap["id"] = recordId

	req := s.client.NewRequest(methodNameserverUpdateRecord, requestMap)

	_, err := s.client.Do(*req)
	if err != nil {
		return err
	}

	return nil
}

func (s *NameserverServiceOp) DeleteRecord(recordId int) error {
	req := s.client.NewRequest(methodNameserverDeleteRecord, map[string]interface{}{
		"id": recordId,
	})

	_, err := s.client.Do(*req)
	if err != nil {
		return err
	}

	return nil
}
