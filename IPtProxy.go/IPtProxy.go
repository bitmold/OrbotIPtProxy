package IPtProxy

import (
	snowflakeclient "git.torproject.org/pluggable-transports/snowflake.git/v2/client"
	"git.torproject.org/pluggable-transports/snowflake.git/v2/common/safelog"
	sfp "git.torproject.org/pluggable-transports/snowflake.git/v2/proxy/lib"
	"gitlab.com/yawning/obfs4.git/obfs4proxy"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

var meekPort = 47000

// MeekPort - Port where Obfs4proxy will provide its Meek service.
// Only use this after calling StartObfs4Proxy! It might have changed after that!
//goland:noinspection GoUnusedExportedFunction
func MeekPort() int {
	return meekPort
}

var obfs2Port = 47100

// Obfs2Port - Port where Obfs4proxy will provide its Obfs2 service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//goland:noinspection GoUnusedExportedFunction
func Obfs2Port() int {
	return obfs2Port
}

var obfs3Port = 47200

// Obfs3Port - Port where Obfs4proxy will provide its Obfs3 service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//goland:noinspection GoUnusedExportedFunction
func Obfs3Port() int {
	return obfs3Port
}

var obfs4Port = 47300

// Obfs4Port - Port where Obfs4proxy will provide its Obfs4 service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//goland:noinspection GoUnusedExportedFunction
func Obfs4Port() int {
	return obfs4Port
}

var scramblesuitPort = 47400

// ScramblesuitPort - Port where Obfs4proxy will provide its Scramblesuit service.
// Only use this property after calling StartObfs4Proxy! It might have changed after that!
//goland:noinspection GoUnusedExportedFunction
func ScramblesuitPort() int {
	return scramblesuitPort
}

var snowflakePort = 52000

// SnowflakePort - Port where Snowflake will provide its service.
// Only use this property after calling StartSnowflake! It might have changed after that!
//goland:noinspection GoUnusedExportedFunction
func SnowflakePort() int {
	return snowflakePort
}

var obfs4ProxyRunning = false
var snowflakeRunning = false
var snowflakeProxy *sfp.SnowflakeProxy

// StateLocation - Override TOR_PT_STATE_LOCATION, which defaults to "$TMPDIR/pt_state".
var StateLocation string

func init() {
	if runtime.GOOS == "android" {
		StateLocation = "/data/local/tmp"
	} else {
		StateLocation = os.Getenv("TMPDIR")
	}

	StateLocation += "/pt_state"
}

// StartObfs4Proxy - Start the Obfs4Proxy.
//
// This will test, if the default ports are available. If not, it will increment them until there is.
// Only use the port properties after calling this, they might have been changed!
//
// @param logLevel Log level (ERROR/WARN/INFO/DEBUG). Defaults to ERROR if empty string.
//
// @param enableLogging Log to TOR_PT_STATE_LOCATION/obfs4proxy.log.
//
// @param unsafeLogging Disable the address scrubber.
//
// @param proxy HTTP, SOCKS4 or SOCKS5 proxy to be used behind Obfs4proxy. E.g. "socks5://127.0.0.1:12345"
//
// @return Port number where Obfs4Proxy will listen on for Obfs4(!), if no error happens during start up.
//	If you need the other ports, check MeekPort, Obfs2Port, Obfs3Port and ScramblesuitPort properties!
//
//
//goland:noinspection GoUnusedExportedFunction
func StartObfs4Proxy(logLevel string, enableLogging, unsafeLogging bool, proxy string) int {
	if obfs4ProxyRunning {
		return obfs4Port
	}

	obfs4ProxyRunning = true

	for !isAvailable(meekPort) {
		meekPort++
	}

	if meekPort >= obfs2Port {
		obfs2Port = meekPort + 1
	}

	for !isAvailable(obfs2Port) {
		obfs2Port++
	}

	if obfs2Port >= obfs3Port {
		obfs3Port = obfs2Port + 1
	}

	for !isAvailable(obfs3Port) {
		obfs3Port++
	}

	if obfs3Port >= obfs4Port {
		obfs4Port = obfs3Port + 1
	}

	for !isAvailable(obfs4Port) {
		obfs4Port++
	}

	if obfs4Port >= scramblesuitPort {
		scramblesuitPort = obfs4Port + 1
	}

	for !isAvailable(scramblesuitPort) {
		scramblesuitPort++
	}

	fixEnv()

	if len(proxy) > 0 {
		_ = os.Setenv("TOR_PT_PROXY", proxy)
	} else {
		_ = os.Unsetenv("TOR_PT_PROXY")
	}

	go obfs4proxy.Start(&meekPort, &obfs2Port, &obfs3Port, &obfs4Port, &scramblesuitPort, &logLevel, &enableLogging, &unsafeLogging)

	return obfs4Port
}

// StopObfs4Proxy - Stop the Obfs4Proxy.
//goland:noinspection GoUnusedExportedFunction
func StopObfs4Proxy() {
	if !obfs4ProxyRunning {
		return
	}

	go obfs4proxy.Stop()

	obfs4ProxyRunning = false
}

// StartSnowflake - Start the Snowflake client.
//
// @param ice Comma-separated list of ICE servers.
//
// @param url URL of signaling broker.
//
// @param front Front domain.
//
// @param ampCache OPTIONAL. URL of AMP cache to use as a proxy for signaling.
//        Only needed when you want to do the rendezvous over AMP instead of a domain fronted server.
//
// @param logFile Name of log file. OPTIONAL. Defaults to no log.
//
// @param logToStateDir Resolve the log file relative to Tor's PT state dir.
//
// @param keepLocalAddresses Keep local LAN address ICE candidates.
//
// @param unsafeLogging Prevent logs from being scrubbed.
//
// @param maxPeers Capacity for number of multiplexed WebRTC peers. DEFAULTs to 1 if less than that.
//
// @return Port number where Snowflake will listen on, if no error happens during start up.
//
//goland:noinspection GoUnusedExportedFunction
func StartSnowflake(ice, url, front, ampCache, logFile string, logToStateDir, keepLocalAddresses, unsafeLogging bool, maxPeers int) int {
	if snowflakeRunning {
		return snowflakePort
	}

	snowflakeRunning = true

	for !isAvailable(snowflakePort) {
		snowflakePort++
	}

	fixEnv()

	go snowflakeclient.Start(&snowflakePort, &ice, &url, &front, &ampCache, &logFile, &logToStateDir, &keepLocalAddresses, &unsafeLogging, &maxPeers)

	return snowflakePort
}

// StopSnowflake - Stop the Snowflake client.
//goland:noinspection GoUnusedExportedFunction
func StopSnowflake() {
	if !snowflakeRunning {
		return
	}

	go snowflakeclient.Stop()

	snowflakeRunning = false
}

// SnowflakeClientConnected - Interface to use when clients connect
// to the snowflake proxy. For use with StartSnowflakeProxy
type SnowflakeClientConnected interface {
	// Connected - callback method to handle snowflake proxy client connections.
	Connected()
}

// StartSnowflakeProxy - Start the Snowflake proxy.
//
// @param capacity Maximum concurrent clients. OPTIONAL. Defaults to 10, if 0.
//
// @param broker Broker URL. OPTIONAL. Defaults to https://snowflake-broker.torproject.net/, if empty.
//
// @param relay WebSocket relay URL. OPTIONAL. Defaults to wss://snowflake.bamsoftware.com/, if empty.
//
// @param stun STUN URL. OPTIONAL. Defaults to stun:stun.stunprotocol.org:3478, if empty.
//
// @param natProbe OPTIONAL. Defaults to https://snowflake-broker.torproject.net:8443/probe, if empty.
//
// @param logFile Name of log file. OPTIONAL. Defaults to STDERR.
//
// @param keepLocalAddresses Keep local LAN address ICE candidates.
//
// @param unsafeLogging Prevent logs from being scrubbed.
//
// @param clientConnected A delegate which is called when a client successfully connected.
//       Will be called on its own thread! You will need to switch to your own UI thread,
//       if you want to do UI stuff!! OPTIONAL
//
//goland:noinspection GoUnusedExportedFunction
func StartSnowflakeProxy(capacity int, broker, relay, stun, natProbe, logFile string, keepLocalAddresses, unsafeLogging bool, clientConnected SnowflakeClientConnected) {
	if snowflakeProxy != nil {
		return
	}

	if capacity < 1 {
		capacity = 0
	}

	snowflakeProxy = &sfp.SnowflakeProxy{
		Capacity:           uint(capacity),
		STUNURL:            stun,
		BrokerURL:          broker,
		KeepLocalAddresses: keepLocalAddresses,
		RelayURL:           relay,
		NATProbeURL:        natProbe,
		ClientConnectedCallback: func() {
			if clientConnected != nil {
				clientConnected.Connected()
			}
		},
	}

	fixEnv()

	go func(snowflakeProxy *sfp.SnowflakeProxy) {
		var logOutput io.Writer = os.Stderr
		log.SetFlags(log.LstdFlags | log.LUTC)

		log.SetFlags(log.LstdFlags | log.LUTC)
		if logFile != "" {
			f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if err != nil {
				log.Fatal(err)
			}
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
			logOutput = io.MultiWriter(os.Stderr, f)
		}
		if unsafeLogging {
			log.SetOutput(logOutput)
		} else {
			log.SetOutput(&safelog.LogScrubber{Output: logOutput})
		}

		err := snowflakeProxy.Start()
		if err != nil {
			log.Fatal(err)
		}
	}(snowflakeProxy)
}

// StopSnowflakeProxy - Stop the Snowflake proxy.
//goland:noinspection GoUnusedExportedFunction
func StopSnowflakeProxy() {
	if snowflakeProxy == nil {
		return
	}

	go func() {
		snowflakeProxy.Stop()
		snowflakeProxy = nil
	}()
}

// Hack: Set some environment variables that are either
// required, or values that we want. Have to do this here, since we can only
// launch this in a thread and the manipulation of environment variables
// from within an iOS app won't end up in goptlib properly.
//
// Note: This might be called multiple times when using different functions here,
// but that doesn't necessarily mean, that the values set are independent each
// time this is called. It's still the ENVIRONMENT, we're changing here, so there might
// be race conditions.
func fixEnv() {
	_ = os.Setenv("TOR_PT_CLIENT_TRANSPORTS", "meek_lite,obfs2,obfs3,obfs4,scramblesuit,snowflake")
	_ = os.Setenv("TOR_PT_MANAGED_TRANSPORT_VER", "1")

	_ = os.Setenv("TOR_PT_STATE_LOCATION", StateLocation)
}

func isAvailable(port int) bool {
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)

	if err != nil {
		return true
	}

	_ = conn.Close()

	return false
}


var lwipStack core.LWIPStack

// PacketFlow should be implemented in Java/Kotlin.
type PacketFlow interface {
	// WritePacket should writes packets to the TUN fd.
	WritePacket(packet []byte)
}

// Write IP packets to the lwIP stack. Call this function in the main loop of
// the VpnService in Java/Kotlin, which should reads packets from the TUN fd.
func InputPacket(data []byte) {
	if lwipStack != nil {
		lwipStack.Write(data)
	}
}

func StartSocks(packetFlow PacketFlow, proxyHost string, proxyPort int) {
	if packetFlow != nil {
		lwipStack = core.NewLWIPStack()
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, uint16(proxyPort)))
		//core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, uint16(proxyPort), 2*time.Minute))
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})
	}
}

