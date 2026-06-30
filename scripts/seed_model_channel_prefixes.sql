-- Seed channel registration prefixes for model public-name stripping.
-- Run once after deploying model_channel_prefixes migration.
-- SQLite example:
--   sqlite3 /data/one-api.db < scripts/seed_model_channel_prefixes.sql

INSERT OR IGNORE INTO model_channel_prefixes (prefix, note, enabled, sort_order, created_time, updated_time)
VALUES
  ('119337-', 'api.119337.xyz Grok Video', 1, 10, strftime('%s','now'), strftime('%s','now')),
  ('aini-', 'Aini 生图', 1, 20, strftime('%s','now'), strftime('%s','now')),
  ('byte-', '字节 Gemini 生图', 1, 30, strftime('%s','now'), strftime('%s','now')),
  ('ctlove-', 'CTLove 字节火山 Seedance', 1, 40, strftime('%s','now'), strftime('%s','now')),
  ('czeq-', 'Czeq 生图', 1, 50, strftime('%s','now'), strftime('%s','now')),
  ('go2api-', 'Go2API 生图', 1, 60, strftime('%s','now'), strftime('%s','now')),
  ('gz-', 'GZ / 冠臻 Seedance', 1, 70, strftime('%s','now'), strftime('%s','now')),
  ('happyhorse-', 'HappyHorse 视频', 1, 80, strftime('%s','now'), strftime('%s','now')),
  ('niming-', '匿名 Gemini 生图', 1, 90, strftime('%s','now'), strftime('%s','now')),
  ('oairegbox-', 'OAIREGBox', 1, 100, strftime('%s','now'), strftime('%s','now')),
  ('tengda-', '腾达 Geeknow Veo（td.geeknow.top）', 1, 105, strftime('%s','now'), strftime('%s','now')),
  ('yunwu-', '云雾 / Apifox 聚合', 1, 110, strftime('%s','now'), strftime('%s','now')),
  ('zeabur-', 'Zeabur 托管 Gemini/GLM', 1, 120, strftime('%s','now'), strftime('%s','now'));
