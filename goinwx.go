package goinwx

import (
	"log"
	"net/url"
	"os"

	"fmt"

	"github.com/andrexus/xmlrpc"
	"github.com/hashicorp/logutils"
)

const (
	libraryVersion    = "0.3.2"
	logLevelEnvName   = "GOINWX_LOG"
	APIBaseUrl        = "https://api.domrobot.com/xmlrpc/"
	APISandboxBaseUrl = "https://api.ote.domrobot.com/xmlrpc/"
	APILanguage       = "eng"
)

func init() {
	logLevel := os.Getenv(logLevelEnvName)
	if logLevel == "" {
		logLevel = "INFO"
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

// Client manages communication with INWX API.
type Client struct {
	// HTTP client used to communicate with the INWX API.
	RPCClient *xmlrpc.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// API username
	Username string

	// API password
	Password string

	// User agent for client
	APILanguage string

	// Services used for communicating with the API
	Account     AccountService
	Domains     DomainService
	Nameservers NameserverService
	Contacts    ContactService
}

type ClientOptions struct {
	Sandbox bool
}

type Request struct {
	ServiceMethod string
	Args          map[string]interface{}
}

// Response is a INWX API response. This wraps the standard http.Response returned from INWX.
type Response struct {
	Code         int                    `xmlrpc:"code"`
	Message      string                 `xmlrpc:"msg"`
	ReasonCode   string                 `xmlrpc:"reasonCode"`
	Reason       string                 `xmlrpc:"reason"`
	ResponseData map[string]interface{} `xmlrpc:"resData"`
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	Code       int    `xmlrpc:"code"`
	Message    string `xmlrpc:"msg"`
	ReasonCode string `xmlrpc:"reasonCode"`
	Reason     string `xmlrpc:"reason"`
}

// NewClient returns a new INWX API client.
func NewClient(username, password string, opts *ClientOptions) *Client {
	useSandbox := false

	if opts != nil {
		useSandbox = opts.Sandbox
	}

	var baseURL *url.URL

	if useSandbox == true {
		baseURL, _ = url.Parse(APISandboxBaseUrl)

	} else {
		baseURL, _ = url.Parse(APIBaseUrl)
	}

	rpcClient, _ := xmlrpc.NewClient(baseURL.String(), nil)

	log.Printf("[DEBUG] Base URL: %s\n", baseURL)

	client := &Client{RPCClient: rpcClient,
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}

	client.Account = &AccountServiceOp{client: client}
	client.Domains = &DomainServiceOp{client: client}
	client.Nameservers = &NameserverServiceOp{client: client}
	client.Contacts = &ContactServiceOp{client: client}

	return client
}

// NewRequest creates an API request.
func (c *Client) NewRequest(serviceMethod string, args map[string]interface{}) *Request {
	if args != nil {
		args["lang"] = APILanguage
	}
	req := &Request{ServiceMethod: serviceMethod, Args: args}
	return req
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req Request) (*map[string]interface{}, error) {
	var resp Response
	err := c.RPCClient.Call(req.ServiceMethod, req.Args, &resp)
	log.Printf("[DEBUG] Response: %v", resp)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(&resp)
	return &resp.ResponseData, err
}

func (r *ErrorResponse) Error() string {
	if r.Reason != "" {
		return fmt.Sprintf("(%d) %s. Reason: (%s) %s",
			r.Code, r.Message, r.ReasonCode, r.Reason)
	}
	return fmt.Sprintf("(%d) %s", r.Code, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *Response) error {
	if c := r.Code; c >= 1000 && c <= 1500 {
		return nil
	}

	errorResponse := &ErrorResponse{Code: r.Code, Message: r.Message, Reason: r.Reason, ReasonCode: r.ReasonCode}

	return errorResponse
}
