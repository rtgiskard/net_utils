package main

import (
	"fmt"
	"time"
)

// ZtNetPost is simplified struct for network settings
type ZtNetPost struct {
	Description string `json:"description,omitempty"`
	Config      struct {
		Name            string `json:"name,omitempty"`
		Private         bool   `json:"private,omitempty"`
		EnableBroadcast bool   `json:"enableBroadcast,omitempty"`
		MulticastLimit  int    `json:"multicastLimit,omitempty"`
		Mtu             int    `json:"mtu,omitempty"`

		Routes []struct {
			Target string      `json:"target,omitempty"`
			Via    interface{} `json:"via,omitempty"`
		} `json:"routes,omitempty"`

		IPAssignmentPools []struct {
			IPRangeStart string `json:"ipRangeStart,omitempty"`
			IPRangeEnd   string `json:"ipRangeEnd,omitempty"`
		} `json:"ipAssignmentPools,omitempty"`

		V4AssignMode struct {
			Zt bool `json:"zt,omitempty"`
		} `json:"v4AssignMode,omitempty"`

		V6AssignMode struct {
			SixPlane bool `json:"6plane,omitempty"`
			Rfc4193  bool `json:"rfc4193,omitempty"`
			Zt       bool `json:"zt,omitempty"`
		} `json:"v6AssignMode,omitempty"`
	} `json:"config,omitempty"`
}

// ZtNetInfo is the generated struct to hold returned network information
type ZtNetInfo struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`

	Config struct {
		Name            string `json:"name,omitempty"`
		Private         bool   `json:"private,omitempty"`
		EnableBroadcast bool   `json:"enableBroadcast,omitempty"`
		MulticastLimit  int    `json:"multicastLimit,omitempty"`
		Mtu             int    `json:"mtu,omitempty"`
		CreationTime    int64  `json:"creationTime,omitempty"`
		LastModified    int64  `json:"lastModified,omitempty"`

		DNS struct {
			Domain  string   `json:"domain,omitempty"`
			Servers []string `json:"servers,omitempty"`
		} `json:"dns,omitempty"`

		IPAssignmentPools []struct {
			IPRangeStart string `json:"ipRangeStart,omitempty"`
			IPRangeEnd   string `json:"ipRangeEnd,omitempty"`
		} `json:"ipAssignmentPools,omitempty"`

		Routes []struct {
			Target string      `json:"target,omitempty"`
			Via    interface{} `json:"via,omitempty"`
		} `json:"routes,omitempty"`

		V4AssignMode struct {
			Zt bool `json:"zt,omitempty"`
		} `json:"v4AssignMode,omitempty"`

		V6AssignMode struct {
			SixPlane bool `json:"6plane,omitempty"`
			Rfc4193  bool `json:"rfc4193,omitempty"`
			Zt       bool `json:"zt,omitempty"`
		} `json:"v6AssignMode,omitempty"`
	} `json:"config,omitempty"`

	OwnerID               string `json:"ownerId,omitempty"`
	OnlineMemberCount     int    `json:"onlineMemberCount,omitempty"`
	AuthorizedMemberCount int    `json:"authorizedMemberCount,omitempty"`
	TotalMemberCount      int    `json:"totalMemberCount,omitempty"`
}

// ZtNetMemberPost is simplified struct for network member settings
type ZtNetMemberPost struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Hidden      bool   `json:"hidden,omitempty"`

	Config struct {
		Authorized      bool     `json:"authorized,omitempty"`
		ActiveBridge    bool     `json:"activeBridge,omitempty"`
		NoAutoAssignIps bool     `json:"noAutoAssignIps,omitempty"`
		IPAssignments   []string `json:"ipAssignments,omitempty"`
	} `json:"config,omitempty"`
}

// ZtNetMemberInfo is the generated struct to hold returned network member information
type ZtNetMemberInfo struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	NetworkID       string `json:"networkId,omitempty"`
	NodeID          string `json:"nodeId,omitempty"`
	Hidden          bool   `json:"hidden,omitempty"`
	PhysicalAddress string `json:"physicalAddress,omitempty"`
	ClientVersion   string `json:"clientVersion,omitempty"`
	ProtocolVersion int    `json:"protocolVersion,omitempty"`
	Clock           int64  `json:"clock,omitempty"`
	LastOnline      int64  `json:"lastOnline,omitempty"`

	Config struct {
		ActiveBridge         bool     `json:"activeBridge,omitempty"`
		Authorized           bool     `json:"authorized,omitempty"`
		NoAutoAssignIps      bool     `json:"noAutoAssignIps,omitempty"`
		IPAssignments        []string `json:"ipAssignments,omitempty"`
		CreationTime         int64    `json:"creationTime,omitempty"`
		LastAuthorizedTime   int64    `json:"lastAuthorizedTime,omitempty"`
		LastDeauthorizedTime int      `json:"lastDeauthorizedTime,omitempty"`
	} `json:"config,omitempty"`
}

// DumpHeader returns header info for data from DumpInfo()
func (i ZtNetInfo) DumpHeader() []interface{} {
	return []interface{}{"NID", "Name", "Route", "Private", "O/T/A", "CreationTime"}
}

// DumpInfo returns a slice of predefined infomation
func (i *ZtNetInfo) DumpInfo() []interface{} {
	createTime := time.Unix(i.Config.CreationTime/1000, 0).Local().Format(time.RFC3339)
	ota := fmt.Sprintf("%d/%d/%d", i.OnlineMemberCount, i.TotalMemberCount, i.AuthorizedMemberCount)

	route := "-"
	if len(i.Config.Routes) > 0 {
		route = i.Config.Routes[0].Target
	}

	return []interface{}{i.ID, i.Config.Name, route, i.Config.Private, ota, createTime}
}

// DumpHeader returns header info for data from DumpInfo()
func (i ZtNetMemberInfo) DumpHeader() []interface{} {
	return []interface{}{"MID", "Name", "IP_assign", "IP_physical", "Version", "LastOnline", "Auth", "Hidden"}
}

// DumpInfo returns a slice of predefined infomation
func (i *ZtNetMemberInfo) DumpInfo() []interface{} {
	lastOnline := time.Unix(i.LastOnline/1000, 0)
	lastduration := time.Since(lastOnline).Truncate(time.Second)

	ipAssigned := "-"
	if len(i.Config.IPAssignments) > 0 {
		ipAssigned = i.Config.IPAssignments[0]
	}

	return []interface{}{i.NodeID, i.Name, ipAssigned, i.PhysicalAddress,
		i.ClientVersion, lastduration, i.Config.Authorized, i.Hidden}
}

func displayNetworks(networks []ZtNetInfo) {
	if len(networks) == 0 {
		fmt.Println("<empty>")
		return
	}

	if args.Verbose {
		for i, v := range networks {
			fmt.Printf("-- net %d: %s\n", i, v.ID)
			fmt.Println(Dumps(v, args.Format))
		}
	} else {
		info := [][]interface{}{ZtNetInfo{}.DumpHeader()}

		for i := range networks {
			info = append(info, networks[i].DumpInfo())
		}

		ShowTable(info)
	}
}

func displayNetworkMembers(members []ZtNetMemberInfo) {
	if len(members) == 0 {
		fmt.Println("<empty>")
		return
	}

	if args.Verbose {
		for i, v := range members {
			fmt.Printf("-- netm %d: %s\n", i, v.NodeID)
			fmt.Println(Dumps(v, args.Format))
		}
	} else {
		info := [][]interface{}{ZtNetMemberInfo{}.DumpHeader()}

		// show brief info
		for i := range members {
			info = append(info, members[i].DumpInfo())
		}

		ShowTable(info)
	}
}

func ztDisplay(o interface{}) {
	switch v := o.(type) {
	case []ZtNetInfo:
		displayNetworks(v)
	case []ZtNetMemberInfo:
		displayNetworkMembers(v)
	default:
		fmt.Println(Dumps(v, args.Format))
	}
}
