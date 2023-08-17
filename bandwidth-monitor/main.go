package main

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/jroimartin/gocui"
)

func displayGui(done chan struct{}, uploadBytes *float64, downloadBytes *float64) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Printf("Error creating GUI: %v\n", err)
		return
	}
	defer g.Close()
	var lastTime time.Time
	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("stats", 0, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Network Statistics"
			v.Autoscroll = true
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						g.Update(func(g *gocui.Gui) error {
							elapsedTime := time.Since(lastTime).Seconds()
							v.Clear()
							fmt.Fprintf(v, "Upload Speed: %.2f bytes/s\n", *uploadBytes/elapsedTime)
							fmt.Fprintf(v, "Download Speed: %.2f bytes/s\n", *downloadBytes/elapsedTime)
							lastTime = time.Now()
							return nil
						})
						time.Sleep(1 * time.Second)
					}
				}
			}()
		}
		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Printf("Error in MainLoop: %v\n", err)
	}
}

func cpaturePacket(iface string, macAddress net.HardwareAddr, done chan struct{}, uploadBytes *float64, downloadBytes *float64) {
	defer close(done)
	handle, err := pcap.OpenLive(iface, 65536, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Error opening interface %s: %s\n", iface, err.Error())
		return
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process captured packet here

		packetSize := float64(len(packet.Data()))

		srcMAC := packet.LinkLayer().LinkFlow().Src()
		isUpload := srcMAC.String() == macAddress.String()

		if isUpload {
			*uploadBytes += packetSize
		} else {
			*downloadBytes += 1
		}
	}

}

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error retrieving network interfaces: %s\n", err.Error())
		return
	}
	done := make(chan struct{})
	uploadBytes := 0.0
	downloadBytes := 0.0

	go displayGui(done, &uploadBytes, &downloadBytes)
	for _, iface := range interfaces {
		addreses, err := iface.Addrs()
		if err != nil {
			return
		}
		if len(iface.Name) > 0 {
			for _, addr := range addreses {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					mac := iface.HardwareAddr
					if len(mac) > 0 {
						fmt.Printf("Capture the packet for %s\n", iface.Name)
						go cpaturePacket(iface.Name, mac, done, &uploadBytes, &downloadBytes)
						break
					}
				}
			}
		}

	}

	time.Sleep(30 * time.Second)
	close(done)

	for range interfaces {
		<-done
	}
}
