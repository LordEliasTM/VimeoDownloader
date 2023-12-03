package main

type Master struct {
	ClipId  string  `json:"clip_id"`
	BaseUrl string  `json:"base_url"`
	Video   []Video `json:"video"`
	Audio   []Audio `json:"audio"`
}

type Media struct {
	Id                 string    `json:"id"`
	BaseUrl            string    `json:"base_url"`
	Format             string    `json:"format"`
	MimeType           string    `json:"mime_type"`
	Codecs             string    `json:"codecs"`
	Bitrate            int       `json:"bitrate"`
	AvgBitrate         int       `json:"avg_bitrate"`
	Duration           float64   `json:"duration"`
	MaxSegmentDuration int       `json:"max_segment_duration"`
	InitSegment        string    `json:"init_segment"`
	IndexSegment       string    `json:"index_segment"`
	Segments           []Segment `json:"segments"`
}

type Video struct {
	Media
	Framerate float64 `json:"framerate"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
}

type Audio struct {
	Media
	Channels   int `json:"channels"`
	SampleRate int `json:"sample_rate"`
}

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Url   string  `json:"url"`
	Size  int     `json:"size"`
}
