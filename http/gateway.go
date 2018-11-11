package http

import "strings"

// ------------------------------------------------------------------------------------------------
// Functions specific to Discord Websockets
// ------------------------------------------------------------------------------------------------

// Gateway returns the websocket Gateway address
func (s *Session) Gateway() (gateway string, err error) {

	response, err := s.requestWithBucketID("GET", EndpointGateway, nil, EndpointGateway)
	if err != nil {
		return
	}

	temp := struct {
		URL string `json:"url"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	gateway = temp.URL

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(gateway, "/") {
		gateway += "/"
	}

	return
}

// GatewayBot returns the websocket Gateway address and the recommended number of shards
func (s *Session) GatewayBot() (st *GatewayBotResponse, err error) {

	response, err := s.requestWithBucketID("GET", EndpointGatewayBot, nil, EndpointGatewayBot)
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	if err != nil {
		return
	}

	// Ensure the gateway always has a trailing slash.
	// MacOS will fail to connect if we add query params without a trailing slash on the base domain.
	if !strings.HasSuffix(st.URL, "/") {
		st.URL += "/"
	}

	return
}
