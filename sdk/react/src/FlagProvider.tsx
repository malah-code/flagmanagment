import React, { createContext, useContext, useEffect, useRef } from 'react';
import { FlagClient } from './client';

interface FlagContextType {
  client: FlagClient;
}

const FlagContext = createContext<FlagContextType | null>(null);

export const useFlagContext = () => {
  const context = useContext(FlagContext);
  if (!context) {
    throw new Error('useFlagContext must be used within a FlagProvider');
  }
  return context;
};

interface FlagProviderProps {
  apiKey: string;
  streamUrl: string;
  context?: Record<string, any>;
  children: React.ReactNode;
}

export const FlagProvider: React.FC<FlagProviderProps> = ({ 
  apiKey, 
  streamUrl, 
  context,
  children 
}) => {
  // Use a ref to ensure the client is instantiated only once
  const clientRef = useRef<FlagClient | null>(null);

  if (!clientRef.current) {
    clientRef.current = new FlagClient(apiKey, streamUrl);
    if (context) {
      clientRef.current.setContext(context);
    }
  }

  useEffect(() => {
    const client = clientRef.current;
    if (client) {
      client.connect();
    }
    return () => {
      if (client) {
        client.disconnect();
      }
    };
  }, []);

  useEffect(() => {
    if (clientRef.current && context) {
      clientRef.current.setContext(context);
    }
  }, [context]);

  return (
    <FlagContext.Provider value={{ client: clientRef.current! }}>
      {children}
    </FlagContext.Provider>
  );
};
