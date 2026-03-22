# Copa.Guru 2026 - Plano de Graficos, Visualizacoes e Animacoes 3D

## Contexto

O site original (2014) usava D3.js v3 com dois graficos principais:
1. **Chord Diagram** - Diagrama circular de corda mostrando confrontos e pontos por grupo
2. **Molecule Graph** - Grafo de forca (force-directed) para o mata-mata com simulacao fisica

O site atual (2026) tem versoes simplificadas em SVG puro. Este plano detalha a
reimplementacao completa com visualizacoes modernas, incluindo 3D com Three.js/React Three Fiber.

---

## Stack de Visualizacao

| Lib | Uso | Bundle |
|-----|-----|--------|
| `@react-three/fiber` | Renderer React para Three.js (cenas 3D) | ~150KB |
| `@react-three/drei` | Helpers (OrbitControls, Text3D, Float, etc) | ~80KB |
| `react-globe.gl` | Globo 3D interativo com selecoes | ~200KB |
| `d3` (modulos) | d3-chord, d3-force, d3-scale, d3-arc | ~30KB (tree-shaken) |
| `gsap` (ja instalado) | ScrollTrigger, Flip, timeline | ja incluso |
| `three` (peer dep) | Engine 3D, shaders, post-processing | ~600KB (tree-shaken) |

Todas as visualizacoes 3D serao lazy-loaded para nao impactar o carregamento inicial.

---

## Visualizacoes do Projeto Original (2014) — Analise Detalhada

### 1. Chord Diagram Original

O chord diagram original era significativamente mais sofisticado que a versao atual:

**Estrutura SVG:**
- SVG dinamico (1024px desktop, 460px mobile)
- Centrado com `translate(width/2, height/2)`
- Drop shadow filter (feGaussianBlur stdDeviation=2, offset dx=2 dy=2)

**Arcos de Pontos (9 aneis concentricos):**
- 9 aneis cinza de fundo representando niveis de pontos (1-9)
- Para cada time: arcos preenchidos com cor primaria (cor1) ate o nivel de pontos
- Arco de fundo com cor secundaria (cor2) a 40% opacity
- Labels de numero de pontos nos aneis (so desktop)
- Calculo de raio: `innerRadius + padding + ((i+1) * padding)` para cada camada

**Chords (Conexoes entre times):**
- Um path por confronto usando `d3.svg.chord()`
- Fill-opacity 0.2 por padrao, 0.8 no hover
- Gradientes lineares SVG por chord (cor time1 → cor time2)
- Direcao do gradiente calculada por quadrante (getSide/getHemisphere)
- Labels de gols posicionados no inicio do chord (innerRadius - 20)

**Icones dos Times:**
- Imagens SVG 48x48px posicionadas ao redor do arco
- Transform: `rotate(angle) translate(distance) rotate(-angle)` (mantém upright)
- Distancia: innerRadius + padding * 4.5
- Opacity 0.7 com drop shadow

**Interacoes:**
- Hover em time: destaca arcos, atualiza `$scope.selected`
- Hover em chord: traz para frente, adiciona classe "active", mostra labels de gols, opacity 0.8
- Mouseout: reset visual completo
- Animacao: transitions de 750ms em opacity ao atualizar dados

### 2. Molecule Graph Original (Grafo de Forca)

**Simulacao Fisica:**
- `d3.layout.force()` com charge -200 (repulsao)
- Link distance: `radius(source.size) + radius(target.size) + edge`
- Edge: 16px desktop, 8px mobile
- Escala sqrt para raios: `d3.scale.sqrt().range([0, 6])`
- Tick function continua atualizando posicoes

**Nodes (34 nos):**
- Indices 0-15: times classificados (tamanho base x1)
- Indices 16-23: quartas de final (tamanho x2)
- Indices 24-27: semifinais (tamanho x3)
- Indices 28-29: finalistas (tamanho x4)
- Indice 30: campeao (tamanho x5)
- Indices 31-33: disputa 3o lugar (tamanhos x2/x3)
- Circulos SVG cinza (#ccc) a 50% opacity
- Imagens 96px centradas no circulo
- Labels de gols abaixo do no

**Links (23 conexoes):**
- Linhas SVG brancas a 30% opacity
- Stroke-width = teams[d.source].size / 3
- Coordenadas atualizadas a cada tick

**Interacoes:**
- Drag: `force.drag()` para reposicionar nos
- Hover: tooltip d3-tip com nome do time
- Click: borda vermelha no no selecionado (stroke-width 4)
- Click: mostra dados do jogo no painel lateral

---

## Plano de Implementacao — Novas Visualizacoes

### V1. Chord Diagram 2.0 (D3.js v7)

**Reimplementar o chord diagram fiel ao original, mas com melhorias modernas.**

Substituir o ChordDiagram.tsx atual (SVG manual) por uma implementacao
completa usando `d3-chord`, `d3-arc` e `d3-ribbon`.

**Features:**
- [ ] 9 aneis concentricos de pontos com preenchimento proporcional
- [ ] Arcos com cor primaria do time (opacity layered)
- [ ] Arco de fundo com cor secundaria a 40%
- [ ] Chords com `d3.ribbon()` e gradientes SVG por confronto
- [ ] Icones de bandeira posicionados ao redor (rotate-translate-rotate)
- [ ] Drop shadow filter SVG nos icones
- [ ] Hover em time: destaca arcos + atualiza info panel
- [ ] Hover em chord: traz para frente, opacity 0.8, mostra gols
- [ ] Transition 750ms em opacity ao mudar placares
- [ ] Responsivo: 3 breakpoints (mobile 300px, tablet 400px, desktop 500px)
- [ ] Labels de pontos nos aneis (so desktop)

**Arquivo:** `src/components/Groups/ChordDiagram.tsx` (rewrite)

---

### V2. Molecule Graph 2.0 (D3.js v7 Force)

**Reimplementar o grafo de forca do mata-mata com D3 force simulation moderna.**

Criar uma alternativa visual ao bracket SVG atual, usando `d3-force`
para layout fisico interativo.

**Features:**
- [ ] `d3.forceSimulation()` com forceManyBody (charge -200), forceLink, forceCenter
- [ ] Nodes com tamanho progressivo por fase (R32 → Final, escala 1x-5x)
- [ ] Circulos com opacity 0.5, imagens de bandeira 96px centradas
- [ ] Links com stroke-width proporcional ao tamanho do source
- [ ] Drag interativo nos nodes (d3.drag)
- [ ] Tooltip no hover com nome do time e placar
- [ ] Click no node: borda dourada, expande info do jogo
- [ ] Node do campeao: tamanho 5x, glow dourado pulsante
- [ ] Animacao de entrada: nodes nascem do centro e se espalham
- [ ] Auto-populate reativo a partir dos standings dos grupos
- [ ] Responsivo com recalculo de dimensoes

**Arquivo:** `src/components/Knockout/MoleculeGraph.tsx` (novo)

---

### V3. Globo 3D de Selecoes (React Three Fiber + react-globe.gl)

**Globo terrestre 3D interativo mostrando todas as selecoes classificadas.**

Nova visualizacao que nao existia no original. Um globo rotativo
onde cada selecao e um ponto no mapa, com arcos conectando paises
do mesmo grupo e cilindros representando FIFA ranking.

**Features:**
- [ ] Globo 3D com textura realista da Terra (dia/noite)
- [ ] Pontos cilindrcos 3D por selecao (altura = FIFA ranking)
- [ ] Cor do ponto = cor primaria do time
- [ ] Arcos entre selecoes do mesmo grupo (cor do grupo)
- [ ] Hover: tooltip com nome, confederacao, grupo
- [ ] Click: zoom suave para o pais, expande card de info
- [ ] Rotacao automatica lenta (auto-rotate)
- [ ] Drag para rotacao manual (OrbitControls)
- [ ] Labels flutuantes com codigo FIFA do time
- [ ] Hexagons por confederacao (hexPolygonsData)
- [ ] Transicao animada ao filtrar por confederacao
- [ ] Atmosfera com glow (custom shader)

**Arquivo:** `src/components/Globe/WorldGlobe.tsx` (novo)
**Dependencias:** `react-globe.gl`, `three`

---

### V4. Stadium 3D (React Three Fiber)

**Estadio 3D low-poly como background decorativo dos jogos.**

Cena 3D de um estadio simplificado que reage aos placares.

**Features:**
- [ ] Modelo low-poly de estadio (geometria procedural com Three.js)
- [ ] Campo verde com linhas brancas (PlaneGeometry + textura)
- [ ] Arquibancadas com cor do time mandante
- [ ] Placar digital 3D (Text3D do drei) mostrando resultado
- [ ] Iluminacao: spotlights nas quatro torres
- [ ] Particulas de papel picado ao inserir gol (points buffer)
- [ ] Camera orbital com zoom
- [ ] Versao mini (200x150px) para usar dentro de MatchCard
- [ ] Versao full para pagina de jogo individual

**Arquivo:** `src/components/3d/Stadium.tsx` (novo)

---

### V5. Trophy 3D (React Three Fiber)

**Taca da Copa do Mundo 3D rotativa no Hero.**

Modelo 3D da taca como elemento visual no hero section.

**Features:**
- [ ] Modelo GLTF da taca (ou geometria procedural dourada)
- [ ] Material metalico (MeshStandardMaterial, metalness 0.9, roughness 0.1)
- [ ] Rotacao automatica suave (useFrame)
- [ ] Reflexo do ambiente (environment map)
- [ ] Float animation sutil do drei (Float component)
- [ ] Glow dourado (bloom post-processing)
- [ ] Mouse parallax: taca segue levemente o cursor
- [ ] Responsivo: escala menor no mobile

**Arquivo:** `src/components/3d/Trophy.tsx` (novo)

---

### V6. Particle Field (Canvas 2D / Three.js Points)

**Campo de particulas interativo como background do hero.**

Particulas flutuantes nas cores da copa que reagem ao mouse.

**Features:**
- [ ] 200-500 particulas flutuantes (Points ou instanced mesh)
- [ ] Cores: dourado, branco, azul, verde (palette da copa)
- [ ] Movimento browniano suave (noise function)
- [ ] Conexoes entre particulas proximas (linhas finas)
- [ ] Mouse repulsion: particulas fogem do cursor
- [ ] Scroll: particulas desaceleram e desaparecem conforme desce
- [ ] Performance: requestAnimationFrame com throttle

**Arquivo:** `src/components/3d/ParticleField.tsx` (novo)

---

### V7. Stats Radar Chart (SVG + GSAP)

**Grafico radar para comparacao de selecoes.**

Visualizacao hexagonal/pentagonal comparando atributos dos times.

**Features:**
- [ ] Eixos: Ataque, Defesa, Posse, Passes, Finalizacoes, Faltas
- [ ] SVG com poligono preenchido por time (2 times sobrepostos)
- [ ] Cores dos times com opacity 0.3
- [ ] Hover nos vertices: tooltip com valor
- [ ] Animacao de entrada: poligono cresce do centro
- [ ] Morph animation ao trocar time (d3.interpolate nos vertices)
- [ ] Labels nos eixos
- [ ] Seletor de times com bandeiras

**Arquivo:** `src/components/Stats/RadarChart.tsx` (novo)

---

### V8. Goal Timeline (SVG + GSAP)

**Timeline horizontal de gols de cada jogo.**

Linha do tempo mostrando minuto a minuto os eventos do jogo.

**Features:**
- [ ] Linha horizontal 0-90+ minutos
- [ ] Dots nos momentos de gol (cor do time que marcou)
- [ ] Hover no dot: tooltip com minuto e placar parcial
- [ ] Animacao: dots aparecem sequencialmente com stagger
- [ ] Indicadores de intervalo (45') e acrescimos
- [ ] Icones de cartao (amarelo/vermelho) se disponivel
- [ ] Scroll-triggered na entrada

**Arquivo:** `src/components/Groups/GoalTimeline.tsx` (novo)

---

### V9. Bracket 3D Tunnel (React Three Fiber)

**Versao 3D do bracket com perspectiva de tunel.**

O bracket atual em 2D reimaginado com profundidade 3D, onde cada rodada
recua no espaco criando efeito de tunel convergente.

**Features:**
- [ ] Cada rodada em um plano Z diferente (R32 z=0, R16 z=-2, QF z=-4, etc)
- [ ] Nodes como cards 3D flutuantes (Float do drei)
- [ ] Linhas de conexao como tubes 3D (TubeGeometry)
- [ ] Particulas fluindo pelas conexoes (instanced points animados)
- [ ] Camera animada: scroll controla zoom no tunel
- [ ] Click em node: camera voa ate ele (gsap camera animation)
- [ ] Team path highlight: tubo ilumina com cor do time
- [ ] Bloom nos nodes vencedores
- [ ] Versao alternativa ao bracket SVG (toggle 2D/3D)

**Arquivo:** `src/components/Knockout/Bracket3D.tsx` (novo)

---

### V10. Animated Counters (GSAP)

**Contadores animados para estatisticas.**

Numeros que contam de 0 ate o valor final com easing.

**Features:**
- [ ] CountUp animation com GSAP (0 → valor, duration 2s, ease power2.out)
- [ ] Formato: numeros inteiros, decimais, porcentagens
- [ ] Scroll-triggered (inicia quando entra na viewport)
- [ ] Roda so uma vez (once: true)
- [ ] Separador de milhares
- [ ] Prefixo/sufixo configuravel

**Arquivo:** `src/components/ui/AnimatedCounter.tsx` (novo)

---

### V11. Group Galaxy View (React Three Fiber)

**Visualizacao alternativa dos grupos como sistemas solares 3D.**

Cada grupo e um cluster onde times sao "planetas" orbitando.

**Features:**
- [ ] 12 clusters (1 por grupo) dispostos em grid ou circulo
- [ ] Times como esferas orbitando o centro (bandeira como textura)
- [ ] Raio da orbita = posicao na tabela (1o mais perto)
- [ ] Tamanho da esfera = pontos acumulados
- [ ] Velocidade orbital diferente por time
- [ ] Rastro de orbita semi-transparente
- [ ] Click no "planeta": zoom + info
- [ ] Estrela central com cor media do grupo
- [ ] Toggle 2D/3D nos group cards

**Arquivo:** `src/components/Groups/GroupGalaxy.tsx` (novo)

---

## Melhorias em Componentes Existentes

### E1. Hero Section
- [ ] Substituir circulos decorativos por ParticleField (V6)
- [ ] Adicionar Trophy 3D (V5) ao lado do titulo
- [ ] Background: gradient + particles em vez de divs blur

### E2. GroupCard
- [ ] Toggle entre ChordDiagram 2D (V1) e GroupGalaxy 3D (V11)
- [ ] GoalTimeline (V8) dentro da aba de jogos
- [ ] AnimatedCounter (V10) nos numeros da tabela

### E3. KnockoutStage
- [ ] Toggle entre Bracket SVG e Bracket3D (V9)
- [ ] MoleculeGraph (V2) como terceira opcao de visualizacao
- [ ] Confetti aprimorado ao definir campeao (3D particles)

### E4. TeamGrid
- [ ] Adicionar secao WorldGlobe (V3) acima do grid
- [ ] RadarChart (V7) no modal de detalhe do time
- [ ] Cards com hover 3D (tilt effect via transform perspective)

### E5. PredictionSection
- [ ] AnimatedCounter (V10) nas estatisticas
- [ ] Confetti 3D (Three.js points) ao selecionar campeao
- [ ] ShareCard com preview do globo

---

## Prioridades de Implementacao

### Fase 1 — Graficos do Original (alta prioridade)
1. **V1** - Chord Diagram 2.0 (rewrite com D3 v7)
2. **V2** - Molecule Graph 2.0 (D3 force simulation)

### Fase 2 — 3D Hero e Atmosfera
3. **V6** - Particle Field (hero background)
4. **V5** - Trophy 3D (hero element)
5. **V10** - Animated Counters (global)

### Fase 3 — Visualizacoes Novas
6. **V3** - Globo 3D de Selecoes
7. **V7** - Stats Radar Chart
8. **V8** - Goal Timeline

### Fase 4 — 3D Avancado
9. **V9** - Bracket 3D Tunnel
10. **V11** - Group Galaxy View
11. **V4** - Stadium 3D

---

## Notas Tecnicas

### Performance
- Todas as cenas 3D devem usar `React.lazy` + `Suspense`
- Three.js nao carrega no bundle principal — so quando o user scrollar ate a secao
- Canvas 3D com `dpr={[1, 1.5]}` para balancear qualidade vs performance
- `frameloop="demand"` nas cenas estaticas (so renderiza quando muda)
- Dispose de geometrias/materiais/texturas no cleanup
- Mobile: reduzir particulas, desabilitar post-processing, usar LOD

### Acessibilidade
- Toda visualizacao 3D tem fallback 2D
- `aria-label` descritivo nos canvas
- Dados acessiveis em tabela HTML escondida (sr-only)
- Reducao de movimento: `prefers-reduced-motion` desabilita animacoes 3D

### Bundle
- Three.js: importar apenas modulos necessarios (tree-shaking)
- react-globe.gl: lazy load completo (so carrega na secao de selecoes)
- Target total de JS: < 500KB gzipped (atualmente 132KB)
- CSS permanece < 10KB gzipped
