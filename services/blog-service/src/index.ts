import { start, stop } from './app';

// Graceful shutdown handlers
process.on('SIGINT', async () => {
  console.log('Received SIGINT, shutting down gracefully...');
  try {
    await stop();
    process.exit(0);
  } catch (error) {
    console.error('Error during shutdown:', error);
    process.exit(1);
  }
});

process.on('SIGTERM', async () => {
  console.log('Received SIGTERM, shutting down gracefully...');
  try {
    await stop();
    process.exit(0);
  } catch (error) {
    console.error('Error during shutdown:', error);
    process.exit(1);
  }
});

// Start the server
async function main() {
  try {
    await start();
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

// Only start if this file is run directly
if (require.main === module) {
  main();
}
