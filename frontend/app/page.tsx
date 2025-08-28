export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-6xl text-center font-bold text-gray-800 mb-8 underline decoration-blue-500 decoration-4">
          Hello World
        </h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white rounded-lg shadow-lg p-6 border border-gray-200">
            <h2 className="text-2xl font-semibold text-gray-700 mb-4">Tailwind Test</h2>
            <p className="text-gray-600">This card should have a white background, rounded corners, and a shadow.</p>
          </div>
          <div className="bg-blue-500 rounded-lg shadow-lg p-6 text-white">
            <h2 className="text-2xl font-semibold mb-4">Blue Card</h2>
            <p>This card should have a blue background and white text.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
