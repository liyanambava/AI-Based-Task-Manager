"use client";

import { useState } from "react";

const TaskForm = () => {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/tasks", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title, description, status: "pending" }),
    });
    if (res.ok) {
      alert("Task created!");
      setTitle("");
      setDescription("");
    } else {
      alert("Failed to create task.");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full max-w-4xl mx-auto p-4 bg-gray-100 shadow-lg text-white rounded-lg bg-opacity-80">
      <h2 className="text-lg font-bold mb-2 text-gray-800">Create New Task</h2>
      <input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className="w-full p-2 border border-gray-450 rounded mb-2 bg-gray-400 text-black placeholder-gray-600"
      />
      <textarea
        placeholder="Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        className="w-full p-2 border border-gray-450 rounded mb-2 bg-gray-400 text-black placeholder-gray-600"
      />
      <button
        type="submit"
        className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
      >
        Add Task
      </button>
    </form>
  );
};

export default TaskForm;
