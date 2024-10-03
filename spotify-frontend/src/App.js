import React, { useState, useEffect } from 'react';
import LoginPage from './LoginPage';

function App() {
  const [user, setUser] = useState(null);
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('localhost:8080/api/user/profile', { method: 'GET' });
        if (response.ok) {
          const data = await response.json();
          setUser(data);
          setIsLoggedIn(true);
        }
      } catch (error) {
        // Handle error
        console.error('Failed to fetch user data:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <div>
      {isLoggedIn && user ? (
        <div>
          <p>Welcome, {user.userName}!</p>
          <img src={user.url} alt="User" />
        </div>
      ) : (
        <LoginPage></LoginPage>
      )}
    </div>
  );
};

export default App;