import React, { useState, useEffect } from "react"
import BloomFilter from "./utils/BloomFilter";

const App = () => {
  const [bloomFilter, setBloomFilter] = useState(null);
  const [password, setPassword] = useState('');
  const [searchResult, setSearchResult] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  

  // const loadBloomFilter = async () => {
  //   try {
  //     setLoading(true);
  //     const response = await fetch('http://localhost:8000/bloom-filter');
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
        const response = await fetch('http://localhost:8000/bloom-filter');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = response.json();
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
    setPassword(e.target.value);
  };
    
  const handleCheck = (password) => {
    if (bloomFilter) {
      const result = bloomFilter.checkPassword(password);
      setSearchResult(result ? 'Password might be common!' : 'Password is not common and can be used.');
    }
    else {
      setError('Bloom Filter not found');
    }
  }
    
  const handleSave = (password) => {
    if (bloomFilter) {
      bloomFilter.addPassword(password);
      setSearchResult("Password is added as common password.")
    }
    else {
      setError('Bloom Filter not found');
    }
  };

  return (
    <div>
      {loading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>Error: {error}</p>}
      <p>Signup page</p>
      {bloomFilter && (
          <div style={{ margin: '20px', fontFamily: 'Arial, sans-serif' }}>
            <h1>Password Check</h1>
            <div>
              <label htmlFor="password">Enter your password:</label>
              <br />
              <input
                type="password"
                id="password"
                value={password}
                onChange={handlePasswordChange}
                style={{ margin: '10px 0', padding: '5px', width: '300px' }}
              />
            </div>
            <div>
              <button onClick={handleCheck} style={{ marginRight: '10px', padding: '5px 10px' }}>
                Check
              </button>
              <button onClick={handleSave} style={{ padding: '5px 10px' }}>
                Save
              </button>
            </div>
          </div>
      )}
    </div>
  );
};

export default App;
