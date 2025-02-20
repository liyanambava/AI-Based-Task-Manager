"use client"; // Ensure it's a Client Component

import { useEffect, useState } from "react";

interface Task {
  _id?: string;
  title: string;
  description: string;
  status: string;
}

const TaskList = () => {
  const [tasks, setTasks] = useState<Task[]>([]);

  useEffect(() => {
    fetch("http://localhost:8080/tasks")
      .then((res) => res.json())
      .then((data) => setTasks(data))
      .catch((err) => console.error("Error fetching tasks:", err));
  }, []);

  return (
    <div className="p-6 bg-gray-100 shadow-lg rounded-lg">
      <h2 className="text-2xl font-bold text-gray-800 mb-4">Task List</h2>
      {tasks.length === 0 ? (
        <p className="text-gray-600">No tasks found.</p>
      ) : (
        <ul className="space-y-3">
          {tasks.map((task, index) => (
            <li
              key={task._id || `task-${index}`}
              className="border border-gray-300 p-4 rounded-lg bg-white shadow-md"
            >
              <strong className="text-gray-900">{task.title}</strong>:{" "}
              <span className="text-gray-700">{task.description}</span> -{" "}
              <span
                className={`${
                  task.status === "completed"
                    ? "text-green-600 font-semibold"
                    : "text-blue-600"
                }`}
              >
                {task.status}
              </span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default TaskList;
