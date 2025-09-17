import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';
import * as schema from './schema';
import config from "./index"

// Create the connection
const client = postgres(config.database.url, {
  max: 10,
  idle_timeout: 20,
  connect_timeout: 10,
});

// Create the database instance
export const db = drizzle(client, { schema });

export default db;
