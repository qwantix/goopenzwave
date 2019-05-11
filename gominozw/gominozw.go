package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/qwantix/goopenzwave"
)

type NodeInfo struct {
	HomeID uint32
	NodeID uint8
	Node   *goopenzwave.Node
	Values map[uint64]*goopenzwave.ValueID
}

var Nodes = make(map[uint8]*NodeInfo)
var initialQueryComplete = make(chan bool)
var sentinitialQueryComplete = false

func main() {
	var controllerPath string
	var configPath string
	flag.StringVar(&controllerPath, "controller", "/dev/ttyUSB0", "the path to your controller device")
	flag.StringVar(&configPath, "config", "/usr/local/etc/openzwave/", "the path to open-zwave config")
	flag.Parse()

	fmt.Println("gominozw started with openzwave version:", goopenzwave.GetVersionLongAsString())

	// Setup the OpenZWave library.
	options := goopenzwave.CreateOptions(configPath, "", "")
	options.AddOptionLogLevel("SaveLogLevel", goopenzwave.LogLevelNone)
	options.AddOptionLogLevel("QueueLogLevel", goopenzwave.LogLevelNone)
	options.AddOptionInt("DumpTrigger", 4)
	options.AddOptionInt("PollInterval", 500)
	options.AddOptionBool("IntervalBetweenPolls", true)
	options.AddOptionBool("ValidateValueChanges", true)
	options.Lock()

	// Start the library and listen for notifications.
	err := goopenzwave.Start(handleNotification)
	if err != nil {
		fmt.Println("ERROR: failed to start goopenzwave package:", err)
		return
	}

	// Add a driver using the supplied controller path.
	err = goopenzwave.AddDriver(controllerPath)
	if err != nil {
		fmt.Println("ERROR: failed to add goopenzwave driver:", err)
		return
	}

	// Wait here until the initial node query has completed.
	<-initialQueryComplete
	fmt.Println("Finished initial scan, now setting up polling...")
	return
	// Now we will enable polling for a variable. In this simple example, it
	// has been hardwired to poll COMMAND_CLASS_BASIC on each node that
	// supports this setting.
	for _, node := range Nodes {
		// Skip the controller (most likely node 1).
		if node.NodeID == 1 {
			continue
		}

		// For each value for this node, set up polling.
		for i := range node.Values {
			valueid := node.Values[i]

			if valueid.CommandClassID == 0x20 {
				// Enable polling with "intensity" of 2. Though, this is
				// irrelevant with only one value polled.
				valueid.EnablePoll(2)
			}
		}
	}

	fmt.Println("Initial scan complete. Now polling for updates...")
	fmt.Println("Hit ctrl-c to quit")

	// Now wait for the user to quit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	// Print out what we know about the network.
	fmt.Println("Nodes:")
	for id, node := range Nodes {
		fmt.Printf("\t%3d: Node: %s\n", id, node.Node)
		fmt.Printf("\t     Values: (%d)\n", len(node.Values))
		for i := range node.Values {
			fmt.Printf("\t\t0x%x: %s\n", i, node.Values[i])
		}
	}

	// All done now finish up.
	goopenzwave.RemoveDriver(controllerPath)
	goopenzwave.Stop()
	goopenzwave.DestroyOptions()
}

func handleNotification(notification *goopenzwave.Notification) {
	if notification.HomeID == 0 {
		return // Ignore
	}
	fmt.Println("Received notification:", notification)

	// Switch based on notification type.
	switch notification.Type {
	case goopenzwave.NotificationTypeNodeAdded:
		// Create a NodeInfo from the notification then add it to the
		// map.
		nodeinfo := &NodeInfo{
			HomeID: notification.HomeID,
			NodeID: notification.NodeID,
			Node:   goopenzwave.NewNode(notification.HomeID, notification.NodeID),
			Values: make(map[uint64]*goopenzwave.ValueID),
		}
		Nodes[nodeinfo.NodeID] = nodeinfo

	case goopenzwave.NotificationTypeNodeRemoved:
		// Remove the NodeInfo from the map.
		if _, found := Nodes[notification.NodeID]; found {
			delete(Nodes, notification.NodeID)
		}

	case goopenzwave.NotificationTypeValueAdded, goopenzwave.NotificationTypeValueChanged:
		// Find the NodeInfo in the map and add/change the ValueID.
		if node, found := Nodes[notification.NodeID]; found {
			node.Values[notification.ValueID.ID] = notification.ValueID
		}

	case goopenzwave.NotificationTypeValueRemoved:
		// Find the NodeInfo in the map and remove the ValueID.
		if node, found := Nodes[notification.NodeID]; found {
			if _, foundVal := node.Values[notification.ValueID.ID]; foundVal {
				delete(node.Values, notification.ValueID.ID)
			}
		}

	case goopenzwave.NotificationTypeAwakeNodesQueried, goopenzwave.NotificationTypeAllNodesQueried, goopenzwave.NotificationTypeAllNodesQueriedSomeDead:
		// The initial node query has completed.
		if sentinitialQueryComplete == false {
			initialQueryComplete <- true
		}
	}
}
