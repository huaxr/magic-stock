package env

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

const (
	UnknownIDC   = "-"
	DC_HY        = "hy"
	DC_LF        = "lf"
	DC_HL        = "hl"
	DC_WJ        = "wj"
	DC_VA        = "va"
	DC_SG        = "sg"
	DC_BR        = "br"
	DC_JP        = "jp"
	DC_IN        = "in"
	DC_CA        = "ca" // West America
	DC_FRAWS     = "fraws"
	DC_ALISG     = "alisg" // Singapore Aliyun
	DC_ALIVA     = "aliva"
	DC_MALIVA    = "maliva"
	DC_ALINC2    = "alinc2"    // aliyun north
	DC_MAWSJP    = "mawsjp"    // musical.ly 东京老机房
	DC_SDQD      = "sdqd"      // sdqd
	DC_BOE       = "boe"       // bytedance offline environment
	DC_IBOE      = "boei18n"   // bytedance offline environment(international)
	DC_LQ        = "lq"        // 灵丘机房
	DC_SYKALARK  = "sykalark"  // 顺义KA机房(lark)
	DC_SYKA2LARK = "syka2lark" // 顺义KA2机房(lark)
	DC_USEAST2A  = "useast2a"  // 美东第二机房
)

const idcFileDefault = "/opt/tiger/chadc/netsegment"
const idcRefreshDur = regionRefreshDur

var (
	idc          atomic.Value // local idc (string)
	idcFile      atomic.Value // file name (string)
	ipv4Segments atomic.Value // ipv4 segments - idc ([]*ipv4Segment)
	ipv6Segments atomic.Value // ipv6 segments - idc ([]*ipv6Segment)
)

// IDC .
func IDC() string {
	if v := idc.Load(); v != nil {
		return v.(string)
	}

	if dc := os.Getenv("RUNTIME_IDC_NAME"); dc != "" {
		idc.Store(dc)
		return dc
	}

	b, err := ioutil.ReadFile("/opt/tmp/consul_agent/datacenter")
	if err == nil {
		if dc := strings.TrimSpace(string(b)); dc != "" {
			idc.Store(dc)
			return dc
		}
	}

	cmd0 := exec.Command("/opt/tiger/consul_deploy/bin/determine_dc.sh")
	output0, err := cmd0.Output()
	if err == nil {
		dc := strings.TrimSpace(string(output0))
		if hasIDC(dc) {
			idc.Store(dc)
			return dc
		}
	}

	cmd := exec.Command(`bash`, `-c`, `sd report|grep "Data center"|awk '{print $3}'`)
	output, err := cmd.Output()
	if err == nil {
		dc := strings.TrimSpace(string(output))
		if hasIDC(dc) {
			idc.Store(dc)
			return dc
		}
	}

	idc.Store(UnknownIDC)
	return UnknownIDC
}

// GetIDCFromHost .
func GetIDCFromHost(ip string) string {
	netIP := net.ParseIP(ip)
	// ipv4
	if strings.Contains(ip, ".") {
		return lookupIpv4(netIP)
	}
	// ipv6
	if strings.Contains(ip, ":") {
		return lookupIpv6(netIP)
	}
	return UnknownIDC
}

// GetIDCList .
func GetIDCList() []string {
	return idcList()
}

// SetIDC .
func SetIDC(v string) {
	idc.Store(v)
}

// SetIDCFile .
func SetIDCFile(file string) {
	idcFile.Store(file)
	updateSegments()
}

func init() {
	idcFile.Store(idcFileDefault)
	refreshIDCs()
}

// ATTENTION: IT COMES WITH A LOOP, DON'T CALL IT AGAIN.
func refreshIDCs() {
	defer time.AfterFunc(idcRefreshDur, refreshIDCs)
	updateSegments()
}

func updateSegments() {
	file, _ := idcFile.Load().(string)
	lines := readFile(file)
	if len(lines) == 0 {
		return
	}
	var ipv4s []*ipv4Segment
	var ipv6s []*ipv6Segment
	for i, _ := range lines {
		s := strings.Split(lines[i], " ")
		if len(s) != 2 {
			continue
		}
		_, ipNet, err := net.ParseCIDR(s[0])
		if err != nil {
			continue
		}
		// ipv4
		if strings.Contains(s[0], ".") {
			ipv4s = append(ipv4s, newIpv4Segment(ipNet, s[1]))
		}
		// ipv6
		if strings.Contains(s[0], ":") {
			ipv6s = append(ipv6s, newIpv6Segment(ipNet, s[1]))
		}
	}
	// sort & merge
	ipv4Segments.Store(mergeIpv4Segment(ipv4s))
	ipv6Segments.Store(mergeIpv6Segment(ipv6s))
}

// readFile return lines
func readFile(file string) (lines []string) {
	f, err := os.Open(file)
	if err != nil {
		return lines
	}
	buf := bufio.NewReader(f)
	line := ""
	for {
		bytes, isPre, err := buf.ReadLine()
		if err != nil {
			return lines
		}
		line = fmt.Sprintf("%s%s", line, string(bytes))
		if !isPre {
			lines = append(lines, line)
			line = ""
		}
	}
}

// ipv4 section
//
// lookupIpv4 .
func lookupIpv4(ip net.IP) (idc string) {
	if len(ip) == net.IPv6len {
		ip = ip[12:16]
	}
	if len(ip) != net.IPv4len {
		return UnknownIDC
	}
	ipv4s, ok := ipv4Segments.Load().([]*ipv4Segment)
	if !ok || len(ipv4s) == 0 {
		return UnknownIDC
	}
	// Binary search
	target := bigEndianUint32(ip)
	search := func(i int) bool {
		return target <= ipv4s[i].end
	}
	index := sort.Search(len(ipv4s), search)
	// search not found
	if index < 0 || index >= len(ipv4s) {
		return UnknownIDC
	}
	if ipv4s[index].start <= target && target <= ipv4s[index].end {
		return ipv4s[index].idc
	} else {
		return UnknownIDC
	}
}

// ipv4Segment .
type ipv4Segment struct {
	start uint32
	end   uint32
	idc   string
}

// newIpv4Segment .
func newIpv4Segment(ipNet *net.IPNet, idc string) *ipv4Segment {
	segment := &ipv4Segment{
		idc:   idc,
		start: bigEndianUint32(ipNet.IP),
	}
	segment.end = segment.start + ^bigEndianUint32(ipNet.Mask)
	return segment
}

// mergeIpv4Segment Merge adjacent segments if idc is same
func mergeIpv4Segment(ipv4s []*ipv4Segment) (merged []*ipv4Segment) {
	if len(ipv4s) == 0 {
		return merged
	}
	sort.Sort(sortIpv4Segment(ipv4s))
	merged = append(merged, ipv4s[0])
	last := ipv4s[0]
	for i := 1; i < len(ipv4s); i++ {
		if ipv4s[i].idc == last.idc && ipv4s[i].start <= last.end+1 {
			if ipv4s[i].end > last.end {
				last.end = ipv4s[i].end
			}
		} else {
			merged = append(merged, ipv4s[i])
			last = ipv4s[i]
		}
	}
	return merged
}

// sortIpv4Segment .
type sortIpv4Segment []*ipv4Segment

// Len .
func (s sortIpv4Segment) Len() int {
	return len(s)
}

// Swap .
func (s sortIpv4Segment) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less .
func (s sortIpv4Segment) Less(i, j int) bool {
	return s[i].start < s[j].start
}

// ipv6 section
//
// lookupIpv6 .
func lookupIpv6(ip net.IP) (idc string) {
	if len(ip) != net.IPv6len {
		return UnknownIDC
	}
	ipv6s, ok := ipv6Segments.Load().([]*ipv6Segment)
	if !ok || len(ipv6s) == 0 {
		return UnknownIDC
	}
	// Binary search
	target := bigEndianUint128(ip)
	search := func(i int) bool {
		// target <= ipv6s[i].end
		return leU128(target, ipv6s[i].end)
	}
	index := sort.Search(len(ipv6s), search)
	// search not found
	if index < 0 || index >= len(ipv6s) {
		return UnknownIDC
	}
	// ipv6s[i].start <= target && target <= ipv6s[i].end
	if leU128(ipv6s[index].start, target) && leU128(target, ipv6s[index].end) {
		return ipv6s[index].idc
	} else {
		return UnknownIDC
	}
}

// ipv6 has 128 bits and uses [4]uint32.
type ipv6Segment struct {
	start uint128
	end   uint128
	idc   string
}

// newIpv6Segment .
func newIpv6Segment(ipNet *net.IPNet, idc string) *ipv6Segment {
	segment := &ipv6Segment{
		idc:   idc,
		start: bigEndianUint128(ipNet.IP),
	}
	segment.end = sumU128(segment.start, negU128(bigEndianUint128(ipNet.Mask)))
	return segment
}

// mergeIpv6Segment Merge adjacent segments if idc is same
func mergeIpv6Segment(ipv6s []*ipv6Segment) (merged []*ipv6Segment) {
	if len(ipv6s) == 0 {
		return merged
	}
	sort.Sort(sortIpv6Segment(ipv6s))
	merged = append(merged, ipv6s[0])
	last := ipv6s[0]
	for i := 1; i < len(ipv6s); i++ {
		// idc == idc && ipv6s[i].start <= last.end + 1
		if ipv6s[i].idc == last.idc && leU128(ipv6s[i].start, sumU128(last.end, uint128{0, 0, 0, 1})) {
			// last.end < ipv6s[i].end
			if ltU128(last.end, ipv6s[i].end) {
				last.end = ipv6s[i].end
			}
		} else {
			merged = append(merged, ipv6s[i])
			last = ipv6s[i]
		}
	}
	return merged
}

// sortIpv6Segment .
type sortIpv6Segment []*ipv6Segment

// Len .
func (s sortIpv6Segment) Len() int {
	return len(s)
}

// Swap .
func (s sortIpv6Segment) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less .
func (s sortIpv6Segment) Less(i, j int) bool {
	return ltU128(s[i].start, s[j].start)
}

// uint128 Express with four uint32s
type uint128 [4]uint32

// geU128 return true when i >= j.
func geU128(i, j uint128) bool {
	return !ltU128(i, j)
}

// leU128 return true when i <= j.
func leU128(i, j uint128) bool {
	return geU128(j, i)
}

// ltU128 return true when i < j.
func ltU128(i, j uint128) bool {
	return i[0] < j[0] ||
		(i[0] == j[0] && i[1] < j[1]) ||
		(i[1] == j[1] && i[2] < j[2]) ||
		(i[2] == j[2] && i[3] < j[3])
}

// sumU128 .
func sumU128(i, j uint128) uint128 {
	var res uint128
	for k := 3; k > 0; k-- {
		tmp := uint64(i[k]) + uint64(j[k]) + uint64(res[k])
		res[k] = uint32(tmp)
		res[k-1] = uint32(tmp >> 32)
	}
	res[0] += i[0] + j[0]
	return res
}

// negU128 Negate number.
func negU128(i uint128) uint128 {
	return uint128{^i[0], ^i[1], ^i[2], ^i[3]}
}

// bigEndianUint32 .
func bigEndianUint32(p []byte) uint32 {
	_ = p[3]
	return uint32(p[0])<<24 | uint32(p[1])<<16 | uint32(p[2])<<8 | uint32(p[3])
}

// bigEndianUint128 .
func bigEndianUint128(p []byte) uint128 {
	_ = p[15]
	return uint128{bigEndianUint32(p[0:4]), bigEndianUint32(p[4:8]), bigEndianUint32(p[8:12]), bigEndianUint32(p[12:16])}
}
