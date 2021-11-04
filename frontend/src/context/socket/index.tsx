import React, { useReducer } from 'react'

import reducer from './reducer'
import { AppContextProvider, initialState } from './Socket'

const AppProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, dispatch] = useReducer(reducer, initialState);

  return (
    <AppContextProvider value={{ state, dispatch }}>
      {children}
    </AppContextProvider>
  );
};

export default AppProvider;
