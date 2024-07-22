package git

import "testing"

func TestRemoteToLink(t *testing.T) {

	expected := "https://github.com/TestUser/TestRepo/commit/"

	remotes := []string{
		"git@github.com:TestUser/TestRepo",
		"git@github.com:TestUser/TestRepo.git",
		"https://github.com/TestUser/TestRepo",
		"https://github.com/TestUser/TestRepo.git",
		"http://github.com/TestUser/TestRepo",
		"http://github.com/TestUser/TestRepo.git",
	}

	for _, remote := range remotes {
		link := remoteToLink(remote)
		if expected != link {
			t.Fatalf("%s --> %s, expected %s", remote, link, expected)
		}
	}
}
