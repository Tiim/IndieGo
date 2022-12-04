package indieauth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"tiim/go-comment-api/mfobjects"

	"willnorris.com/go/microformats"
)

type appInfo struct {
	ClientId     string
	Name         string
	Logo         string
	Summary      string
	Author       string
	RedirectUris []string
}

func getAppInfo(clientId string, client http.Client, ctx context.Context) (*appInfo, error) {
	cidUrl, err := url.ParseRequestURI(clientId)

	if err != nil {
		return nil, err
	}

	if cidUrl.Scheme != "https" && cidUrl.Scheme != "http" {
		return nil, fmt.Errorf("clientId must have http or https scheme")
	}

	ips, err := net.LookupHost(cidUrl.Host)

	for _, ip := range ips {
		ipp := net.ParseIP(ip)
		if ipp.IsLoopback() {
			return &appInfo{ClientId: clientId, Name: "localhost"}, nil
		}
	}

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", clientId, nil)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	mfData := microformats.Parse(res.Body, cidUrl)
	if len(mfData.Items) == 0 {
		return &appInfo{ClientId: clientId}, nil
	}

	happ := mfobjects.GetHApp(mfData)

	return &appInfo{
		ClientId:     clientId,
		Name:         happ.Name,
		Logo:         happ.Logo,
		Summary:      happ.Summary,
		Author:       happ.Author.Name,
		RedirectUris: happ.RedirectUris,
	}, nil
}
