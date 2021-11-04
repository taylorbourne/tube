import { createContext, FC, useState } from 'react'

interface WizardContextProps {
  isXEPGEnabled: boolean;
  payload: undefined | Object;
  setIsXEPGEnabled:
    | React.Dispatch<React.SetStateAction<boolean>>
    | (() => null);
  setPayload: React.Dispatch<React.SetStateAction<any>> | (() => null);
}

export const WizardContext = createContext<WizardContextProps>({
  isXEPGEnabled: true,
  payload: undefined,
  setIsXEPGEnabled: () => null,
  setPayload: () => null,
});

const WizardProvider: FC = ({ children }) => {
  const [payload, setPayload] = useState(undefined);
  const [isXEPGEnabled, setIsXEPGEnabled] = useState(true);

  return (
    <WizardContext.Provider
      value={{ isXEPGEnabled, payload, setIsXEPGEnabled, setPayload }}
    >
      {children}
    </WizardContext.Provider>
  );
};

export default WizardProvider;
