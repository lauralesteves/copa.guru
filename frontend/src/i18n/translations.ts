export type Locale = 'pt' | 'en' | 'es';

export interface AlternateSite {
  locale: Locale;
  flag: string;
  label: string;
}

export interface Translations {
  lang: string;
  siteTitle: string;
  metaDescription: string;
  nav: {
    groups: string;
    knockout: string;
    worldMap: string;
    teams: string;
    compare: string;
    predictions: string;
  };
  hero: {
    subtitle: string;
    viewGroups: string;
    makePredictions: string;
  };
  groups: {
    title: string;
    subtitle: string;
    tabTable: string;
    tabMatches: string;
    tabChart: string;
    tab3D: string;
    hintTable: string;
    hintMatches: string;
    hintChart: string;
    hint3D: string;
    tbd: string;
    host: string;
    debutant: string;
    noMatches: string;
  };
  standings: {
    team: string;
    pts: string;
    played: string;
    won: string;
    drawn: string;
    lost: string;
    gd: string;
  };
  knockout: {
    title: string;
    subtitle: string;
    bracket: string;
    graph: string;
    threeD: string;
    champion: string;
    tbd: string;
    roundOf16: string;
    quarters: string;
    semi: string;
    final: string;
  };
  globe: {
    title: string;
    subtitle: string;
    group: string;
  };
  teams: {
    title: string;
    subtitle: string;
    filterAll: string;
  };
  compare: {
    title: string;
    subtitle: string;
    vs: string;
    teams: string;
    groups: string;
    matches: string;
    stats: string[];
  };
  predictions: {
    title: string;
    subtitle: string;
    progress: string;
    instruction: string;
    goToGroups: string;
    whoWins: string;
    yourPick: string;
    change: string;
    share: string;
    clearAll: string;
    confirmClear: string;
  };
  shareCard: {
    title: string;
    champion: string;
    games: string;
    totalGoals: string;
    homeWins: string;
    draws: string;
    progress: string;
    copyText: string;
    close: string;
    shareLines: {
      header: string;
      champion: string;
      filled: string;
      goals: string;
      results: string;
      cta: string;
    };
  };
  footer: {
    hintDesktop: string;
    hintMobile: string;
    hintPrefix: string;
    privacy: string;
    terms: string;
    tagline: string;
  };
  legal: {
    back: string;
  };
  scrollTop: string;
  alternates: AlternateSite[];
}

export const translations: Record<Locale, Translations> = {
  pt: {
    lang: 'pt-BR',
    siteTitle: 'Copa.Guru',
    metaDescription: 'O jogo que poucos enxergam. Copa.Guru usa dados e inteligência para revelar o que realmente importa em cada partida.',
    nav: {
      groups: 'Grupos',
      knockout: 'Mata-Mata',
      worldMap: 'Mapa Mundial',
      teams: 'Seleções',
      compare: 'Comparar',
      predictions: 'Palpites',
    },
    hero: {
      subtitle: 'Copa do Mundo 2026',
      viewGroups: 'Ver Grupos',
      makePredictions: 'Fazer Palpites!',
    },
    groups: {
      title: 'FASE DE GRUPOS',
      subtitle: '12 grupos · 48 seleções · Insira seus palpites',
      tabTable: 'Tabela',
      tabMatches: 'Jogos',
      tabChart: 'Gráfico',
      tab3D: '3D',
      hintTable: 'classificação ao vivo',
      hintMatches: 'insira seus placares',
      hintChart: 'diagrama de confrontos',
      hint3D: 'visualização orbital',
      tbd: 'A definir',
      host: 'SEDE',
      debutant: 'NOVO',
      noMatches: 'Nenhum jogo disponível',
    },
    standings: {
      team: 'Seleção',
      pts: 'PG',
      played: 'J',
      won: 'V',
      drawn: 'E',
      lost: 'D',
      gd: 'SG',
    },
    knockout: {
      title: 'MATA-MATA',
      subtitle: 'Dos 16 avos até a grande final',
      bracket: 'Chave',
      graph: 'Grafo',
      threeD: '3D',
      champion: 'CAMPEÃO',
      tbd: 'A definir...',
      roundOf16: 'Oitavas',
      quarters: 'Quartas',
      semi: 'Semi',
      final: 'Final',
    },
    globe: {
      title: 'MAPA MUNDIAL',
      subtitle: 'Seleções classificadas ao redor do mundo',
      group: 'Grupo',
    },
    teams: {
      title: 'SELEÇÕES',
      subtitle: 'seleções classificadas',
      filterAll: 'Todas',
    },
    compare: {
      title: 'COMPARAR SELEÇÕES',
      subtitle: 'Radar de atributos lado a lado',
      vs: 'VS',
      teams: 'Seleções',
      groups: 'Grupos',
      matches: 'Jogos',
      stats: ['Ataque', 'Defesa', 'Posse', 'Passes', 'Finalizações', 'Faltas'],
    },
    predictions: {
      title: 'PALPITES',
      subtitle: 'Faça seus palpites e compartilhe',
      progress: 'Progresso dos palpites',
      instruction: 'e inserir placares na aba "Jogos" de cada grupo',
      goToGroups: 'Ir para Fase de Grupos',
      whoWins: 'QUEM SERÁ O CAMPEÃO?',
      yourPick: 'Seu palpite de campeão',
      change: 'Trocar',
      share: 'Compartilhar Palpites',
      clearAll: 'Limpar tudo',
      confirmClear: 'Confirmar reset?',
    },
    shareCard: {
      title: 'Meus Palpites 2026',
      champion: 'Campeão',
      games: 'Jogos',
      totalGoals: 'Gols totais',
      homeWins: 'Vit. Mandante',
      draws: 'Empates',
      progress: 'Progresso',
      copyText: 'Copiar texto',
      close: 'Fechar',
      shareLines: {
        header: 'Meus palpites Copa 2026 - Copa.Guru',
        champion: 'Campeão',
        filled: 'Jogos preenchidos',
        goals: 'Total de gols',
        results: 'Vitórias mandante: %h | Empates: %d | Vitórias visitante: %a',
        cta: 'Faça seus palpites em https://copa.guru',
      },
    },
    footer: {
      hintPrefix: 'Se eu fosse você, eu tentaria',
      hintDesktop: 'digitar o Konami Code (↑↑↓↓←→←→BA)',
      hintMobile: 'clicar 5x rápido no logo do menu',
      privacy: 'Política de Privacidade',
      terms: 'Termos',
      tagline: 'Copa do Mundo 2026 - EUA, Canadá e México',
    },
    legal: {
      back: '← Voltar para Copa.Guru',
    },
    scrollTop: 'Voltar ao topo',
    alternates: [
      { locale: 'pt', flag: '/images/flags/48/BRA.png', label: 'Português' },
      { locale: 'en', flag: '/images/flags/48/USA.png', label: 'English' },
      { locale: 'es', flag: '/images/flags/48/ESP.png', label: 'Español' },
    ],
  },

  en: {
    lang: 'en',
    siteTitle: 'Copa.Guru',
    metaDescription: 'The game few can see. Copa.Guru uses data and intelligence to reveal what truly matters in every match.',
    nav: {
      groups: 'Groups',
      knockout: 'Knockout',
      worldMap: 'World Map',
      teams: 'Teams',
      compare: 'Compare',
      predictions: 'Predictions',
    },
    hero: {
      subtitle: 'World Cup 2026',
      viewGroups: 'View Groups',
      makePredictions: 'Make Predictions!',
    },
    groups: {
      title: 'GROUP STAGE',
      subtitle: '12 groups · 48 teams · Enter your predictions',
      tabTable: 'Table',
      tabMatches: 'Matches',
      tabChart: 'Chart',
      tab3D: '3D',
      hintTable: 'live standings',
      hintMatches: 'enter your scores',
      hintChart: 'match diagram',
      hint3D: 'orbital view',
      tbd: 'TBD',
      host: 'HOST',
      debutant: 'NEW',
      noMatches: 'No matches available',
    },
    standings: {
      team: 'Team',
      pts: 'PTS',
      played: 'P',
      won: 'W',
      drawn: 'D',
      lost: 'L',
      gd: 'GD',
    },
    knockout: {
      title: 'KNOCKOUT',
      subtitle: 'From the Round of 16 to the Grand Final',
      bracket: 'Bracket',
      graph: 'Graph',
      threeD: '3D',
      champion: 'CHAMPION',
      tbd: 'TBD...',
      roundOf16: 'Round of 16',
      quarters: 'Quarters',
      semi: 'Semi',
      final: 'Final',
    },
    globe: {
      title: 'WORLD MAP',
      subtitle: 'Qualified teams around the world',
      group: 'Group',
    },
    teams: {
      title: 'TEAMS',
      subtitle: 'qualified teams',
      filterAll: 'All',
    },
    compare: {
      title: 'COMPARE TEAMS',
      subtitle: 'Side-by-side attribute radar',
      vs: 'VS',
      teams: 'Teams',
      groups: 'Groups',
      matches: 'Matches',
      stats: ['Attack', 'Defense', 'Possession', 'Passing', 'Shots', 'Fouls'],
    },
    predictions: {
      title: 'PREDICTIONS',
      subtitle: 'Make your predictions and share',
      progress: 'Prediction progress',
      instruction: 'and enter scores in the "Matches" tab of each group',
      goToGroups: 'Go to Group Stage',
      whoWins: 'WHO WILL BE THE CHAMPION?',
      yourPick: 'Your champion pick',
      change: 'Change',
      share: 'Share Predictions',
      clearAll: 'Clear all',
      confirmClear: 'Confirm reset?',
    },
    shareCard: {
      title: 'My Predictions 2026',
      champion: 'Champion',
      games: 'Matches',
      totalGoals: 'Total goals',
      homeWins: 'Home wins',
      draws: 'Draws',
      progress: 'Progress',
      copyText: 'Copy text',
      close: 'Close',
      shareLines: {
        header: 'My World Cup 2026 predictions - Copa.Guru',
        champion: 'Champion',
        filled: 'Matches filled',
        goals: 'Total goals',
        results: 'Home wins: %h | Draws: %d | Away wins: %a',
        cta: 'Make your predictions at https://copa.guru/en',
      },
    },
    footer: {
      hintPrefix: 'If I were you, I would try',
      hintDesktop: 'typing the Konami Code (↑↑↓↓←→←→BA)',
      hintMobile: 'tapping the menu logo 5 times quickly',
      privacy: 'Privacy Policy',
      terms: 'Terms',
      tagline: 'World Cup 2026 - USA, Canada & Mexico',
    },
    legal: {
      back: '← Back to Copa.Guru',
    },
    scrollTop: 'Back to top',
    alternates: [
      { locale: 'pt', flag: '/images/flags/48/BRA.png', label: 'Português' },
      { locale: 'en', flag: '/images/flags/48/USA.png', label: 'English' },
      { locale: 'es', flag: '/images/flags/48/ESP.png', label: 'Español' },
    ],
  },

  es: {
    lang: 'es',
    siteTitle: 'Copa.Guru',
    metaDescription: 'El juego que pocos ven. Copa.Guru usa datos e inteligencia para revelar lo que realmente importa en cada partido.',
    nav: {
      groups: 'Grupos',
      knockout: 'Eliminatorias',
      worldMap: 'Mapa Mundial',
      teams: 'Selecciones',
      compare: 'Comparar',
      predictions: 'Pronósticos',
    },
    hero: {
      subtitle: 'Copa del Mundo 2026',
      viewGroups: 'Ver Grupos',
      makePredictions: '¡Hacer Pronósticos!',
    },
    groups: {
      title: 'FASE DE GRUPOS',
      subtitle: '12 grupos · 48 selecciones · Ingresa tus pronósticos',
      tabTable: 'Tabla',
      tabMatches: 'Partidos',
      tabChart: 'Gráfico',
      tab3D: '3D',
      hintTable: 'clasificación en vivo',
      hintMatches: 'ingresa tus marcadores',
      hintChart: 'diagrama de enfrentamientos',
      hint3D: 'vista orbital',
      tbd: 'Por definir',
      host: 'SEDE',
      debutant: 'NUEVO',
      noMatches: 'No hay partidos disponibles',
    },
    standings: {
      team: 'Selección',
      pts: 'PTS',
      played: 'PJ',
      won: 'G',
      drawn: 'E',
      lost: 'P',
      gd: 'DG',
    },
    knockout: {
      title: 'ELIMINATORIAS',
      subtitle: 'De los 16avos hasta la gran final',
      bracket: 'Llave',
      graph: 'Grafo',
      threeD: '3D',
      champion: 'CAMPEÓN',
      tbd: 'Por definir...',
      roundOf16: 'Octavos',
      quarters: 'Cuartos',
      semi: 'Semi',
      final: 'Final',
    },
    globe: {
      title: 'MAPA MUNDIAL',
      subtitle: 'Selecciones clasificadas alrededor del mundo',
      group: 'Grupo',
    },
    teams: {
      title: 'SELECCIONES',
      subtitle: 'selecciones clasificadas',
      filterAll: 'Todas',
    },
    compare: {
      title: 'COMPARAR SELECCIONES',
      subtitle: 'Radar de atributos lado a lado',
      vs: 'VS',
      teams: 'Selecciones',
      groups: 'Grupos',
      matches: 'Partidos',
      stats: ['Ataque', 'Defensa', 'Posesión', 'Pases', 'Tiros', 'Faltas'],
    },
    predictions: {
      title: 'PRONÓSTICOS',
      subtitle: 'Haz tus pronósticos y comparte',
      progress: 'Progreso de pronósticos',
      instruction: 'e ingresa marcadores en la pestaña "Partidos" de cada grupo',
      goToGroups: 'Ir a Fase de Grupos',
      whoWins: '¿QUIÉN SERÁ EL CAMPEÓN?',
      yourPick: 'Tu pronóstico de campeón',
      change: 'Cambiar',
      share: 'Compartir Pronósticos',
      clearAll: 'Borrar todo',
      confirmClear: '¿Confirmar reset?',
    },
    shareCard: {
      title: 'Mis Pronósticos 2026',
      champion: 'Campeón',
      games: 'Partidos',
      totalGoals: 'Goles totales',
      homeWins: 'Vic. Local',
      draws: 'Empates',
      progress: 'Progreso',
      copyText: 'Copiar texto',
      close: 'Cerrar',
      shareLines: {
        header: 'Mis pronósticos Copa 2026 - Copa.Guru',
        champion: 'Campeón',
        filled: 'Partidos completados',
        goals: 'Total de goles',
        results: 'Victorias local: %h | Empates: %d | Victorias visitante: %a',
        cta: 'Haz tus pronósticos en https://copa.guru/es',
      },
    },
    footer: {
      hintPrefix: 'Si yo fuera tú, intentaría',
      hintDesktop: 'escribir el Código Konami (↑↑↓↓←→←→BA)',
      hintMobile: 'tocar el logo del menú 5 veces rápido',
      privacy: 'Política de Privacidad',
      terms: 'Términos',
      tagline: 'Copa del Mundo 2026 - EE.UU., Canadá y México',
    },
    legal: {
      back: '← Volver a Copa.Guru',
    },
    scrollTop: 'Volver arriba',
    alternates: [
      { locale: 'pt', flag: '/images/flags/48/BRA.png', label: 'Português' },
      { locale: 'en', flag: '/images/flags/48/USA.png', label: 'English' },
      { locale: 'es', flag: '/images/flags/48/ESP.png', label: 'Español' },
    ],
  },
};
