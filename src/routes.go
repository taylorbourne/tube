package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// Index : Web Server /
func Index(w http.ResponseWriter, r *http.Request) {

	var err error
	var response []byte
	var path = r.URL.Path
	var debug string

	debug = fmt.Sprintf("Web Server Request:Path: %s", path)
	showDebug(debug, 2)

	switch path {

	case "/discover.json":
		response, err = getDiscover()
		w.Header().Set("Content-Type", "application/json")

	case "/lineup_status.json":
		response, err = getLineupStatus()
		w.Header().Set("Content-Type", "application/json")

	case "/lineup.json":
		response, err = getLineup()
		w.Header().Set("Content-Type", "application/json")

	case "/device.xml", "/capability":
		response, err = getCapability()
		w.Header().Set("Content-Type", "application/xml")

	default:
		response, err = getCapability()
		w.Header().Set("Content-Type", "application/xml")
	}

	if err == nil {

		w.WriteHeader(200)
		w.Write(response)
		return

	}

	httpStatusError(w, r, 500)
}

// Stream : Web Server /stream/
func Stream(w http.ResponseWriter, r *http.Request) {

	var path = strings.Replace(r.RequestURI, "/stream/", "", 1)
	//var stream = strings.SplitN(path, "-", 2)

	streamInfo, err := getStreamInfo(path)
	if err != nil {
		ShowError(err, 1203)
		httpStatusError(w, r, 404)
		return
	}

	// If an UDPxy host is set, and the stream URL is multicast (i.e. starts with 'udp://@'),
	// then streamInfo.URL needs to be rewritten to point to UDPxy.
	if Settings.UDPxy != "" && strings.HasPrefix(streamInfo.URL, "udp://@") {
		streamInfo.URL = fmt.Sprintf("http://%s/udp/%s/", Settings.UDPxy, strings.TrimPrefix(streamInfo.URL, "udp://@"))
	}

	switch Settings.Buffer {

	case "-":
		showInfo(fmt.Sprintf("Buffer:false [%s]", Settings.Buffer))

	case "xteve":
		if strings.Contains(streamInfo.URL, "rtsp://") || strings.Contains(streamInfo.URL, "rtp://") {
			err = errors.New("RTSP and RTP streams are not supported")
			ShowError(err, 2004)

			showInfo("Streaming URL:" + streamInfo.URL)
			http.Redirect(w, r, streamInfo.URL, http.StatusFound)

			showInfo("Streaming Info:URL was passed to the client")
			return
		}

		showInfo(fmt.Sprintf("Buffer:true [%s]", Settings.Buffer))

	default:
		showInfo(fmt.Sprintf("Buffer:true [%s]", Settings.Buffer))

	}

	if Settings.Buffer != "-" {
		showInfo(fmt.Sprintf("Buffer Size:%d KB", Settings.BufferSize))
	}

	showInfo(fmt.Sprintf("Channel Name:%s", streamInfo.Name))
	showInfo(fmt.Sprintf("Client User-Agent:%s", r.Header.Get("User-Agent")))

	// Prüfen ob der Buffer verwendet werden soll
	switch Settings.Buffer {

	case "-":
		showInfo("Streaming URL:" + streamInfo.URL)
		http.Redirect(w, r, streamInfo.URL, http.StatusFound)

		showInfo("Streaming Info:URL was passed to the client.")
		showInfo("Streaming Info:xTeVe is no longer involved, the client connects directly to the streaming server.")

	default:
		bufferingStream(streamInfo.PlaylistID, streamInfo.URL, streamInfo.Name, w, r)

	}

}

// Auto : HDHR routing (wird derzeit nicht benutzt)
func Auto(w http.ResponseWriter, r *http.Request) {

	var channelID = strings.Replace(r.RequestURI, "/auto/v", "", 1)
	fmt.Println(channelID)

	/*
		switch Settings.Buffer {

		case true:
			var playlistID, streamURL, err = getStreamByChannelID(channelID)
			if err == nil {
				bufferingStream(playlistID, streamURL, w, r)
			} else {
				httpStatusError(w, r, 404)
			}

		case false:
			httpStatusError(w, r, 423)
		}
	*/
}

// xTeVe : Web Server /xmltv/ und /m3u/
func Tube(w http.ResponseWriter, r *http.Request) {

	var groupTitle, file, content, contentType string
	var err error
	var path = strings.TrimPrefix(r.URL.Path, "/")
	var groups = []string{}

	setGlobalDomain(r.Host)

	// XMLTV Datei
	if strings.Contains(path, "xmltv/") {

		file = System.Folder.Data + getFilenameFromPath(path)

		content, err = readStringFromFile(file)
		if err != nil {
			httpStatusError(w, r, 404)
			return
		}

	}

	// M3U Datei
	if strings.Contains(path, "m3u/") {

		groupTitle = r.URL.Query().Get("group-title")

		if !System.Dev {
			// false: Dateiname wird im Header gesetzt
			// true: M3U wird direkt im Browser angezeigt
			w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(path))
		}

		if len(groupTitle) > 0 {
			groups = strings.Split(groupTitle, ",")
		}

		content, err = buildM3U(groups)
		if err != nil {
			ShowError(err, 000)
		}

	}

	contentType = http.DetectContentType([]byte(content))
	if strings.Contains(strings.ToLower(contentType), "xml") {
		contentType = "application/xml; charset=utf-8"
	}

	w.Header().Set("Content-Type", contentType)

	if err == nil {
		w.Write([]byte(content))
	}
}

// Images : Image Cache /images/
func Images(w http.ResponseWriter, r *http.Request) {

	var path = strings.TrimPrefix(r.URL.Path, "/")
	var filePath = System.Folder.ImagesCache + getFilenameFromPath(path)

	content, err := readByteFromFile(filePath)
	if err != nil {
		httpStatusError(w, r, 404)
		return
	}

	w.Header().Add("Content-Type", getContentType(filePath))
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(content)))
	w.WriteHeader(200)
	w.Write(content)
}

// DataImages : Image Pfad für Logos / Bilder die hochgeladen wurden /data_images/
func DataImages(w http.ResponseWriter, r *http.Request) {

	var path = strings.TrimPrefix(r.URL.Path, "/")
	var filePath = System.Folder.ImagesUpload + getFilenameFromPath(path)

	content, err := readByteFromFile(filePath)
	if err != nil {
		httpStatusError(w, r, 404)
		return
	}

	w.Header().Add("Content-Type", getContentType(filePath))
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(content)))
	w.WriteHeader(200)
	w.Write(content)
}

// Download : Datei Download
func Download(w http.ResponseWriter, r *http.Request) {

	var path = r.URL.Path
	var file = System.Folder.Temp + getFilenameFromPath(path)
	w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(file))

	content, err := readStringFromFile(file)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	os.RemoveAll(System.Folder.Temp + getFilenameFromPath(path))
	w.Write([]byte(content))
}

func httpStatusError(w http.ResponseWriter, r *http.Request, httpStatusCode int) {
	http.Error(w, fmt.Sprintf("%s [%d]", http.StatusText(httpStatusCode), httpStatusCode), httpStatusCode)
}

func getContentType(filename string) (contentType string) {

	if strings.HasSuffix(filename, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(filename, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(filename, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(filename, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(filename, ".jpg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(filename, ".gif") {
		contentType = "image/gif"
	} else if strings.HasSuffix(filename, ".svg") {
		contentType = "image/svg+xml"
	} else if strings.HasSuffix(filename, ".mp4") {
		contentType = "video/mp4"
	} else if strings.HasSuffix(filename, ".webm") {
		contentType = "video/webm"
	} else if strings.HasSuffix(filename, ".ogg") {
		contentType = "video/ogg"
	} else if strings.HasSuffix(filename, ".mp3") {
		contentType = "audio/mp3"
	} else if strings.HasSuffix(filename, ".wav") {
		contentType = "audio/wav"
	} else {
		contentType = "text/plain"
	}

	return
}

func Backup(w http.ResponseWriter, r *http.Request) {
	_, err := xteveBackup()
	if err != nil {
		httpStatusError(w, r, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Config(w http.ResponseWriter, r *http.Request) {
	var request ConfigRequest

	_ = json.NewDecoder(r.Body).Decode(&request)

	err := saveWizard(request)

	if err == nil {
		System.ConfigurationWizard = false
		httpStatusError(w, r, 500)
	}
	w.WriteHeader(http.StatusOK)
}

func Restore(w http.ResponseWriter, r *http.Request) {
	var request RestoreRequest
	var response RestoreResponse

	_ = json.NewDecoder(r.Body).Decode(&request)

	WebScreenLog.Log = make([]string, 0)
	WebScreenLog.Errors = 0
	WebScreenLog.Warnings = 0

	if len(request.Base64) > 0 {

		newWebURL, err := xteveRestoreFromWeb(request.Base64)
		if err != nil {
			ShowError(err, 000)
			response.Alert = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}

		if err == nil {

			if len(newWebURL) > 0 {
				response.Alert = "Backup was successfully restored.\nThe port of the sTeVe URL has changed, you have to restart xTeVe.\nAfter a restart, xTeVe can be reached again at the following URL:\n" + newWebURL
			} else {
				response.Alert = "Backup was successfully restored."
			}
			showInfo("xTeVe:" + "Backup successfully restored.")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	}
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	response := GetStatusResponse{
		EpgSource: Settings.EpgSource,
		Endpoints: GetStatusResponseEndpoints{
			URLDvr:  System.Domain,
			URLM3U:  System.ServerProtocol.M3U + "://" + System.Domain + "/m3u/xteve.m3u",
			URLXepg: System.ServerProtocol.XML + "://" + System.Domain + "/xmltv/xteve.xml",
		},
		Streams: GetStatusResponseStreams{
			StreamsActive: int64(len(Data.Streams.Active)),
			StreamsAll:    int64(len(Data.Streams.All)),
			StreamsXepg:   int64(Data.XEPG.XEPGCount),
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func HDHRUpdate(w http.ResponseWriter, r *http.Request) {

	var request FileRequest

	_ = json.NewDecoder(r.Body).Decode(&request)

	err := updateFile(request, "hdhr")

	if err != nil {
		httpStatusError(w, r, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {

	if System.ScanInProgress == 0 {
		if len(Settings.Files.M3U) == 0 && len(Settings.Files.HDHR) == 0 {
			System.ConfigurationWizard = true
		}
	}

	response := InfoResponse{
		ARCH:           System.ARCH,
		EpgSource:      Settings.EpgSource,
		DVR:            System.Addresses.DVR,
		M3U:            System.Addresses.M3U,
		XML:            System.Addresses.XML,
		OS:             System.OS,
		Streams:        fmt.Sprintf("%d / %d", len(Data.Streams.Active), len(Data.Streams.All)),
		UUID:           Settings.UUID,
		Errors:         WebScreenLog.Errors,
		Warnings:       WebScreenLog.Warnings,
		Notification:   System.Notification,
		ScanInProgress: System.ScanInProgress,
		// @TODO – does this make sense somewhere else or on it's own request?
		ConfigurationWizard: System.ConfigurationWizard,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetXEPG(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var XEPG = make(map[string]interface{})

	if len(Data.Streams.Active) > 0 {
		XEPG["epgMapping"] = Data.XEPG.Channels
		XEPG["xmltvMap"] = Data.XMLTV.Mapping
	} else {
		XEPG["epgMapping"] = make(map[string]interface{})
		XEPG["xmltvMap"] = make(map[string]interface{})
	}

	response = XEPG

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Files(w http.ResponseWriter, r *http.Request) {
	fileType := mux.Vars(r)["type"]

	switch fileType {
	case "m3u":
	}

	var request FileRequest
	var response FileResponse
	if len(request.Base64) > 0 {
		LogoURL, err := uploadLogo(request.Base64, request.Filename)
		if err != nil {
			httpStatusError(w, r, 500)
			return
		}
		response.LogoURL = LogoURL
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}

func SaveXEPG(w http.ResponseWriter, r *http.Request) {
	var request XEPGRequest

	err := saveXEpgMapping(request)
	if err != nil {
		httpStatusError(w, r, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetLog(w http.ResponseWriter, r *http.Request) {
	response := WebScreenLog.Log
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteLog(w http.ResponseWriter, r *http.Request) {
	response := WebScreenLog.Log
	WebScreenLog.Log = make([]string, 0)
	WebScreenLog.Errors = 0
	WebScreenLog.Warnings = 0
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateFile(w http.ResponseWriter, r *http.Request) {
	var request FileRequest
	fileType := mux.Vars(r)["type"]
	err := updateFile(request, fileType)
	if err != nil {
		httpStatusError(w, r, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SaveFile(w http.ResponseWriter, r *http.Request) {
	var request FileRequest
	var response FileResponse
	fileType := mux.Vars(r)["type"]
	_ = json.NewDecoder(r.Body).Decode(&request)

	if fileType == "logo" {
		if len(request.Base64) > 0 {
			LogoURL, err := uploadLogo(request.Base64, request.Filename)
			if err != nil {
				httpStatusError(w, r, 500)
				return
			}
			response.LogoURL = LogoURL
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	err := saveFiles(request, fileType)

	if err != nil {
		httpStatusError(w, r, 500)
	}
	w.WriteHeader(http.StatusOK)
}

func GetPlaylistStreams(w http.ResponseWriter, r *http.Request) {
	response := StreamResponse{
		Active:   Data.StreamPreviewUI.Active,
		Inactive: Data.StreamPreviewUI.Inactive,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GetPlaylistInfo(w http.ResponseWriter, r *http.Request) {
	response := PlaylistReponse{
		M3U: M3U{
			Groups: Groups{
				Text:  Data.Playlist.M3U.Groups.Text,
				Value: Data.Playlist.M3U.Groups.Value,
			},
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	var response SettingsResponse = Settings
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SaveSettings(w http.ResponseWriter, r *http.Request) {
	var request SettingsRequest
	var response SettingsResponse = Settings
	_ = json.NewDecoder(r.Body).Decode(&request)
	settings, err := updateServerSettings(request)
	response = settings
	if err != nil {
		httpStatusError(w, r, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func PlaylistFilter(w http.ResponseWriter, r *http.Request) {
	var request map[int64]interface{}

	_ = json.NewDecoder(r.Body).Decode(&request)

	response, err := saveFilter(request)

	if err != nil {
		httpStatusError(w, r, 500)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var err error
	updateType := mux.Vars(r)["type"]

	switch updateType {
	case "m3u":
		err = getProviderData("m3u", "")
		if err != nil {
			break
		}
		err = buildDatabaseDVR()
		if err != nil {
			break
		}

	case "hdhr":
		err = getProviderData("hdhr", "")
		if err != nil {
			break
		}

		err = buildDatabaseDVR()
		if err != nil {
			break
		}

	case "xmltv":
		err = getProviderData("xmltv", "")
		if err != nil {
			break
		}

	case "update.xepg":
		buildXEPG(false)

	default:
		err = errors.New(getErrMsg(5000))
	}

	if err != nil {
		httpStatusError(w, r, 500)
		return
	}

	w.WriteHeader(http.StatusOK)

}
