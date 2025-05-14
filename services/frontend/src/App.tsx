import './App.css'

import React from 'react';
import SignIn from './components/auth/SignIn.tsx'
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./services/AuthContext.js";
import SignUp from './components/auth/SignUp.js';
import PrivateRoute from './services/PrivateRoute.js';
import CreateEvent from './components/event/CreateEvent.js';
import Events from './components/tabs/Events.js';
import Calendar from './components/tabs/Calendar.tsx';
import Menu from './components/menu/Menu.js';
import Profile from './components/tabs/Profile.tsx';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/signIn" element={<SignIn />} />
          <Route path="/" element={
            <PrivateRoute>
              <Menu />
              <Events />
            </PrivateRoute>
          }/>
          <Route path="/calendar" element={
            <PrivateRoute>
              <Menu />
              <Calendar />
            </PrivateRoute>
          }/>
          <Route path="/events" element={
            <PrivateRoute>
              <CreateEvent />
            </PrivateRoute>
          }/>
          <Route path="/signUp" element={<SignUp />}/>
          <Route path="/profile" element={
            <PrivateRoute>
              <Menu />
              <Profile />
            </PrivateRoute>
          }/>
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
