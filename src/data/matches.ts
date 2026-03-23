import type { GroupName, Match } from '../types/worldcup';
import { groups } from './groups';

// === DADOS FAKE PARA TESTE — REMOVER ANTES DO DEPLOY ===
// Resultados simulados para visualizar o site completo funcionando.
// Cada grupo tem 3 jogos (times validos), com placares verossimeis.
const FAKE_SCORES: Record<string, [number, number]> = {
  // Grupo A: MEX vs KOR vs RSA (TBD1 ignorado)
  'MEX-KOR': [2, 1],
  'MEX-RSA': [3, 0],
  'KOR-RSA': [1, 1],

  // Grupo B: CAN vs SUI vs QAT (TBD2 ignorado)
  'CAN-SUI': [0, 1],
  'CAN-QAT': [2, 0],
  'SUI-QAT': [3, 1],

  // Grupo C: BRA vs MAR vs SCO vs HAI
  'BRA-MAR': [2, 1],
  'BRA-SCO': [3, 0],
  'BRA-HAI': [5, 0],
  'MAR-SCO': [1, 0],
  'MAR-HAI': [4, 1],
  'SCO-HAI': [2, 2],

  // Grupo D: USA vs PAR vs AUS (TBD3 ignorado)
  'USA-PAR': [1, 0],
  'USA-AUS': [2, 2],
  'PAR-AUS': [0, 1],

  // Grupo E: GER vs ECU vs CIV vs CUW
  'GER-ECU': [2, 0],
  'GER-CIV': [1, 1],
  'GER-CUW': [4, 0],
  'ECU-CIV': [1, 2],
  'ECU-CUW': [3, 0],
  'CIV-CUW': [2, 0],

  // Grupo F: NED vs JPN vs TUN (TBD4 ignorado)
  'NED-JPN': [1, 2],
  'NED-TUN': [3, 0],
  'JPN-TUN': [2, 1],

  // Grupo G: BEL vs IRN vs EGY vs NZL
  'BEL-IRN': [3, 1],
  'BEL-EGY': [2, 0],
  'BEL-NZL': [4, 0],
  'IRN-EGY': [0, 0],
  'IRN-NZL': [2, 1],
  'EGY-NZL': [1, 0],

  // Grupo H: ESP vs URU vs KSA vs CPV
  'ESP-URU': [1, 1],
  'ESP-KSA': [3, 0],
  'ESP-CPV': [4, 0],
  'URU-KSA': [2, 0],
  'URU-CPV': [3, 1],
  'KSA-CPV': [1, 1],

  // Grupo I: FRA vs SEN vs NOR (TBD5 ignorado)
  'FRA-SEN': [2, 1],
  'FRA-NOR': [1, 0],
  'SEN-NOR': [1, 1],

  // Grupo J: ARG vs AUT vs ALG vs JOR
  'ARG-AUT': [2, 0],
  'ARG-ALG': [3, 1],
  'ARG-JOR': [4, 0],
  'AUT-ALG': [1, 1],
  'AUT-JOR': [2, 0],
  'ALG-JOR': [1, 0],

  // Grupo K: POR vs COL vs UZB (TBD6 ignorado)
  'POR-COL': [1, 1],
  'POR-UZB': [3, 0],
  'COL-UZB': [2, 1],

  // Grupo L: ENG vs CRO vs PAN vs GHA
  'ENG-CRO': [2, 1],
  'ENG-PAN': [3, 0],
  'ENG-GHA': [1, 0],
  'CRO-PAN': [2, 0],
  'CRO-GHA': [1, 1],
  'PAN-GHA': [0, 2],
};
// === FIM DADOS FAKE ===

function generateGroupMatches(
  groupName: GroupName,
  teams: string[],
  startId: number,
): Match[] {
  const valid = teams.filter((t) => !t.startsWith('TBD'));
  const matches: Match[] = [];
  let id = startId;

  for (let i = 0; i < valid.length; i++) {
    for (let j = i + 1; j < valid.length; j++) {
      const key = `${valid[i]}-${valid[j]}`;
      const fake = FAKE_SCORES[key];

      matches.push({
        id: id++,
        team1: valid[i],
        team2: valid[j],
        goals1: fake ? fake[0] : null,
        goals2: fake ? fake[1] : null,
        date: '',
        time: '',
        stadium: '',
        city: '',
        group: groupName,
        stage: 'group',
      });
    }
  }

  return matches;
}

export function createInitialMatches(): Match[] {
  const all: Match[] = [];
  let nextId = 1;

  for (const group of groups) {
    const gm = generateGroupMatches(group.name, group.teams, nextId);
    all.push(...gm);
    nextId += gm.length;
  }

  return all;
}
