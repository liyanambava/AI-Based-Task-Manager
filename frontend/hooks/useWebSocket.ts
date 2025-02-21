import { useEffect, useState } from "react";

const useWebSocket = (url: string) => {
  const [messages, setMessages] = useState<string[]>([]);
  const [connected, setConnected] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket(url);

    socket.onopen = () => {
      setConnected(true);
      console.log("✅ WebSocket Connected");
    };

    socket.onmessage = (event) => {
      console.log("📩 Message received:", event.data);
      setMessages((prev) => [...prev, event.data]);
    };

    socket.onclose = () => {
      setConnected(false);
      console.log("❌ WebSocket Disconnected. Reconnecting...");
      setTimeout(() => setWs(new WebSocket(url)), 3000); // Auto-reconnect
    };

    socket.onerror = (error) => {
      console.error("⚠️ WebSocket Error:", error);
    };

    setWs(socket);

    return () => {
      socket.close();
    };
  }, [url]);

  return { connected, messages, ws };
};

export default useWebSocket;
