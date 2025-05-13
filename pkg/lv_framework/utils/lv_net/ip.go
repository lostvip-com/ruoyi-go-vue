package lv_net

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

// GET preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return localAddr.IP.String()
}

func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

// IsPrivateIP 判断是否为内网IP
func IsPrivateIP(ipStr string) bool {
	// 将字符串解析为 net.IP
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false // 如果解析失败，返回 false
	}

	// 定义内网IP地址段
	privateIPBlocks := []*net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
		{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("169.254.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	// 遍历所有内网地址段，检查IP是否在其中
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func GetRemoteClientIp(r *http.Request) string {
	remoteIp := r.RemoteAddr

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteIp = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(remoteIp)
	}

	//本地ip
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}

	return remoteIp
}

// 获取外网ip地址
func GetLocation(ip string) string {
	return ip
	//if ip == "127.0.0.1" || ip == "localhost" {
	//	return "内部IP"
	//}
	//resp, err := http.Get("https://restapi.amap.com/v3/ip?ip=" + ip + "&key=3fabc36c20379fbb9300c79b19d5d05e")
	//if err != nil {
	//	panic(err)
	//
	//}
	//defer resp.Body.Close()
	//s, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf(string(s))
	//
	//m := make(map[string]string)
	//
	//err = json.Unmarshal(s, &m)
	//if err != nil {
	//	fmt.Println("Umarshal failed:", err)
	//}
	//if m["province"] == "" {
	//	return "未知位置"
	//}
	//return m["province"] + "-" + m["city"]
}

// 获取局域网ip地址
func GetLocaHost() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}

type Tunit struct {
	Pro  string `json:"pro"`
	City string `json:"city"`
}

func GetRealAddressByIP(ip string) string {
	url := "http://whois.pconline.com.cn/ipJson.jsp?ip=" + ip + "&json=true"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	var result string
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		dws := new(Tunit)
		json.Unmarshal(body, &dws)
		result = dws.Pro + " " + dws.City
	} else {
		result = ip
	}
	return result
}
