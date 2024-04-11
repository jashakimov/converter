package api

type Task struct {
	Path    string `json:"path" validate:"required"`
	ID      string `json:"id" validate:"required"`
	StartMs int    `json:"startMs"`
	EndMs   int    `json:"endMs"`
}

type Preset string

const (
	AVC720    Preset = "AVC 720 HD (1280x720)"
	AVC1080   Preset = "AVC 1080 FHD (1920x1080)"
	AVC240    Preset = "AVC 240 LD (426x240)"
	AVC480    Preset = "AVC 480 SD (854x480)"
	AVC360    Preset = "AVC 360 SD (640x360)"
	AVC144    Preset = "AVC 144 LD (256x144)"
	AVC4K     Preset = "AVC 4K UHD (3840x2160)"
	HEVC4K    Preset = "HEVC 4K UHD (3840x2160)"
	HEVC720   Preset = "HEVC 720 HD (1280x720)"
	HEVC480   Preset = "HEVC 480 SD (854x480)"
	HEVC360   Preset = "HEVC 360 SD (640x360)"
	HEVC144   Preset = "HEVC 144 LD (256x144)"
	HEVC240   Preset = "HEVC 240 LD (426x240)"
	HEVC1080  Preset = "HEVC 1080 FHD (1920x1080)"
	Mpeg2720  Preset = "Mpeg2 720 HD (1280x720)"
	Mpeg21080 Preset = "Mpeg2 1080 FHD (1920x1080)"
	Mpeg2480  Preset = "Mpeg2 480 SD (854x480)"
	Mpeg2360  Preset = "Mpeg2 360 SD (640x360)"
	Mpeg2144  Preset = "Mpeg2 144 LD (256x144)"
	Mpeg2240  Preset = "Mpeg2 240 LD (426x240)"
	AAC       Preset = "AAC (1 ch)"
	DASH      Preset = "DASH"
	FR_MP4    Preset = "File Render (mp4 mux)"
	FR_TS     Preset = "File Render (ts mux) HLS"
)

var AllowedPresets = []Preset{
	AVC720,
	AVC1080,
	AVC240,
	AVC480,
	AVC360,
	AVC144,
	AVC4K,
	HEVC4K,
	HEVC720,
	HEVC480,
	HEVC360,
	HEVC144,
	HEVC240,
	HEVC1080,
	Mpeg2720,
	Mpeg21080,
	Mpeg2480,
	Mpeg2360,
	Mpeg2144,
	Mpeg2240,
	AAC,
	DASH,
	FR_MP4,
	FR_TS,
}

type ChanResult struct {
	Error   error
	Message string
}
