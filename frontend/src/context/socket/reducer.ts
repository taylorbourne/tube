import { ServerConfigType, StateType } from 'types'

type ActionMap<M extends { [index: string]: any }> = {
  [Key in keyof M]: M[Key] extends undefined
    ? {
        type: Key;
      }
    : {
        type: Key;
        payload: M[Key];
      }
};

export enum Types {
  GET_SYSTEM_DATA = "GET_SYSTEM_DATA",
}

type Payload = {
  [Types.GET_SYSTEM_DATA]: ServerConfigType;
};

export type Actions = ActionMap<Payload>[keyof ActionMap<
  Payload
>];

export const reducer = (
  state: StateType,
  action: Actions
) => {
  switch (action.type) {
    case Types.GET_SYSTEM_DATA:
      return {
        ...state,
        serverConfig: {
          ...state.serverConfig,
          ...action.payload,
        },
      };
    default:
      return state;
  }
};

export default reducer