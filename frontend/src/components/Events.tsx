import { useEffect, useState } from "react";
import { Link } from "react-router-dom"; // 确保安装 react-router-dom

type Event = {
  id: number;
  name: string;
  location: string;
  start_time: string;
  end_time: string;
  current_participants: number;
  max_participants: number;
};

export default function EventList() {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const pageSize = 20;
  const token = localStorage.getItem("token");

  const headers = {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };

  const fetchEvents = async () => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/events?page=${page}&page_size=${pageSize}`, {
        method: "GET",
        headers,
      });

      if (!response.ok) throw new Error("无法加载活动信息");

      const data = await response.json();
      setEvents(data.events);
      setHasMore(data.events.length >= pageSize);
    } catch (err) {
      setError("无法加载活动信息");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
  }, [page]);

  const handleJoinEvent = async (eventId: number) => {
    if (!token) {
      alert("请先登录");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/join", {
        method: "POST",
        headers,
        body: JSON.stringify({ event_id: eventId }),
      });

      if (!response.ok) throw new Error("加入失败，请重试");

      setEvents((prevEvents) =>
        prevEvents.map((event) =>
          event.id === eventId
            ? { ...event, current_participants: event.current_participants + 1 }
            : event
        )
      );
      alert("加入成功！");
    } catch (error) {
      alert(error instanceof Error ? error.message : "发生未知错误");
    }
  };

  if (loading) return <p className="text-center mt-5">加载中...</p>;
  if (error) return <p className="text-red-500 text-center mt-5">{error}</p>;

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">即将举行的活动</h1>
      {events.length === 0 ? (
        <p className="text-center">暂无即将举行的活动</p>
      ) : (
        <div>
          <ul className="space-y-4">
            {events.map((event) => (
              <li key={event.id} className="bg-white p-4 rounded-lg shadow-md border">
                {/* 点击跳转到活动详情页面 */}
                <h2 className="text-lg font-semibold">
                  <Link to={`/events/${event.id}`} className="text-blue-500 hover:underline">
                    {event.name}
                  </Link>
                </h2>
                <p className="text-gray-500">{event.location}</p>
                <p className="text-gray-500">
                  {new Date(event.start_time).toLocaleString()} - {new Date(event.end_time).toLocaleString()}
                </p>
                <p className="text-gray-600">
                  参加人数: {event.current_participants} / {event.max_participants}
                </p>
                <button
                  onClick={() => handleJoinEvent(event.id)}
                  disabled={event.current_participants >= event.max_participants}
                  className={`mt-2 px-4 py-2 rounded-lg text-white font-medium transition ${
                    event.current_participants >= event.max_participants
                      ? "bg-gray-400 cursor-not-allowed"
                      : "bg-blue-500 hover:bg-blue-600"
                  }`}
                >
                  {event.current_participants >= event.max_participants ? "已满员" : "加入活动"}
                </button>
              </li>
            ))}
          </ul>

          {/* 分页按钮 */}
          <div className="flex justify-between mt-6">
            <button
              onClick={() => setPage((prev) => Math.max(prev - 1, 1))}
              disabled={page === 1}
              className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg disabled:opacity-50"
            >
              上一页
            </button>
            <span className="text-gray-700">当前页: {page}</span>
            <button
              onClick={() => setPage((prev) => (hasMore ? prev + 1 : prev))}
              disabled={!hasMore}
              className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg disabled:opacity-50"
            >
              下一页
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
