package tmp

// Guild returns a Guild structure of a specific Guild.
// guildID   : The ID of a Guild
func (s *Session) Guild(guildID string) (st *Guild, err error) {
	if s.StateEnabled {
		// Attempt to grab the guild from State first.
		st, err = s.State.Guild(guildID)
		if err == nil && !st.Unavailable {
			return
		}
	}
	/*
		body, err := s.RequestWithBucketID("GET", EndpointGuild(guildID), nil, EndpointGuild(guildID))
		if err != nil {
			return
		}

		err = unmarshal(body, &st)
	*/
	return
}
