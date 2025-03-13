import React from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Navbar from "./components/Navbar";
import Login from "./components/Login";
import EventList from "./components/Events";
import EventDetail from "./components/EventDetail";
import RegisterPage from "./components/Register";
import ProtectedRoute from "./components/ProtectedRoute"; // 导入保护组件

const App: React.FC = () => {
  return (
    <Router>
      <Navbar />
      <Routes>
        {/* 公开页面 */}
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* 受保护页面 */}
        <Route element={<ProtectedRoute />}>
          <Route path="/events" element={<EventList />} />
          <Route path="/events/:id" element={<EventDetail />} />
        </Route>

        {/* 默认重定向 */}
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </Router>
  );
};

export default App;
