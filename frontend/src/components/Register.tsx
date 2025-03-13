import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

const RegisterPage: React.FC = () => {
    const [step, setStep] = useState<"preregister" | "register">("preregister");
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [verificationCode, setVerificationCode] = useState("");
    const [password, setPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");
    const [countdown, setCountdown] = useState(3);

    const navigate = useNavigate();

    // 预注册：获取验证码
    const handlePreregister = async (event: React.FormEvent) => {
        event.preventDefault();
        setError("");
        setSuccess("");
        setLoading(true);

        if (!email.endsWith("@sjtu.edu.cn")) {
            setError("请输入 @sjtu.edu.cn 结尾的邮箱");
            setLoading(false);
            return;
        }

        try {
            const response = await fetch("http://localhost:8080/preregister", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, email }),
            });

            const result = await response.json();
            if (!response.ok) throw new Error(result.error || "预注册失败");

            setStep("register");
            setSuccess("验证码已发送，请检查邮箱");
        } catch (error: any) {
            setError(error.message);
        } finally {
            setLoading(false);
        }
    };

    // 最终注册
    const handleRegister = async (event: React.FormEvent) => {
        event.preventDefault();
        setError("");
        setSuccess("");
        setLoading(true);

        try {
            const response = await fetch("http://localhost:8080/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password, email, verification_code: verificationCode }),
            });

            const result = await response.json();
            if (!response.ok) throw new Error(result.error || "注册失败");

            setSuccess("注册成功！3 秒后跳转到登录页面...");
        } catch (error: any) {
            setError(error.message);
        } finally {
            setLoading(false);
        }
    };

    // 监听 success 变化，触发 3 秒倒计时跳转
    useEffect(() => {
        if (success.includes("注册成功")) {
            const interval = setInterval(() => {
                setCountdown((prev) => prev - 1);
            }, 1000);

            setTimeout(() => {
                clearInterval(interval);
                navigate("/login"); // 跳转到登录页面
            }, 3000);

            return () => clearInterval(interval);
        }
    }, [success, navigate]);

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
            <div className="bg-white p-6 rounded-lg shadow-md w-96">
                <h2 className="text-2xl font-bold text-center text-gray-800 mb-4">
                    {step === "preregister" ? "用户注册" : "输入验证码"}
                </h2>

                {step === "preregister" ? (
                    <form onSubmit={handlePreregister} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">用户名</label>
                            <input
                                type="text"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                required
                                className="mt-1 p-2 w-full border rounded-md focus:ring focus:ring-blue-300"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">邮箱</label>
                            <input
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                placeholder="name@sjtu.edu.cn"
                                className="mt-1 p-2 w-full border rounded-md focus:ring focus:ring-blue-300"
                            />
                        </div>

                        <button
                            type="submit"
                            className="w-full p-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition disabled:bg-gray-400"
                            disabled={loading}
                        >
                            {loading ? "发送中..." : "获取验证码"}
                        </button>
                    </form>
                ) : (
                    <form onSubmit={handleRegister} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">验证码</label>
                            <input
                                type="text"
                                value={verificationCode}
                                onChange={(e) => setVerificationCode(e.target.value)}
                                required
                                className="mt-1 p-2 w-full border rounded-md focus:ring focus:ring-blue-300"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">密码</label>
                            <input
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                className="mt-1 p-2 w-full border rounded-md focus:ring focus:ring-blue-300"
                            />
                        </div>

                        <button
                            type="submit"
                            className="w-full p-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition disabled:bg-gray-400"
                            disabled={loading}
                        >
                            {loading ? "注册中..." : "注册"}
                        </button>
                    </form>
                )}

                {error && <p className="text-red-600 text-sm mt-3">{error}</p>}
                {success && (
                    <p className="text-green-600 text-sm mt-3">
                        {success} {success.includes("注册成功") && countdown > 0 && `(${countdown}s)`}
                    </p>
                )}
            </div>
        </div>
    );
};

export default RegisterPage;
