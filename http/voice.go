package http

// VoiceRegions returns the voice server regions
func (s *Session) VoiceRegions() (st []*VoiceRegion, err error) {

	body, err := s.requestWithBucketID("GET", EndpointVoiceRegions, nil, EndpointVoiceRegions)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}
