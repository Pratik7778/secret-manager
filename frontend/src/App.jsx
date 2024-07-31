import { useState } from 'react'
import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';

import { Route, Routes } from "react-router-dom"

import Homepage from './components/Welcome'
import LoginRegisterPage from './components/LoginRegister'
import ListSecretsPage from './components/ListSecrets';
// import { defineConfig, loadEnv } from 'vite';

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <div>
        <Routes>
          <Route path="/" element={<Homepage />} />
          <Route path="/login" element={<LoginRegisterPage />} />
          <Route path="/secrets" element={<ListSecretsPage />} />
        </Routes>
      </div>
    </>
  )
}

export default App
