// src/types.ts - TypeScript types
export interface Player {
  name: string;
  isReady: boolean;
}

export interface GameData {
  text: string;
  StartTime: string;
  IsActive: string;
}

export interface WordCompletePayload {
  type: 'wordComplete';
  content: {
    word: string;
  };
}

export interface JoinPayload {
  type: 'join';
  content: {
    room: string;
  };
}

export interface ReadyPayload {
  type: 'ready';
}
