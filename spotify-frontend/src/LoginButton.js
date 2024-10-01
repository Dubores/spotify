// LoginButton.js
import React, { useState, useEffect } from 'react';

const LoginButton = () => {
  const [authURL, setAuthURL] = useState('');

  useEffect(() => {
    fetch('http://localhost:8080/spotify/auth-url')
      .then(response => response.json())
      .then(data => setAuthURL(data.authURL));
  }, []);

  const handleLoginClick = () => {
    if (authURL) {
      window.location.href = authURL;
    }
  };

  return (
    <button onClick={handleLoginClick}>Log In to Spotify</button>
  );
};

export default LoginButton;