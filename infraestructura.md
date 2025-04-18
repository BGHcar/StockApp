stock-analyzer/
│
├── backend/                     # Código fuente del backend en Go
│   ├── api/                    # Controladores y rutas HTTP
│   │   └── handlers.go
│   ├── client/                 # Cliente para consumir la API externa
│   │   └── apiConsumer.go
│   ├── config/                 # Configuración de la aplicación
│   │   └── config.go
│   ├── db/                     # Lógica de conexión y queries a la base de datos
│   │   └── db.go
│   ├── factory/                # Fábricas para crear instancias de servicios/repositorios
│   │   └── factory.go
│   ├── interfaces/             # Definición de interfaces
│   │   └── interfaces.go
│   ├── models/                 # Estructuras de datos y modelos
│   │   └── stock.go
│   ├── repositories/           # Operaciones de persistencia
│   │   └── stock_repository.go
│   ├── services/               # Lógica de negocio
│   │   └── stock_service.go
│   ├── main.go                 # Punto de entrada de la API
│   ├── go.mod                  # Dependencias de Go
│   └── go.sum                  # Dependencias de Go
│
frontend/
├── public/                   # Archivos estáticos
│   └── favicon.ico           # Ícono de la aplicación
│
├── src/                      # Código fuente principal
│   ├── components/           # Componentes reutilizables
│   │   ├── StockTable.vue    # Tabla para mostrar datos de stocks
│   │   ├── StockSearch.vue   # Buscador de stocks
│   │   ├── StockRecommendations.vue # Componente de recomendaciones
│   │   │
│   │   ├── base/             # Componentes base fundamentales
│   │   │   ├── BaseTable.vue # Componente base para tablas
│   │   │   └── BaseScroll.vue # Componente con scroll personalizado
│   │   │
│   │   └── icons/            # Componentes de íconos SVG
│   │
│   ├── pages/                # Vistas/páginas principales
│   │   └── Home.vue          # Página principal de la aplicación
│   │
│   ├── router/               # Configuración de rutas
│   │   └── index.ts          # Definición de rutas de la aplicación
│   │
│   ├── services/             # Servicios para API y lógica de negocio
│   │   ├── stockService.ts   # Servicio para consumir la API de stocks
│   │   └── analysisService.ts # Servicio para análisis de datos
│   │
│   ├── stores/               # Estado global con Pinia
│   │   └── stockStore.ts     # Store para manejo de datos de stocks
│   │
│   ├── types/                # Definiciones de tipos TypeScript
│   │   ├── stock.ts          # Interfaces para datos de stocks
│   │   └── table.ts          # Interfaces para componentes de tabla
│   │
│   ├── App.vue               # Componente raíz de la aplicación
│   └── main.ts               # Punto de entrada - inicialización
│
├── .env                      # Variables de entorno (APIs, etc.)
├── tsconfig.json             # Configuración de TypeScript
├── tsconfig.app.json         # Config TypeScript específica para la app
├── tsconfig.node.json        # Config TypeScript para entorno Node
├── vite.config.ts            # Configuración del bundler Vite
├── tailwind.config.js        # Configuración de TailwindCSS
├── index.html                # Plantilla HTML principal
└── package.json              # Dependencias y scripts
│
├── infra/                      # Infraestructura como código (opcional)
│   ├── cockroachdb/
│   ├── docker/
│   └── main.tf
│
├── docs/                       # Documentación del sistema, API, esquemas, etc.
│   └── architecture.pdf
│
├── .env                        # Variables de entorno del proyecto
├── .gitignore                  # Archivos a ignorar por Git
└── README.md                   # Documentación principal del proyecto
