import { FontDemo, LoginButton } from "@/components";

export default function Home() {
  return (
    <main className="min-h-screen">
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">GymLog</h1>
          <LoginButton />
        </div>
      <FontDemo />
      </div>
    </main>
  );
}
