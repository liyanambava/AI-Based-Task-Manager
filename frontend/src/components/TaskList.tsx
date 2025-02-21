"use client"; // ✅ Ensure this is a Client Component

import { useEffect, useState } from "react";

interface Task {
  _id?: string; // MongoDB `_id`
  id?: string; // Some APIs might return `id`
  title: string;
  description: string;
  status: string;
}

const TaskList = () => {
  const [tasks, setTasks] = useState<Task[] | null>(null); // ✅ Use `null` to prevent hydration mismatch
  const [ws, setWs] = useState<WebSocket | null>(null);

  // Fetch tasks on client mount
  useEffect(() => {
    fetchTasks();

    if (typeof window !== "undefined") {
      const socket = new WebSocket("ws://localhost:8080/ws");

      socket.onopen = () => {
        console.log("✅ WebSocket connected!");
      };

      socket.onmessage = (event) => {
        console.log("📩 WebSocket Message Received:", event.data);
        fetchTasks(); // ✅ Refresh tasks after WebSocket update
      };

      socket.onclose = () => {
        console.warn("⚠️ WebSocket Disconnected. Reconnecting in 3s...");
        setTimeout(() => setWs(new WebSocket("ws://localhost:8080/ws")), 3000);
      };

      setWs(socket);

      return () => {
        socket.close();
      };
    }
  }, []);

  // Function to fetch all tasks
  const fetchTasks = () => {
    fetch("http://localhost:8080/tasks")
      .then((res) => res.json())
      .then((data) => {
        console.log("🔍 Tasks fetched:", data);
        setTasks(data);
      })
      .catch((err) => console.error("❌ Error fetching tasks:", err));
  };

  // Function to update task status
  const handleCheckboxChange = async (task: Task) => {
    const taskId = task._id || task.id; // ✅ Use `_id` if available, otherwise fallback to `id`

    if (!taskId) {
      console.error("❌ Task ID is undefined! Cannot update status.");
      return;
    }

    console.log(`🔄 Updating task status for ID: ${taskId}`);

    const updatedStatus = task.status === "completed" ? "pending" : "completed";

    try {
      const res = await fetch(`http://localhost:8080/tasks/${taskId}/status`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: updatedStatus }),
      });

      if (res.ok) {
        const updatedTask = await res.json();
        setTasks((prevTasks) =>
          prevTasks
            ? prevTasks.map((t) => (t._id === updatedTask._id ? updatedTask : t))
            : []
        );
        console.log(`✅ Task status updated successfully: ${taskId}`);
      } else {
        console.error("❌ Failed to update task status");
      }
    } catch (error) {
      console.error("❌ Error updating task:", error);
    }
  };

  const handleDeleteTask = async (taskId: string | undefined) => {
    if (!taskId) {
      console.error("❌ Task ID is undefined! Cannot delete.");
      return;
    }

    try {
      const res = await fetch(`http://localhost:8080/tasks/${taskId}`, {
        method: "DELETE",
      });

      if (res.ok) {
        setTasks((prevTasks) =>
          prevTasks ? prevTasks.filter((task) => task._id !== taskId) : []
        );
        console.log(`🗑️ Task deleted: ${taskId}`);
      } else {
        console.error("❌ Failed to delete task");
      }
    } catch (error) {
      console.error("❌ Error deleting task:", error);
    }
  };

  return (
    <div className="w-full max-w-4xl mx-auto p-4 bg-gray-100 shadow-lg text-white rounded-lg bg-opacity-80">
      <h2 className="text-2xl font-bold text-gray-800 mb-4">Task List</h2>

      {tasks === null ? (
        <p className="text-gray-600">Loading tasks...</p>
      ) : tasks.length === 0 ? (
        <p className="text-gray-600">No tasks found.</p>
      ) : (
        <ul className="space-y-3">
          {tasks.map((task, index) => (
            <li
              key={task._id || task.id || `task-${index}`}
              className="flex justify-between items-center border border-gray-300 p-4 rounded-lg bg-white shadow-md"
            >
              <div className="flex items-center">
                <input
                  type="checkbox"
                  checked={task.status === "completed"}
                  onChange={() => handleCheckboxChange(task)}
                  className="mr-3"
                />
                <div>
                  <strong className="text-gray-900">{task.title}</strong>:{" "}
                  <span className="text-gray-700">{task.description}</span> -{" "}
                  <span
                    className={
                      task.status === "completed"
                        ? "text-green-600 font-semibold"
                        : "text-blue-600"
                    }
                  >
                    {task.status}
                  </span>
                </div>
              </div>
              <button
                onClick={() => handleDeleteTask(task._id)}
                className="ml-4 px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700"
              >
                Delete
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};


export default TaskList;
