package main

func login (login string, pass string) string {
	if login == "admin" && pass == "admin" {
		return "admin"
	} else if login == "user" && pass == "user" {
		return "user"
	} else {
		return "wrong"
	}
}
