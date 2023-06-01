package main

import (
	scrappy_gin "github.com/akakou/scrappy/gin"
)

func main() {
	scrappy_gin.RunIssuer("../templates/*.html")
}
