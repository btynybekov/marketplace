BEGIN;

-- UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1) Пользователи
CREATE TABLE IF NOT EXISTS app_user (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  phone         TEXT UNIQUE,
  display_name  TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2) Категории (иерархия)
CREATE TABLE IF NOT EXISTS category (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  parent_id     UUID REFERENCES category(id) ON DELETE RESTRICT,
  name          TEXT NOT NULL,
  slug          TEXT NOT NULL,
  path          TEXT NOT NULL,          -- например: electronics/phones/smartphones
  is_active     BOOLEAN NOT NULL DEFAULT TRUE,
  sort_order    INT NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(parent_id, name),
  UNIQUE(path),
  UNIQUE(slug)
);
CREATE INDEX IF NOT EXISTS idx_category_parent ON category(parent_id);

-- 3) Бренды
CREATE TABLE IF NOT EXISTS brand (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name          TEXT NOT NULL UNIQUE,
  slug          TEXT NOT NULL UNIQUE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 4) Товары (каталог/модель)
CREATE TABLE IF NOT EXISTS product (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  category_id   UUID NOT NULL REFERENCES category(id) ON DELETE RESTRICT,
  brand_id      UUID REFERENCES brand(id),
  model         TEXT,                    -- напр. "iPhone 13"
  title         TEXT NOT NULL,           -- заголовок модели
  specs         JSONB NOT NULL DEFAULT '{}'::jsonb,  -- характеристики модели
  is_active     BOOLEAN NOT NULL DEFAULT TRUE,
  created_by    UUID REFERENCES app_user(id),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_product_category ON product(category_id);
CREATE INDEX IF NOT EXISTS idx_product_brand ON product(brand_id);
CREATE INDEX IF NOT EXISTS idx_product_specs_gin ON product USING GIN (specs);

-- 5) Медиа для товаров (1:N)
CREATE TABLE IF NOT EXISTS product_media (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  product_id    UUID NOT NULL REFERENCES product(id) ON DELETE CASCADE,
  url           TEXT NOT NULL,
  type          TEXT NOT NULL DEFAULT 'image',      -- image|video|doc
  is_cover      BOOLEAN NOT NULL DEFAULT FALSE,
  alt           TEXT,
  width         INT,
  height        INT,
  bytes         INT,
  hash          TEXT,
  variants      JSONB NOT NULL DEFAULT '[]'::jsonb, -- [{label:'thumb',url:'...'}]
  sort_order    INT NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_product_media_product ON product_media(product_id, sort_order);
CREATE UNIQUE INDEX IF NOT EXISTS uniq_product_media_cover
  ON product_media(product_id)
  WHERE is_cover = TRUE;

-- 6) Объявления (конкретные предложения)
CREATE TABLE IF NOT EXISTS listing (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  seller_id     UUID NOT NULL REFERENCES app_user(id) ON DELETE RESTRICT,
  product_id    UUID REFERENCES product(id) ON DELETE SET NULL,
  category_id   UUID NOT NULL REFERENCES category(id) ON DELETE RESTRICT, -- дублируем для быстрого фильтра
  title         TEXT NOT NULL,
  description   TEXT NOT NULL,
  price_amount  NUMERIC(12,2) NOT NULL,
  currency_code CHAR(3) NOT NULL DEFAULT 'KGS',
  condition     TEXT NOT NULL CHECK (condition IN ('new','used')),
  location_text TEXT,
  attrs         JSONB NOT NULL DEFAULT '{}'::jsonb, -- частные атрибуты объявления
  status        TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','paused','sold','deleted')),
  expires_at    TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_listing_category ON listing(category_id, status);
CREATE INDEX IF NOT EXISTS idx_listing_price ON listing(currency_code, price_amount);
CREATE INDEX IF NOT EXISTS idx_listing_attrs_gin ON listing USING GIN (attrs);

-- 7) Медиа для объявлений (1:N)
CREATE TABLE IF NOT EXISTS listing_media (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  listing_id    UUID NOT NULL REFERENCES listing(id) ON DELETE CASCADE,
  url           TEXT NOT NULL,
  type          TEXT NOT NULL DEFAULT 'image',
  is_cover      BOOLEAN NOT NULL DEFAULT FALSE,
  alt           TEXT,
  width         INT,
  height        INT,
  bytes         INT,
  hash          TEXT,
  variants      JSONB NOT NULL DEFAULT '[]'::jsonb,
  sort_order    INT NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_listing_media_listing ON listing_media(listing_id, sort_order);
CREATE UNIQUE INDEX IF NOT EXISTS uniq_listing_media_cover
  ON listing_media(listing_id)
  WHERE is_cover = TRUE;

-- 8) Диалоги с ИИ
CREATE TABLE IF NOT EXISTS conversation (
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id         UUID REFERENCES app_user(id) ON DELETE SET NULL,
  channel         TEXT NOT NULL DEFAULT 'web' CHECK (channel IN ('web','mobile','tg','api')),
  intent          TEXT CHECK (intent IN ('smalltalk','search','post')),
  summary         TEXT,
  slots           JSONB NOT NULL DEFAULT '{}'::jsonb,
  state           JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_message_at TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_conversation_user ON conversation(user_id);
CREATE INDEX IF NOT EXISTS idx_conversation_updated ON conversation(updated_at DESC);

-- 9) Сообщения
CREATE TABLE IF NOT EXISTS message (
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  conversation_id UUID NOT NULL REFERENCES conversation(id) ON DELETE CASCADE,
  role            TEXT NOT NULL CHECK (role IN ('user','assistant','system')),
  content         TEXT NOT NULL,
  meta            JSONB NOT NULL DEFAULT '{}'::jsonb, -- {model, usage, finish_reason}
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_message_conv_created ON message(conversation_id, created_at);

-- 10) Логи запусков ИИ
CREATE TABLE IF NOT EXISTS ai_run (
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  conversation_id UUID REFERENCES conversation(id) ON DELETE SET NULL,
  kind            TEXT NOT NULL CHECK (kind IN ('intent','parse','search_rank','seller_draft')),
  model           TEXT NOT NULL,
  temperature     NUMERIC(3,2) NOT NULL,
  prompt          TEXT NOT NULL,
  response        TEXT NOT NULL,
  usage           JSONB NOT NULL DEFAULT '{}'::jsonb,
  status          TEXT NOT NULL DEFAULT 'ok' CHECK (status IN ('ok','error')),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_ai_run_conversation ON ai_run(conversation_id, created_at);

-- 11) Логи поисковых запросов к бэкенду
CREATE TABLE IF NOT EXISTS search_request (
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  conversation_id UUID REFERENCES conversation(id) ON DELETE SET NULL,
  filters         JSONB NOT NULL,
  result_count    INT,
  filter_url      TEXT,
  backend_meta    JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_search_request_conv ON search_request(conversation_id, created_at);

COMMIT;

