export type ClientInfoType = {
  arch: string;
  branch: string;
  DVR: string;
  epgSource: string;
  errors: number;
  "m3u-url": string;
  os: string;
  streams: string;
  uuid: string;
  version: string;
  warnings: number;
  xepg: number;
  "xepg-url": string;
}

export type DataType = {
  playlist: {
    m3u: {
      groups: {
        text: Array<string>;
        value: Array<string>;
      }
    }
  },
  StreamPreviewUI: {
    activeStreams: Array<any>,
    inactiveStreams: Array<any>,
  }
}

export type LogType = {
  errors: number;
  log: Array<string>;
  warnings: number;
}

export type Notification = {
  "headline": string;
  "message": string;
  "new": boolean;
  "time": Date;
  "type": "info" | "warning" | "error";
}

export type SettingsType = {
  api: boolean;
  'authentication.api': boolean;
  'authentication.m3u': boolean;
  'authentication.pms': boolean;
  'authentication.web': boolean;
  'authentication.xml': boolean;
  'backup.keep': number,
  'backup.path': string;
  'git.branch': string;
  buffer: string;
  'buffer.size.kb': number;
  'buffer.timeout': number;
  'cache.images': boolean;
  epgSource: string;
  'ffmpeg.options'?: string;
  'ffmpeg.path'?: string;
  'vlc.options': string;
  'vlc.path': string;
  files: {
      hdhr: Object;
      m3u: Object;
      xmltv: Object;
  },
  'files.update': boolean;
  filter: Object;
  language: string;
  'log.entries.ram': number;
  'm3u8.adaptive.bandwidth.mbps': number;
  'mapping.first.channel': number;
  port: string;
  ssdp: boolean;
  'temp.path': string;
  tuner: number;
  update: Array<string>;
  'user.agent': string;
  uuid: string;
  udpxy: string;
  version: string;
  'xepg.replace.missing.images': boolean;
  xteveAutoUpdate: boolean;}

export type ServerConfigType = {
  "clientInfo": ClientInfoType,
  "data": DataType,
  "configurationWizard": boolean,
  "log": LogType,
  "settings": SettingsType
  "status": boolean;
  "xepg": {
      "epgMapping": Object;
      "xmltvMap": Object;
  },
}

export interface StateType {
  serverConfig: ServerConfigType;
  "notification": {
    [key: string]: Notification;
  }
}