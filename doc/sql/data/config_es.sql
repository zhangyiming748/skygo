DROP sequence IF EXISTS "public"."config_es_id_seq" ;

CREATE SEQUENCE "public"."config_es_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for config_es
-- ----------------------------
DROP TABLE IF EXISTS "public"."config_es";
CREATE TABLE "public"."config_es" (
  "id" int4 NOT NULL DEFAULT nextval('config_es_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "brokers" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "version" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "username" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "password" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "category" int2 NOT NULL DEFAULT 2,
  "auth_type" int2 NOT NULL DEFAULT 2,
  "create_time" int4 NOT NULL
)
;
COMMENT ON COLUMN "public"."config_es"."id" IS '自增主键';
COMMENT ON COLUMN "public"."config_es"."name" IS 'es集群名称';
COMMENT ON COLUMN "public"."config_es"."brokers" IS 'es集群地址，格式ip:port,ip:port';
COMMENT ON COLUMN "public"."config_es"."version" IS 'es集群版本号，比如6.5';
COMMENT ON COLUMN "public"."config_es"."username" IS '用户名';
COMMENT ON COLUMN "public"."config_es"."password" IS '密码';
COMMENT ON COLUMN "public"."config_es"."category" IS '分类：1预定义，2自定义';
COMMENT ON COLUMN "public"."config_es"."auth_type" IS '认证方式，1无 2用户密码';
COMMENT ON COLUMN "public"."config_es"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."config_es" IS 'kafka数据源配置';

-- ----------------------------
-- Records of config_es
-- ----------------------------
INSERT INTO "public"."config_es" VALUES (1, 'kafka集群1', '127.0.0.1:8911', '', 'soc', 'gwc-2345', 2, 1, 0);

-- ----------------------------
-- Primary Key structure for table config_es
-- ----------------------------
ALTER TABLE "public"."config_es" ADD CONSTRAINT "_copy_1_copy_1" PRIMARY KEY ("id");
