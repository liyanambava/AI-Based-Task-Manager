import TaskList from "../components/TaskList";
import TaskForm from "../components/TaskForm";

export default function Home() {
  return (
    <div className="min-h-screen w-full bg-cover bg-center bg-[url('https://happywall-img-gallery.imgix.net/69281/minimalist_arches_evergreen_limited.jpg')]">

      <h1 className="flex flex-col items-center justify-center text-5xl font-bold mb-4 p-4">Full-Stack Task Manager</h1>
      <TaskForm />
      <TaskList />
    </div>
  );
}

