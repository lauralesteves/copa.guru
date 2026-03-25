import { useCallback, useEffect, useState } from 'react';
import type { Match } from '../types/worldcup';

const STORAGE_KEY = 'copa-guru-predictions';

export interface Prediction {
  matchId: number;
  goals1: number;
  goals2: number;
  lockedAt: number;
}

interface PredictionsState {
  predictions: Prediction[];
  champion: string | null;
  createdAt: number;
}

function loadState(): PredictionsState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) return JSON.parse(raw);
  } catch {
    // ignore
  }
  return { predictions: [], champion: null, createdAt: Date.now() };
}

function saveState(state: PredictionsState) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch {
    // ignore
  }
}

export function usePredictions(matches: Match[]) {
  const [state, setState] = useState<PredictionsState>(loadState);

  useEffect(() => {
    saveState(state);
  }, [state]);

  const setPrediction = useCallback(
    (matchId: number, goals1: number, goals2: number) => {
      setState((prev) => {
        const existing = prev.predictions.findIndex(
          (p) => p.matchId === matchId,
        );
        const prediction: Prediction = {
          matchId,
          goals1,
          goals2,
          lockedAt: Date.now(),
        };

        const predictions =
          existing >= 0
            ? prev.predictions.map((p, i) => (i === existing ? prediction : p))
            : [...prev.predictions, prediction];

        return { ...prev, predictions };
      });
    },
    [],
  );

  const setChampion = useCallback((teamCode: string | null) => {
    setState((prev) => ({ ...prev, champion: teamCode }));
  }, []);

  const getPrediction = useCallback(
    (matchId: number): Prediction | undefined => {
      return state.predictions.find((p) => p.matchId === matchId);
    },
    [state.predictions],
  );

  const clearAll = useCallback(() => {
    setState({ predictions: [], champion: null, createdAt: Date.now() });
  }, []);

  const stats = {
    total: matches.length,
    filled: state.predictions.length,
    percentage: matches.length > 0
      ? Math.round((state.predictions.length / matches.length) * 100)
      : 0,
  };

  return {
    predictions: state.predictions,
    champion: state.champion,
    setPrediction,
    setChampion,
    getPrediction,
    clearAll,
    stats,
  };
}
