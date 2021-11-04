import React, { createContext, Dispatch } from 'react'

import type { Actions } from './reducer'
import type { StateType } from 'types'

export const initialState = {
    serverConfig: {"clientInfo": {
        "arch": "",
        "branch": "",
        "DVR": "",
        "epgSource": "",
        "errors": 0,
        "m3u-url": "",
        "os": "",
        "streams": "0 / 0",
        "uuid": "",
        "version": "",
        "warnings": 0,
        "xepg": 0,
        "xepg-url": ""
    },
    "data": {
        "playlist": {
            "m3u": {
                "groups": {
                    "text": [],
                    "value": []
                }
            }
        },
        "StreamPreviewUI": {
            "activeStreams": [],
            "inactiveStreams": []
        }
    },
    "configurationWizard": false,
    "log": {
        "errors": 0,
        "log": [],
        "warnings": 0
    },
    "settings": {
        "api": false,
        "authentication.api": false,
        "authentication.m3u": false,
        "authentication.pms": false,
        "authentication.web": false,
        "authentication.xml": false,
        "backup.keep": 10,
        "backup.path": "",
        "buffer": "-",
        "buffer.size.kb": 1024,
        "buffer.timeout": 500,
        "cache.images": false,
        "epgSource": "",
        "ffmpeg.options": "",
        "ffmpeg.path": "",
        "vlc.options": "",
        "vlc.path": "",
        "files": {
            "hdhr": {},
            "m3u": {},
            "xmltv": {}
        },
        "files.update": true,
        "filter": {},
        "language": "en",
        "log.entries.ram": 500,
        "m3u8.adaptive.bandwidth.mbps": 10,
        "mapping.first.channel": 1000,
        "port": "34400",
        "ssdp": true,
        "temp.path": "/var/folders/34/f7s_0_s56qnfq76vb5s8rgmh0000gn/T/xteve/",
        "tuner": 1,
        "update": [
            "0000"
        ],
        "user.agent": "xTeVe",
        "uuid": "2019-01-DEV-xTeVe!",
        "udpxy": "",
        "version": "2.1.0",
        "xepg.replace.missing.images": true,
    },
    "status": true,
    "xepg": {
        "epgMapping": {},
        "xmltvMap": {}
    },
  },
    "notification": {},
}


const ctxt = createContext<{
  state: StateType;
  dispatch: Dispatch<Actions>;
}>({
  state: initialState,
  dispatch: () => null
});

export const AppContextProvider = ctxt.Provider;
  
export const AppContextConsumer = ctxt.Consumer;