package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Index : Web Server /
func Index(w http.ResponseWriter, r *http.Request) {

	var err error
	var response []byte
	var path = r.URL.Path
	var debug string

	setGlobalDomain(r.Host)

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

// Web : Web Server /web/
func Web(w http.ResponseWriter, r *http.Request) {

	var lang = make(map[string]interface{})
	var err error

	var requestFile = strings.Replace(r.URL.Path, "/web", "html", -1)
	var content, contentType string

	var language LanguageUI

	setGlobalDomain(r.Host)

	// if System.Dev {

	// 	lang, err = loadJSONFileToMap(fmt.Sprintf("html/lang/%s.json", Settings.Language))
	// 	if err != nil {
	// 		ShowError(err, 000)
	// 	}

	// } else {

	// 	var languageFile = "html/lang/en.json"

	// 	if _, ok := webUI[languageFile].(string); ok {
	// 		content = ""
	// 		lang = jsonToMap(content)
	// 	}

	// }

	err = json.Unmarshal([]byte(mapToJSON(lang)), &language)
	if err != nil {
		ShowError(err, 000)
		return
	}

	if getFilenameFromPath(requestFile) == "html" {

		if System.ScanInProgress == 0 {

			if len(Settings.Files.M3U) == 0 && len(Settings.Files.HDHR) == 0 {
				System.ConfigurationWizard = true
			}

		}

		// switch System.ConfigurationWizard {

		// case true:
		// 	file = requestFile + "configuration.html"
		// 	Settings.AuthenticationWEB = false

		// case false:
		// 	file = requestFile + "index.html"

		// }

		// if System.ScanInProgress == 1 {
		// 	file = requestFile + "maintenance.html"
		// }

		// 	requestFile = requestFile + "index.html"

		// 	if _, ok := webUI[requestFile]; ok {

		// 		if contentType == "text/plain" {
		// 			w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(requestFile))
		// 		}

		// 	} else {

		// 		httpStatusError(w, r, 404)
		// 		return
		// 	}

	}

	// if _, ok := webUI[requestFile].(string); ok {

	// 	content = ""
	// 	contentType = getContentType(requestFile)

	// 	if contentType == "text/plain" {
	// 		w.Header().Set("Content-Disposition", "attachment; filename="+getFilenameFromPath(requestFile))
	// 	}

	// } else {
	// 	httpStatusError(w, r, 404)
	// 	return
	// }

	contentType = getContentType(requestFile)

	if System.Dev {
		// Lokale Webserver Dateien werden geladen, nur für die Entwicklung
		content, _ = readStringFromFile(requestFile)
	}

	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(200)

	// if contentType == "text/html" || contentType == "application/javascript" {
	// 	content = parseTemplate(content, lang)
	// }

	w.Write([]byte(content))
}

// API : API request /api/
func API(w http.ResponseWriter, r *http.Request) {

	/*
			API Bedingungen (ohne Authentifizierung):
			- API muss in den Einstellungen aktiviert sein

			Beispiel API Request mit curl
			Status:
			curl -X POST -H "Content-Type: application/json" -d '{"cmd":"status"}' http://localhost:34400/api/

			- - - - -

			API Bedingungen (mit Authentifizierung):
			- API muss in den Einstellungen aktiviert sein
			- API muss bei den Authentifizierungseinstellungen aktiviert sein
			- Benutzer muss die Berechtigung API haben

			Nach jeder API Anfrage wird ein Token generiert, dieser ist einmal in 60 Minuten gültig.
			In jeder Antwort ist ein neuer Token enthalten

			Beispiel API Request mit curl
			Login:
			curl -X POST -H "Content-Type: application/json" -d '{"cmd":"login","username":"plex","password":"123"}' http://localhost:34400/api/

			Antwort:
			{
		  	"status": true,
		  	"token": "U0T-NTSaigh-RlbkqERsHvUpgvaaY2dyRGuwIIvv"
			}

			Status mit Verwendung eines Tokens:
			curl -X POST -H "Content-Type: application/json" -d '{"cmd":"status","token":"U0T-NTSaigh-RlbkqERsHvUpgvaaY2dyRGuwIIvv"}' http://localhost:4400/api/

			Antwort:
			{
			  "epg.source": "XEPG",
			  "status": true,
			  "streams.active": 7,
			  "streams.all": 63,
			  "streams.xepg": 2,
			  "token": "mXiG1NE1MrTXDtyh7PxRHK5z8iPI_LzxsQmY-LFn",
			  "url.dvr": "localhost:34400",
			  "url.m3u": "http://localhost:34400/m3u/xteve.m3u",
			  "url.xepg": "http://localhost:34400/xmltv/xteve.xml",
			  "version.api": "1.1.0",
			  "version.xteve": "1.3.0"
			}
	*/

	setGlobalDomain(r.Host)
	var request APIRequestStruct
	var response APIResponseStruct

	var responseAPIError = func(err error) {

		var response APIResponseStruct

		response.Status = false
		response.Error = err.Error()
		w.Write([]byte(mapToJSON(response)))
	}

	response.Status = true

	if !Settings.API {
		httpStatusError(w, r, 423)
		return
	}

	if r.Method == "GET" {
		httpStatusError(w, r, 404)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		httpStatusError(w, r, 400)
		return

	}

	err = json.Unmarshal(b, &request)
	if err != nil {
		httpStatusError(w, r, 400)
		return
	}

	w.Header().Set("content-type", "application/json")

	switch request.Cmd {
	case "login": // Muss nichts übergeben werden

	case "status":

		response.VersionXteve = System.Version
		response.VersionAPI = System.APIVersion
		response.StreamsActive = int64(len(Data.Streams.Active))
		response.StreamsAll = int64(len(Data.Streams.All))
		response.StreamsXepg = int64(Data.XEPG.XEPGCount)
		response.EpgSource = Settings.EpgSource
		response.URLDvr = System.Domain
		response.URLM3U = System.ServerProtocol.M3U + "://" + System.Domain + "/m3u/xteve.m3u"
		response.URLXepg = System.ServerProtocol.XML + "://" + System.Domain + "/xmltv/xteve.xml"

	case "update.m3u":
		err = getProviderData("m3u", "")
		if err != nil {
			break
		}

		err = buildDatabaseDVR()
		if err != nil {
			break
		}

	case "update.hdhr":

		err = getProviderData("hdhr", "")
		if err != nil {
			break
		}

		err = buildDatabaseDVR()
		if err != nil {
			break
		}

	case "update.xmltv":
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
		responseAPIError(err)
	}

	w.Write([]byte(mapToJSON(response)))
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

func PlaylistFilter(w http.ResponseWriter, r *http.Request) {
	var request map[int64]interface{}

	_ = json.NewDecoder(r.Body).Decode(&request)

	response, err := saveFilter(request)

	if err == nil {
		httpStatusError(w, r, 500)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
