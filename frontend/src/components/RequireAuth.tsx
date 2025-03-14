import { Navigate, useLocation } from "react-router-dom";
import { JSX } from "react/jsx-dev-runtime";

const RequireAuth: React.FC<{ children: JSX.Element }> = ({ children }) => {
  const token = localStorage.getItem("token"); // 获取 Token
  const location = useLocation(); // 获取当前访问路径

  if (!token) {
    localStorage.setItem("redirectPath", location.pathname); // 存储用户访问路径
    return <Navigate to="/login" replace />; // 跳转到登录页
  }

  return children;
};

export default RequireAuth;
