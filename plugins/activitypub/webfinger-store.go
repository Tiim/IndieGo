package activitypub

import (
	"fmt"

	"github.com/go-ap/activitypub"
)

type webfingerActorStore struct {
	baseUrl         string
	actorProfileUrl string
	actorName       string
	host            string
}

func (d *webfingerActorStore) Load(iri activitypub.IRI) (activitypub.Item, error) {
	fmt.Printf("iri: %s\n", iri)
	url, err := iri.URL()
	if err != nil {
		return nil, err
	}
	if url.Path == "/" && url.Hostname() == d.host {
		fmt.Println("Service actor")
		return d.buildServiceActor()
	} else {
		fmt.Printf("path %s, host: %s ref %s\n", url.Path, url.Hostname(), d.host)
		fmt.Println("User actor")
		return d.buildPersonActor()
	}

}

func (d *webfingerActorStore) buildServiceActor() (activitypub.Item, error) {
	actor := activitypub.ActorNew(
		activitypub.IRI(d.baseUrl+"/ap"),
		activitypub.ServiceType,
	)

	actor.PreferredUsername.Set(activitypub.DefaultLang, activitypub.Content("IndieGo Server"))
	actor.URL = actor.ID

	actor.Name.Set(activitypub.DefaultLang, activitypub.Content("indiego"))

	return actor, nil
}

func (d *webfingerActorStore) buildPersonActor() (activitypub.Item, error) {
	actor := activitypub.ActorNew(
		activitypub.IRI(d.baseUrl+"/ap/users/"+d.actorName),
		activitypub.PersonType,
	)

	actor.PreferredUsername.Set(activitypub.DefaultLang, activitypub.Content(d.actorName))
	if d.actorProfileUrl != "" {
		actor.URL = activitypub.IRI(d.actorProfileUrl)
	}
	actor.Name.Set(activitypub.DefaultLang, activitypub.Content(d.actorName))

	return actor, nil
}
