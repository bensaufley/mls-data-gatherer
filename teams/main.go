package teams

import "errors"

// AbbrevIsValid returns an error if the passed abbreviation isn't in the list
func AbbrevIsValid(team string) error {
	for name := range Teams {
		if name == team {
			return nil
		}
	}
	return errors.New("No such team abbreviation " + team)
}

// NameIsValid returns an error if the passed name isn't in the list
func NameIsValid(team string) error {
	for _, name := range Teams {
		if name == team {
			return nil
		}
	}
	return errors.New("No such team name " + team)
}

// AbbrevFor returns the abbreviation for a team by its name
func AbbrevFor(team string) (string, error) {
	for abbrv, name := range Teams {
		if name == team {
			return abbrv, nil
		}
	}
	return "", errors.New("No team found")
}
