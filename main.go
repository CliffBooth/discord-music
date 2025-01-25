package main

import (
	"fmt"
	"time"
)

func main() {
	var snowflake uint = 1332348323797012590
	fmt.Printf("snoflake hex = %x\n", snowflake)

	timestmp := snowflake >> 22 // timestamp == bits from 64 to 22
	fmt.Printf("timestmp hex = %x\n", timestmp)
	var discordEpoch uint = 1420070400000

	ageTime := discordEpoch + timestmp

	t := time.UnixMilli(int64(ageTime))
	fmt.Println(t)
}
