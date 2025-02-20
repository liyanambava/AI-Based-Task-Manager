import TaskList from "../components/TaskList";
import TaskForm from "../components/TaskForm";

export default function Home() {
  return (
    <div className="max-w-2xl mx-auto mt-10">
      <h1 className="text-2xl font-bold mb-4">AI Task Manager</h1>
      <TaskForm />
      <TaskList />
    </div>
  );
}
