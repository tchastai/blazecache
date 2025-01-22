package util

import "testing"

func TestRemoveStringFromList(t *testing.T) {
	adress := []string{":8080", ":8081", ":8082"}
	ownAdress := ":8080"
	newAdress := RemoveStringFromList(adress, ownAdress)
	for _, v := range newAdress {
		if v == ownAdress {
			t.Error("string in the list")
		}
	}
}
