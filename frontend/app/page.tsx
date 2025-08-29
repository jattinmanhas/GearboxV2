import { ThemeToggle } from '@/components/ui/theme-toggle'

export default function Home() {
  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-6xl font-bold text-foreground underline decoration-primary decoration-4">
            Hello World
          </h1>
          <ThemeToggle />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-card rounded-lg shadow-lg p-6 border border-border">
            <h2 className="text-2xl font-semibold text-card-foreground mb-4">Tailwind Test</h2>
            <p className="text-muted-foreground">This card should have a white background, rounded corners, and a shadow.</p>
          </div>
          <div className="bg-primary rounded-lg shadow-lg p-6 text-primary-foreground">
            <h2 className="text-2xl font-semibold mb-4">Blue Card</h2>
            <p>This card should have a blue background and white text.</p>
          </div>
        </div>
      </div>
    </div>
  );
}
