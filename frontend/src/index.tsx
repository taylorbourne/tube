import './index.css'

import { ChakraProvider, extendTheme, ThemeConfig } from '@chakra-ui/react'
import { StepsStyleConfig as Steps } from 'chakra-ui-steps'
import React from 'react'
import ReactDOM from 'react-dom'

import App from './App'
import reportWebVitals from './reportWebVitals'

const config: ThemeConfig = {
  initialColorMode: "dark",
  useSystemColorMode: true,
};

const theme = extendTheme({
  components: {
    Steps,
  },
  colors: {
    teal: {
      "50": "#E5FDFF",
      "100": "#B8F8FF",
      "200": "#8AF4FF",
      "300": "#5CEFFF",
      "400": "#2EEBFF",
      "500": "#00E6FF",
      "600": "#00B8CC",
      "700": "#008A99",
      "800": "#005C66",
      "900": "#002E33",
    },
  },
  config,
});

ReactDOM.render(
  <React.StrictMode>
    <ChakraProvider theme={theme}>
      <App />
    </ChakraProvider>
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
