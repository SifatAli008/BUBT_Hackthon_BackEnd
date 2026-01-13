/**
 * Configuration interface for application settings
 */
export interface IConfig {
  port: number;
  nodeEnv: string;
  database: {
    url: string;
    host: string;
    port: number;
    user: string;
    password: string;
    database: string;
    maxPoolSize: number;
    idleTimeoutMs: number;
  };
  jwt: {
    secret: string;
    expiresIn: string;
    refreshSecret: string;
    refreshExpiresIn: string;
  };
  logging: {
    level: string;
  };
}
