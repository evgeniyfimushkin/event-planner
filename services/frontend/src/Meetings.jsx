import './Meetings.css'
import Bar from "./components/menu/Bar.jsx"
import Grid from "./components/cards/Grid.jsx"
import Login from './components/auth/Login.jsx'

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./services/AuthContext.jsx";
import Register from './components/auth/Register.jsx';
import PrivateRoute from './services/PrivateRoute.jsx';
import CreateEvent from './components/event/CreateEvent.jsx';
import Events from './components/tabs/Events.jsx';
import Calendar from './components/tabs/Calendar.jsx';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={
            <PrivateRoute>
              <Bar />
              <Events />
            </PrivateRoute>
          }/>
          <Route path="/calendar" element={
            <PrivateRoute>
              <Bar />
              <Calendar />
            </PrivateRoute>
          }/>
          <Route path="/events" element={
            <PrivateRoute>
              <CreateEvent />
            </PrivateRoute>
          }/>
          <Route path="/register" element={<Register />}/>
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
