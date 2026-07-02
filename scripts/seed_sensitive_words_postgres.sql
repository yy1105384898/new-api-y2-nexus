-- Seed public sensitive-word options for PostgreSQL (production new-api).
-- Safe to re-run: ON CONFLICT refreshes toggles and the baseline word list.
--
-- Usage (inside newapi-postgres container):
--   psql -U root -d new-api -f /path/to/seed_sensitive_words_postgres.sql
--
-- Categories align with video/image content moderation guide:
-- minors, sexual/nudity, violence, drugs, gambling, IP/copyright, face-swap, and common upstream triggers.

INSERT INTO options (key, value) VALUES
  ('CheckSensitiveEnabled', 'true'),
  ('CheckSensitiveOnPromptEnabled', 'true'),
  ('StopOnSensitiveEnabled', 'true'),
  ('SensitiveWords', $words$nude
naked
nsfw
porn
pornography
pornographic
xxx
hentai
erotic
seductive
cleavage
onsen
ryokan
topless
undress
stripper
prostitute
zombie
corpse
fallen bodies
child porn
loli
shota
pedophile
paedophile
underage sex
minor sex
bestiality
zoophilia
rape
gangbang
gore
bloodbath
bodies
torture
beheading
decapitation
suicide
self-harm
self harm
terrorism
terrorist attack
school shooting
mass shooting
heroin
cocaine
methamphetamine
lsd
marijuana
gambling
casino
色情
裸体
淫秽
性交
做爱
口交
肛交
强奸
轮奸
幼女
儿童色情
恋童
兽交
血腥
自杀
自残
恐怖主义
爆炸物
炸弹
雷管
爆炸
丧尸
尸
制毒
冰毒
海洛因
可卡因
赌博
换脸
锁定人脸
人脸锁定
半透明
情色
裸照
刀斧手
血渍
疯狂动物城
朱迪
尼克
迪士尼公主$words$)
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;
