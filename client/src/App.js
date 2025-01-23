import React, { useState, useEffect } from "react"
import BloomFilter from "./utils/BloomFilter";
import "./App.css";

const App = () => {
  const [bloomFilter, setBloomFilter] = useState(null);
  const [password, setPassword] = useState("");
  const [searchAndSaveResult, setSearchAndSaveResult] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  

  // const loadBloomFilter = async () => {
  //   try {
  //     setLoading(true);
  //     const response = await fetch(`${process.env.BACKEND_BASE_URL}/bloom-filter`);
  //     if (!response.ok) {
  //       throw new Error('Failed to fetch bloom filter');
  //     }
  //     const data = await response.json();
      
  //     // Store in localStorage for offline use
  //     localStorage.setItem('bloomFilter', JSON.stringify(data));
      
  //     setBloomFilter(data);
  //     setError(null);
  //   } catch (error) {
  //     console.error('Error loading bloom filter:', error);
  //     setError('Failed to load bloom filter');
      
  //     const cachedData = localStorage.getItem('bloomFilter');
  //     if (cachedData) {
  //       setBloomFilter(JSON.parse(cachedData));
  //     }
  //   } finally {
  //     setLoading(false);
  //   }
  // };

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        // TODO: Get this url from env file
        const response = await fetch(`${process.env.BACKEND_BASE_URL}/bloom-filter`);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        const { bitArray, bitArraySize, hashFunctionCount } = result;
        const filter = new BloomFilter(bitArray, bitArraySize, hashFunctionCount);
        setBloomFilter(filter);
      } catch (err) {
        setError(err.message)
      }
      finally {
        setLoading(false)
      }
    };

    fetchData();
  }, []);

  const handlePasswordChange = (e) => {
    setPassword(String(e.target.value));
  };
    
  const handleCheck = () => {
    if (bloomFilter) {
      const result = bloomFilter.checkPassword(password);
      const message = result ? 'Password is too common. Choose a stronger one.' : 'This password is safe to use.'
      setSearchAndSaveResult(message);
    }
    else {
      setError('Bloom Filter not found');
    }
  }

  const postToBackend = async () => {
    const data = {
      "password": password
    }
    try {
      // TODO: Get this url from env file
      const response = await fetch(`${process.env.BACKEND_BASE_URL}/add`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data)
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      console.log("Password added as common password to backend as well...", result)
    } catch (err) {
      setError(err.message)
    }
  }
    
  const handleSave = () => {
    if (bloomFilter) {
      bloomFilter.addPassword(password);
      postToBackend();
      setSearchAndSaveResult("Password is added as common password.")
    }
    else {
      setError('Bloom Filter not found');
    }
  };

  return (
    <div className="screen">
      {loading && <p className="loading">Loading...</p>}
      {error && <p className="error">Error: {error}</p>}
      {bloomFilter && (
        <div className="card">
          <h1>Check your password</h1>
          <div className="input-container">
            <label htmlFor="password">Enter your password:</label>
            <br />
            <input
              type="text"
              id="password"
              value={password}
              onChange={handlePasswordChange}
            />
          </div>
          <div className="button-container">
            <button onClick={handleCheck} className="button">
              Check
            </button>
            <button onClick={handleSave} className="button">
              Save
            </button>
          </div>
          {searchAndSaveResult && (
            <div className="message-container">
              <p className={`result-message ${searchAndSaveResult.includes("too common") ? "error" : "success"}`}>{searchAndSaveResult}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default App;
