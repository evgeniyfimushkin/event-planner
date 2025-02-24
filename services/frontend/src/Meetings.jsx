import './Meetings.css'
import Bar from "./components/menu/Bar.jsx"
import Grid from "./components/cards/Grid.jsx"
import Login from './components/auth/Login.jsx'

import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./services/AuthContext.jsx";

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
            <Route path="/" element={
              <>
                <Bar />
                <Grid />
              </>
            }/>
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
