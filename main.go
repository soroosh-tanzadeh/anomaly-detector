package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/pcap"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"github.com/soroosh-tanzadeh/anormaly_detector/internal/detectors"
	"github.com/soroosh-tanzadeh/anormaly_detector/internal/streams"
	"github.com/soroosh-tanzadeh/anormaly_detector/internal/trafficplot"
)

var plotLock *sync.Mutex = &sync.Mutex{}

const unit = 125000

func main() {
	threshold, err := strconv.ParseFloat(os.Getenv("SMA_DETECTOR_THERHSOLD"), 64)
	if err != nil {
		log.Fatal(err)
	}
	threshold = threshold * unit

	windowSize, err := strconv.ParseFloat(os.Getenv("SMA_DETECTOR_WINDOW_SIZE"), 64)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("SMA_DETECTOR_REDIS_ADDR"),
		Password: os.Getenv("SMA_DETECTOR_REDIS_PASS"),
		DB:       5,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis Error: %s", err.Error())
	}

	trafficplot.CreateTrafficPlot()

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

	stream := streams.NewRedisTrafficStream("status:traffic", redisClient)

	go packetLogger(stream, *packetSource)

	// Change Thershold]
	go windowTracker(stream, windowSize, threshold)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
	trafficplot.Save()
	fmt.Println("exiting")
}

// stream Traffic Stream
// thershold max moving average in bits
func windowTracker(stream streams.TrafficStream, windowSize, thershold float64) {
	runtime.LockOSThread()
	prevTime := time.Now()
	for {
		nowTime := time.Now()
		if nowTime.UnixMilli()-prevTime.UnixMilli() >= int64(windowSize*1000) {
			traffics, err := stream.Range(context.Background(), prevTime, nowTime)
			if err != nil {
				fmt.Printf("Redis Error: %s", err.Error())
			}

			go func() {
				if len(traffics) >= 10 {
					detections := detectors.DetectAnomalyWithSMA(traffics, int(math.Floor(windowSize/4)), thershold)
					for detectTime, detection := range detections {
						if detection.Detected {
							dt := nowTime.Add(-time.Second * time.Duration(windowSize)).Add(time.Second * time.Duration(detectTime+1))
							plotLock.Lock()
							trafficplot.CaptureAnomaly(dt, detection.Traffic/unit)
							plotLock.Unlock()
							fmt.Printf("Anormaly detected in %s \n", dt.Format("15:04:05"))
						}
					}
				}
			}()
			prevTime = nowTime
		}
	}
}

func packetLogger(stream streams.TrafficStream, packetSource gopacket.PacketSource) {
	runtime.LockOSThread()
	sumInSecond := 0
	prevTime := time.Now()
	for packet := range packetSource.Packets() {
		nowTime := time.Now()
		if nowTime.UnixMilli()-prevTime.UnixMilli() >= 1000 {
			err := stream.Add(context.Background(), float64(sumInSecond))
			if err != nil {
				fmt.Println(err.Error())
			}

			plotLock.Lock()
			trafficplot.Capture(((float64(sumInSecond)) / unit))
			plotLock.Unlock()

			fmt.Printf("%.8f MB/Sec - %s\n", ((float64(sumInSecond)) / unit), nowTime.Format("15:04:05"))
			sumInSecond = 0
			prevTime = nowTime
		}
		sumInSecond += len(packet.Data())
	}
	runtime.UnlockOSThread()
}
