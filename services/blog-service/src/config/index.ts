import dotenv from "dotenv";

dotenv.config(); // load .env here (single source of truth)

export interface Config {
  port: number;
  host: string;
  database: {
    url: string;
  };
  jwt: {
    secret: string;
    expiresIn: string;
  };
  cors: {
    origin: string[];
  };
}

const config: Config = {
  port: parseInt(process.env.PORT || "3009"),
  host: process.env.HOST || "0.0.0.0",
  database: {
    url: process.env.DATABASE_URL || "postgresql://localhost:5432/blog_service",
  },
  jwt: {
    secret: process.env.JWT_SECRET || "your-secret-key",
    expiresIn: process.env.JWT_EXPIRES_IN || "7d",
  },
  cors: {
    origin: process.env.CORS_ORIGIN
      ? process.env.CORS_ORIGIN.split(",")
      : ["http://localhost:3000"],
  },
};

export default config;