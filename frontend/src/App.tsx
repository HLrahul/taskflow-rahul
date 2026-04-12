import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { ProtectedRoute } from "./components/ProtectedRoute";

import { Login } from "./pages/Login";
import Register from "./pages/Register";
import { Navbar } from "./components/layout/Navbar";
import { Dashboard } from "./pages/Dashboard";
import { ProjectDetail } from "./pages/ProjectDetail";

function App() {
  return (
    <Router>
      <Routes>
        {/* PUBLIC ROUTES */}
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        {/* PROTECTED ROUTES - Wrapped by our Auth Guard Bouncer! */}
        <Route element={<ProtectedRoute />}>
          <Route
            path="/"
            element={
              <>
                <Navbar />
                <Dashboard />
              </>
            }
          />
          <Route path="/project/:id" element={<><Navbar /><ProjectDetail /></>} />
        </Route>

        {/* CATCH ALL - Redirect unknown URLs to the dashboard */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
