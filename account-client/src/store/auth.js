import { reactive, inject, toRefs, readonly, watchEffect } from 'vue';
import { storeTokens, getTokens, doRequest, getTokenPayload } from '../util';
import { useRouter } from 'vue-router';

const state = reactive({
  currentUser: null,
  idToken: null,
  isLoading: false,
  error: null,
});

const storeSymbol = Symbol();

const signin = async (email, password) =>
  await authenticate(email, password, '/api/account/signin');

const signup = async (email, password) =>
  await authenticate(email, password, '/api/account/signup');

const initializeUser = async () => {
  state.isLoading = true;
  state.error = null;

  const [idToken, refreshToken] = getTokens();

  const idTokenClaims = getTokenPayload(idToken);
  const refreshTokenClaims = getTokenPayload(refreshToken);

  if (idTokenClaims) {
    state.idToken = idToken;
    state.currentUser = idTokenClaims.user;
  }

  state.isLoading = false;

  // silently refresh tokens in local storage
  // if we for some reason don't have refresh token (e.g., if the user deleted it manually)
  // then we don't proceed
  if (!refreshTokenClaims) {
    return;
  }

  const { data, error } = await doRequest({
    url: '/api/account/tokens',
    method: 'post',
    data: {
      refreshToken,
    },
  });

  if (error) {
    console.error('Error refreshing tokens\n', error);
    return;
  }

  const { tokens } = data;
  storeTokens(tokens.idToken, tokens.refreshToken);

  const updatedIdTokenClaims = getTokenPayload(tokens.idToken);

  state.currentUser = updatedIdTokenClaims.user;
  state.idToken = tokens.idToken;
};

export const createAuthStore = (authStoreOptions) => {
  const { onAuthRoute, requireAuthRoute } = authStoreOptions || {};

  const authStore = {
    ...toRefs(readonly(state)),
    signin,
    signup,
    initializeUser,
    onAuthRoute,
    requireAuthRoute,
  };

  return {
    authStore,
    install: (app) => {
      app.provide(storeSymbol, authStore);
    },
  };
};

export function useAuth() {
  const store = inject(storeSymbol);

  if (!store) {
    throw new Error('Auth store has not been instantiated!');
  }

  const router = useRouter();

  watchEffect(() => {
    if (store.currentUser.value && store.onAuthRoute) {
      router.push(store.onAuthRoute);
    }

    if (!store.currentUser.value && store.requireAuthRoute) {
      router.push(store.requireAuthRoute);
    }
  });

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
