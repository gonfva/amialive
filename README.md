# AmIAlive

A very simple status bar app for MacOS to check if we have Internet.

![example](example.png)

It pings to a DNS server (it can be changed) and outputs the current round trip time.

You can then choose to be alerted on of two conditions.

+ The most recent ping exceeds a certain threshold (by default 250ms).
+ The most recent ping exceeds a certain multiple of the moving average of the last 10 pings.

It only works on Mac. And I have only tested on a Mac x86.

You can download the most recent release, but it might complain about some file coming from the Internet.

You can also build the code yourself. Just download the repo and do

`go build .`

Thanks to https://github.com/caseymrm/menuet
