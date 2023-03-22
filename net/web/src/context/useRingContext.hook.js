import { useEffect, useContext, useState, useRef } from 'react';

export function useRingContext() {
  const [state, setState] = useState({
    ringing: new Map(),
  });
  const access = useRef(null);

  const EXPIRE = 3000000
  const ringing = useRef(new Map());

  const updateState = (value) => {
    setState((s) => ({ ...s, ...value }))
  }

  const actions = {
    clear: () => {
      ringing.current = new Map();
      updateState({ ringing: ringing.current });
    },
    ring: (cardId, callId, calleeToken) => {
      const key = `${cardId}:${callId}`
      const call = ringing.current.get(key) || { cardId, calleeToken, callId }
      call.expires = Date.now() + EXPIRE;
      ringing.current.set(key, call);
      updateState({ ringing: ringing.current });
      setTimeout(() => {
        updateState({ ringing: ringing.current });
      }, 3000);
    },
    ignore: (cardId, callId) => {
      const key = `${cardId}:${callId}`
console.log("IGNORE", key);
      const call = ringing.current.get(key);
      if (call) {
        call.status = 'ignored'
        ringing.current.set(key, call);
        updateState({ ringing: ringing.current });
console.log(ringing.current);

      }
    },
    decline: (cardId, callId) => {
      const key = `${cardId}:${callId}`
      const call = ringing.current.get(key);
      if (call) {
        call.status = 'declined'
        ringing.current.set(key, call);
        updateState({ ringing: ringing.current });
      }
    },
    accept: (cardId, callId) => {
      const key = `${cardId}:${callId}`
      const call = ringing.current.get(key);
      if (call) {
        call.status = 'accepted'
        ringing.current.set(key, call);
        updateState({ ringing: ringing.current });
      }
    },
  }

  return { state, actions }
}

