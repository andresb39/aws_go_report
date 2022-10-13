package main

import (
	cost "github.com/andresb39/aws_go_cost/cost"
	mail "github.com/andresb39/aws_go_cost/email"
	"os"
)

func main() {
	from := os.Getenv("From")
	to := os.Getenv("To")
	regions := os.Getenv("Region")

	cost.GenerateCost(regions)
	defer mail.SendMail(regions, from, to)

}
