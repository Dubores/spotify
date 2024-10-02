import React from 'react';
import LoginButton from './LoginButton';

const LoginPage = () => {

  return (
    <div className="App">
      <header className="App-header">
        <div>
          <h1>Welcome to My Spotify App. Please log in</h1>
          <LoginButton />
        </div>
      </header>
    </div>
  );
};

export default LoginPage;