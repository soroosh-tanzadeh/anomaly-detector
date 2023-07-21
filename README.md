# SMA and Autoregressive model Anomaly Detector
**For Educational Purposes Only**
## Usage
create .env file from .env.example
```
mv .env.example .env
```

Set Redis Address and Port for your enviroment

you can also change thershold with updating ``SMA_DETECTOR_THERHSOLD`` 

**Enter Thereshold in Megabits**

#### Install dependencies
```
go mod vendor
```

This application require libpcap for network packets inspection.
```
sudo apt-get install libpcap-dev
```


#### Run Application
```
go run main.go {interface}
```

replace {interface} with your network interface name

#### Example
```
go run main.go eth0
```