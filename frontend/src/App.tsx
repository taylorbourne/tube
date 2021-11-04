import './App.css'

import React, { useCallback, useEffect, useState } from 'react'
import useWebSocket, { ReadyState } from 'react-use-websocket'

import SetupWizard from './components/Wizard'
import { GET_SYSTEM_DATA } from './data'

const WebSocketDemo = () => {
  var protocolWS;
  switch (window.location.protocol) {
    case "http:":
      protocolWS = "ws://";
      break;
    case "https:":
      protocolWS = "wss://";
      break;
  }

  //Public API that will echo messages sent to it back to the client

  const [serverConfig, setServerConfig] = useState([]);

  const { sendMessage, lastMessage, readyState } = useWebSocket(
    protocolWS + window.location.hostname + ":" + 34400 + "/data/"
  );

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  useEffect(() => {
    if (connectionStatus === "Open") {
      sendMessage(GET_SYSTEM_DATA);
    }
  }, [connectionStatus]);

  useEffect(() => {
    if (lastMessage !== null) {
      console.log("what am i", JSON.parse(lastMessage.data));
      setServerConfig(JSON.parse(lastMessage.data));
    }
  }, [lastMessage]);

  return <SetupWizard />;
};

export default WebSocketDemo;
