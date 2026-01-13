import { createLogger, format, transports, Logger } from 'winston';

/**
 * Winston logger configuration
 * Provides structured logging with different levels and formats
 */
export const createAppLogger = (context?: string): Logger => {
  const logger = createLogger({
    level: process.env.LOG_LEVEL || 'info',
    format: format.combine(
      format.timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }),
      format.errors({ stack: true }),
      format.splat(),
      format.json(),
    ),
    defaultMeta: { service: 'foodlink-backend', context },
    transports: [
      // Console transport with colorized output for development
      new transports.Console({
        format: format.combine(
          format.colorize(),
          format.printf(({ timestamp, level, message, context, stack }) => {
            const contextStr = context ? `[${context}]` : '';
            const stackStr = stack ? `\n${stack}` : '';
            return `${timestamp} ${level} ${contextStr} ${message}${stackStr}`;
          }),
        ),
      }),
      // File transport for errors
      new transports.File({
        filename: 'logs/error.log',
        level: 'error',
        format: format.combine(format.timestamp(), format.json()),
      }),
      // File transport for all logs
      new transports.File({
        filename: 'logs/combined.log',
        format: format.combine(format.timestamp(), format.json()),
      }),
    ],
  });

  // In production, also log to a rotating file
  if (process.env.NODE_ENV === 'production') {
    logger.add(
      new transports.File({
        filename: 'logs/combined.log',
        maxsize: 5242880, // 5MB
        maxFiles: 5,
      }),
    );
  }

  return logger;
};
