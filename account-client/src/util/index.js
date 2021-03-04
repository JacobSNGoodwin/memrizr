import axios from 'axios';
import jwt_decode from 'jwt-decode';

// doRequest is a helper function for
// handling axios responses - reqOptions follow axios req config
export const doRequest = async (reqOptions) => {
  let error;
  let data;

  try {
    const response = await axios.request(reqOptions);
    data = response.data;
  } catch (e) {
    if (e.response) {
      error = e.response.data.error;
    } else if (e.request) {
      error = e.request;
    } else {
      error = e;
    }
  }

  if (error) {
    console.error(error);
  }

  return {
    data,
    error,
  };
};

const idTokenKey = '__malcorpId';
const refreshTokenKey = '__malcorpRf';

// storeTokens utility for storing idAndRefreshToken
export const storeTokens = (idToken, refreshToken) => {
  localStorage.setItem(idTokenKey, idToken);
  localStorage.setItem(refreshTokenKey, refreshToken);
};

export const getTokens = () => {
  return [
    localStorage.getItem(idTokenKey),
    localStorage.getItem(refreshTokenKey),
  ];
};

// gets the token's payload, and returns null
// if invalid
export const getTokenPayload = (token) => {
  if (!token) {
    return null;
  }

  const tokenClaims = jwt_decode(token);

  if (Date.now() / 1000 >= tokenClaims.exp) {
    return null;
  }

  return tokenClaims;
};
