import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import Login from "./components/Login";
import EventList from "./components/Events";
import EventDetail from "./components/EventDetail";
import RegisterPage from "./components/Register";
import RequireAuth from "./components/RequireAuth"; // 导入 RequireAuth 组件

const App: React.FC = () => {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* 受保护的路由：未登录的用户会被重定向到 /login */}
        <Route
          path="/events"
          element={
            <RequireAuth>
              <EventList />
            </RequireAuth>
          }
        />
        <Route
          path="/events/:id"
          element={
            <RequireAuth>
              <EventDetail />
            </RequireAuth>
          }
        />
      </Routes>
    </Router>
  );
};

export default App;
