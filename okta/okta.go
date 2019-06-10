package okta

import (
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	log "github.com/sirupsen/logrus"
)

const (
	ReadOnlyOktaGroup = "global-read-only "
)

func AddGroup(client *okta.Client, name string) {
	gp := &okta.GroupProfile{
		Name: name,
	}
	g := &okta.Group{
		Profile: gp,
	}
	_, _, err := client.Group.CreateGroup(*g, nil)

	if err != nil {
		log.Printf("Error creating group : %v in Okta", name)
	} else {
		log.Printf("Created group : %v in Okta", name)
	}
}

func DeleteGroup(client *okta.Client, name string) {
	groups, _, err := client.Group.ListGroups(query.NewQueryParams(query.WithQ(name)))
	if err != nil || len(groups) == 0 {
		log.Printf("Error searching group : %v in Okta", name)
	}

	group := groups[0]
	_, err = client.Group.DeleteGroup(group.Id, nil)

	if err != nil {
		log.Printf("Error deleting group : %v in Okta", name)
	} else {
		log.Printf("Deleted group : %v in Okta", name)
	}

}
