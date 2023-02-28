package web

import (
	"fmt"
	"runtime"
)

type FileInfo struct {
	Name  string `json:"name"`
	Pool  string `json:"pool"`
	Size  int    `json:"size"`
	Parts int    `json:"parts"`
}

type mimetypes struct {
	Videos map[string]string
	Audios map[string]string
}

var HttpMimes = &mimetypes{
	Videos: map[string]string{},
	Audios: map[string]string{},
}

func (*mimetypes) Populate() {
	HttpMimes.Videos["1d-interleaved-parityfec"] = "video/1d-interleaved-parityfec"
	HttpMimes.Videos["3gpp"] = "video/3gpp"
	HttpMimes.Videos["3gpp2"] = "video/3gpp2"
	HttpMimes.Videos["3gpp-tt"] = "video/3gpp-tt"
	HttpMimes.Videos["av1"] = "video/av1"
	HttpMimes.Videos["bmpeg"] = "video/bmpeg"
	HttpMimes.Videos["bt656"] = "video/bt656"
	HttpMimes.Videos["celb"] = "video/celb"
	HttpMimes.Videos["dv"] = "video/dv"
	HttpMimes.Videos["encaprtp"] = "video/encaprtp"
	HttpMimes.Videos["example"] = "video/example"
	HttpMimes.Videos["ffv1"] = "video/ffv1"
	HttpMimes.Videos["flexfec"] = "video/flexfec"
	HttpMimes.Videos["h261"] = "video/h261"
	HttpMimes.Videos["h263"] = "video/h263"
	HttpMimes.Videos["h263-1998"] = "video/h263-1998"
	HttpMimes.Videos["h263-2000"] = "video/h263-2000"
	HttpMimes.Videos["h264"] = "video/h264"
	HttpMimes.Videos["h264-rcdo"] = "video/h264-rcdo"
	HttpMimes.Videos["h264-svc"] = "video/h264-svc"
	HttpMimes.Videos["h265"] = "video/h265"
	HttpMimes.Videos["h266"] = "video/h266"
	HttpMimes.Videos["iso.segment"] = "video/iso.segment"
	HttpMimes.Videos["jxsv"] = "video/jxsv"
	HttpMimes.Videos["mj2"] = "video/mj2"
	HttpMimes.Videos["mp1s"] = "video/mp1s"
	HttpMimes.Videos["mp2p"] = "video/mp2p"
	HttpMimes.Videos["mp2t"] = "video/mp2t"
	HttpMimes.Videos["mp4"] = "video/mp4"
	HttpMimes.Videos["mp4v-es"] = "video/mp4v-es"
	HttpMimes.Videos["mpv"] = "video/mpv"
	HttpMimes.Videos["mpeg"] = "video/mpeg"
	HttpMimes.Videos["mpeg4-generic"] = "video/mpeg4-generic"
	HttpMimes.Videos["mkv"] = "video/x-matroska"
	HttpMimes.Videos["nv"] = "video/nv"
	HttpMimes.Videos["parityfec"] = "video/parityfec"
	HttpMimes.Videos["pointer"] = "video/pointer"
	HttpMimes.Videos["quicktime"] = "video/quicktime"
	HttpMimes.Videos["raptorfec"] = "video/raptorfec"
	HttpMimes.Videos["raw"] = "video/raw"
	HttpMimes.Videos["rtp-enc-aescm128"] = "video/rtp-enc-aescm128"
	HttpMimes.Videos["rtploopback"] = "video/rtploopback"
	HttpMimes.Videos["rtx"] = "video/rtx"
	HttpMimes.Videos["scip"] = "video/scip"
	HttpMimes.Videos["smpte291"] = "video/smpte291"
	HttpMimes.Videos["smpte292m"] = "video/smpte292m"
	HttpMimes.Videos["ulpfec"] = "video/ulpfec"
	HttpMimes.Videos["vc1"] = "video/vc1"
	HttpMimes.Videos["vc2"] = "video/vc2"
	HttpMimes.Videos["vnd.cctv"] = "video/vnd.cctv"
	HttpMimes.Videos["vnd.dece.hd"] = "video/vnd.dece.hd"
	HttpMimes.Videos["vnd.dece.mobile"] = "video/vnd.dece.mobile"
	HttpMimes.Videos["vnd.dece.mp4"] = "video/vnd.dece.mp4"
	HttpMimes.Videos["vnd.dece.pd"] = "video/vnd.dece.pd"
	HttpMimes.Videos["vnd.dece.sd"] = "video/vnd.dece.sd"
	HttpMimes.Videos["vnd.dece.video"] = "video/vnd.dece.video"
	HttpMimes.Videos["vnd.directv.mpeg"] = "video/vnd.directv.mpeg"
	HttpMimes.Videos["vnd.directv.mpeg-tts"] = "video/vnd.directv.mpeg-tts"
	HttpMimes.Videos["vnd.dlna.mpeg-tts"] = "video/vnd.dlna.mpeg-tts"
	HttpMimes.Videos["vnd.dvb.file"] = "video/vnd.dvb.file"
	HttpMimes.Videos["vnd.fvt"] = "video/vnd.fvt"
	HttpMimes.Videos["vnd.hns.video"] = "video/vnd.hns.video"
	HttpMimes.Videos["vnd.iptvforum.1dparityfec-1010"] = "video/vnd.iptvforum.1dparityfec-1010"
	HttpMimes.Videos["vnd.iptvforum.1dparityfec-2005"] = "video/vnd.iptvforum.1dparityfec-2005"
	HttpMimes.Videos["vnd.iptvforum.2dparityfec-1010"] = "video/vnd.iptvforum.2dparityfec-1010"
	HttpMimes.Videos["vnd.iptvforum.2dparityfec-2005"] = "video/vnd.iptvforum.2dparityfec-2005"
	HttpMimes.Videos["vnd.iptvforum.ttsavc"] = "video/vnd.iptvforum.ttsavc"
	HttpMimes.Videos["vnd.iptvforum.ttsmpeg2"] = "video/vnd.iptvforum.ttsmpeg2"
	HttpMimes.Videos["vnd.motorola.video"] = "video/vnd.motorola.video"
	HttpMimes.Videos["vnd.motorola.videop"] = "video/vnd.motorola.videop"
	HttpMimes.Videos["vnd.mpegurl"] = "video/vnd.mpegurl"
	HttpMimes.Videos["vnd.ms-playready.media.pyv"] = "video/vnd.ms-playready.media.pyv"
	HttpMimes.Videos["vnd.nokia.interleaved-multimedia"] = "video/vnd.nokia.interleaved-multimedia"
	HttpMimes.Videos["vnd.nokia.mp4vr"] = "video/vnd.nokia.mp4vr"
	HttpMimes.Videos["vnd.nokia.videovoip"] = "video/vnd.nokia.videovoip"
	HttpMimes.Videos["vnd.objectvideo"] = "video/vnd.objectvideo"
	HttpMimes.Videos["vnd.radgamettools.bink"] = "video/vnd.radgamettools.bink"
	HttpMimes.Videos["vnd.radgamettools.smacker"] = "video/vnd.radgamettools.smacker"
	HttpMimes.Videos["vnd.sealed.mpeg1"] = "video/vnd.sealed.mpeg1"
	HttpMimes.Videos["vnd.sealed.mpeg4"] = "video/vnd.sealed.mpeg4"
	HttpMimes.Videos["vnd.sealed.swf"] = "video/vnd.sealed.swf"
	HttpMimes.Videos["vnd.sealedmedia.softseal.mov"] = "video/vnd.sealedmedia.softseal.mov"
	HttpMimes.Videos["vnd.uvvu.mp4"] = "video/vnd.uvvu.mp4"
	HttpMimes.Videos["vnd.youtube.yt"] = "video/vnd.youtube.yt"
	HttpMimes.Videos["vnd.vivo"] = "video/vnd.vivo"
	HttpMimes.Videos["vp8"] = "video/vp8"
	HttpMimes.Videos["vp9"] = "video/vp9"
	HttpMimes.Videos["webm"] = "video/webm"
	HttpMimes.Audios["wav"] = "audio/x-wav"
	HttpMimes.Audios["aifc"] = "audio/x-aifc"
	HttpMimes.Audios["aiff"] = "audio/x-aiff"
	HttpMimes.Audios["mp3"] = "audio/mpeg"
	HttpMimes.Audios["ogg"] = "application/ogg"
	HttpMimes.Audios["m4a"] = "audio/mp4"
	HttpMimes.Audios["mp2"] = "audio/mpeg"

}

func (*mimetypes) Lookup(mt string) (string, bool) {
	mvi, vok := HttpMimes.Videos[mt]
	if vok {
		return mvi, true
	}
	mau, mak := HttpMimes.Audios[mt]
	if mak {
		return mau, true
	}
	return "application/octet-stream", false
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
