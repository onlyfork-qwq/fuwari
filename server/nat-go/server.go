package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Msg struct {
	SDP           string `json:"sdp,omitempty"`
	ICECandidate  string `json:"ice-candidate,omitempty"`
	ProbeSDP      string `json:"probe-sdp,omitempty"`
	ProbeAnswer   string `json:"probe-answer,omitempty"`
	ProbeCandidate string `json:"probe-candidate,omitempty"`
}

type Session struct {
	ws       *websocket.Conn
	pcA      *webrtc.PeerConnection
	pcB      *webrtc.PeerConnection
	publicIP string
	ports    map[string]bool
	mu       sync.Mutex
	pcBReady bool
	pcBRecv  bool
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("NAT检测服务器启动在 :9000")
	http.ListenAndServe(":9000", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket升级失败:", err)
		return
	}
	defer ws.Close()

	session := &Session{
		ws:    ws,
		ports: make(map[string]bool),
	}

	// 创建pcA（主连接）
	pcA, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		log.Println("创建pcA失败:", err)
		return
	}
	session.pcA = pcA
	defer pcA.Close()

	// 创建pcB（探测连接，使用不同STUN）
	pcB, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun1.l.google.com:19302"}},
		},
	})
	if err != nil {
		log.Println("创建pcB失败:", err)
		return
	}
	session.pcB = pcB
	defer pcB.Close()

	// pcA接收DataChannel
	pcA.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.Println("pcA收到DataChannel:", dc.Label())

		dc.OnOpen(func() {
			log.Println("pcA DataChannel已打开")
			go session.startProbe()
		})
	})

	pcA.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		ws.WriteJSON(map[string]string{"ice-candidate": c.ToJSON().Candidate})
	})

	pcA.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Println("pcA ICE状态:", state.String())
	})

	// pcB配置
	pcB.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		log.Println("pcB ICE候选者:", c.ToJSON().Candidate[:60])
		ws.WriteJSON(map[string]string{"probe-candidate": c.ToJSON().Candidate})
	})

	pcB.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Println("pcB ICE状态:", state.String())
	})

	pcB.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.Println("pcB收到DataChannel:", dc.Label())
		dc.OnOpen(func() {
			log.Println("pcB DataChannel已打开 - Full Cone!")
			session.mu.Lock()
			session.pcBRecv = true
			session.mu.Unlock()
		})
	})

	// 消息循环
	for {
		var msg Msg
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("读取消息失败:", err)
			break
		}

		if msg.SDP != "" {
			session.handleSDP(msg.SDP)
		}
		if msg.ICECandidate != "" {
			session.handleICE(msg.ICECandidate)
		}
		if msg.ProbeAnswer != "" {
			session.handleProbeAnswer(msg.ProbeAnswer)
		}
		if msg.ProbeCandidate != "" {
			session.handleProbeCandidate(msg.ProbeCandidate)
		}
	}
}

func (s *Session) handleSDP(sdp string) {
	log.Println("收到SDP Offer")
	s.pcA.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp})
	answer, _ := s.pcA.CreateAnswer(nil)
	s.pcA.SetLocalDescription(answer)
	s.ws.WriteJSON(map[string]string{"sdp": answer.SDP})
	log.Println("已发送SDP Answer")
}

func (s *Session) handleICE(candidate string) {
	s.pcA.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate})

	if strings.Contains(candidate, "srflx") && strings.Contains(candidate, "udp") {
		parts := strings.Split(candidate, " ")
		if len(parts) >= 6 {
			ip, port := parts[4], parts[5]
			if !strings.Contains(ip, ":") {
				s.mu.Lock()
				s.publicIP = ip
				s.ports[port] = true
				s.mu.Unlock()
				log.Printf("记录: IP=%s Port=%s\n", ip, port)
			}
		}
	}
}

func (s *Session) handleProbeAnswer(answer string) {
	log.Println("收到探测Answer")
	s.pcB.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer})
	s.mu.Lock()
	s.pcBReady = true
	s.mu.Unlock()
}

func (s *Session) handleProbeCandidate(candidate string) {
	log.Println("收到探测候选者")
	s.pcB.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate})
}

func (s *Session) startProbe() {
	time.Sleep(1 * time.Second)

	s.mu.Lock()
	ip := s.publicIP
	portCount := len(s.ports)
	s.mu.Unlock()

	if ip == "" {
		s.sendResult("Blocked", "")
		return
	}

	if portCount > 1 {
		s.sendResult("Symmetric", ip)
		return
	}

	// 使用pcB探测Full Cone
	log.Println("开始Full Cone探测...")

	dcB, _ := s.pcB.CreateDataChannel("probe", nil)
	dcB.OnOpen(func() {
		log.Println("pcB DataChannel已打开 - Full Cone!")
		s.mu.Lock()
		s.pcBRecv = true
		s.mu.Unlock()
	})

	offer, _ := s.pcB.CreateOffer(nil)
	s.pcB.SetLocalDescription(offer)

	// 发送探测offer给客户端
	s.ws.WriteJSON(map[string]string{"probe-offer": offer.SDP})

	// 等待连接建立
	time.Sleep(3 * time.Second)

	s.mu.Lock()
	connected := s.pcBRecv
	s.mu.Unlock()

	if connected {
		s.sendResult("Full Cone", ip)
	} else {
		s.sendResult("Port Restricted Cone", ip)
	}
}

func (s *Session) sendResult(natType, ip string) {
	result := map[string]string{"nat_type": natType, "public_ip": ip}
	data, _ := json.Marshal(result)
	log.Println("发送结果:", string(data))
	s.ws.WriteJSON(result)
}
