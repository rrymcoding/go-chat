package main

func main() {

	s := newServer()
	go s.run()
	s.start()

}
