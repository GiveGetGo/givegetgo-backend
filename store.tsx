// store.ts
import { createStore, combineReducers, Action } from 'redux';

// Define the action types
const SET_AVATAR_URI = 'SET_AVATAR_URI';

// Define the action creators
export const setAvatarUri = (uri: string) => ({
  type: SET_AVATAR_URI,
  payload: uri
});

// Define the state and action types
type AvatarState = {
  selectedAvatarUri: string;
};

type AvatarAction = ReturnType<typeof setAvatarUri>;

// Reducer
const avatarReducer = (state: AvatarState = { selectedAvatarUri: '' }, action: AvatarAction) => {
  switch (action.type) {
    case SET_AVATAR_URI:
      return { ...state, selectedAvatarUri: action.payload };
    default:
      return state;
  }
};

// Combine reducers if you have more than one reducer
const rootReducer = combineReducers({
  avatar: avatarReducer
});

// Create the store
const store = createStore(rootReducer);

export default store;
