# Copa.Guru

Plataforma interativa de palpites para a Copa do Mundo FIFA 2026. Faça previsões de placares, explore visualizações de dados e compartilhe seus palpites com amigos.

**Acesse em [copa.guru](https://copa.guru)**

<p align="center">
  <img src="frontend/public/images/logo.webp" alt="Copa.Guru Logo" width="300" />
</p>

## Funcionalidades

### Fase de Grupos
- Palpites de placares para todas as 48 seleções em 12 grupos
- Tabela de classificação com critérios de desempate automáticos
- Diagrama de acordes mostrando relações entre partidas
- Visualização 3D orbital dos grupos

### Mata-Mata
- Chave completa das oitavas de final até a final
- Visualização interativa em grafo/rede
- Renderização 3D do chaveamento
- Previsão do campeão

### Mapa Mundial
- Globo 3D interativo exibindo todas as seleções classificadas
- Seleções identificadas por cor de grupo

### Seleções
- Grade com todas as 48 seleções classificadas
- Filtro por confederação (UEFA, CONMEBOL, CONCACAF, CAF, AFC, OFC)

### Estatísticas & Comparação
- Gráficos radar comparando atributos das seleções (Ataque, Defesa, Posse, Passes, Finalizações, Faltas)

### Palpites
- Cards de palpites compartilháveis com estatísticas das partidas

### Easter Egg
- Konami Code no desktop (&#8593;&#8593;&#8595;&#8595;&#8592;&#8594;&#8592;&#8594;BA) ou 5 toques rápidos no logo no mobile

## Tech Stack

| Camada         | Tecnologia                                       |
|----------------|--------------------------------------------------|
| Framework      | React 19 + TypeScript 5.7                        |
| Build          | Vite 6 (SWC plugin) + Bun                        |
| Estilização    | Tailwind CSS 4                                   |
| 3D             | Three.js + React Three Fiber + Drei              |
| Visualizações  | D3.js (chord, force, scale, shape, drag)         |
| Animações      | GSAP + ScrollTrigger                             |
| Rotas          | React Router DOM 7                               |
| Linting        | Biome                                            |
| Testes         | Vitest                                           |

## Internacionalização

Três idiomas suportados, cada um com seu próprio prefixo de rota:

| Idioma     | Rota  |
|------------|-------|
| Português  | `/`   |
| Inglês     | `/en` |
| Espanhol   | `/es` |

As traduções ficam em `frontend/src/i18n/translations.ts`. Arquivos HTML por idioma com meta tags apropriadas são gerados no build via `frontend/scripts/generate-html.ts`.

## Como Rodar

### Pré-requisitos

- [Bun](https://bun.sh/) >= 1.x

### Instalar dependências

```bash
cd frontend
bun install
```

### Servidor de desenvolvimento

```bash
make frontend-server
# ou
cd frontend && bun run dev --port 5006
```

Acesse [http://localhost:5006](http://localhost:5006).

### Build

```bash
make frontend-build
# ou
cd frontend && bun run build
```

### Lint & Formatação

```bash
cd frontend && bun run check
```

### Testes

```bash
cd frontend && bun run test
```

## Estrutura do Projeto

```
copa.guru/
├── frontend/
│   ├── public/
│   │   ├── data/           # worldcup2026.json
│   │   └── flags/          # Bandeiras SVG dos países
│   ├── scripts/            # Geração de HTML no build
│   ├── src/
│   │   ├── assets/         # Ícones e assets estáticos
│   │   ├── components/     # Componentes React por feature
│   │   ├── data/           # Fetch de dados e constantes
│   │   ├── hooks/          # Hooks customizados
│   │   ├── i18n/           # Traduções e contexto de locale
│   │   ├── lib/            # Utilitários
│   │   └── types/          # Definições de tipos TypeScript
│   ├── biome.json
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── backend/                # (em breve)
└── Makefile
```

## Deploy

Hospedado em **AWS S3 + CloudFront**.

```bash
make frontend-deploy
```

Faz o build, sincroniza com o S3 com headers de cache apropriados (1 ano para assets imutáveis, 1 hora para HTML/metadados) e invalida a distribuição do CloudFront.

CI/CD via GitHub Actions: PRs disparam build de verificação, merges na `main` disparam deploy automático.
