package models

import "testing"

func TestUsernameCensor(t *testing.T) {
	// Create a simple instance of APIUser for testing
	regularUser := APIUser{Role: Regular}
	moderator := APIUser{Role: Moderator}

	usernames := []string{"u", "us", "use", "user", "usern", "userna", "usernam", "username"}
	censoredUsernames := []string{"u", "u*", "u**", "u***", "us***", "use***", "use****", "use*****"}

	for i, username := range usernames {
		banEntry := Log{ToUser: &User{Username: username}}
		censoredUsername := banEntry.BannedDisplayName(regularUser)
		if censoredUsername != censoredUsernames[i] {
			t.Errorf("Not censored correctly. Expected %s, got %s", censoredUsernames[i], censoredUsername)
		}
		cleanUsername := banEntry.BannedDisplayName(moderator)
		if cleanUsername != username {
			t.Errorf("Moderators should see full usernames. Expected %s, got %s", censoredUsernames[i], censoredUsername)
		}
	}

}
