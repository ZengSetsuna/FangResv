import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

const Login: React.FC = () => {
  const [username, setUsername] = useState(""); 
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "登录失败");
      }

      localStorage.setItem("token", data.token); // 存储 Token

      // 获取存储的重定向路径
      const redirectPath = localStorage.getItem("redirectPath") || "/events";
      localStorage.removeItem("redirectPath"); // 登录成功后删除记录

      navigate(redirectPath); // 登录成功后跳转到之前访问的页面
    } catch (err: any) {
      console.error("登录请求出错:", err);
      setError(err.message);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-100">
      <div className="w-full max-w-md bg-white p-8 rounded-lg shadow-lg">
        <h2 className="text-2xl font-bold text-center text-gray-700">登录</h2>

        {error && <p className="mt-2 text-center text-red-500">{error}</p>}

        <form onSubmit={handleSubmit} className="mt-6">
          <div>
            <label className="block text-gray-600">用户名</label>
            <input
              type="text"
              name="username"
              className="w-full mt-1 p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="请输入用户名"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>

          <div className="mt-4">
            <label className="block text-gray-600">密码</label>
            <input
              type="password"
              name="password"
              className="w-full mt-1 p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="请输入密码"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <button
            type="submit"
            className="w-full mt-6 bg-blue-500 text-white py-3 rounded-lg hover:bg-blue-600 transition"
          >
            登录
          </button>
        </form>

        {/* 注册按钮 */}
        <div className="mt-4 text-center">
          <p className="text-gray-600">还没有账号？</p>
          <button
            onClick={() => navigate("/register")}
            className="mt-2 w-full bg-gray-300 text-gray-700 py-3 rounded-lg hover:bg-gray-400 transition"
          >
            注册
          </button>
        </div>
      </div>
    </div>
  );
};

export default Login;
