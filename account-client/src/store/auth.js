import { reactive, provide, inject, toRefs, readonly } from 'vue';
import { storeTokens, doRequest, getTokenPayload } from '../util';

const state = reactive({
  currentUser: 'user',
  idToken: null,
  isLoading: false,
  error: null,
});

const signin = async (email, password) =>
  await authenticate(email, password, '/api/account/signin');

const signup = async (email, password) =>
  await authenticate(email, password, '/api/account/signup');

export const authStore = {
  ...toRefs(readonly(state)),
  signin,
  signup,
};

const storeSymbol = Symbol();

export function provideAuth() {
  provide(storeSymbol, authStore);
}

export function useAuth() {
  const store = inject(storeSymbol);

  if (!store) {
    throw new Error('Auth store has not been instantiated!');
  }

  return store;
}

// authenticate implements common code between signin and signup
const authenticate = async (email, password, url) => {
  state.isLoading = true;
  state.error = null;

  const { data, error } = await doRequest({
    url,
    method: 'post',
    data: {
      email,
      password,
    },
  });

  if (error) {
    state.error = error;
    state.isLoading = false;
    return;
  }

  const { tokens } = data;

  storeTokens(tokens.idToken, tokens.refreshToken);

  const tokenClaims = getTokenPayload(tokens.idToken);

  // set tokens to local storage with expiry (separate function)
  state.idToken = tokens.idToken;
  state.currentUser = tokenClaims.user;
  state.isLoading = false;
};
