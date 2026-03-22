# Copa.Guru - Rebuild do Site para Copa do Mundo 2026

## Contexto

O Copa.Guru original foi criado em 2014 pela Viking Labs para a Copa do Mundo no Brasil. Era uma plataforma interativa de visualizacao de dados esportivos construida com AngularJS + D3.js, que permitia simular resultados e acompanhar a fase de grupos e mata-mata com diagramas de corda (chord) e grafos de forca (molecule).

O novo site sera uma reconstrucao completa usando stack moderna, mantendo o espirito de interatividade e visualizacao de dados do original, mas com animacoes e experiencia de usuario muito superiores.

---

## Stack Tecnica

Baseado na estrutura do projeto `quecoreagora.lauraesteves.com`:

- **Framework:** React 19 + TypeScript
- **Build:** Vite 6 + SWC
- **Styling:** Tailwind CSS 4
- **Animacoes:** GSAP (ScrollTrigger, timeline) + Motion (Framer Motion) para React
- **Visualizacao de Dados:** D3.js v7 (evolucao do D3 v3 usado no original)
- **Linting:** Biome
- **Testes:** Vitest
- **Package Manager:** Bun
- **Deploy:** AWS S3 + CloudFront (mesmo pipeline do quecoreagora)

---

## Estrutura do Projeto

Copiar padrao do `quecoreagora.lauraesteves.com`:

```
copa.guru/
├── .ai/                           # Issues e documentacao de projeto
├── .claude/                       # Config Claude
├── .github/                       # Workflows
├── public/
│   ├── images/
│   │   ├── flags/                 # Logos/bandeiras das selecoes (SVG)
│   │   │   ├── 96/               # 96px (nodes grandes)
│   │   │   └── 48/               # 48px (diagramas, tabelas)
│   │   ├── og/                    # OpenGraph images
│   │   └── icons/                 # Icones do site
│   ├── data/
│   │   └── worldcup2026.json     # Dados da copa (selecoes, jogos, grupos, estadios)
│   ├── favicon.ico
│   ├── site.webmanifest
│   ├── robots.txt
│   └── ai.txt
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── index.css
│   ├── assets/
│   │   ├── fonts/
│   │   └── icons/
│   ├── components/
│   │   ├── Navbar/
│   │   │   ├── Navbar.tsx
│   │   │   ├── NavLinks.tsx
│   │   │   ├── MobileMenu.tsx
│   │   │   └── SocialIcons.tsx
│   │   ├── Footer/
│   │   │   └── Footer.tsx
│   │   ├── Hero/
│   │   │   └── Hero.tsx           # Secao de abertura com animacao
│   │   ├── Groups/
│   │   │   ├── GroupStage.tsx     # Container da fase de grupos
│   │   │   ├── GroupCard.tsx      # Card individual de cada grupo
│   │   │   ├── ChordDiagram.tsx   # Diagrama de corda D3 (legado reimaginado)
│   │   │   ├── StandingsTable.tsx # Tabela de classificacao
│   │   │   └── MatchCard.tsx      # Card de jogo individual
│   │   ├── Knockout/
│   │   │   ├── KnockoutStage.tsx  # Container do mata-mata
│   │   │   ├── Bracket.tsx        # Chave do torneio (evolucao do molecule)
│   │   │   └── MatchNode.tsx      # Node de jogo no bracket
│   │   ├── Teams/
│   │   │   ├── TeamGrid.tsx       # Grid de todas as selecoes
│   │   │   ├── TeamCard.tsx       # Card de selecao com bandeira
│   │   │   └── TeamDetail.tsx     # Modal/pagina de detalhe
│   │   ├── Predictions/
│   │   │   ├── PredictionForm.tsx # Formulario de palpites
│   │   │   └── ShareCard.tsx      # Card compartilhavel de palpites
│   │   ├── Stats/
│   │   │   ├── StatsSection.tsx   # Secao de estatisticas
│   │   │   └── AnimatedCounter.tsx
│   │   └── ui/
│   │       ├── ExternalLink.tsx
│   │       ├── Flag.tsx           # Componente reutilizavel de bandeira
│   │       ├── ScoreInput.tsx     # Input de placar (legado reimaginado)
│   │       └── Tooltip.tsx
│   ├── hooks/
│   │   ├── useWorldCupData.ts     # Hook para dados da copa
│   │   ├── useGroupStandings.ts   # Calculo de classificacao
│   │   ├── useAnimations.ts       # Hook wrapper GSAP
│   │   └── usePredictions.ts      # Gerenciamento de palpites (localStorage)
│   ├── lib/
│   │   ├── calculations.ts        # Logica de classificacao e desempate
│   │   ├── animations.ts          # Configuracoes de animacao reutilizaveis
│   │   └── constants.ts           # Cores, configuracoes, enums
│   ├── types/
│   │   └── worldcup.ts            # Tipos TypeScript (Team, Match, Group, etc.)
│   └── data/
│       └── teams.ts               # Dados estaticos das selecoes
├── scripts/
│   └── pre-push                   # Hook de pre-push
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
├── biome.json
├── Makefile
└── bun.lock
```

---

## Selecoes Classificadas - Copa 2026

### 42 Confirmadas + 6 via Playoff (26-31 de Marco 2026)

#### CONCACAF (6)
| Selecao | Codigo | Status |
|---------|--------|--------|
| Canada | CAN | Sede |
| Mexico | MEX | Sede |
| Estados Unidos | USA | Sede |
| Haiti | HAI | Classificado |
| Panama | PAN | Classificado |
| Curacao | CUW | Classificado |

#### AFC - Asia (8)
| Selecao | Codigo | Status |
|---------|--------|--------|
| Australia | AUS | Classificado |
| Ira | IRN | Classificado |
| Japao | JPN | Classificado |
| Jordania | JOR | Estreante |
| Catar | QAT | Classificado |
| Arabia Saudita | KSA | Classificado |
| Coreia do Sul | KOR | Classificado |
| Uzbequistao | UZB | Estreante |

#### CAF - Africa (9)
| Selecao | Codigo | Status |
|---------|--------|--------|
| Argelia | ALG | Classificado |
| Cabo Verde | CPV | Estreante |
| Egito | EGY | Classificado |
| Gana | GHA | Classificado |
| Costa do Marfim | CIV | Classificado |
| Marrocos | MAR | Classificado |
| Senegal | SEN | Classificado |
| Africa do Sul | RSA | Classificado |
| Tunisia | TUN | Classificado |

#### CONMEBOL (6)
| Selecao | Codigo | Status |
|---------|--------|--------|
| Argentina | ARG | Classificado |
| Brasil | BRA | Classificado |
| Colombia | COL | Classificado |
| Equador | ECU | Classificado |
| Paraguai | PAR | Classificado |
| Uruguai | URU | Classificado |

#### OFC - Oceania (1)
| Selecao | Codigo | Status |
|---------|--------|--------|
| Nova Zelandia | NZL | Classificado |

#### UEFA - Europa (12 + 4 via playoff)
| Selecao | Codigo | Status |
|---------|--------|--------|
| Austria | AUT | Classificado |
| Belgica | BEL | Classificado |
| Croacia | CRO | Classificado |
| Inglaterra | ENG | Classificado |
| Franca | FRA | Classificado |
| Alemanha | GER | Classificado |
| Holanda | NED | Classificado |
| Noruega | NOR | Classificado |
| Portugal | POR | Classificado |
| Escocia | SCO | Classificado |
| Espanha | ESP | Classificado |
| Suica | SUI | Classificado |
| Playoff A | TBD | Italia/Irlanda do Norte vs Pais de Gales/Bosnia → Grupo B |
| Playoff B | TBD | Ucrania/Suecia vs Polonia/Albania → Grupo F |
| Playoff C | TBD | Eslovaquia/Kosovo vs Turquia/Romenia → Grupo D |
| Playoff D | TBD | Tchequi/Irlanda vs Dinamarca/Macedonia do Norte → Grupo A |

#### Playoffs Intercontinentais (2 vagas)
| Chave | Semifinal (26/03) | Final (31/03) | Grupo |
|-------|-------------------|---------------|-------|
| 1 | Nova Caledonia vs Jamaica | Vencedor vs RD Congo | K |
| 2 | Bolivia vs Suriname | Vencedor vs Iraque | I |

---

## Grupos da Copa 2026

| Grupo | Pot 1 | Pot 2 | Pot 3 | Pot 4 |
|-------|-------|-------|-------|-------|
| **A** | Mexico | Coreia do Sul | Africa do Sul | UEFA Playoff D |
| **B** | Canada | Suica | Catar | UEFA Playoff A |
| **C** | Brasil | Marrocos | Escocia | Haiti |
| **D** | Estados Unidos | Paraguai | Australia | UEFA Playoff C |
| **E** | Alemanha | Equador | Costa do Marfim | Curacao |
| **F** | Holanda | Japao | Tunisia | UEFA Playoff B |
| **G** | Belgica | Ira | Egito | Nova Zelandia |
| **H** | Espanha | Uruguai | Arabia Saudita | Cabo Verde |
| **I** | Franca | Senegal | Noruega | Intercon. Playoff 2 |
| **J** | Argentina | Austria | Argelia | Jordania |
| **K** | Portugal | Colombia | Uzbequistao | Intercon. Playoff 1 |
| **L** | Inglaterra | Croacia | Panama | Gana |

---

## Animacoes

### Animacoes do Projeto Original (2014)

O site original possuia as seguintes animacoes:

1. **Sidebar Slide** - CSS transition 0.4s ease para menu mobile
2. **Chord Diagram Transitions** - D3.js opacity transition 750ms ao atualizar resultados
3. **Chord Hover** - Opacity 0.2 → 0.8 com highlight e reordenacao de z-index
4. **Force-Directed Graph** - D3 force layout com simulacao fisica continua, charge -200, drag interativo
5. **Tooltips D3** - Show/hide com offset -10px no mouseover
6. **Table Row Hover** - Background-color com cor do time a 30% opacity
7. **SVG Drop Shadows** - feGaussianBlur filter nos logos
8. **Smooth Scrolling** - Angular smooth-scroll directive
9. **Loading Spinner** - Angular spinner com spin.js
10. **Easter Egg Konami Code** (↑↑↓↓←→←→BA):
    - Personagem "Allejo" sobe do bottom (jQuery animate)
    - Bounce de 100ms (-10px)
    - Travessia horizontal 2200ms
    - Shake de gol (jQuery UI, 7x, 10px, 100ms)
    - Background fadeIn 100ms / fadeOut 2200ms

### Novas Animacoes Propostas

#### Hero / Abertura
- **Parallax de entrada** - Background do estadio com parallax multi-camada (GSAP ScrollTrigger)
- **Logo reveal** - Copa.Guru logo aparece com morphing SVG ou text-reveal animado
- **Countdown timer** - Contagem regressiva animada para o inicio da copa (se antes do torneio)
- **Particle field** - Particulas flutuantes nas cores da copa que reagem ao mouse (Canvas/WebGL)

#### Fase de Grupos
- **Scroll-triggered group cards** - Cards dos grupos entram com stagger animation conforme scroll (GSAP ScrollTrigger + stagger)
- **Chord diagram 2.0** - Reimaginacao do chord diagram original com:
  - Animacao de entrada: arcos crescem do centro para fora
  - Transicao suave entre estados ao mudar placar
  - Glow effect nas cordas ao hover (SVG filter + CSS)
  - Pulsacao sutil nos arcos da selecao selecionada
- **Tabela animada** - Linhas da tabela reorganizam com flip animation (FLIP technique) ao mudar classificacao
- **Score input micro-interactions** - Numeros do placar sobem/descem com rolling number animation
- **Flag wave** - Bandeiras com efeito sutil de ondulacao CSS (transform skew + animation)

#### Mata-Mata / Bracket
- **Bracket reveal** - Chave do torneio se constroi progressivamente (path drawing animation com SVG stroke-dasharray)
- **Match connection lines** - Linhas conectando jogos com efeito de "fluxo de energia" (animated gradient along path)
- **Winner celebration** - Ao definir vencedor: confetti burst (canvas particles), glow pulse no node, shake sutil
- **Bracket zoom** - Transicao suave de zoom in/out nas fases do bracket (GSAP transform)
- **Team path highlight** - Ao clicar numa selecao, ilumina todo o caminho dela no bracket com trail animation

#### Selecoes / Teams
- **Team grid entrance** - Cards entram em wave/stagger por confederacao (GSAP stagger + ScrollTrigger)
- **Card flip** - Card de selecao vira (3D CSS transform perspective) para mostrar detalhes no verso
- **Flag morph** - Transicao suave entre bandeiras ao navegar (crossfade com scale)
- **Hover lift** - Cards levantam com sombra dinamica no hover (transform translateY + box-shadow transition)

#### Palpites / Predictions
- **Share card generation** - Card de palpites animado para compartilhar (html2canvas + animacao de "carimbo")
- **Prediction lock** - Animacao de "cadeado fechando" ao confirmar palpite
- **Score reveal** - Numeros do placar aparecem tipo slot machine (vertical scroll de digitos)

#### Estatisticas
- **Animated counters** - Numeros sobem de 0 ao valor final com easing (countUp animation)
- **Chart transitions** - Graficos crescem/morpham suavemente ao trocar visualizacao
- **Data pulse** - Dados atualizados piscam sutilmente com glow animation

#### Micro-interacoes Globais
- **Page transitions** - Transicao entre secoes com crossfade + slide (Motion AnimatePresence)
- **Scroll progress** - Barra de progresso no topo que acompanha o scroll
- **Loading skeleton** - Skeleton screens animados durante carregamento
- **Toast notifications** - Notificacoes entram/saem com slide + fade
- **Button ripple** - Efeito material ripple nos botoes
- **Cursor trail** - Sutil trail effect no cursor em areas interativas

#### Easter Eggs (Tradicao do Original)
- **Konami Code 2026** - Novo easter egg com referencia a Copa 2026 (reimaginar o Allejo)
- **Goal celebration** - Ao inserir um placar de goleada (5+), animacao especial de celebracao
- **Hidden team facts** - Clique longo em bandeira revela fato curioso com animacao de "cartao"

---

## Novas Ideias para Exibicao de Resultados

### 1. Timeline Interativa
Em vez de mostrar todos os jogos de uma vez, uma **timeline horizontal scrollavel** onde cada dia da copa e um ponto. Ao clicar/scrollar, os jogos daquele dia se expandem com animacao. Inspirado em timelines de redes sociais.

### 2. Mapa Interativo das Sedes
Mapa dos EUA/Canada/Mexico com os **estadios marcados**. Ao clicar num estadio, mostra os jogos que acontecem ali com animacao de "zoom into stadium". Os jogos em andamento pulsam no mapa.

### 3. Match Cards Estilo Scoreboard
Cards individuais por jogo com:
- Bandeiras grandes das selecoes
- Placar central grande e animado
- Timeline de gols (minuto a minuto) como dots numa linha
- Status do jogo com badge animado (AO VIVO pulsando, ENCERRADO com checkmark)

### 4. Group Galaxy View
Visualizacao alternativa dos grupos como um "sistema solar" onde:
- Cada grupo e uma galaxia/cluster
- Times sao planetas orbitando
- Tamanho do planeta = pontos acumulados
- Proximidade ao centro = posicao na tabela
- Evolui conforme os resultados mudam

### 5. Head-to-Head Comparator
Ferramenta para comparar duas selecoes lado a lado:
- Estatisticas em barras animadas que "duelam"
- Historico de confrontos
- Radar chart com atributos (ataque, defesa, posse, etc.)

### 6. Prediction Mode Gamificado
- Usuario faz palpites para todos os jogos
- Sistema calcula pontuacao conforme resultados reais
- Ranking compartilhavel
- Badge/achievements por acertos (ex: "Vidente" se acertar placar exato)
- Share card personalizado para redes sociais

### 7. Match Story
Cada jogo tem uma "historia" contada em cards verticais (estilo Instagram Stories):
- Card 1: Apresentacao dos times
- Card 2: Escalacoes
- Card 3: Momentos-chave (gols, cartoes)
- Card 4: Resultado final com stats
- Navegacao por swipe/tap

### 8. Live Pulse Dashboard
Durante os jogos, um dashboard que mostra:
- Todos os jogos simultaneos lado a lado
- Indicador visual de "intensidade" do jogo (baseado em gols recentes)
- Notificacao animada quando sai gol em qualquer jogo
- Mini-mapa de posicao no bracket/grupo

### 9. Champion Path Retrospective
Apos a copa, uma animacao que mostra o caminho do campeao:
- Cada jogo como um "capitulo"
- Scroll vertical com parallax
- Estatisticas acumulativas que crescem
- Momento decisivo de cada jogo destacado

### 10. Social Share Cards
Cards bonitos e animados para compartilhar:
- Resultado de um jogo especifico
- Classificacao de um grupo
- Palpite do usuario
- Caminho de uma selecao no torneio
- Gerados como imagem/video curto

---

## Dados Necessarios (worldcup2026.json)

```typescript
interface Team {
  name: string;               // "Brasil"
  code: string;               // "BRA"
  confederation: Confederation;
  group: Group;
  flag96: string;             // path para logo 96px
  flag48: string;             // path para logo 48px
  primaryColor: string;       // "#FFDF00"
  secondaryColor: string;     // "#009739"
  isHost: boolean;
  isDebutant: boolean;
  fifaRanking: number;
}

interface Match {
  id: number;
  team1: string;              // code
  team2: string;              // code
  goals1: number | null;
  goals2: number | null;
  date: string;               // "2026-06-11"
  time: string;               // "17:00"
  stadium: string;
  city: string;
  group: Group | null;
  stage: 'group' | 'round32' | 'round16' | 'quarter' | 'semi' | 'third' | 'final';
}

interface Group {
  name: string;               // "A" ... "L"
  teams: string[];            // team codes
  matches: number[];          // match ids
}

type Confederation = 'AFC' | 'CAF' | 'CONCACAF' | 'CONMEBOL' | 'OFC' | 'UEFA';
```

---

## Bandeiras / Logos das Selecoes

### Task: Baixar logos de todas as 48 selecoes

Precisamos dos logos/bandeiras em dois tamanhos (48px e 96px) para as 42 selecoes confirmadas. Os 6 times de playoff serao adicionados apos definicao (31 de marco de 2026).

**Formato preferido:** SVG (escalavel) com fallback PNG
**Fonte sugerida:** API de bandeiras (flagcdn.com, countryflagsapi.com) ou FIFA assets

**Lista de arquivos necessarios (42 confirmados):**
```
flags/ARG.svg  flags/ALG.svg  flags/AUS.svg  flags/AUT.svg
flags/BEL.svg  flags/BRA.svg  flags/CAN.svg  flags/CPV.svg
flags/CIV.svg  flags/COL.svg  flags/CRO.svg  flags/CUW.svg
flags/ECU.svg  flags/EGY.svg  flags/ENG.svg  flags/FRA.svg
flags/GER.svg  flags/GHA.svg  flags/HAI.svg  flags/IRN.svg
flags/JOR.svg  flags/JPN.svg  flags/KOR.svg  flags/KSA.svg
flags/MAR.svg  flags/MEX.svg  flags/NED.svg  flags/NOR.svg
flags/NZL.svg  flags/PAN.svg  flags/PAR.svg  flags/POR.svg
flags/QAT.svg  flags/RSA.svg  flags/SCO.svg  flags/SEN.svg
flags/ESP.svg  flags/SUI.svg  flags/TUN.svg  flags/URU.svg
flags/USA.svg  flags/UZB.svg
```

**Pendentes (playoff - definicao ate 31/03/2026):**
```
# UEFA Playoffs
flags/ITA.svg OR flags/NIR.svg OR flags/WAL.svg OR flags/BIH.svg
flags/UKR.svg OR flags/SWE.svg OR flags/POL.svg OR flags/ALB.svg
flags/SVK.svg OR flags/KOS.svg OR flags/TUR.svg OR flags/ROU.svg
flags/CZE.svg OR flags/IRL.svg OR flags/DEN.svg OR flags/MKD.svg

# Intercontinental Playoffs
flags/COD.svg OR flags/JAM.svg OR flags/NCL.svg
flags/IRQ.svg OR flags/BOL.svg OR flags/SUR.svg
```

---

## Navbar e Footer

Copiar estrutura do `quecoreagora.lauraesteves.com` e adaptar:

### Navbar
- Logo "Copa.Guru" no canto esquerdo (fonte customizada)
- Links de navegacao: Grupos | Mata-Mata | Selecoes | Palpites | Stats
- Social icons (mesmo padrao do quecoreagora)
- Menu mobile hamburger com dropdown animado (max-h transition 300ms)
- Fixed top com backdrop-blur
- Scroll behavior: navbar fica mais compacta ao scrollar (GSAP ScrollTrigger)

### Footer
- Fixed bottom com fundo semi-transparente
- "Copa.Guru 2026 - Todos os dados da Copa do Mundo"
- Copyright + creditos Laura Esteves
- Links para redes sociais

---

## Milestones

### M1 - Setup & Estrutura ✅
- [x] Inicializar projeto (Vite + React + TS + Tailwind + Biome)
- [x] Copiar estrutura de Navbar/Footer do quecoreagora
- [x] Configurar build pipeline (Makefile, scripts)
- [x] Configurar deploy S3 + CloudFront
- [x] Criar tipos TypeScript (worldcup.ts)
- [x] Montar dados dos grupos (src/data/groups.ts, src/data/teams.ts)
- [x] Hero section com animacao de entrada e scroll indicator
- [x] Scroll progress bar global
- [x] Navbar com comportamento de scroll (compacta ao rolar)

### M2 - Bandeiras & Assets ✅
- [x] Baixar bandeiras SVG de todas as 42 selecoes confirmadas + 22 playoff
- [x] Gerar versoes 48px e 96px (PNG fallback)
- [x] Criar componente Flag reutilizavel (src/components/ui/Flag.tsx)
- [x] Definir cores primarias/secundarias de cada selecao (src/data/teams.ts)

### M3 - Fase de Grupos ✅
- [x] Implementar GroupStage container com scroll-triggered stagger (GSAP)
- [x] Implementar GroupCard com abas (Tabela/Jogos/Grafico) e listagem com badges
- [x] Implementar StandingsTable reativa (atualiza ao vivo com cores, saldo +/-, highlight classificados)
- [x] Logica de calculo de classificacao e desempate (src/lib/calculations.ts)
- [x] Implementar ChordDiagram SVG (arcos de pontos por time, chords gradiente entre confrontos, hover interativo, labels de placar)
- [x] Implementar MatchCard com score input interativo (botoes +/-, scale animation, highlight vencedor)
- [x] Scroll-triggered animations com GSAP ScrollTrigger (titulos, subtitulos, grid stagger)
- [x] Hook useWorldCupData para gerenciar estado global de partidas e standings
- [x] Geracao automatica de partidas por grupo (src/data/matches.ts)
- [x] FLIP animation na reordenacao da tabela (GSAP Flip plugin, anima rows ao mudar posicao)

### M3.5 - Selecoes ✅
- [x] TeamGrid com filtro por confederacao + GSAP stagger animation na troca de filtro
- [x] TeamCard com flip animation 3D (frente: bandeira/nome, verso: detalhes)

### M4 - Mata-Mata ✅
- [x] Bracket SVG completo: R32 (16 jogos) → Oitavas (8) → Quartas (4) → Semi (2) → Final
- [x] Layout em duas metades (top/bottom) com bracket de chaves convergente
- [x] Path drawing animation nas conexoes (stroke-dasharray, ScrollTrigger, stagger 0.04s)
- [x] MatchNode interativo com bandeiras, nomes, highlight dourado
- [x] Team path highlight (hover numa selecao ilumina todos os seus nodes)
- [x] Label "CAMPEAO" acima da final
- [x] Curvas bezier conectando rodadas com gradiente dourado
- [x] Auto-populate R32 com 1o e 2o de cada grupo (reativo aos palpites)
- [x] Confetti celebration animation (canvas particles, 150 particulas, gravidade, fade out)

### M5 - Palpites & Social ✅
- [x] Sistema de palpites com localStorage (usePredictions hook, persist/restore)
- [x] Secao de palpites com barra de progresso (jogos preenchidos/total)
- [x] Seletor de campeao (grid de 42 selecoes, confetti ao selecionar)
- [x] ShareCard modal (resumo visual: campeao, stats, gols totais, copiar texto)
- [x] Botao limpar tudo com confirmacao de seguranca (timeout 3s)
- [x] Confetti burst reutilizavel (canvas, 150 particulas, cores copa, gravidade, fade)

### M6 - Hero & Polish ✅
- [x] Hero section com GSAP entrance timeline (stagger children)
- [x] Parallax multi-camada (bg, circulos decorativos, conteudo movem em velocidades diferentes via ScrollTrigger scrub)
- [x] Floating flags decorativas (8 bandeiras com gsap.to random y/x/rotation, yoyo infinito)
- [x] Scroll progress bar
- [x] SEO e meta tags
- [x] Easter egg Konami Code (↑↑↓↓←→←→BA → overlay "GOOOL!" com GSAP timeline + confetti)
- [x] Loading skeletons (SectionSkeleton com shimmer animation)
- [x] Code splitting via React.lazy + Suspense (GroupStage 40KB, Knockout 7KB, Predictions 7KB, TeamGrid 4KB)
- [x] will-change-transform nas camadas parallax para GPU acceleration

### M7 - Dados ao Vivo (pos-inicio da copa)
- [ ] Integracao com API de resultados ao vivo (se disponivel)
- [ ] Live pulse dashboard
- [ ] Notificacoes de gol
- [ ] Atualizacao automatica de classificacao
