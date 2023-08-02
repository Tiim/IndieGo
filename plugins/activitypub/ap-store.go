package activitypub

import (
	"fmt"
	"strings"

	"github.com/go-ap/activitypub"
)

type apStore struct {
	baseUrl         string
	actorName       string
	actorProfileUrl string
}

func (s *apStore) getActorFromName(name string) (*activitypub.Actor, error) {

	if strings.ToLower(name) != strings.ToLower(s.actorName) {

		fmt.Printf("expected: '%s', gotten: '%s'", s.actorName, name)
		return nil, fmt.Errorf("Actor not found: %s", name)
	}

	actor := activitypub.ActorNew(
		activitypub.IRI(s.baseUrl+"/ap/users/"+s.actorName),
		activitypub.PersonType,
	)

	actor.PreferredUsername.Set(activitypub.DefaultLang, activitypub.Content(s.actorName))
	if s.actorProfileUrl != "" {
		actor.URL = activitypub.IRI(s.actorProfileUrl)
	}
	actor.Name.Set(activitypub.DefaultLang, activitypub.Content(s.actorName))

	actor.Inbox = activitypub.IRI(s.baseUrl + "/ap/users/" + s.actorName + "/inbox")
	actor.Outbox = activitypub.IRI(s.baseUrl + "/ap/users/" + s.actorName + "/outbox")
	actor.Following = activitypub.IRI(s.baseUrl + "/ap/users/" + s.actorName + "/following")
	actor.Followers = activitypub.IRI(s.baseUrl + "/ap/users/" + s.actorName + "/followers")

	//set actor.PublicKey
    actor.Endpoints = &activitypub.Endpoints{}
	actor.Endpoints.SharedInbox = activitypub.IRI(s.baseUrl + "/ap/shared-inbox")
	return actor, nil
}
