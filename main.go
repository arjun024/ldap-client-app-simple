package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-ldap/ldap"
)

func ldap_connect() string {
	// Free test server
	ldapServer := "ldap.forumsys.com:389"
	bindUsername := "uid=tesla,dc=example,dc=com"
	bindPassword := "password"

	conn, err := ldap.Dial("tcp", ldapServer)
	if err != nil {
		log.Fatalf("Error connecting to LDAP server: %v", err)
	}
	defer conn.Close()

	if err := conn.Bind(bindUsername, bindPassword); err != nil {
		log.Fatalf("Error binding to LDAP server: %v", err)
	}

	// Search for entries in the specified base DN
	searchBaseDN := "dc=example,dc=com"
	searchFilter := "(objectClass=*)"
	searchRequest := ldap.NewSearchRequest(
		searchBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		nil,
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("Error searching LDAP directory: %v", err)
	}

	var buf bytes.Buffer
	for _, entry := range searchResult.Entries {
		fmt.Fprintf(&buf, "DN: %s\n", entry.DN)
		for _, attr := range entry.Attributes {
			fmt.Fprintf(&buf, "%s: %s\n", attr.Name, attr.Values)
		}
		fmt.Fprintln(&buf)
	}
	return buf.String()
}

func main() {
	http.HandleFunc("/", hello)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	ldap_response := ldap_connect()
	fmt.Fprintln(res, "go server start")
	fmt.Fprintln(res, ldap_response)
}
