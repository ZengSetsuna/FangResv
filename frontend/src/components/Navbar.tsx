import { Link } from "react-router-dom";
import React from "react";

const Navbar: React.FC = () => {
  return (
    <nav style={styles.navbar}>
      <h1 style={styles.logo}>坊预约</h1>
      <ul style={styles.navLinks}>
        <li><Link to="/events" style={styles.link}>活动列表</Link></li>
        <li><Link to="/account" style={styles.link}>我的账户</Link></li>
      </ul>
    </nav>
  );
};


const styles: { [key: string]: React.CSSProperties } = {
  navbar: {
    position: "fixed",  // 固定在顶部
    top: 0,
    left: 0,
    width: "100%",      // 宽度占满整个页面
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    padding: "1rem 2rem",
    background: "#282c34",
    color: "white",
    zIndex: 1000,       // 确保在最前面
  },
  logo: { fontSize: "1.5rem" },
  navLinks: { listStyle: "none", display: "flex", gap: "1.5rem", padding: 0 },
  link: { color: "white", textDecoration: "none", fontSize: "1.2rem" },
};

export default Navbar;
