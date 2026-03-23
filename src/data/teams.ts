import type { Team } from '../types/worldcup';

export const teams: Team[] = [
  // CONCACAF
  { name: 'Canadá', code: 'CAN', confederation: 'CONCACAF', group: 'B', primaryColor: '#FF0000', secondaryColor: '#FFFFFF', isHost: true, isDebutant: false },
  { name: 'México', code: 'MEX', confederation: 'CONCACAF', group: 'A', primaryColor: '#006847', secondaryColor: '#CE1126', isHost: true, isDebutant: false },
  { name: 'Estados Unidos', code: 'USA', confederation: 'CONCACAF', group: 'D', primaryColor: '#002868', secondaryColor: '#BF0A30', isHost: true, isDebutant: false },
  { name: 'Haiti', code: 'HAI', confederation: 'CONCACAF', group: 'C', primaryColor: '#00209F', secondaryColor: '#D21034', isHost: false, isDebutant: false },
  { name: 'Panamá', code: 'PAN', confederation: 'CONCACAF', group: 'L', primaryColor: '#DA121A', secondaryColor: '#003DA5', isHost: false, isDebutant: false },
  { name: 'Curaçao', code: 'CUW', confederation: 'CONCACAF', group: 'E', primaryColor: '#002B7F', secondaryColor: '#F9E814', isHost: false, isDebutant: true },

  // AFC
  { name: 'Austrália', code: 'AUS', confederation: 'AFC', group: 'D', primaryColor: '#FFCD00', secondaryColor: '#008751', isHost: false, isDebutant: false },
  { name: 'Irã', code: 'IRN', confederation: 'AFC', group: 'G', primaryColor: '#239F40', secondaryColor: '#DA0000', isHost: false, isDebutant: false },
  { name: 'Japão', code: 'JPN', confederation: 'AFC', group: 'F', primaryColor: '#000080', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Jordânia', code: 'JOR', confederation: 'AFC', group: 'J', primaryColor: '#007A33', secondaryColor: '#CE1126', isHost: false, isDebutant: true },
  { name: 'Catar', code: 'QAT', confederation: 'AFC', group: 'B', primaryColor: '#8A1538', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Arábia Saudita', code: 'KSA', confederation: 'AFC', group: 'H', primaryColor: '#006C35', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Coréia do Sul', code: 'KOR', confederation: 'AFC', group: 'A', primaryColor: '#CD2E3A', secondaryColor: '#0047A0', isHost: false, isDebutant: false },
  { name: 'Uzbequistão', code: 'UZB', confederation: 'AFC', group: 'K', primaryColor: '#1EB53A', secondaryColor: '#0099B5', isHost: false, isDebutant: true },

  // CAF
  { name: 'Argélia', code: 'ALG', confederation: 'CAF', group: 'J', primaryColor: '#006233', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Cabo Verde', code: 'CPV', confederation: 'CAF', group: 'H', primaryColor: '#003893', secondaryColor: '#CF2027', isHost: false, isDebutant: true },
  { name: 'Egito', code: 'EGY', confederation: 'CAF', group: 'G', primaryColor: '#CE1126', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Gana', code: 'GHA', confederation: 'CAF', group: 'L', primaryColor: '#006B3F', secondaryColor: '#FCD116', isHost: false, isDebutant: false },
  { name: 'Costa do Marfim', code: 'CIV', confederation: 'CAF', group: 'E', primaryColor: '#FF8200', secondaryColor: '#009A44', isHost: false, isDebutant: false },
  { name: 'Marrocos', code: 'MAR', confederation: 'CAF', group: 'C', primaryColor: '#C1272D', secondaryColor: '#006233', isHost: false, isDebutant: false },
  { name: 'Senegal', code: 'SEN', confederation: 'CAF', group: 'I', primaryColor: '#00853F', secondaryColor: '#FDEF42', isHost: false, isDebutant: false },
  { name: 'África do Sul', code: 'RSA', confederation: 'CAF', group: 'A', primaryColor: '#007749', secondaryColor: '#FFB81C', isHost: false, isDebutant: false },
  { name: 'Tunísia', code: 'TUN', confederation: 'CAF', group: 'F', primaryColor: '#E70013', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },

  // CONMEBOL
  { name: 'Argentina', code: 'ARG', confederation: 'CONMEBOL', group: 'J', primaryColor: '#75AADB', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Brasil', code: 'BRA', confederation: 'CONMEBOL', group: 'C', primaryColor: '#FFDF00', secondaryColor: '#009739', isHost: false, isDebutant: false },
  { name: 'Colômbia', code: 'COL', confederation: 'CONMEBOL', group: 'K', primaryColor: '#FCD116', secondaryColor: '#003893', isHost: false, isDebutant: false },
  { name: 'Equador', code: 'ECU', confederation: 'CONMEBOL', group: 'E', primaryColor: '#FFD100', secondaryColor: '#0038A8', isHost: false, isDebutant: false },
  { name: 'Paraguai', code: 'PAR', confederation: 'CONMEBOL', group: 'D', primaryColor: '#DA121A', secondaryColor: '#0038A8', isHost: false, isDebutant: false },
  { name: 'Uruguai', code: 'URU', confederation: 'CONMEBOL', group: 'H', primaryColor: '#5CBFEB', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },

  // OFC
  { name: 'Nova Zelândia', code: 'NZL', confederation: 'OFC', group: 'G', primaryColor: '#000000', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },

  // UEFA
  { name: 'Áustria', code: 'AUT', confederation: 'UEFA', group: 'J', primaryColor: '#ED2939', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Bélgica', code: 'BEL', confederation: 'UEFA', group: 'G', primaryColor: '#ED2939', secondaryColor: '#000000', isHost: false, isDebutant: false },
  { name: 'Croácia', code: 'CRO', confederation: 'UEFA', group: 'L', primaryColor: '#FF0000', secondaryColor: '#0000FF', isHost: false, isDebutant: false },
  { name: 'Inglaterra', code: 'ENG', confederation: 'UEFA', group: 'L', primaryColor: '#FFFFFF', secondaryColor: '#CF081F', isHost: false, isDebutant: false },
  { name: 'França', code: 'FRA', confederation: 'UEFA', group: 'I', primaryColor: '#002395', secondaryColor: '#ED2939', isHost: false, isDebutant: false },
  { name: 'Alemanha', code: 'GER', confederation: 'UEFA', group: 'E', primaryColor: '#000000', secondaryColor: '#DD0000', isHost: false, isDebutant: false },
  { name: 'Holanda', code: 'NED', confederation: 'UEFA', group: 'F', primaryColor: '#FF6600', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Noruega', code: 'NOR', confederation: 'UEFA', group: 'I', primaryColor: '#BA0C2F', secondaryColor: '#00205B', isHost: false, isDebutant: false },
  { name: 'Portugal', code: 'POR', confederation: 'UEFA', group: 'K', primaryColor: '#006600', secondaryColor: '#FF0000', isHost: false, isDebutant: false },
  { name: 'Escócia', code: 'SCO', confederation: 'UEFA', group: 'C', primaryColor: '#003399', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
  { name: 'Espanha', code: 'ESP', confederation: 'UEFA', group: 'H', primaryColor: '#AA151B', secondaryColor: '#F1BF00', isHost: false, isDebutant: false },
  { name: 'Suíça', code: 'SUI', confederation: 'UEFA', group: 'B', primaryColor: '#FF0000', secondaryColor: '#FFFFFF', isHost: false, isDebutant: false },
];

export function getTeamByCode(code: string): Team | undefined {
  return teams.find((t) => t.code === code);
}

export function getTeamsByGroup(group: string): Team[] {
  return teams.filter((t) => t.group === group);
}

export function getTeamsByConfederation(conf: string): Team[] {
  return teams.filter((t) => t.confederation === conf);
}
