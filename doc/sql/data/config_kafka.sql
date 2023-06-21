DROP sequence IF EXISTS "public"."config_kafka_id_seq";

CREATE SEQUENCE "public"."config_kafka_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

SELECT setval('"public"."config_kafka_id_seq"', 1, true);

-- ----------------------------
-- Table structure for config_kafka
-- ----------------------------
DROP TABLE IF EXISTS "public"."config_kafka";
CREATE TABLE "public"."config_kafka" (
  "id" int4 NOT NULL DEFAULT nextval('config_kafka_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "brokers" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "version" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "character_encoding" int2 NOT NULL DEFAULT 1,
  "auth_type" int2 NOT NULL DEFAULT 0,
  "key_tab" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "krb" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "server_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "user_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "category" int2 NOT NULL DEFAULT 2,
  "create_time" int4 NOT NULL
)
;
COMMENT ON COLUMN "public"."config_kafka"."id" IS '自增ID';
COMMENT ON COLUMN "public"."config_kafka"."name" IS 'kafka集群名称';
COMMENT ON COLUMN "public"."config_kafka"."brokers" IS 'kafka集群地址，格式ip:port,ip:port';
COMMENT ON COLUMN "public"."config_kafka"."version" IS 'kafka集群版本号，比如1.0.0';
COMMENT ON COLUMN "public"."config_kafka"."character_encoding" IS '字符编码，1utf-8';
COMMENT ON COLUMN "public"."config_kafka"."auth_type" IS '权限校验，0无 1kerberos';
COMMENT ON COLUMN "public"."config_kafka"."key_tab" IS 'keytab文件路径';
COMMENT ON COLUMN "public"."config_kafka"."krb" IS 'krb配置文件路径';
COMMENT ON COLUMN "public"."config_kafka"."server_name" IS '访问服务器名称';
COMMENT ON COLUMN "public"."config_kafka"."url" IS '访问域名';
COMMENT ON COLUMN "public"."config_kafka"."user_name" IS '用户名';
COMMENT ON COLUMN "public"."config_kafka"."category" IS '分类：1预定义，2自定义';
COMMENT ON COLUMN "public"."config_kafka"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."config_kafka" IS '存储配置kafka
每条记录对应一个kafka集群的配置';

-- ----------------------------
-- Records of config_kafka
-- ----------------------------
INSERT INTO "public"."config_kafka" VALUES (1, 'iisop集群', '10.99.99.60:9092,10.99.99.60:9093,10.99.99.60:9094', '2.1.1', 1, 0, '', '', '', '', '', 2, 1614586876);

-- ----------------------------
-- Primary Key structure for table config_kafka
-- ----------------------------
ALTER TABLE "public"."config_kafka" ADD CONSTRAINT "_copy_1" PRIMARY KEY ("id");
