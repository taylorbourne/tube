package src

import "tube/src/internal/imgcache"

// SystemStruct : Beinhaltet alle Systeminformationen
type SystemStruct struct {
	Addresses struct {
		DVR string
		M3U string
		XML string
	}

	APIVersion             string
	AppName                string
	ARCH                   string
	BackgroundProcess      bool
	Branch                 string
	Build                  string
	Compatibility          string
	ConfigurationWizard    bool
	DBVersion              string
	Dev                    bool
	DeviceID               string
	Domain                 string
	PlexChannelLimit       int
	UnfilteredChannelLimit int

	FFmpeg struct {
		DefaultOptions string
		Path           string
	}

	VLC struct {
		DefaultOptions string
		Path           string
	}

	File struct {
		M3U      string
		PMS      string
		Settings string
		URLS     string
		XEPG     string
		XML      string
	}

	Compressed struct {
		GZxml string
	}

	Flag struct {
		Branch   string
		Debug    int
		Info     bool
		Port     string
		Restore  string
		SSDP     bool
		Mode     string
		Proxy    string
		Dir      string
		Embed    string
		Packaged string
	}

	Folder struct {
		Backup       string
		Cache        string
		Config       string
		Data         string
		ImagesCache  string
		ImagesUpload string
		Temp         string
	}

	Hostname               string
	ImageCachingInProgress int
	IPAddress              string
	IPAddressesList        []string
	IPAddressesV4          []string
	IPAddressesV6          []string
	Name                   string
	OS                     string
	ScanInProgress         int

	Notification map[string]Notification

	ServerProtocol struct {
		API string
		DVR string
		M3U string
		WEB string
		XML string
	}

	Update struct {
		Git  string
		Name string
	}

	URLBase string
	UDPxy   string
	Version string
	WEB     struct {
		Menu []string
	}
}

// DataStruct : Alle Daten werden hier abgelegt. (Lineup, XMLTV)
type DataStruct struct {
	Cache struct {
		Images      *imgcache.Cache
		ImagesCache []string
		ImagesFiles []string
		ImagesURLS  []string
		PMS         map[string]string

		StreamingURLS map[string]StreamInfo
		XMLTV         map[string]XMLTV

		Streams struct {
			Active []string
		}
	}

	Filter []Filter

	Playlist struct {
		M3U struct {
			Groups struct {
				Text  []string
				Value []string
			}
		}
	}

	StreamPreviewUI struct {
		Active   []string
		Inactive []string
	}

	Streams struct {
		Active   []interface{}
		All      []interface{}
		Inactive []interface{}
	}

	XMLTV struct {
		Files   []string
		Mapping map[string]interface{}
	}

	XEPG struct {
		Channels  map[string]interface{}
		XEPGCount int64
	}
}

// Filter : Wird f??r die Filterregeln verwendet
type Filter struct {
	CaseSensitive bool
	Rule          string
	Type          string
}

// XEPGChannelStruct : XEPG Struktur
type XEPGChannelStruct struct {
	FileM3UID          string `json:"_file.m3u.id,required"`
	FileM3UName        string `json:"_file.m3u.name,required"`
	FileM3UPath        string `json:"_file.m3u.path,required"`
	GroupTitle         string `json:"group-title,required"`
	Name               string `json:"name,required"`
	TvgID              string `json:"tvg-id,required"`
	TvgLogo            string `json:"tvg-logo,required"`
	TvgName            string `json:"tvg-name,required"`
	URL                string `json:"url,required"`
	UUIDKey            string `json:"_uuid.key,required"`
	UUIDValue          string `json:"_uuid.value,omitempty"`
	Values             string `json:"_values,required"`
	XActive            bool   `json:"x-active,required"`
	XCategory          string `json:"x-category,required"`
	XChannelID         string `json:"x-channelID,required"`
	XEPG               string `json:"x-epg,required"`
	XGroupTitle        string `json:"x-group-title,required"`
	XMapping           string `json:"x-mapping,required"`
	XmltvFile          string `json:"x-xmltv-file,required"`
	XName              string `json:"x-name,required"`
	XUpdateChannelIcon bool   `json:"x-update-channel-icon,required"`
	XUpdateChannelName bool   `json:"x-update-channel-name,required"`
	XDescription       string `json:"x-description,required"`
}

// M3UChannelStructXEPG : M3U Struktur f??r XEPG
type M3UChannelStructXEPG struct {
	FileM3UID   string `json:"_file.m3u.id,required"`
	FileM3UName string `json:"_file.m3u.name,required"`
	FileM3UPath string `json:"_file.m3u.path,required"`
	GroupTitle  string `json:"group-title,required"`
	Name        string `json:"name,required"`
	TvgID       string `json:"tvg-id,required"`
	TvgLogo     string `json:"tvg-logo,required"`
	TvgName     string `json:"tvg-name,required"`
	URL         string `json:"url,required"`
	UUIDKey     string `json:"_uuid.key,required"`
	UUIDValue   string `json:"_uuid.value,required"`
	Values      string `json:"_values,required"`
}

// FilterStruct : Filter Struktur
type FilterStruct struct {
	Active        bool   `json:"active,required"`
	CaseSensitive bool   `json:"caseSensitive,required"`
	Description   string `json:"description,required"`
	Exclude       string `json:"exclude,required"`
	Filter        string `json:"filter,required"`
	Include       string `json:"include,required"`
	Name          string `json:"name,required"`
	Rule          string `json:"rule,omitempty"`
	Type          string `json:"type,required"`
}

// StreamingURLS : Informationen zu allen streaming URL's
type StreamingURLS struct {
	Streams map[string]StreamInfo `json:"channels,required"`
}

// StreamInfo : Informationen zum Kanal f??r die streaming URL
type StreamInfo struct {
	ChannelNumber string `json:"channelNumber,required"`
	Name          string `json:"name,required"`
	PlaylistID    string `json:"playlistID,required"`
	URL           string `json:"url,required"`
	URLid         string `json:"urlID,required"`
}

// Notification : Notifikationen im Webinterface
type Notification struct {
	Headline string `json:"headline,required"`
	Message  string `json:"message,required"`
	New      bool   `json:"new,required"`
	Time     string `json:"time,required"`
	Type     string `json:"type,required"`
}

// SettingsStruct : Inhalt der settings.json
type SettingsStruct struct {
	API           bool     `json:"api"`
	BackupKeep    int      `json:"backupKeep"`
	BackupPath    string   `json:"backupPath"`
	Buffer        string   `json:"buffer"`
	BufferSize    int      `json:"bufferSizeKb"`
	BufferTimeout float64  `json:"bufferTimeout"`
	CacheImages   bool     `json:"cacheImages"`
	EpgSource     string   `json:"epgSource"`
	FFmpegOptions string   `json:"ffmpegOptions"`
	FFmpegPath    string   `json:"ffmpegPath"`
	VLCOptions    string   `json:"vlcOptions"`
	VLCPath       string   `json:"vlcPath"`
	FileM3U       []string `json:"file,omitempty"`  // Beim Wizard wird die M3U in ein Slice gespeichert
	FileXMLTV     []string `json:"xmltv,omitempty"` // Altes Speichersystem der Provider XML Datei Slice (Wird f??r die Umwandlung auf das neue ben??tigt)

	Files struct {
		HDHR  map[string]interface{} `json:"hdhr"`
		M3U   map[string]interface{} `json:"m3u"`
		XMLTV map[string]interface{} `json:"xmltv"`
	} `json:"files"`

	FilesUpdate               bool                  `json:"filesUpdate"`
	Filter                    map[int64]interface{} `json:"filter"`
	Key                       string                `json:"key,omitempty"`
	Language                  string                `json:"language"`
	LogEntriesRAM             int                   `json:"logEntriesRam"`
	M3U8AdaptiveBandwidthMBPS int                   `json:"m3u8BandwidthMbps"`
	MappingFirstChannel       float64               `json:"mappingFirstChannel"`
	Port                      string                `json:"port"`
	SSDP                      bool                  `json:"ssdp"`
	TempPath                  string                `json:"tempPath"`
	Tuner                     int                   `json:"tuner"`
	Update                    []string              `json:"update"`
	UpdateURL                 string                `json:"updateUrl,omitempty"`
	UserAgent                 string                `json:"userAgent"`
	UUID                      string                `json:"uuid"`
	UDPxy                     string                `json:"udpxy"`
	Version                   string                `json:"version"`
	XepgReplaceMissingImages  bool                  `json:"xepgReplaceMissingImages"`
}

// LanguageUI : Sprache f??r das WebUI
type LanguageUI struct {
	Login struct {
		Failed string
	}
}
