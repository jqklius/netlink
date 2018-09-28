package nl

import (
	"encoding/binary"
	"unsafe"
	"fmt"
	"net"
	"strconv"
)

// LinkLayer
const (
	LINKLAYER_UNSPEC = iota
	LINKLAYER_ETHERNET
	LINKLAYER_ATM
)

// ATM
const (
	ATM_CELL_PAYLOAD = 48
	ATM_CELL_SIZE    = 53
)

const TC_LINKLAYER_MASK = 0x0F

// Police
const (
	TCA_POLICE_UNSPEC = iota
	TCA_POLICE_TBF
	TCA_POLICE_RATE
	TCA_POLICE_PEAKRATE
	TCA_POLICE_AVRATE
	TCA_POLICE_RESULT
	TCA_POLICE_MAX = TCA_POLICE_RESULT
)

// Message types
const (
	TCA_UNSPEC = iota
	TCA_KIND
	TCA_OPTIONS
	TCA_STATS
	TCA_XSTATS
	TCA_RATE
	TCA_FCNT
	TCA_STATS2
	TCA_STAB
	TCA_MAX = TCA_STAB
)

const (
	TCA_ACT_TAB = 1
	TCAA_MAX    = 1
)

const (
	TCA_ACT_UNSPEC = iota
	TCA_ACT_KIND
	TCA_ACT_OPTIONS
	TCA_ACT_INDEX
	TCA_ACT_STATS
	TCA_ACT_MAX
)

const (
	TCA_PRIO_UNSPEC = iota
	TCA_PRIO_MQ
	TCA_PRIO_MAX = TCA_PRIO_MQ
)

const (
	TCA_STATS_UNSPEC = iota
	TCA_STATS_BASIC
	TCA_STATS_RATE_EST
	TCA_STATS_QUEUE
	TCA_STATS_APP
	TCA_STATS_MAX = TCA_STATS_APP
)

const (
	SizeofTcMsg          = 0x14
	SizeofTcActionMsg    = 0x04
	SizeofTcPrioMap      = 0x14
	SizeofTcRateSpec     = 0x0c
	SizeofTcNetemQopt    = 0x18
	SizeofTcNetemCorr    = 0x0c
	SizeofTcNetemReorder = 0x08
	SizeofTcNetemCorrupt = 0x08
	SizeofTcTbfQopt      = 2*SizeofTcRateSpec + 0x0c
	SizeofTcHtbCopt      = 2*SizeofTcRateSpec + 0x14
	SizeofTcHtbGlob      = 0x14
	SizeofTcU32Key       = 0x10
	SizeofTcU32Sel       = 0x10 // without keys
	SizeofTcGen          = 0x14
	SizeofTcMirred       = SizeofTcGen + 0x08
	SizeofTcPolice       = 2*SizeofTcRateSpec + 0x20
	SizeofTcTunnelKey    = SizeofTcGen + 0x04
)

// struct tcmsg {
//   unsigned char tcm_family;
//   unsigned char tcm__pad1;
//   unsigned short  tcm__pad2;
//   int   tcm_ifindex;
//   __u32   tcm_handle;
//   __u32   tcm_parent;
//   __u32   tcm_info;
// };

type TcMsg struct {
	Family  uint8
	Pad     [3]byte
	Ifindex int32
	Handle  uint32
	Parent  uint32
	Info    uint32
}

func (msg *TcMsg) Len() int {
	return SizeofTcMsg
}

func DeserializeTcMsg(b []byte) *TcMsg {
	return (*TcMsg)(unsafe.Pointer(&b[0:SizeofTcMsg][0]))
}

func (x *TcMsg) Serialize() []byte {
	return (*(*[SizeofTcMsg]byte)(unsafe.Pointer(x)))[:]
}

// struct tcamsg {
//   unsigned char tca_family;
//   unsigned char tca__pad1;
//   unsigned short  tca__pad2;
// };

type TcActionMsg struct {
	Family uint8
	Pad    [3]byte
}

func (msg *TcActionMsg) Len() int {
	return SizeofTcActionMsg
}

func DeserializeTcActionMsg(b []byte) *TcActionMsg {
	return (*TcActionMsg)(unsafe.Pointer(&b[0:SizeofTcActionMsg][0]))
}

func (x *TcActionMsg) Serialize() []byte {
	return (*(*[SizeofTcActionMsg]byte)(unsafe.Pointer(x)))[:]
}

const (
	TC_PRIO_MAX = 15
)

// struct tc_prio_qopt {
// 	int bands;      /* Number of bands */
// 	__u8  priomap[TC_PRIO_MAX+1]; /* Map: logical priority -> PRIO band */
// };

type TcPrioMap struct {
	Bands   int32
	Priomap [TC_PRIO_MAX + 1]uint8
}

func (msg *TcPrioMap) Len() int {
	return SizeofTcPrioMap
}

func DeserializeTcPrioMap(b []byte) *TcPrioMap {
	return (*TcPrioMap)(unsafe.Pointer(&b[0:SizeofTcPrioMap][0]))
}

func (x *TcPrioMap) Serialize() []byte {
	return (*(*[SizeofTcPrioMap]byte)(unsafe.Pointer(x)))[:]
}

const (
	TCA_TBF_UNSPEC = iota
	TCA_TBF_PARMS
	TCA_TBF_RTAB
	TCA_TBF_PTAB
	TCA_TBF_RATE64
	TCA_TBF_PRATE64
	TCA_TBF_BURST
	TCA_TBF_PBURST
	TCA_TBF_MAX = TCA_TBF_PBURST
)

// struct tc_ratespec {
//   unsigned char cell_log;
//   __u8    linklayer; /* lower 4 bits */
//   unsigned short  overhead;
//   short   cell_align;
//   unsigned short  mpu;
//   __u32   rate;
// };

type TcRateSpec struct {
	CellLog   uint8
	Linklayer uint8
	Overhead  uint16
	CellAlign int16
	Mpu       uint16
	Rate      uint32
}

func (msg *TcRateSpec) Len() int {
	return SizeofTcRateSpec
}

func DeserializeTcRateSpec(b []byte) *TcRateSpec {
	return (*TcRateSpec)(unsafe.Pointer(&b[0:SizeofTcRateSpec][0]))
}

func (x *TcRateSpec) Serialize() []byte {
	return (*(*[SizeofTcRateSpec]byte)(unsafe.Pointer(x)))[:]
}

/**
* NETEM
 */

const (
	TCA_NETEM_UNSPEC = iota
	TCA_NETEM_CORR
	TCA_NETEM_DELAY_DIST
	TCA_NETEM_REORDER
	TCA_NETEM_CORRUPT
	TCA_NETEM_LOSS
	TCA_NETEM_RATE
	TCA_NETEM_ECN
	TCA_NETEM_RATE64
	TCA_NETEM_MAX = TCA_NETEM_RATE64
)

// struct tc_netem_qopt {
//	__u32	latency;	/* added delay (us) */
//	__u32   limit;		/* fifo limit (packets) */
//	__u32	loss;		/* random packet loss (0=none ~0=100%) */
//	__u32	gap;		/* re-ordering gap (0 for none) */
//	__u32   duplicate;	/* random packet dup  (0=none ~0=100%) */
// 	__u32	jitter;		/* random jitter in latency (us) */
// };

type TcNetemQopt struct {
	Latency   uint32
	Limit     uint32
	Loss      uint32
	Gap       uint32
	Duplicate uint32
	Jitter    uint32
}

func (msg *TcNetemQopt) Len() int {
	return SizeofTcNetemQopt
}

func DeserializeTcNetemQopt(b []byte) *TcNetemQopt {
	return (*TcNetemQopt)(unsafe.Pointer(&b[0:SizeofTcNetemQopt][0]))
}

func (x *TcNetemQopt) Serialize() []byte {
	return (*(*[SizeofTcNetemQopt]byte)(unsafe.Pointer(x)))[:]
}

// struct tc_netem_corr {
//  __u32   delay_corr; /* delay correlation */
//  __u32   loss_corr;  /* packet loss correlation */
//  __u32   dup_corr;   /* duplicate correlation  */
// };

type TcNetemCorr struct {
	DelayCorr uint32
	LossCorr  uint32
	DupCorr   uint32
}

func (msg *TcNetemCorr) Len() int {
	return SizeofTcNetemCorr
}

func DeserializeTcNetemCorr(b []byte) *TcNetemCorr {
	return (*TcNetemCorr)(unsafe.Pointer(&b[0:SizeofTcNetemCorr][0]))
}

func (x *TcNetemCorr) Serialize() []byte {
	return (*(*[SizeofTcNetemCorr]byte)(unsafe.Pointer(x)))[:]
}

// struct tc_netem_reorder {
//  __u32   probability;
//  __u32   correlation;
// };

type TcNetemReorder struct {
	Probability uint32
	Correlation uint32
}

func (msg *TcNetemReorder) Len() int {
	return SizeofTcNetemReorder
}

func DeserializeTcNetemReorder(b []byte) *TcNetemReorder {
	return (*TcNetemReorder)(unsafe.Pointer(&b[0:SizeofTcNetemReorder][0]))
}

func (x *TcNetemReorder) Serialize() []byte {
	return (*(*[SizeofTcNetemReorder]byte)(unsafe.Pointer(x)))[:]
}

// struct tc_netem_corrupt {
//  __u32   probability;
//  __u32   correlation;
// };

type TcNetemCorrupt struct {
	Probability uint32
	Correlation uint32
}

func (msg *TcNetemCorrupt) Len() int {
	return SizeofTcNetemCorrupt
}

func DeserializeTcNetemCorrupt(b []byte) *TcNetemCorrupt {
	return (*TcNetemCorrupt)(unsafe.Pointer(&b[0:SizeofTcNetemCorrupt][0]))
}

func (x *TcNetemCorrupt) Serialize() []byte {
	return (*(*[SizeofTcNetemCorrupt]byte)(unsafe.Pointer(x)))[:]
}

// struct tc_tbf_qopt {
//   struct tc_ratespec rate;
//   struct tc_ratespec peakrate;
//   __u32   limit;
//   __u32   buffer;
//   __u32   mtu;
// };

type TcTbfQopt struct {
	Rate     TcRateSpec
	Peakrate TcRateSpec
	Limit    uint32
	Buffer   uint32
	Mtu      uint32
}

func (msg *TcTbfQopt) Len() int {
	return SizeofTcTbfQopt
}

func DeserializeTcTbfQopt(b []byte) *TcTbfQopt {
	return (*TcTbfQopt)(unsafe.Pointer(&b[0:SizeofTcTbfQopt][0]))
}

func (x *TcTbfQopt) Serialize() []byte {
	return (*(*[SizeofTcTbfQopt]byte)(unsafe.Pointer(x)))[:]
}

const (
	TCA_HTB_UNSPEC = iota
	TCA_HTB_PARMS
	TCA_HTB_INIT
	TCA_HTB_CTAB
	TCA_HTB_RTAB
	TCA_HTB_DIRECT_QLEN
	TCA_HTB_RATE64
	TCA_HTB_CEIL64
	TCA_HTB_MAX = TCA_HTB_CEIL64
)

//struct tc_htb_opt {
//	struct tc_ratespec	rate;
//	struct tc_ratespec	ceil;
//	__u32	buffer;
//	__u32	cbuffer;
//	__u32	quantum;
//	__u32	level;		/* out only */
//	__u32	prio;
//};

type TcHtbCopt struct {
	Rate    TcRateSpec
	Ceil    TcRateSpec
	Buffer  uint32
	Cbuffer uint32
	Quantum uint32
	Level   uint32
	Prio    uint32
}

func (msg *TcHtbCopt) Len() int {
	return SizeofTcHtbCopt
}

func DeserializeTcHtbCopt(b []byte) *TcHtbCopt {
	return (*TcHtbCopt)(unsafe.Pointer(&b[0:SizeofTcHtbCopt][0]))
}

func (x *TcHtbCopt) Serialize() []byte {
	return (*(*[SizeofTcHtbCopt]byte)(unsafe.Pointer(x)))[:]
}

type TcHtbGlob struct {
	Version      uint32
	Rate2Quantum uint32
	Defcls       uint32
	Debug        uint32
	DirectPkts   uint32
}

func (msg *TcHtbGlob) Len() int {
	return SizeofTcHtbGlob
}

func DeserializeTcHtbGlob(b []byte) *TcHtbGlob {
	return (*TcHtbGlob)(unsafe.Pointer(&b[0:SizeofTcHtbGlob][0]))
}

func (x *TcHtbGlob) Serialize() []byte {
	return (*(*[SizeofTcHtbGlob]byte)(unsafe.Pointer(x)))[:]
}

// HFSC

type Curve struct {
	m1 uint32
	d  uint32
	m2 uint32
}

type HfscCopt struct {
	Rsc Curve
	Fsc Curve
	Usc Curve
}

func (c *Curve) Attrs() (uint32, uint32, uint32) {
	return c.m1, c.d, c.m2
}

func (c *Curve) Set(m1 uint32, d uint32, m2 uint32) {
	c.m1 = m1
	c.d = d
	c.m2 = m2
}

func DeserializeHfscCurve(b []byte) *Curve {
	return &Curve{
		m1: binary.LittleEndian.Uint32(b[0:4]),
		d:  binary.LittleEndian.Uint32(b[4:8]),
		m2: binary.LittleEndian.Uint32(b[8:12]),
	}
}

func SerializeHfscCurve(c *Curve) (b []byte) {
	t := make([]byte, binary.MaxVarintLen32)
	binary.LittleEndian.PutUint32(t, c.m1)
	b = append(b, t[:4]...)
	binary.LittleEndian.PutUint32(t, c.d)
	b = append(b, t[:4]...)
	binary.LittleEndian.PutUint32(t, c.m2)
	b = append(b, t[:4]...)
	return b
}

type TcHfscOpt struct {
	Defcls uint16
}

func (x *TcHfscOpt) Serialize() []byte {
	return (*(*[2]byte)(unsafe.Pointer(x)))[:]
}

const (
	TCA_U32_UNSPEC = iota
	TCA_U32_CLASSID
	TCA_U32_HASH
	TCA_U32_LINK
	TCA_U32_DIVISOR
	TCA_U32_SEL
	TCA_U32_POLICE
	TCA_U32_ACT
	TCA_U32_INDEV
	TCA_U32_PCNT
	TCA_U32_MARK
	TCA_U32_MAX = TCA_U32_MARK
)

// struct tc_u32_key {
//   __be32    mask;
//   __be32    val;
//   int   off;
//   int   offmask;
// };

type TcU32Key struct {
	Mask    uint32 // big endian
	Val     uint32 // big endian
	Off     int32
	OffMask int32
}

func (msg *TcU32Key) Len() int {
	return SizeofTcU32Key
}

func DeserializeTcU32Key(b []byte) *TcU32Key {
	return (*TcU32Key)(unsafe.Pointer(&b[0:SizeofTcU32Key][0]))
}

func (x *TcU32Key) Serialize() []byte {
	return (*(*[SizeofTcU32Key]byte)(unsafe.Pointer(x)))[:]
}

// struct tc_u32_sel {
//   unsigned char   flags;
//   unsigned char   offshift;
//   unsigned char   nkeys;
//
//   __be16      offmask;
//   __u16     off;
//   short     offoff;
//
//   short     hoff;
//   __be32      hmask;
//   struct tc_u32_key keys[0];
// };

const (
	TC_U32_TERMINAL  = 1 << iota
	TC_U32_OFFSET    = 1 << iota
	TC_U32_VAROFFSET = 1 << iota
	TC_U32_EAT       = 1 << iota
)

type TcU32Sel struct {
	Flags    uint8
	Offshift uint8
	Nkeys    uint8
	Pad      uint8
	Offmask  uint16 // big endian
	Off      uint16
	Offoff   int16
	Hoff     int16
	Hmask    uint32 // big endian
	Keys     []TcU32Key
}

func (msg *TcU32Sel) Len() int {
	return SizeofTcU32Sel + int(msg.Nkeys)*SizeofTcU32Key
}

func DeserializeTcU32Sel(b []byte) *TcU32Sel {
	x := &TcU32Sel{}
	copy((*(*[SizeofTcU32Sel]byte)(unsafe.Pointer(x)))[:], b)
	next := SizeofTcU32Sel
	var i uint8
	for i = 0; i < x.Nkeys; i++ {
		x.Keys = append(x.Keys, *DeserializeTcU32Key(b[next:]))
		next += SizeofTcU32Key
	}
	return x
}

func (x *TcU32Sel) Serialize() []byte {
	// This can't just unsafe.cast because it must iterate through keys.
	buf := make([]byte, x.Len())
	copy(buf, (*(*[SizeofTcU32Sel]byte)(unsafe.Pointer(x)))[:])
	next := SizeofTcU32Sel
	for _, key := range x.Keys {
		keyBuf := key.Serialize()
		copy(buf[next:], keyBuf)
		next += SizeofTcU32Key
	}
	return buf
}

type TcGen struct {
	Index   uint32
	Capab   uint32
	Action  int32
	Refcnt  int32
	Bindcnt int32
}

func (msg *TcGen) Len() int {
	return SizeofTcGen
}

func DeserializeTcGen(b []byte) *TcGen {
	return (*TcGen)(unsafe.Pointer(&b[0:SizeofTcGen][0]))
}

func (x *TcGen) Serialize() []byte {
	return (*(*[SizeofTcGen]byte)(unsafe.Pointer(x)))[:]
}

// #define tc_gen \
//   __u32                 index; \
//   __u32                 capab; \
//   int                   action; \
//   int                   refcnt; \
//   int                   bindcnt

const (
	TCA_ACT_GACT = 5
)

const (
	TCA_GACT_UNSPEC = iota
	TCA_GACT_TM
	TCA_GACT_PARMS
	TCA_GACT_PROB
	TCA_GACT_MAX = TCA_GACT_PROB
)

type TcGact TcGen

const (
	TCA_ACT_BPF = 13
)

const (
	TCA_ACT_BPF_UNSPEC = iota
	TCA_ACT_BPF_TM
	TCA_ACT_BPF_PARMS
	TCA_ACT_BPF_OPS_LEN
	TCA_ACT_BPF_OPS
	TCA_ACT_BPF_FD
	TCA_ACT_BPF_NAME
	TCA_ACT_BPF_MAX = TCA_ACT_BPF_NAME
)

const (
	TCA_BPF_FLAG_ACT_DIRECT uint32 = 1 << iota
)

const (
	TCA_BPF_UNSPEC = iota
	TCA_BPF_ACT
	TCA_BPF_POLICE
	TCA_BPF_CLASSID
	TCA_BPF_OPS_LEN
	TCA_BPF_OPS
	TCA_BPF_FD
	TCA_BPF_NAME
	TCA_BPF_FLAGS
	TCA_BPF_MAX = TCA_BPF_FLAGS
)

type TcBpf TcGen

const (
	TCA_ACT_MIRRED = 8
)

const (
	TCA_MIRRED_UNSPEC = iota
	TCA_MIRRED_TM
	TCA_MIRRED_PARMS
	TCA_MIRRED_MAX = TCA_MIRRED_PARMS
)

// struct tc_mirred {
// 	tc_gen;
// 	int                     eaction;   /* one of IN/EGRESS_MIRROR/REDIR */
// 	__u32                   ifindex;  /* ifindex of egress port */
// };

type TcMirred struct {
	TcGen
	Eaction int32
	Ifindex uint32
}

func (msg *TcMirred) Len() int {
	return SizeofTcMirred
}

func DeserializeTcMirred(b []byte) *TcMirred {
	return (*TcMirred)(unsafe.Pointer(&b[0:SizeofTcMirred][0]))
}

func (x *TcMirred) Serialize() []byte {
	return (*(*[SizeofTcMirred]byte)(unsafe.Pointer(x)))[:]
}

// tunnel_key begin

const (
	TCA_TUNNEL_KEY_UNSPEC = iota
	TCA_TUNNEL_KEY_TM
	TCA_TUNNEL_KEY_PARMS
	TCA_TUNNEL_KEY_ENC_IPV4_SRC	/* be32 */
	TCA_TUNNEL_KEY_ENC_IPV4_DST	/* be32 */
	TCA_TUNNEL_KEY_ENC_IPV6_SRC	/* struct in6_addr */
	TCA_TUNNEL_KEY_ENC_IPV6_DST	/* struct in6_addr */
	TCA_TUNNEL_KEY_ENC_KEY_ID	/* be64 */
	TCA_TUNNEL_KEY_PAD
	TCA_TUNNEL_KEY_ENC_DST_PORT	/* be16 */
	TCA_TUNNEL_KEY_NO_CSUM		/* u8 */
	TCA_TUNNEL_KEY_MAX = TCA_TUNNEL_KEY_NO_CSUM
)

type TcTunnelKey struct {
	TcGen
	T_ACTION int32
}

func (msg *TcTunnelKey) Len() int {
	return SizeofTcTunnelKey
}

func DeserializeTcTunnelKey(b []byte) *TcTunnelKey {
	return (*TcTunnelKey)(unsafe.Pointer(&b[0:SizeofTcTunnelKey][0]))
}

func (x *TcTunnelKey) Serialize() []byte {
	return (*(*[SizeofTcTunnelKey]byte)(unsafe.Pointer(x)))[:]
}

// tunnel_key end


// tc pedit begin

const (
	PeditDebug = true
	MAX_OFFS = 128
)

const (
	TIPV4 = 1
	TIPV6 = 2
	TINT = 3
	TU32 = 4
	TMAC = 5
)

const (
	RU32 = 0xFFFFFFFF
	RU16 = 0xFFFF
	RU8 = 0xFF
)

const (
	TCA_PEDIT_UNSPEC = iota
	TCA_PEDIT_TM
	TCA_PEDIT_PARMS
	TCA_PEDIT_PAD
	TCA_PEDIT_PARMS_EX
	TCA_PEDIT_KEYS_EX
	TCA_PEDIT_KEY_EX
	TCA_PEDIT_MAX = TCA_PEDIT_KEY_EX
)

const (
	TCA_PEDIT_KEY_EX_CMD_SET = 0
	TCA_PEDIT_KEY_EX_CMD_ADD = 1
	TCA_PEDIT_CMD_MAX = TCA_PEDIT_KEY_EX_CMD_ADD
)

const (
	TCA_PEDIT_KEY_EX_HTYPE = 1
	TCA_PEDIT_KEY_EX_CMD = 2
	TCA_PEDIT_KEY_EX_MAX = TCA_PEDIT_KEY_EX_CMD
)

const (
	TCA_PEDIT_KEY_EX_HDR_TYPE_NETWORK = iota
	TCA_PEDIT_KEY_EX_HDR_TYPE_ETH
	TCA_PEDIT_KEY_EX_HDR_TYPE_IP4
	TCA_PEDIT_KEY_EX_HDR_TYPE_IP6
	TCA_PEDIT_KEY_EX_HDR_TYPE_TCP
	TCA_PEDIT_KEY_EX_HDR_TYPE_UDP
	TCA_PEDIT_HDR_TYPE_MAX = TCA_PEDIT_KEY_EX_HDR_TYPE_UDP
)

type TcPeditKey struct {
   Mask  	uint32 /* AND */
   Val   	uint32 /*XOR */
   Off   	uint32 /*offset */
   At		uint32
   Offmask	uint32
   Shift	uint32
}

type TcPeditSel struct {
	TcGen
	Nkeys	uint8
	Flags 	uint8
	Keys	[MAX_OFFS]TcPeditKey
}

func (tcsel *TcPeditSel) Len() int {
	return (int)(unsafe.Offsetof(tcsel.Keys)) + (int)(tcsel.Nkeys) * (int)(unsafe.Sizeof(TcPeditKey{}))
}

func DeserializeTcPeditSel(b []byte) *TcPeditSel {
	//return (*TcTunnelKey)(unsafe.Pointer(&b[0:SizeofTcTunnelKey][0]))
	return nil
}

func (tcsel *TcPeditSel) Serialize() []byte {
	const size = (uint32)(unsafe.Sizeof(*tcsel))
	datalen := tcsel.Len()
	fmt.Printf("size:%v, datalen:%v\n", size, datalen)
	return (*(*[size]byte)(unsafe.Pointer(tcsel)))[:datalen]
}

type PeditCMD struct {
	Key		string
	Action	string
	Val		string
}

type MPeditKey struct {
    Mask 	uint32  /* AND */
    Val		uint32	  /*XOR */
    Off		uint32 /*offset */
	At		uint32
	Offmask	uint32
	Shift	uint32

	Htype 	uint32
	Cmd		uint32
}

type MPeditKeyEx struct {
	Htype	uint32
	Cmd		uint32
}

type MPeditSel struct {
	Sel TcPeditSel
	//Keys[MAX_OFFS] TcPeditKey
	Keys_ex[MAX_OFFS] MPeditKeyEx
	Extended bool
}

func PackKey(_sel *MPeditSel, tkey *MPeditKey) int {
	var sel *TcPeditSel = nil
	var keys_ex *[MAX_OFFS] MPeditKeyEx = nil
	var hwm uint8
	
	sel = &_sel.Sel
	keys_ex = &_sel.Keys_ex
	hwm = sel.Nkeys
	
	if (hwm >= MAX_OFFS) {
		return -1
	}
	
	if ((tkey.Off % 4) > 0) {
		fmt.Printf("offsets MUST be in 32 bit boundaries\n")
		return -1
	}

	sel.Keys[hwm].Val = tkey.Val;
	sel.Keys[hwm].Mask = tkey.Mask;
	sel.Keys[hwm].Off = tkey.Off;
	sel.Keys[hwm].At = tkey.At;
	sel.Keys[hwm].Offmask = tkey.Offmask;
	sel.Keys[hwm].Shift = tkey.Shift;
	
	if (_sel.Extended) {
		keys_ex[hwm].Htype = tkey.Htype;
		keys_ex[hwm].Cmd = tkey.Cmd;
		fmt.Printf("--- tkey.Cmd:%v\n", tkey.Cmd)
	} else {
		if (tkey.Htype != TCA_PEDIT_KEY_EX_HDR_TYPE_NETWORK ||
		    tkey.Cmd != TCA_PEDIT_KEY_EX_CMD_SET) {
			fmt.Printf(
				"Munge parameters not supported. Use 'pedit ex munge ...'.\n")
			return -1
		}
	}

	sel.Nkeys++
	return 0
	
}

func Btou32l (b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func Htonl(val uint32) []byte {
	fmt.Printf("Htonl: %v\n", val)
	
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, val)
	
	fmt.Printf("htonl_ret: %v\n", Btou32l(bytes))
	return bytes
}

func Ntohl(buf []byte) uint32 {
	fmt.Printf("Ntohl: %v\n", Btou32l(buf))

	ret := binary.BigEndian.Uint32(buf)
	
	fmt.Printf("ntohl_ret: %v\n", ret)
	return ret
}

func Ntohs(buf []byte) uint16 {
	fmt.Printf("Ntohs: %v\n", binary.LittleEndian.Uint16(buf))
	
	ret := binary.BigEndian.Uint16(buf)
	fmt.Printf("ntohs_ret: %v\n", ret)
	return ret
}

func PackKey32(retain uint32, sel *MPeditSel, tkey *MPeditKey) int {
	if (tkey.Off > (tkey.Off & (^(uint32)(3)))) {
		fmt.Printf(
			"PackKey32: 32 bit offsets must begin in 32bit boundaries\n")
		return -1
	}

	var ret []byte
	ret = Htonl(tkey.Val & retain)
	fmt.Printf("key32:val_ret:%v, \n", Btou32l(ret))
	
	tkey.Val = Btou32l(ret)
	
	ret = Htonl(tkey.Mask | ^retain)
	fmt.Printf("key32:mask_ret:%v, \n", Btou32l(ret))
	
	tkey.Mask = Btou32l(ret)
	return PackKey(sel, tkey)
}

func PackKey16(retain uint32, sel *MPeditSel, tkey *MPeditKey) int {
	var ind, stride uint32
	var m = [4] uint32 { 0x0000FFFF, 0xFF0000FF, 0xFFFF0000 }

	if (tkey.Val > 0xFFFF || tkey.Mask > 0xFFFF) {
		fmt.Printf("PackKey16 bad value\n")
		return -1
	}

	ind = (uint32)(tkey.Off & 3)

	if (ind == 3) {
		fmt.Printf("PackKey16 bad index value %d\n", ind)
		return -1
	}

	stride = 8 * (2 - ind)
	
	var ret []byte
	ret = Htonl((tkey.Val & retain) << stride)
	tkey.Val = Btou32l(ret)
	
	ret = Htonl(((tkey.Mask | ^retain) << stride) | m[ind])
	tkey.Mask = Btou32l(ret)

	tkey.Off &= ^(uint32)(3)

	if (PeditDebug) {
		fmt.Printf("PackKey16: Final val %08x mask %08x\n",
		       tkey.Val, tkey.Mask)
	}
	return PackKey(sel, tkey)

}

func PackKey8(retain uint32, sel *MPeditSel, tkey *MPeditKey) int {
	var ind, stride uint32
	var m = [4] uint32 { 0x00FFFFFF, 0xFF00FFFF, 0xFFFF00FF, 0xFFFFFF00 }
	
	if (tkey.Val > 0xFF || tkey.Mask > 0xFF) {
		fmt.Printf("PackKey8 bad value (val %x mask %x\n",
			tkey.Val, tkey.Mask)
		return -1
	}

	ind = (uint32)(tkey.Off & 3)

	stride = 8 * (3 - ind)
	
	var ret []byte
	ret = Htonl((tkey.Val & retain) << stride)
	tkey.Val = Btou32l(ret)
	
	ret = Htonl(((tkey.Mask | ^retain) << stride) | m[ind])
	tkey.Mask = Btou32l(ret)

	tkey.Off &= ^(uint32)(3)

	if (PeditDebug) {
		fmt.Printf("PackKey8: Final word off %d  val %08x mask %08x\n",
		       tkey.Off, tkey.Val, tkey.Mask)
	}
	return PackKey(sel, tkey)
}

func PackMac(sel *MPeditSel, tkey *MPeditKey, mac [6]byte) int {
	var ret int = 0

	if (0 == (tkey.Off & 0x3)) {
		fmt.Printf("tkey.Off & 0x3 ...");
		
		tkey.Mask = 0
		tkey.Val = Ntohl(mac[0:4])
		ret |= PackKey32(^(uint32)(0), sel, tkey)

		tkey.Off += 4
		tkey.Mask = 0
		tkey.Val = (uint32)(Ntohs(mac[4:6]))
		
		fmt.Printf("off & 0x3, val:%d\n", tkey.Val);
		
		ret |= PackKey16(^(uint32)(0), sel, tkey)
	} else if ( 0 == (tkey.Off & 0x1)) {
		fmt.Printf("tkey.Off & 0x1 ...");
		
		tkey.Mask = 0
		tkey.Val = (uint32)(Ntohs(mac[0:2]))
		ret |= PackKey16(^(uint32)(0), sel, tkey)

		tkey.Off += 4
		tkey.Mask = 0
		tkey.Val = Ntohl((mac[2:6]))
		
		fmt.Printf("off & 0x1, val:%d\n", tkey.Val);
		
		ret |= PackKey32(^(uint32)(0), sel, tkey)
	} else {
		fmt.Printf(
			"PackMac: mac offsets must begin in 32bit or 16bit boundaries\n")
		return -1
	}

	return ret
}

func GetInteger(s string) (int, error) {
	ret, err := strconv.Atoi(s)
	return ret, err
}

func GetU32(s string) (uint32, error) {
	uret64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	uret32 := (uint32)(uret64)
	return uret32, err
}

func GetIpv4(s string) net.IP {
	return net.ParseIP(s)
}

func GetMac(s string) (hw net.HardwareAddr, err error) {
	return net.ParseMAC(s)
}

func ParseCmd(cmd PeditCMD, cmd_len uint32, cmd_type int, retain uint32, sel *MPeditSel, tkey *MPeditKey) int {
	var mask,val uint32 = 0, 0
	/*
	var m *uint32
	var v *uint32
	var o uint32 = 0xFF
	*/
	var res int = -1
	
	if len(cmd.Key) <= 0 {
		return -1
	}
	
	if (PeditDebug){ 
		fmt.Printf("ParseCmd %v length %d, val:%d, mask:%d\n",
			       cmd, cmd_len, val, mask)
	}

	/*
	// TODO 
	if (len == 2)
		o = 0xFFFF;
	if (len == 4)
		o = 0xFFFFFFFF;
	*/
	
	if (cmd.Action == "add") {
		tkey.Cmd = TCA_PEDIT_KEY_EX_CMD_ADD
	}
	
	if (cmd_type == TMAC) {
		hardware_addr, err := GetMac(cmd.Val)
		if err != nil {
			fmt.Printf("get mac from %v failed\n", cmd.Val)
			return -1
		}
		
		var mac [6]byte
		for index, val := range hardware_addr {
			mac[index] = val
		}
		
		res = PackMac(sel, tkey, mac)
		goto done;
	}
	
	if (cmd_type == TIPV4) {
		netip := GetIpv4(cmd.Val)
		netip = netip.To4()
		tkey.Val = Ntohl(netip)
		res = 0
	}

	if (cmd_type == TINT) {
		ret, err := GetInteger(cmd.Val)
		if err != nil {
			goto done
		}
		tkey.Val = (uint32)(ret)
		res = 0
	}

	if (cmd_type == TU32) {
		ret, err := GetU32(cmd.Val)		
		if err != nil {
			goto done
		}
		tkey.Val = ret
		res = 0
	}

	if (cmd_len == 1) {
		res = PackKey8(retain, sel, tkey);
		goto done;
	}
	if (cmd_len == 2) {
		res = PackKey16(retain, sel, tkey);
		goto done;
	}
	if (cmd_len == 4) {
		res = PackKey32(retain, sel, tkey);
		goto done;
	}

	return -1;
done:
	if (PeditDebug) {
		fmt.Printf("ParseCmd done: offset %d length %d\n",
			        tkey.Off, cmd_len)
		fmt.Printf("ParseCmd done: word off %d val %08x mask %08x\n",
		       tkey.Off, tkey.Val, tkey.Mask)
	}

	return res;
}

// tc pedit end


// struct tc_police {
// 	__u32			index;
// 	int			action;
// 	__u32			limit;
// 	__u32			burst;
// 	__u32			mtu;
// 	struct tc_ratespec	rate;
// 	struct tc_ratespec	peakrate;
// 	int				refcnt;
// 	int				bindcnt;
// 	__u32			capab;
// };

type TcPolice struct {
	Index    uint32
	Action   int32
	Limit    uint32
	Burst    uint32
	Mtu      uint32
	Rate     TcRateSpec
	PeakRate TcRateSpec
	Refcnt   int32
	Bindcnt  int32
	Capab    uint32
}

func (msg *TcPolice) Len() int {
	return SizeofTcPolice
}

func DeserializeTcPolice(b []byte) *TcPolice {
	return (*TcPolice)(unsafe.Pointer(&b[0:SizeofTcPolice][0]))
}

func (x *TcPolice) Serialize() []byte {
	return (*(*[SizeofTcPolice]byte)(unsafe.Pointer(x)))[:]
}

const (
	TCA_FW_UNSPEC = iota
	TCA_FW_CLASSID
	TCA_FW_POLICE
	TCA_FW_INDEV
	TCA_FW_ACT
	TCA_FW_MASK
	TCA_FW_MAX = TCA_FW_MASK
)

const (
	TCA_MATCHALL_UNSPEC = iota
	TCA_MATCHALL_CLASSID
	TCA_MATCHALL_ACT
	TCA_MATCHALL_FLAGS
)

const (
	TCA_FQ_UNSPEC             = iota
	TCA_FQ_PLIMIT             // limit of total number of packets in queue
	TCA_FQ_FLOW_PLIMIT        // limit of packets per flow
	TCA_FQ_QUANTUM            // RR quantum
	TCA_FQ_INITIAL_QUANTUM    // RR quantum for new flow
	TCA_FQ_RATE_ENABLE        // enable/disable rate limiting
	TCA_FQ_FLOW_DEFAULT_RATE  // obsolete do not use
	TCA_FQ_FLOW_MAX_RATE      // per flow max rate
	TCA_FQ_BUCKETS_LOG        // log2(number of buckets)
	TCA_FQ_FLOW_REFILL_DELAY  // flow credit refill delay in usec
	TCA_FQ_ORPHAN_MASK        // mask applied to orphaned skb hashes
	TCA_FQ_LOW_RATE_THRESHOLD // per packet delay under this rate
)

const (
	TCA_FQ_CODEL_UNSPEC = iota
	TCA_FQ_CODEL_TARGET
	TCA_FQ_CODEL_LIMIT
	TCA_FQ_CODEL_INTERVAL
	TCA_FQ_CODEL_ECN
	TCA_FQ_CODEL_FLOWS
	TCA_FQ_CODEL_QUANTUM
	TCA_FQ_CODEL_CE_THRESHOLD
	TCA_FQ_CODEL_DROP_BATCH_SIZE
	TCA_FQ_CODEL_MEMORY_LIMIT
)

const (
	TCA_HFSC_UNSPEC = iota
	TCA_HFSC_RSC
	TCA_HFSC_FSC
	TCA_HFSC_USC
)

/* Flower classifier */

const (
	TCA_FLOWER_UNSPEC = iota
	TCA_FLOWER_CLASSID
	TCA_FLOWER_INDEV
	TCA_FLOWER_ACT
	TCA_FLOWER_KEY_ETH_DST		/* ETH_ALEN */
	TCA_FLOWER_KEY_ETH_DST_MASK	/* ETH_ALEN */
	TCA_FLOWER_KEY_ETH_SRC		/* ETH_ALEN */
	TCA_FLOWER_KEY_ETH_SRC_MASK	/* ETH_ALEN */
	TCA_FLOWER_KEY_ETH_TYPE	/* be16 */
	TCA_FLOWER_KEY_IP_PROTO	/* u8 */
	TCA_FLOWER_KEY_IPV4_SRC	/* be32 */
	TCA_FLOWER_KEY_IPV4_SRC_MASK	/* be32 */
	TCA_FLOWER_KEY_IPV4_DST	/* be32 */
	TCA_FLOWER_KEY_IPV4_DST_MASK	/* be32 */
	TCA_FLOWER_KEY_IPV6_SRC	/* struct in6_addr */
	TCA_FLOWER_KEY_IPV6_SRC_MASK	/* struct in6_addr */
	TCA_FLOWER_KEY_IPV6_DST	/* struct in6_addr */
	TCA_FLOWER_KEY_IPV6_DST_MASK	/* struct in6_addr */
	TCA_FLOWER_KEY_TCP_SRC		/* be16 */
	TCA_FLOWER_KEY_TCP_DST		/* be16 */
	TCA_FLOWER_KEY_UDP_SRC		/* be16 */
	TCA_FLOWER_KEY_UDP_DST		/* be16 */

	TCA_FLOWER_FLAGS
	TCA_FLOWER_KEY_VLAN_ID		/* be16 */
	TCA_FLOWER_KEY_VLAN_PRIO	/* u8   */
	TCA_FLOWER_KEY_VLAN_ETH_TYPE	/* be16 */

	TCA_FLOWER_KEY_ENC_KEY_ID	/* be32 */
	TCA_FLOWER_KEY_ENC_IPV4_SRC	/* be32 */
	TCA_FLOWER_KEY_ENC_IPV4_SRC_MASK/* be32 */
	TCA_FLOWER_KEY_ENC_IPV4_DST	/* be32 */
	TCA_FLOWER_KEY_ENC_IPV4_DST_MASK/* be32 */
	TCA_FLOWER_KEY_ENC_IPV6_SRC	/* struct in6_addr */
	TCA_FLOWER_KEY_ENC_IPV6_SRC_MASK/* struct in6_addr */
	TCA_FLOWER_KEY_ENC_IPV6_DST	/* struct in6_addr */
	TCA_FLOWER_KEY_ENC_IPV6_DST_MASK/* struct in6_addr */

	TCA_FLOWER_KEY_TCP_SRC_MASK	/* be16 */
	TCA_FLOWER_KEY_TCP_DST_MASK	/* be16 */
	TCA_FLOWER_KEY_UDP_SRC_MASK	/* be16 */
	TCA_FLOWER_KEY_UDP_DST_MASK	/* be16 */
	TCA_FLOWER_KEY_SCTP_SRC_MASK	/* be16 */
	TCA_FLOWER_KEY_SCTP_DST_MASK	/* be16 */

	TCA_FLOWER_KEY_SCTP_SRC	/* be16 */
	TCA_FLOWER_KEY_SCTP_DST	/* be16 */

	TCA_FLOWER_KEY_ENC_UDP_SRC_PORT	/* be16 */
	TCA_FLOWER_KEY_ENC_UDP_SRC_PORT_MASK	/* be16 */
	TCA_FLOWER_KEY_ENC_UDP_DST_PORT	/* be16 */
	TCA_FLOWER_KEY_ENC_UDP_DST_PORT_MASK	/* be16 */

	TCA_FLOWER_KEY_FLAGS		/* be32 */
	TCA_FLOWER_KEY_FLAGS_MASK	/* be32 */

	TCA_FLOWER_KEY_ICMPV4_CODE	/* u8 */
	TCA_FLOWER_KEY_ICMPV4_CODE_MASK/* u8 */
	TCA_FLOWER_KEY_ICMPV4_TYPE	/* u8 */
	TCA_FLOWER_KEY_ICMPV4_TYPE_MASK/* u8 */
	TCA_FLOWER_KEY_ICMPV6_CODE	/* u8 */
	TCA_FLOWER_KEY_ICMPV6_CODE_MASK/* u8 */
	TCA_FLOWER_KEY_ICMPV6_TYPE	/* u8 */
	TCA_FLOWER_KEY_ICMPV6_TYPE_MASK/* u8 */

	TCA_FLOWER_KEY_ARP_SIP		/* be32 */
	TCA_FLOWER_KEY_ARP_SIP_MASK	/* be32 */
	TCA_FLOWER_KEY_ARP_TIP		/* be32 */
	TCA_FLOWER_KEY_ARP_TIP_MASK	/* be32 */
	TCA_FLOWER_KEY_ARP_OP		/* u8 */
	TCA_FLOWER_KEY_ARP_OP_MASK	/* u8 */
	TCA_FLOWER_KEY_ARP_SHA		/* ETH_ALEN */
	TCA_FLOWER_KEY_ARP_SHA_MASK	/* ETH_ALEN */
	TCA_FLOWER_KEY_ARP_THA		/* ETH_ALEN */
	TCA_FLOWER_KEY_ARP_THA_MASK	/* ETH_ALEN */

	TCA_FLOWER_KEY_MPLS_TTL	/* u8 - 8 bits */
	TCA_FLOWER_KEY_MPLS_BOS	/* u8 - 1 bit */
	TCA_FLOWER_KEY_MPLS_TC		/* u8 - 3 bits */
	TCA_FLOWER_KEY_MPLS_LABEL	/* be32 - 20 bits */

	TCA_FLOWER_KEY_TCP_FLAGS	/* be16 */
	TCA_FLOWER_KEY_TCP_FLAGS_MASK	/* be16 */

	TCA_FLOWER_KEY_IP_TOS		/* u8 */
	TCA_FLOWER_KEY_IP_TOS_MASK	/* u8 */
	TCA_FLOWER_KEY_IP_TTL		/* u8 */
	TCA_FLOWER_KEY_IP_TTL_MASK	/* u8 */

	TCA_FLOWER_MAX = TCA_FLOWER_KEY_IP_TTL_MASK
)
