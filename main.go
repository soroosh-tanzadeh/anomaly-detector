package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/pcap"
	"github.com/redis/go-redis/v9"
	"github.com/soroosh-tanzadeh/anormaly_detector/internal/streams"
)

func main() {
	// Check if file argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a pcap file to read")
		os.Exit(1)
	}

	// Open up the pcap file for reading
	handle, err := pcap.OpenLive(os.Args[1], 16000, true, time.Microsecond)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	c := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       5,
	})
	if err := c.Ping(context.Background()).Err(); err != nil {
		fmt.Errorf("Redis Error", err)
	}

	stream := streams.NewRedisTrafficStream("status:traffic", c)

	sumInSecond := 0
	prevTime := time.Now().UnixMilli()
	for packet := range packetSource.Packets() {
		nowTime := time.Now().UnixMilli()
		if nowTime-prevTime >= 1000 {
			err := stream.Add(int64(sumInSecond))
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Printf("%v MB/Sec\n", ((float64(sumInSecond)) / 1e+6))
			sumInSecond = 0
			prevTime = nowTime
		}
		sumInSecond += len(packet.Data())
	}
}
