package goinwx

import (
	"fmt"
	"os"
	"testing"
)

const (
	testClientUsernameEnvName = "GOINWX_USERNAME"
	testClientPasswordEnvName = "GOINWX_PASSWORD"
)

var (
	client *Client
)

func setup(t *testing.T) {
	client = NewClient(os.Getenv(testClientUsernameEnvName), os.Getenv(testClientPasswordEnvName))
	err := client.Login()
	if err != nil {
		t.Error(err)
	}
}

func teardown() {
	client.Logout()
}

func TestDomainCheck(t *testing.T) {
	setup(t)

	defer teardown()

	items, err := client.Domains.Check([]string{"foobar.com", "golang-meets-inwx.com"})
	if err != nil {
		t.Error(err)
	}

	for _, item := range items {
		fmt.Printf("Domain: %s\n", item.Domain)
		fmt.Printf("Available: %d\n", item.Available)
		fmt.Printf("Status: %s\n", item.Status)
		fmt.Printf("Check time: %f\n", item.CheckTime)
		fmt.Printf("Name: %s\n", item.Name)
		fmt.Printf("TLD: %s\n", item.TLD)
		fmt.Printf("Check method: %s\n", item.CheckMethod)
		fmt.Printf("Price: %f\n", item.Price)

		fmt.Println()
	}
}

func TestDomainRegister(t *testing.T) {
	setup(t)

	defer teardown()

	request := &DomainRegisterRequest{
		Domain: "golang-meets-inwx.com",
		Registrant: 1081859,
		Admin: 1081859,
		Tech: 1081859,
		Billing: 1081859,
		Nameservers: []string{"ns.ote.inwx.de", "ns2.ote.inwx.de"},
	}
	result, err := client.Domains.Register(request)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Domain register result", result)
}
