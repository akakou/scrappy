package main

func main() {
	err := setup()

	if err != nil {
		panic(err)
	}

	runServer()
}
