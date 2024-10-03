// LoginButton.js
import React, { useState, useEffect } from 'react';

const LoginButton = () => {
  const [authURL, setAuthURL] = useState('');

  const handleLoginClick = () => {
    window.location.href = "http://localhost:8080/login";
  };

  return (
    <button onClick={handleLoginClick}>Log In to Spotify</button>
  );
};

export default LoginButton;