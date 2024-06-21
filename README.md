# AmIAlive

A very simple status bar app for MacOS to check if we have Internet. 

![example](example.png)

It pings to a DNS server (it can be changed) and outputs the current round trip time. 

It also keeps an internal record of the last 10 times (it can be changed) and if the RTT is bigger than double the average (or if we have packet loss) it beeps.

It only works on Mac. And I have only tested on a Mac x86.

To build, download the code and do

`go build .`

Thanks to https://github.com/caseymrm/menuet
