package goinwx

import (
	"time"

	"fmt"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

const (
	methodDomainCheck       = "domain.check"
	methodDomainCreate      = "domain.create"
	methodDomainDelete      = "domain.delete"
	methodDomainGetPrices   = "domain.getPrices"
	methodDomainGetRules    = "domain.getRules"
	methodDomainInfo        = "domain.info"
	methodDomainList        = "domain.list"
	methodDomainLog         = "domain.log"
	methodDomainPush        = "domain.push"
	methodDomainRenew       = "domain.renew"
	methodDomainRestore     = "domain.restore"
	methodDomainStats       = "domain.stats"
	methodDomainTrade       = "domain.trade"
	methodDomainTransfer    = "domain.transfer"
	methodDomainTransferOut = "domain.transferOut"
	methodDomainUpdate      = "domain.update"
	methodDomainWhois       = "domain.whois"
)

type DomainService interface {
	Check(domains []string) ([]DomainCheckResponse, error)
	Register(request *DomainRegisterRequest) (*DomainRegisterResponse, error)
	Delete(domain string, scheduledDate time.Time) error
	Info(domain string, roId int) (*DomainInfoResponse, error)
}

type DomainServiceOp struct {
	client *Client
}

var _ DomainService = &DomainServiceOp{}

type domainCheckResponseRoot struct {
	Domains []DomainCheckResponse `mapstructure:"domain"`
}
type DomainCheckResponse struct {
	Available   int     `mapstructure:"avail"`
	Status      string  `mapstructure:"status"`
	Name        string  `mapstructure:"name"`
	Domain      string  `mapstructure:"domain"`
	TLD         string  `mapstructure:"tld"`
	CheckMethod string  `mapstructure:"checkmethod"`
	Price       float32 `mapstructure:"price"`
	CheckTime   float32 `mapstructure:"checktime"`
}

type DomainRegisterRequest struct {
	Domain        string   `structs:"domain"`
	Period        string   `structs:"period,omitempty"`
	Registrant    int      `structs:"registrant"`
	Admin         int      `structs:"admin"`
	Tech          int      `structs:"tech"`
	Billing       int      `structs:"billing"`
	Nameservers   []string `structs:"ns,omitempty"`
	TransferLock  string   `structs:"transferLock,omitempty"`
	RenewalMode   string   `structs:"renewalMode,omitempty"`
	WhoisProvider string   `structs:"whoisProvider,omitempty"`
	WhoisUrl      string   `structs:"whoisUrl,omitempty"`
	ScDate        string   `structs:"scDate,omitempty"`
	ExtDate       string   `structs:"extDate,omitempty"`
	Asynchron     string   `structs:"asynchron,omitempty"`
	Voucher       string   `structs:"voucher,omitempty"`
	Testing       string   `structs:"testing,omitempty"`
}

type DomainRegisterResponse struct {
	RoId     int
	Price    float32
	Currency string
}

type DomainInfoResponse struct {
	RoId         int                `mapstructure:"roId"`
	Domain       string             `mapstructure:"domain"`
	DomainAce    string             `mapstructure:"domainAce"`
	Period       string             `mapstructure:"period"`
	CrDate       time.Time          `mapstructure:"crDate"`
	ExDate       time.Time          `mapstructure:"exDate"`
	UpDate       time.Time          `mapstructure:"upDate"`
	ReDate       time.Time          `mapstructure:"reDate"`
	ScDate       time.Time          `mapstructure:"scDate"`
	TransferLock int                `mapstructure:"transferLock"`
	Status       string             `mapstructure:"status"`
	AuthCode     string             `mapstructure:"authCode"`
	RenewalMode  string             `mapstructure:"renewalMode"`
	TransferMode string             `mapstructure:"transferMode"`
	Registrant   int                `mapstructure:"registrant"`
	Admin        int                `mapstructure:"admin"`
	Tech         int                `mapstructure:"tech"`
	Billing      int                `mapstructure:"billing"`
	Ns           []string           `mapstructure:"ns"`
	NoDelegation string             `mapstructure:"noDelegation"`
	Contacts     map[string]Contact `mapstructure:"contact"`
}

type Contact struct {
	RoId          int
	Id            int
	Type          string
	Name          string
	Org           string
	Street        string
	City          string
	PostalCode    string `mapstructure:"pc"`
	StateProvince string `mapstructure:"sp"`
	Country       string `mapstructure:"cc"`
	Phone         string `mapstructure:"voice"`
	Fax           string
	Email         string
	Remarks       string
	Protection    int
}

func (s *DomainServiceOp) Check(domains []string) ([]DomainCheckResponse, error) {
	req := s.client.NewRequest(methodDomainCheck, map[string]interface{}{
		"domain": domains,
		"wide":   "2",
	})

	resp, err := s.client.Do(*req)
	if err != nil {
		return nil, err
	}

	root := new(domainCheckResponseRoot)
	err = mapstructure.Decode(*resp, &root)
	if err != nil {
		return nil, err
	}

	return root.Domains, nil
}

func (s *DomainServiceOp) Register(request *DomainRegisterRequest) (*DomainRegisterResponse, error) {
	req := s.client.NewRequest(methodDomainCreate, structs.Map(request))

	//fmt.Println("Args", req.Args)
	resp, err := s.client.Do(*req)
	if err != nil {
		return nil, err
	}

	var result DomainRegisterResponse
	err = mapstructure.Decode(*resp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *DomainServiceOp) Delete(domain string, scheduledDate time.Time) error {
	req := s.client.NewRequest(methodDomainDelete, map[string]interface{}{
		"domain": domain,
		"scDate": scheduledDate.Format(time.RFC3339),
	})

	_, err := s.client.Do(*req)

	return err
}

func (s *DomainServiceOp) Info(domain string, roId int) (*DomainInfoResponse, error) {
	req := s.client.NewRequest(methodDomainInfo, map[string]interface{}{
		"domain": domain,
		"wide":   "2",
	})
	if roId != 0 {
		req.Args["roId"] = roId
	}

	resp, err := s.client.Do(*req)
	if err != nil {
		return nil, err
	}

	var result DomainInfoResponse
	err = mapstructure.Decode(*resp, &result)
	if err != nil {
		return nil, err
	}
	fmt.Println("Response", result)

	return &result, nil
}
