package main

func main() {
	// dgClient := newClient()
	mux := Route()
	server := NewServer(mux)
	server.Run()
}
