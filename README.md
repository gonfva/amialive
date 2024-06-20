# AmIAlive

A very simple status bar app for MacOS to check if we have Internet. 

It pings to one of three DNS servers (it can be changed) and outputs the current round trip time. 

It also keeps a record of the last 10 times (it can be changed) and if the RTT is bigger than double the average (or if we have packet loss) it beeps.
