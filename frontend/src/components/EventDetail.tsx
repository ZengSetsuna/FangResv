import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function EventDetail() {
  const { id } = useParams();

  interface Event {
    event_name: string;
    description: string | null;
    location: string;
    start_time: string;
    end_time: string;
    organizer: string;
    current_participants: number;
    max_participants: number;
    can_join: boolean;
    participants: string | string[];
  }

  const [event, setEvent] = useState<Event | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const token = localStorage.getItem("token");

  const headers = {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
  

  useEffect(() => {
    const fetchEventDetails = async () => {
      try {
        const response = await fetch(`http://localhost:8080/events/${id}`, {
            method: "GET",
            headers,
            }
        );
        if (!response.ok) throw new Error("无法加载活动信息");
        const data = await response.json();

        // 如果 participants 是 Base64 编码的字符串，则解码
        let participants: string[] = [];
        if (typeof data.participants === "string") {
          try {
            const decodedStr = atob(data.participants); // Base64 解码
            participants = JSON.parse(decodedStr); // 解析 JSON
          } catch (err) {
            console.error("无法解析 participants:", err);
          }
        } else {
          participants = data.participants || [];
        }

        setEvent({ ...data, participants });
      } catch (err) {
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError("未知错误");
        }
      } finally {
        setLoading(false);
      }
    };
    fetchEventDetails();
  }, [id]);

  if (loading) return <p>加载中...</p>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div className="max-w-2xl mx-auto p-6 bg-white shadow-md rounded-lg">
      {event && (
        <>
          <h1 className="text-2xl font-bold mb-4">{event.event_name}</h1>
          <p className="text-gray-600 mb-2">{event.description || "无描述"}</p>
          <p className="text-gray-700">地点: {event.location}</p>
          <p className="text-gray-700">开始时间: {new Date(event.start_time).toLocaleString()}</p>
          <p className="text-gray-700">结束时间: {new Date(event.end_time).toLocaleString()}</p>
          <p className="text-gray-700">组织者: {event.organizer}</p>
          <p className="text-gray-700">
            参加人数: {event.current_participants} / {event.max_participants}
          </p>
          <h2 className="text-lg font-semibold mt-4">参与者名单</h2>
          {event.participants.length > 0 ? (
            <ul>
              {Array.isArray(event.participants) && event.participants.map((p, index) => (
                <li key={index}>{p}</li>
              ))}
            </ul>
          ) : (
            <p>暂无参与者</p>
          )}
        </>
      )}
    </div>
  );
}
